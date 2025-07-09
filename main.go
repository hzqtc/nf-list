package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type nfIcon struct {
	Name string `json:"class"`
	Hex  string `json:"hex"`
	R    rune   `json:"char"`
}

var (
	search     []string
	group      []string
	random     bool
	compact    bool
	outputJson bool
	initShell  string
)

const (
	fishFunc = `function nfzf
		nf-list $argv | fzf --style minimal -m --ansi --preview '
			set parts (string split " -> " -- {})
			set class $parts[1]

			set right (string split " | " -- $parts[2])
			set hex $right[1]
			set char $right[2]

			echo -e "Name:   $class"
			echo -e "Symbol: $char"
			echo -e "Hex:    $hex\n"
			echo -e "$char $char $char"
		'
	end`
	bashFunc = `nfzf() {
		nf-list "$@" | fzf --style minimal -m --ansi --preview '
			line="{}"
			class="${line%% -> *}"
			right="${line#* -> }"
			hex="${right%% | *}"
			char="${right#* | }"

			echo -e "Name:   $class"
			echo -e "Symbol: $char"
			echo -e "Hex:    $hex\n"
			echo -e "$char $char $char"
		'
	}`
)

const (
	url         = "https://raw.githubusercontent.com/ryanoasis/nerd-fonts/refs/heads/master/css/nerd-fonts-generated.css"
	cacheTTL    = 24 * time.Hour
	cacheFile   = "nerd-fonts-generated.css"
	cacheSubdir = "nf-list"
)

var cacheDir = func() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".cache"
	}
	return filepath.Join(home, ".cache")
}()

var rootCmd = &cobra.Command{
	Use:   "nf-list",
	Short: "List Nerd Font glyphs",
	Run: func(cmd *cobra.Command, args []string) {
		if initShell != "" {
			switch initShell {
			case "fish":
				fmt.Println(fishFunc)
			case "bash", "zsh":
				fmt.Println(bashFunc)
			default:
				fmt.Fprintf(os.Stderr, "Unsupported shell: %s\n", initShell)
			}
			return
		}

		icons := loadNfIcons()
		if random && len(icons) > 0 {
			rand.Seed(time.Now().UnixNano())
			icons = []nfIcon{icons[rand.Intn(len(icons))]}
		} else {
			if len(search) > 0 {
				icons = filterIcons(icons, search)
			}
			if len(group) > 0 {
				icons = filterGroup(icons, group)
			}
		}
		printIcons(icons)
	},
}

func printIcons(icons []nfIcon) {
	if outputJson {
		_ = json.NewEncoder(os.Stdout).Encode(icons)
	} else {
		for _, icon := range icons {
			if compact {
				fmt.Printf("%c\n", icon.R)
			} else {
				fmt.Printf("%s -> %s | %c\n", icon.Name, icon.Hex, icon.R)
			}
		}
	}
}

func filterIcons(icons []nfIcon, keywords []string) []nfIcon {
	var filtered []nfIcon
	for _, icon := range icons {
		// Needs to match all keywords
		matches := true
		for _, kw := range keywords {
			if !strings.Contains(icon.Name, kw) && !strings.Contains(icon.Hex, kw) {
				matches = false
				continue
			}
		}
		if matches {
			filtered = append(filtered, icon)
		}
	}
	return filtered
}

func filterGroup(icons []nfIcon, prefixes []string) []nfIcon {
	var filtered []nfIcon
	for _, icon := range icons {
		for _, prefix := range prefixes {
			if strings.HasPrefix(icon.Name, "nf-"+prefix+"-") {
				filtered = append(filtered, icon)
			}
		}
	}
	return filtered
}

func getCachePath() string {
	return filepath.Join(cacheDir, cacheSubdir, cacheFile)
}

func isCacheExpired(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return true
	}
	return time.Since(info.ModTime()) > cacheTTL
}

func fetchCSS() string {
	path := getCachePath()

	if _, err := os.Stat(path); err == nil && !isCacheExpired(path) {
		data, err := os.ReadFile(path)
		if err == nil {
			return string(data)
		}
	}

	resp, err := http.Get(url)
	if err != nil {
		panic("failed to fetch CSS: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("failed to read CSS body: " + err.Error())
	}

	// write to cache
	cacheDir := filepath.Dir(path)
	os.MkdirAll(cacheDir, 0755)
	os.WriteFile(path, body, 0644)

	return string(body)
}

func runeFromHex(hex string) (rune, error) {
	val, err := strconv.ParseInt(hex, 16, 32)
	return rune(val), err
}

func loadNfIcons() []nfIcon {
	css := fetchCSS()
	// Match lines like `.nf-xyz:before { content: "\abcd"; }`
	re := regexp.MustCompile(`(?m)\.(?P<class>[^\s:]+):before\s*\{\s*content:\s*"\\(?P<hex>[a-fA-F0-9]+)";\s*\}`)
	matches := re.FindAllStringSubmatch(css, -1)

	icons := make([]nfIcon, 0, len(matches))
	for _, match := range matches {
		class := match[1]
		hex := match[2]
		r, err := runeFromHex(hex)
		if err == nil {
			icons = append(icons, nfIcon{class, hex, r})
		}
	}

	return icons
}

func init() {
	rootCmd.Flags().StringArrayVar(&search, "search", []string{}, "Filter icons by substring")
	rootCmd.Flags().StringArrayVar(
		&group,
		"group",
		[]string{},
		"Filter icons by group prefix (cod, custom, dev, extra, fa, fae, iec, indent, indentation, linux, md, oct, pl, ple, pom, seti, weather)",
	)
	rootCmd.Flags().BoolVar(&random, "random", false, "Output one random icon")
	rootCmd.Flags().BoolVar(&compact, "compact", false, "Only print the icon character")
	rootCmd.Flags().BoolVar(&outputJson, "json", false, "Output as JSON")
	rootCmd.Flags().StringVar(&initShell, "init", "", "Print shell integration with fzf for [fish|bash|zsh]")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
