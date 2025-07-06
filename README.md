# nf-list

A simple command-line tool to list all Nerd Font icons with their names and codepoints.

The tool fetches the Nerd Fonts CSS file and parses it to extract the icon information. The CSS file is cached locally to avoid repeated downloads.

![](https://raw.github.com/hzqtc/nf-list/master/demo.gif)

## Usage

To run the tool, use the following command:

```sh
make install
nf-list
```

This will output a list of all Nerd Font icons in the following format:

```
<icon-name> -> <hex-code> | <icon>
```

For example:

```
nf-dev-git -> f1d3 | 
```

## fzf integration

The repository also includes wrapper scripts for `bash`, `zsh`, and `fish` to integrate with `fzf` for interactive searching.

*   `nfzf.bash`
*   `nfzf.zsh`
*   `nfzf.fish`

These scripts allow you to search for Nerd Font icons and copy the selected icon to your clipboard.
