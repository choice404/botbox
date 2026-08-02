/*
Copyright © 2025 Austin Choi austinch20@protonmail.com
See end of file for extended copyright information
*/

package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// projectTemplateData holds the values rendered into the project templates
type projectTemplateData struct {
	Version     string
	Name        string
	Prefix      string
	Author      string
	Description string
	LicenseType string
	HasLicense  bool
	Doppler     bool
	Token       string
	Guild       string
	Project     string
	Config      string
	// HelpStyle is written into botbox.conf and read by the generated help cog
	HelpStyle string
	// EnvProvider is written into botbox.conf so later commands know how secrets are supplied
	EnvProvider string
}

// dockerTemplateData holds the values rendered into the docker templates
type dockerTemplateData struct {
	PythonVersion string
	Doppler       bool
}

// optionalValue reads a value off the model values bus, falling back when the key is absent
func optionalValue(values Values, key, fallback string) string {
	if v, ok := values.Map[key]; ok && v != nil && *v != "" {
		return *v
	}
	return fallback
}

// renderToFile renders the named template and writes the result to path
func renderToFile(path string, templateName string, data any) error {
	content, err := RenderTemplate(templateName, data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func CreateProject(rootDir string, values Values, force bool) error {
	directories := []string{
		"src",
		"src/cogs",
		"src/utils",
	}

	for _, dir := range directories {
		fullPath := filepath.Join(rootDir, dir)
		err := os.MkdirAll(fullPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory %s: %w", fullPath, err)
		}
	}

	licenseType := *values.Map["licenseType"]
	// The conf provider mirrors the env choice, anything but doppler stores env
	envProvider := "env"
	if *values.Map["envChoice"] == "doppler" {
		envProvider = "doppler"
	}
	data := projectTemplateData{
		Version:     Version,
		Name:        *values.Map["botName"],
		Prefix:      *values.Map["botPrefix"],
		Author:      *values.Map["botAuthor"],
		Description: *values.Map["botDescription"],
		LicenseType: licenseType,
		HasLicense:  licenseType != "no-license" && licenseType != "",
		Doppler:     *values.Map["envChoice"] == "doppler",
		Token:       *values.Map["botTokenDopplerProject"],
		Guild:       *values.Map["botGuildDopplerEnv"],
		Project:     *values.Map["botTokenDopplerProject"],
		Config:      *values.Map["botGuildDopplerEnv"],
		HelpStyle:   NormalizeHelpStyle(optionalValue(values, "helpStyle", DefaultHelpStyle)),
		EnvProvider: envProvider,
	}

	if confOpt, err := CreateFileOption(filepath.Join(rootDir, "botbox.conf"), force); err == nil && confOpt {
		err := renderToFile(filepath.Join(rootDir, "botbox.conf"), "botbox.conf.tmpl", data)
		if err != nil {
			return fmt.Errorf("error creating botbox.conf file: %w", err)
		}
	} else if err == nil && !confOpt {
		fmt.Println("Not overriding botbox.conf file.")
	} else {
		return fmt.Errorf("error creating botbox.conf file: %w", err)
	}

	if readmeOpt, err := CreateFileOption(filepath.Join(rootDir, "README.md"), force); err == nil && readmeOpt {
		err := renderToFile(filepath.Join(rootDir, "README.md"), "readme.md.tmpl", data)
		if err != nil {
			return fmt.Errorf("error creating README.md file: %w", err)
		}
	} else if err == nil && !readmeOpt {
		fmt.Println("Not overriding README.md file.")
	} else {
		return fmt.Errorf("error creating README.md file: %w", err)
	}

	if data.HasLicense {
		if licenseOpt, err := CreateFileOption(filepath.Join(rootDir, "LICENSE"), force); err == nil && licenseOpt {
			licenseText, err := FetchLicense(licenseType)
			if err != nil {
				return fmt.Errorf("Error fetching license %s: %v\n", licenseType, err)
			}

			err = os.WriteFile(filepath.Join(rootDir, "LICENSE"), []byte(licenseText), 0644)
			if err != nil {
				return fmt.Errorf("Error writing to LICENSE file: %v\n", err)
			}
		} else if err == nil && !licenseOpt {
			fmt.Println("Not overriding LICENSE file.")
		} else {
			return fmt.Errorf("Error creating LICENSE file: %v\n", err)
		}
	}

	if *values.Map["envChoice"] == "doppler" {
		if dopplerOpt, err := CreateFileOption(filepath.Join(rootDir, "doppler.yaml"), force); err == nil && dopplerOpt {
			err := renderToFile(filepath.Join(rootDir, "doppler.yaml"), "doppler.yaml.tmpl", data)
			if err != nil {
				return fmt.Errorf("Error creating doppler.yaml file: %v\n", err)
			}
		} else if err == nil && !dopplerOpt {
			fmt.Println("Not overriding doppler.yaml file.")
		} else {
			return fmt.Errorf("Error creating doppler.yaml file: %v\n", err)
		}
	} else if *values.Map["envChoice"] == "env" {
		if envOpt, err := CreateFileOption(filepath.Join(rootDir, ".env"), force); err == nil && envOpt {
			err := renderToFile(filepath.Join(rootDir, ".env"), "env.tmpl", data)
			if err != nil {
				return fmt.Errorf("Error creating .env file: %v\n", err)
			}
		} else if err == nil && !envOpt {
			fmt.Println("Not overriding .env file.")
		} else {
			return fmt.Errorf("Error creating .env file: %v\n", err)
		}
	} else if *values.Map["envChoice"] == "none" {
		fmt.Println("No environment file will be created.")
	} else {
		return fmt.Errorf("Invalid environment choice: %s", *values.Map["envChoice"])
	}

	if reqOpt, err := CreateFileOption(filepath.Join(rootDir, "requirements.txt"), force); err == nil && reqOpt {
		err := renderToFile(filepath.Join(rootDir, "requirements.txt"), "requirements.txt.tmpl", data)
		if err != nil {
			return fmt.Errorf("error creating requirements.txt file: %w", err)
		}
	} else if err == nil && !reqOpt {
		fmt.Println("Not overriding requirements.txt file.")
	} else {
		return fmt.Errorf("error creating requirements.txt file: %w", err)
	}

	if gitignoreOpt, err := CreateFileOption(filepath.Join(rootDir, ".gitignore"), force); err == nil && gitignoreOpt {
		err := renderToFile(filepath.Join(rootDir, ".gitignore"), "gitignore.tmpl", data)
		if err != nil {
			return fmt.Errorf("error creating .gitignore file: %w", err)
		}
	} else if err == nil && !gitignoreOpt {
		fmt.Println("Not overriding .gitignore file.")
	} else {
		return fmt.Errorf("error creating .gitignore file: %w", err)
	}

	if runOpt, err := CreateFileOption(filepath.Join(rootDir, "run.sh"), force); err == nil && runOpt {
		err := renderToFile(filepath.Join(rootDir, "run.sh"), "run.sh.tmpl", data)
		if err != nil {
			return fmt.Errorf("Error creating run.sh file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "run.sh"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for run.sh file: %v\n", err)
		}
	} else if err == nil && !runOpt {
		fmt.Println("Not overriding run.sh file.")
	} else {
		return fmt.Errorf("Error creating run.sh file: %v\n", err)
	}

	if mainOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "main.py"), force); err == nil && mainOpt {
		err := renderToFile(filepath.Join(rootDir, "src", "main.py"), "main.py.tmpl", data)
		if err != nil {
			return fmt.Errorf("Error creating main.py file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "src", "main.py"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for main.py file: %v\n", err)
		}
	}

	if helloWorldOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "cogs", "helloWorld.py"), force); err == nil && helloWorldOpt {
		err := renderToFile(filepath.Join(rootDir, "src", "cogs", "helloWorld.py"), "helloworld.py.tmpl", data)
		if err != nil {
			return fmt.Errorf("Error creating helloWorld.py file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "src", "cogs", "helloWorld.py"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for helloWorld.py file: %v\n", err)
		}
	}

	if helpOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "cogs", "help.py"), force); err == nil && helpOpt {
		err := renderToFile(filepath.Join(rootDir, "src", "cogs", "help.py"), "help.py.tmpl", data)
		if err != nil {
			return fmt.Errorf("Error creating help.py file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "src", "cogs", "help.py"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for help.py file: %v\n", err)
		}
	}

	if adminOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "cogs", "admin.py"), force); err == nil && adminOpt {
		err := renderToFile(filepath.Join(rootDir, "src", "cogs", "admin.py"), "admin.py.tmpl", data)
		if err != nil {
			return fmt.Errorf("Error creating admin.py file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "src", "cogs", "admin.py"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for admin.py file: %v\n", err)
		}
	}

	if cogsOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "cogs", "cogs.py"), force); err == nil && cogsOpt {
		err := renderToFile(filepath.Join(rootDir, "src", "cogs", "cogs.py"), "cogs.py.tmpl", data)
		if err != nil {
			return fmt.Errorf("Error creating cogs.py file: %v\n", err)
		}
		err = os.Chmod(filepath.Join(rootDir, "src", "cogs", "cogs.py"), 0755)
		if err != nil {
			return fmt.Errorf("Error setting permissions for cogs.py file: %v\n", err)
		}
	}

	if initOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "cogs", "__init__.py"), force); err == nil && initOpt {
		err := renderToFile(filepath.Join(rootDir, "src", "cogs", "__init__.py"), "init.py.tmpl", data)
		if err != nil {
			return fmt.Errorf("error creating __init__.py file: %w", err)
		}
	}

	if loggerOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "utils", "logger.py"), force); err == nil && loggerOpt {
		err := renderToFile(filepath.Join(rootDir, "src", "utils", "logger.py"), "logger.py.tmpl", data)
		if err != nil {
			return fmt.Errorf("error creating logger.py file: %w", err)
		}
	}

	if utilsInitOpt, err := CreateFileOption(filepath.Join(rootDir, "src", "utils", "__init__.py"), force); err == nil && utilsInitOpt {
		err := renderToFile(filepath.Join(rootDir, "src", "utils", "__init__.py"), "utils_init.py.tmpl", data)
		if err != nil {
			return fmt.Errorf("error creating utils __init__.py file: %w", err)
		}
	}

	// Docker files are opt in, the tui confirm and the --docker flag both store yes here
	if optionalValue(values, "dockerize", "no") == "yes" {
		// The docker init command reads the version the same way, global default then 3.11
		pythonVersion := DefaultPythonVersion
		if conf, err := LoadGlobalConfig(); err == nil && conf.Defaults.PythonVersion != "" {
			pythonVersion = conf.Defaults.PythonVersion
		}
		if _, err := GenerateDockerFiles(rootDir, pythonVersion, envProvider, force); err != nil {
			return fmt.Errorf("error creating docker files: %w", err)
		}
	}

	return nil
}

