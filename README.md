ir.binGO
====
ir.binGO is a tool to manage infrared devices in the bytewerk hackerspace using the terminal. It gives you the possibility to define shortcuts to make music louder/quieter or to turn on and off the boxes using just a keystroke. It is written in Go. Obviously.

Quickstart
===
[![asciicast](https://asciinema.org/a/5yuea6wrhe9331w5592yz3w5b.png)](https://asciinema.org/a/5yuea6wrhe9331w5592yz3w5b)

Manual
===
    Usage of ir.binGO:
      -c string
            The command you want to execute. See a list of commands in the config file
      -config string
        Provide the config file to read from. Defaults to /etc/ir.bingo/config.toml (default "/etc/ir.bingo/config.toml")
      -d string
        The device you want to send your commands to. Defaults to the amplifier (default "amplifier")
      -delay int
        The delay in milliseconds to wait between multiple commands (default 200)
      -i    Interactive mode: Use a shell to interact with the receiver
      -p string
            The preset you want to execute. See a list of presets in the config file

Installation
===
1. Either download a pre-built copy of ir.binGO, or follow the build instructions below.
2. *[Optional] Copy the binary to a location that is in your PATH variable so you can issue commands without the need to provide the exact location each time*
3. Download the out-of-the-box configuration (which should be working in the bytewerk), create the `/etc/ir.bingo` directory and copy the configuration to `/etc/ir.bingo/config.toml` *(or just create an own one)*
4. *[Optional] Create own presets and/or commands in the config file. It should be pretty self-explainatory.*

Build
===
Set up a Go environment and make sure your GOPATH is set correctly. Then execute the following commands:

    go get github.com/BurntSushi/toml
    go get github.com/Bytewerk/ir.binGO
    go build -o $GOPATH/bin/ir.binGO github.com/Bytewerk/ir.binGO

License
===
I publish this code under the [CC0-1.0](https://creativecommons.org/publicdomain/zero/1.0/) license.
