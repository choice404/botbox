/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"encoding/json"
	"fmt"
)

type BotConfig struct {
	Name          string `json:"name"`
	CommandPrefix string `json:"command_prefix"`
	Author        string `json:"author"`
	Description   string `json:"description"`
}

type CogConfig struct {
	Name           string        `json:"name"`
	Env            string        `json:"env"`
	File           string        `json:"file"`
	SlashCommands  []CommandInfo `json:"slash_commands"`
	PrefixCommands []CommandInfo `json:"prefix_commands"`
}

func CogConfigSliceToJSON(slice []CogConfig) (string, error) {
	jsonData, err := json.Marshal(slice)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CogConfig slice to JSON: %w", err)
	}
	return string(jsonData), nil
}

func JSONToCogConfigSlice(jsonString string) ([]CogConfig, error) {
	var slice []CogConfig
	err := json.Unmarshal([]byte(jsonString), &slice)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to CogConfig slice: %w", err)
	}
	return slice, nil
}

func JSONToCogConfig(jsonString string) (*CogConfig, error) {
	var cogConfig CogConfig
	err := json.Unmarshal([]byte(jsonString), &cogConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to CogConfig: %w", err)
	}
	return &cogConfig, nil
}

func (c *CogConfig) ToJSON() (string, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CogConfig to JSON: %w", err)
	}
	return string(jsonData), nil
}

type Config struct {
	BotInfo BotConfig   `json:"bot"`
	Cogs    []CogConfig `json:"cogs"`
}

type CommandInfo struct {
	Name        string
	Scope       string
	Type        string
	Description string
	Args        []ArgInfo
	ReturnType  string
}

func CmdInfoSliceToJSON(slice []CommandInfo) (string, error) {
	jsonData, err := json.Marshal(slice)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CommandInfo slice to JSON: %w", err)
	}
	return string(jsonData), nil
}

func JSONToCmdInfoSlice(jsonString string) ([]CommandInfo, error) {
	var slice []CommandInfo
	err := json.Unmarshal([]byte(jsonString), &slice)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to CommandInfo slice: %w", err)
	}
	return slice, nil
}

func JSONToCmdInfo(jsonString string) (*CommandInfo, error) {
	var cmdInfo CommandInfo
	err := json.Unmarshal([]byte(jsonString), &cmdInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to CommandInfo: %w", err)
	}
	return &cmdInfo, nil
}

func (c *CommandInfo) ToJSON() (string, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal CommandInfo to JSON: %w", err)
	}
	return string(jsonData), nil
}

type ArgInfo struct {
	Name        string
	Type        string
	Description string
}

func ArgInfoSliceToJSON(slice []ArgInfo) (string, error) {
	jsonData, err := json.Marshal(slice)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ArgInfo slice to JSON: %w", err)
	}
	return string(jsonData), nil
}

func JSONToArgInfoSlice(jsonString string) ([]ArgInfo, error) {
	var slice []ArgInfo
	err := json.Unmarshal([]byte(jsonString), &slice)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to ArgInfo slice: %w", err)
	}
	return slice, nil
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

MIT License

Copyright (c) 2025 Austin Choi

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
