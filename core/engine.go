package core

import (
	"fmt"
	"os/exec"
)

type Engine struct {
	Binary          string
	EnvVars         map[string]string
	AdditionalFlags []string
}

func NewEngine(name, storageRoot string, storageDriver string) (*Engine, error) {
	binary, err := getEngineBinary(name)
	if err != nil {
		return nil, err
	}

	return &Engine{
		Binary: binary,
		AdditionalFlags: []string{
			"--storage-driver", storageDriver,
			"--root", storageRoot,
		},
	}, nil
}

func DetectEngine() (*Engine, error) {
	for _, binary := range engineBinaries {
		location, err := exec.LookPath(binary)
		if err != nil {
			continue
		}

		return &Engine{
			Binary: location,
		}, nil
	}

	return nil, fmt.Errorf("no supported engine found")
}

func getEngineBinary(name string) (string, error) {
	for _, binary := range engineBinaries {
		if binary == name {
			location, err := exec.LookPath(binary)
			if err != nil {
				return "", err
			}

			return location, nil
		}
	}

	return "", fmt.Errorf("engine %s not found", name)
}

var engineBinaries = []string{
	"docker",
	"podman",
	// "prometheus",
}
