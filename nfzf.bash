nfzf() {
  nf-list | fzf -m --ansi --preview '
    line="{}"
    class="${line%% -> *}"
    right="${line#* -> }"
    hex="${right%% | *}"
    char="${right#* | }"

    echo -e "\e[1;32mName:\e[0m   $class"
    echo -e "\e[1;33mSymbol:\e[0m $char"
    echo -e "\e[1;36mHex:\e[0m    $hex\n"
    echo -e "$char $char $char"
  '
}