package cmd

type BotConfig struct {
	Name          string `json:"name"`
	CommandPrefix string `json:"command_prefix"`
	Author        string `json:"author"`
	Description   string `json:"description"`
}

type CogConfig struct {
	Name     string   `json:"name"`
	File     string   `json:"file"`
	Commands []string `json:"commands"`
}

type Config struct {
	BotInfo BotConfig   `json:"bot"`
	Cogs    []CogConfig `json:"cogs"`
}

type CommandInfo struct {
	Name        string
	Description string
	Args        []ArgInfo
	ReturnType  string
}

type ArgInfo struct {
	Name        string
	Type        string
	Description string
}

type LicenseResponse struct {
	Body string `json:"body"`
}
