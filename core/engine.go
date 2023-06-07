package core

import (
	"fmt"
	"os"
	"os/exec"
)

var engineBinaries = []string{
	"podman",
	"docker",
	// "prometheus",
}

type Engine struct {
	Name            string
	EnvVars         map[string]string
	AdditionalFlags []string
}

func NewPodmanEngine(storageRoot string, storageDriver string) (*Engine, error) {
	_, err := exec.LookPath("podman")
	if err != nil {
		return nil, err
	}

	storageRoot, err = getFullStorageRoot(storageRoot)
	if err != nil {
		return nil, err
	}

	err = createStorageIfMissing(storageRoot)
	if err != nil {
		return nil, err
	}

	return &Engine{
		Name: "podman",
		AdditionalFlags: []string{
			"--storage-driver", storageDriver,
			"--root", storageRoot,
		},
		EnvVars: map[string]string{
			"XDG_DATA_HOME":  storageRoot,
			"STORAGE_DRIVER": storageDriver,
		},
	}, nil
}

func NewCustomEngine(name string, envVars map[string]string, additionalFlags []string) (*Engine, error) {
	_, err := exec.LookPath(name)
	if err != nil {
		return nil, err
	}

	return &Engine{
		Name:            name,
		EnvVars:         envVars,
		AdditionalFlags: additionalFlags,
	}, nil
}

func DetectEngine() (*Engine, error) {
	for _, binary := range engineBinaries {
		_, err := exec.LookPath(binary)
		if err != nil {
			continue
		}

		return &Engine{
			Name: binary,
		}, nil
	}

	return nil, fmt.Errorf("no supported engine found. Supported engines: %s", engineBinaries)
}

func createStorageIfMissing(storageRoot string) error {
	if _, err := os.Stat(storageRoot); os.IsNotExist(err) {
		err := os.MkdirAll(storageRoot, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func getFullStorageRoot(storageRoot string) (string, error) {
	curDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", curDir, storageRoot), nil
}
