package main

const HELPTEXT =
`The following commands are available:
    - execute [command|preset] [name]    Executes a command or a preset
    - list [commands|presets]            Lists available commands or presets
    - inspect [command|preset] [name]    Inspects a command or a preset showing either
                                         the code sent out by the command or the commands
                                         a preset executes
    - raw [data]                         Used to send raw codes to the receiver
    - help                               Displays this help message
    - exit                               Exits interactive mode, dropping you backc to your
                                         terminal`