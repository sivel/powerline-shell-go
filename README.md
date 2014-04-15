# powerline-shell-go

Attempted fork of [powerline-shell](https://github.com/milkbikis/powerline-shell) into [Go](http://golang.org/)

This application does not cover all features of powerline-shell, only those that I currently use, and it is not configurable.

## Performance

```shell
$ time ~/git/milkbikis/powerline-shell/powerline-shell.py > /dev/null
real	0m0.092s
user	0m0.027s
sys	0m0.046s
```

```shell
$ time ~/go/src/github.com/sivel/powerline-shell/powerline-shell > /dev/null
real	0m0.007s
user	0m0.002s
sys	0m0.004s
```
