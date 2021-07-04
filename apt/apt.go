package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/anircodes/configo/commands"
	"github.com/anircodes/configo/core"
	"github.com/anircodes/configo/sarray"
	"gopkg.in/yaml.v2"
)

type Apt struct {
	Pkgs          []string `yaml:pkg`
	ExpectedState string   `yaml:state`
	Restart       bool     `yaml:restart`
	Upgrade       bool     `yaml:upgrade`
	Update        bool     `yaml:update`
}

func (uapt *Apt) Validate() core.ExecOutput {
	STATES := [3]string{"latest", "absent", "present"}
	execOutput := core.ExecOutput{}
	execOutput.ErrorCode = 0
	execOutput.ErrorMessage = ""

	inputYamlFile, err := ioutil.ReadFile("apt.yaml")
	if err != nil {
		execOutput.ErrorCode = -1
		execOutput.ErrorMessage = err.Error()
		return execOutput
	}
	var apt Apt
	err = yaml.Unmarshal(inputYamlFile, &apt)
	if err != nil {
		execOutput.ErrorCode = -1
		execOutput.ErrorMessage = err.Error()
		return execOutput
	}

	if !sarray.Contains(STATES, apt.ExpectedState) {
		execOutput.ErrorCode = -1
		execOutput.ErrorMessage = "Incorrect state value"
		return execOutput
	}
	apt.ExpectedState = strings.ToLower(apt.ExpectedState)
	return execOutput
}

func (apt *Apt) Run() core.Result {
	output := ""
	if apt.Update {
		result, err := commands.Run("apt-get", "update", "-y")
		if err != nil {
			return result
		}
	}
	if apt.Upgrade {
		result, err := commands.Run("apt-get", "upgrade", "-y")
		if err != nil {
			return result
		}
	}

	for _, pkg := range apt.Pkgs {
		isPresnt := false
		isLatest := false
		result, err := commands.Run("apt-get", "policy", pkg, "|", "grep", "Installed:")
		output = output + "\n" + result.Output
		if err != nil {
			fmt.Println("Failed to get package details - ", pkg)
			return result
		}
		installedVersion := strings.TrimSpace(strings.Split(result.Output, ":")[1])

		if installedVersion != "(none)" {
			isPresnt = true
			result, err := commands.Run("apt-get", "policy", pkg, "|", "grep", "Candidate:")
			output = output + "\n" + result.Output
			if err != nil {
				fmt.Println("Failed to get package details - ", pkg)
				return result
			}
			candidateVersion := strings.TrimSpace(strings.Split(result.Output, ":")[1])

			if installedVersion == candidateVersion {
				isLatest = true
			}
		}

		switch apt.ExpectedState {
		case "latest":
			if !isLatest {
				result, err := commands.Run("apt-get", "install", pkg, "-y")
				output = output + "\n" + result.Output
				if err != nil {
					fmt.Println("Failed to get package details - ", pkg)
					return result
				}
			} else {
				fmt.Println("Package already at latest version. Package - ", pkg, " Version - ", installedVersion)
			}
		case "present":
			if !isPresnt {
				result, err := commands.Run("apt-get", "install", pkg, "-y")
				output = output + "\n" + result.Output
				if err != nil {
					fmt.Println("Failed to get package details - ", pkg)
					return result
				}
			} else {
				fmt.Println("Package already at present. Package - ", pkg, " Version - ", installedVersion)
			}
		case "absent":
			if isPresnt {
				result, err := commands.Run("apt-get", "remove", pkg, "-y")
				output = output + "\n" + result.Output
				if err != nil {
					fmt.Println("Failed to get package details - ", pkg)
					return result
				}
			} else {
				fmt.Println("Package already at absent. Package - ", pkg, " Version - ", installedVersion)
			}
		}
	}
	return core.Result{Output: output}
}

func main() {
	overvallResult := 0
	if len(os.Args) == 1 {
		overvallResult = -1
		fmt.Printf("Invalid parameter. ")
		os.Exit(overvallResult)
	}
	var plugin core.Plugin = &Apt{}
	switch action := os.Args[1]; action {
	case "validate":
		execOutput := plugin.Validate()
		fmt.Println(execOutput)
	case "run":
		result := plugin.Run()
		fmt.Println(result)
	default:
		overvallResult = -1
		fmt.Printf("Invalid parameter. ")
	}
	os.Exit(overvallResult)
}
