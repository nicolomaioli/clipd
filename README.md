# clipd

![Go](https://github.com/nicolomaioli/clipd/workflows/Go/badge.svg?branch=master)

A simple clipboard written in Go with support for multiple registries. It
includes:

- An http server to read from and write to an embedded cache , data is received
  and transmitted as `text/plain`;
- A fully configurable cli to start the server, yank (with pipes), and paste.

This project was inspired by
[this post](https://andrewbrookins.com/technology/synchronizing-the-ios-clipboard-with-a-remote-server-using-command-line-tools).
I have been spending some time with Go, and it seemed like a good idea to get
some hands-on experience.

It probably needs some refactoring, and the test coverage leaves a lot to be
desired, but it's a start.

## Get/install

If you have Go installed, just `go get -u github.com/nicolomaioli/clipd`. You
can also clone the repository and `go build` or `go install`. Binaries are not
available at this time.

## Configuration

Checkout the `--help` flag. You can also create a global config file
`$HOME/.clipd.yaml` (or `json`, or `toml`, or any format supported by
[viper](https://github.com/spf13/viper)). Here's an example config:

```yml
server:
    address: ":8891"
    develop: false
    logLevel: 3
client:
    address: ":8891"
```

## Start with Systemd

An example `clipd.service` (you will probably want to change `User` and the
path to the executable in `ExecStart`):

```
[Unit]
Description=clipd server
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=pi
ExecStart=/home/pi/go/bin/clipd start -l 0

[Install]
WantedBy=multi-user.target
```

Put it in `/etc/systemd/system/clipd.service`, then:

```sh
# Set correct permissions
sudo chmod 644 /etc/systemd/system/clipd.service

# Start it
sudo systemctl start clipd.service

# Get the status
systemctl status clipd.service

# Enable it
sudo systemctl enable clipd.service

# Get the logs
journalctl -u clipd
```

## Over SSH with port forwarding

If `clipd` is running on the target machine on port `8891`, you can access it
from your local machine with port forwarding:

```sh
ssh -L 9188:localhost:8891 user@host
```

Then, on the client, you can access the clipboard with:

```sh
echo "Hello world!" | clipd yank -a 9188
clipd paste -a 9188
# Hello world!
```

## Neovim

`:help g:clipboard` for all the good stuff. Here is a minimal `init.vim`:

```vim
let g:clipboard = {
      \   'name': 'clipd',
      \   'copy': {
      \      '+': 'clipd yank',
      \      '*': 'clipd yank',
      \    },
      \   'paste': {
      \      '+': 'clipd paste',
      \      '*': 'clipd paste',
      \   },
      \ }

set clipboard+=unnamedplus
```

## Tmux

You can specify your own clipboard in Tmux:

```tmux
bind-key -T copy-mode-vi y send-keys -X copy-pipe-and-cancel "clipd yank"
bind-key -T copy-mode-vi Enter send-keys -X copy-pipe-and-cancel "clipd paste"
```

## Web UI

Adding a basic web UI is probably going to be the next step in the development
of `clipd`. One consideration here is that the server reads and return
`text/plain` content, so the output in particular should be properly sanitized
before it hits the browser.
