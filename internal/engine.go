package internal

import (
	"fmt"
	"os/exec"
)

type Engine struct {
	Binary          string
	StorageRoot     string
	AdditionalFlags []string
}

func NewEngine(name, storageRoot string, additionalFlags []string) (*Engine, error) {
	binary, err := getEngineBinary(name)
	if err != nil {
		return nil, err
	}

	return &Engine{
		Binary:          binary,
		StorageRoot:     storageRoot,
		AdditionalFlags: additionalFlags,
	}, nil
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
