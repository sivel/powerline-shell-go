# powerline-shell-go

Attempted fork of [powerline-shell](https://github.com/milkbikis/powerline-shell) into [Go](http://golang.org/)

This application does not cover all features of powerline-shell, only those that I currently use, and it is not configurable.

## Usage

Install the binary with

    go get github.com/sivel/powerline-shell-go
    go install github.com/sivel/powerline-shell-go

### Bash

Install powerline-shell-go and add the following to your `~/.bashrc`

    function _update_ps1() {
       export PS1="$(powerline-shell-go bash $? 2> /dev/null)"
    }

    export PROMPT_COMMAND="_update_ps1; $PROMPT_COMMAND"

### Zsh

Install powerline-shell-go and add the following to your `~/.zshrc`

    function powerline_precmd() {
      export PS1="$(powerline-shell-go zsh $? 2> /dev/null)"
    }

    function install_powerline_precmd() {
      for s in "${precmd_functions[@]}"; do
        if [ "$s" = "powerline_precmd" ]; then
          return
        fi
      done
      precmd_functions+=(powerline_precmd)
    }

    install_powerline_precmd

## Performance

```
$ time ~/git/milkbikis/powerline-shell/powerline-shell.py > /dev/null
real    0m0.092s
user    0m0.027s
sys     0m0.046s
```

```
$ time ~/go/src/github.com/sivel/powerline-shell/powerline-shell > /dev/null
real    0m0.007s
user    0m0.002s
sys     0m0.004s
```
