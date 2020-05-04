# clipd

![Go](https://github.com/nicolomaioli/clipd/workflows/Go/badge.svg?branch=master)

A simple clipboard written in Go with support for multiple registries

## Example usage

In `zsh`, start the server in the background, redirecting STDERR to STDOUT to `/tmp/clipd.logs`. Then kill the process gracefully:

```sh
clipd start &> /tmp/clipd.logs &! CLIPD_PID=$!
kill -SIGINT $CLIPD_PID
```
