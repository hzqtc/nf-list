package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
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

func main() {
	css := fetchCSS()
	// Match lines like `.nf-xyz:before { content: "\abcd"; }`
	re := regexp.MustCompile(`(?m)\.(?P<class>[^\s:]+):before\s*\{\s*content:\s*"\\(?P<hex>[a-fA-F0-9]+)";\s*\}`)
	matches := re.FindAllStringSubmatch(css, -1)
	for _, match := range matches {
		class := match[1]
		hex := match[2]
		r, err := runeFromHex(hex)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hex parsing failed: %s -> %s\n", class, hex)
		} else {
			fmt.Printf("%s -> %s | %c\n", class, hex, r)
		}
	}
}
