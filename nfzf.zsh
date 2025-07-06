function nfzf() {
  nf-list | fzf --style minimal -m --ansi --preview '
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
}
