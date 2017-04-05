package main

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"log"
	"flag"
	"strings"
	"fmt"
	"os"
	"bufio"
	"time"
	"strconv"
)

var (
	AMPLIFIER = "amplifier"
	PROJECTOR = "projector"
)

func main() {
	flagConfigFile := flag.String("config", "/etc/ir.bingo/config.toml", "Provide the config file to read from. Defaults to /etc/ir.bingo/config.toml")
	flagDevice := flag.String("d", "amplifier", "The device you want to send your commands to. Defaults to the amplifier")
	flagCommand := flag.String("c", "", "The command you want to execute. See a list of commands in the config file")
	flagPreset := flag.String("p", "", "The preset you want to execute. See a list of presets in the config file")
	flagInteractive := flag.Bool("i", false, "Interactive mode: Use a shell to interact with the receiver")
	flagDelay := flag.Int("delay", 200, "The delay in milliseconds to wait between multiple commands")

	flag.Parse()

	rawConfig, err := ioutil.ReadFile(*flagConfigFile)
	if err != nil {
		log.Fatal(err)
	}
	var tomlConfig TOMLConfig
	_, err = toml.Decode(string(rawConfig[:]), &tomlConfig)

	if len(*flagCommand) == 0 && len(*flagPreset) == 0 && !*flagInteractive {
		flag.Usage()
		return
	}

	infoDevice, ok := tomlConfig.Devices[*flagDevice]
	if !ok {
		log.Fatal("No such device: " + *flagDevice)
		return
	}

	deviceName := *flagDevice
	delay := *flagDelay
	if *flagInteractive {
		for {
			fmt.Print("[ " + deviceName + " ] >> ")
			reader := bufio.NewReader(os.Stdin)
			rawCommand, _ := reader.ReadSlice('\n')
			rawCommandString := strings.Trim(string(rawCommand[:]), "\n")

			commandSplit := strings.Split(rawCommandString, " ")
			if len(commandSplit) == 0 {
				continue
			}
			switch strings.ToLower(commandSplit[0]) {
			case "exit":
				return
			case "set":
				if len(commandSplit) != 3 {
					fmt.Println("Syntax: \"set [key] [value]\"")
					fmt.Println("Use \"list variables\" to see a list of available variables")
				}
				switch commandSplit[1] {
				case "device":
					found := false
					for key := range tomlConfig.Devices {
						if commandSplit[2] == key {
							found = true
							break
						}
					}
					if found {
						deviceName = commandSplit[2]
					} else {
						fmt.Println("No such device: " + commandSplit[2])
					}
					break
				case "delay":
					d, err := strconv.ParseInt(commandSplit[2], 10, 32)
					if err != nil {
						fmt.Println("Couldn't parse the number! Did you really enter an integer?")
						continue
					}
					delay = int(d)
					break
				}
				break
			case "execute":
				if len(commandSplit) != 3 {
					fmt.Println("Syntax: \"execute [command|preset] [name]\"")
					fmt.Println("Use \"list commands\" to see a list of available commands or \"list presets\" to see a list of available presets")
					continue
				}
				switch commandSplit[1] {
				case "preset":
					preset, ok := infoDevice.Presets[commandSplit[2]]
					if !ok {
						fmt.Println("No such preset: " + commandSplit[2])
						continue
					}
					for _, command := range preset.Commands {
						cmd, ok := infoDevice.Commands[command]
						if !ok {
							fmt.Println("No such command: " + command)
							continue
						}
						_, err := SendCommand(tomlConfig.Config.Requestpath, ConvertString(tomlConfig.Config.Getprefix + " " + infoDevice.Settings.Prot + " " + cmd.Command + " 00"))
						if err != nil {
							fmt.Println(err)
						}
						time.Sleep(time.Duration(delay) * time.Millisecond)
					}
					break
				case "command":
					command, ok := infoDevice.Commands[commandSplit[2]]
					if !ok {
						fmt.Println("No such command: " + commandSplit[2])
						continue
					}
					_, err := SendCommand(tomlConfig.Config.Requestpath, ConvertString(tomlConfig.Config.Getprefix + " " + infoDevice.Settings.Prot + " " + command.Command + " 00"))
					if err != nil {
						fmt.Println(err)
					}
					break
				}
				break
			case "list":
				if len(commandSplit) != 2 {
					fmt.Println("Syntax: \"list [commands|presets]\"")
					continue
				}
				switch commandSplit[1] {
				case "commands":
					fmt.Println("The following commands are available:")
					for key := range tomlConfig.Devices[deviceName].Commands {
						fmt.Println("    - " + key)
					}
					break
				case "presets":
					fmt.Println("The following presets are available:")
					for key := range tomlConfig.Devices[deviceName].Presets {
						fmt.Println("    - " + key)
					}
					break
				case "devices":
					fmt.Println("The following devices are available:")
					for key := range tomlConfig.Devices {
						fmt.Println("    - " + key)
					}
					break
				case "variables":
					fmt.Println("The following variables are available:")
					fmt.Println("    - device  The device you want to be sending commands to")
					fmt.Println("    - delay   The delay you want to use between multiple commands")
					break
				default:
					fmt.Println("Syntax: \"list [commands|presets]\"")
					continue
				}
				break
			case "inspect":
				if len(commandSplit) != 3 {
					fmt.Println("Syntax: \"list [command|preset] [name]\"")
					continue
				}
				switch commandSplit[1] {
				case "command":
					command, ok := tomlConfig.Devices[deviceName].Commands[commandSplit[2]]
					if !ok {
						fmt.Println("No such command: " + commandSplit[2])
						continue
					}
					fmt.Println("This command sends out the following code:")
					fmt.Println("    " + tomlConfig.Devices[deviceName].Settings.Prot + " " + command.Command + " 00")
					break
				case "preset":
					preset, ok := tomlConfig.Devices[deviceName].Presets[commandSplit[2]]
					if !ok {
						fmt.Println("No such preset: " + commandSplit[2])
						continue
					}
					fmt.Println("This preset sends out the following commands:")
					for _, command := range preset.Commands {
						fmt.Println("    - " + command)
					}
					break
				}
				break
			case "raw":
				if len(commandSplit) <= 1 {
					fmt.Println("Syntax: \"raw [data]\"")
					continue
				}
				rawData := ""
				for _, data := range commandSplit[1:] {
					rawData += data + " "
				}
				resp, err := SendCommand(tomlConfig.Config.Requestpath, ConvertString(tomlConfig.Config.Getprefix + " " + rawData))
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Command executed successfully!")
					fmt.Println("Response:")
					fmt.Println(resp)
				}
				break
			case "help":
				if len(commandSplit) != 1 {
					fmt.Println("Syntax: \"help\"")
					continue
				}
				fmt.Println(HELPTEXT)
				break
			default:
				fmt.Println(commandSplit[0] + ": command not found. Please use \"help\" to see a list of available commands.")
				break
			}
		}
		return
	}

	if len(*flagCommand) > 0 {
		command, ok := infoDevice.Commands[*flagCommand]
		if !ok {
			log.Fatal("No such command: " + *flagCommand)
			return
		}
		_, err := SendCommand(tomlConfig.Config.Requestpath, ConvertString(tomlConfig.Config.Getprefix + " " + infoDevice.Settings.Prot + " " + command.Command + " 00"))
		if err != nil {
			log.Fatal(err)
		}
	}

	if len(*flagPreset) > 0 {
		preset, ok := infoDevice.Presets[*flagPreset]
		if !ok {
			log.Fatal("No such preset: " + *flagPreset)
			return
		}
		for _, command := range preset.Commands {
			cmd, ok := infoDevice.Commands[command]
			if !ok {
				log.Fatal("No such command: " + *flagCommand)
				return
			}
			_, err := SendCommand(tomlConfig.Config.Requestpath, ConvertString(tomlConfig.Config.Getprefix + " " + infoDevice.Settings.Prot + " " + cmd.Command + " 00"))
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Duration(*flagDelay) * time.Millisecond)
		}
	}
}

func ConvertString(str string) string {
	split := strings.Split(str, " ")
	res := ""
	for _, s := range split {
		res += s + "%20"
	}
	return res
}
