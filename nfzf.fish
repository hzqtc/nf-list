function nfzf
  nf-list | fzf --style minimal -m --ansi --preview '
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
end

