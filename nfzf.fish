function nfzf
  nf-list | fzf -m --ansi --preview '
      set parts (string split " -> " -- {})
      set class $parts[1]

      set right (string split " | " -- $parts[2])
      set hex $right[1]
      set char $right[2]

      echo -e "\e[1;32mName:\e[0m   $class"
      echo -e "\e[1;33mSymbol:\e[0m $char"
      echo -e "\e[1;36mHex:\e[0m    $hex\n"
      echo -e "$char $char $char"
    '
end