// DefaultPythonVersion seeds the docker base image when the global config has no default
const DefaultPythonVersion = "3.11"

// GenerateDockerFiles renders the Dockerfile, docker-compose.yml, and .dockerignore
// into rootDir and returns the paths it actually wrote
func GenerateDockerFiles(rootDir string, pythonVersion string, envProvider string, force bool) ([]string, error) {
	data := dockerTemplateData{
		PythonVersion: pythonVersion,
		Doppler:       envProvider == "doppler",
	}

	// Each output file pairs with the template that renders it
	files := []struct {
		name     string
		template string
	}{
		{"Dockerfile", "dockerfile.tmpl"},
		{"docker-compose.yml", "docker-compose.yml.tmpl"},
		{".dockerignore", "dockerignore.tmpl"},
	}

	var written []string
	for _, file := range files {
		path := filepath.Join(rootDir, file.name)
		opt, err := CreateFileOption(path, force)
		if err != nil {
			return written, fmt.Errorf("error creating %s file: %w", file.name, err)
		}
		// CreateFileOption already reports skipped files, keep stdout clean for the path list
		if !opt {
			continue
		}
		if err := renderToFile(path, file.template, data); err != nil {
			return written, fmt.Errorf("error creating %s file: %w", file.name, err)
		}
		written = append(written, path)
	}

	return written, nil
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
