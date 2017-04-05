package main

type TOMLConfig struct {
	Config InfoConfig
	Devices map[string]InfoDevice
}

type InfoConfig struct {
	Requestpath string
	Getprefix string
}

type InfoDevice struct {
	Settings InfoSettings
	Presets map[string]InfoPreset
	Commands map[string]InfoCommand
}

type InfoSettings struct {
	Prot string
	Dev string
}

type InfoPreset struct {
	Commands []string
}

type InfoCommand struct {
	Command string
}