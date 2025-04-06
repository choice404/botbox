/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

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

/*
Copyright © 2025 Austin "Choice404" Choi

https://github.com/choice404/botbox

Bot Box
A discord bot template generator to help create discord bots quickly and easily

This code is licensed under the MIT License.
Please see the LICENSE file in the root directory of this project for the full license details.
*/
