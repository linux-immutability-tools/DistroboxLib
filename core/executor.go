package core

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mattn/go-shellwords"
)

type Executor struct {
	dboxBinary string
	engine     *Engine
	envVars    map[string]string
}

func NewExecutor(engine *Engine, dboxBinaryPath string) (*Executor, error) {
	var dboxBinary string

	if dboxBinaryPath == "" {
		var err error
		dboxBinary, err = getDboxBinary()
		if err != nil {
			return nil, err
		}
	} else {
		dboxBinary = dboxBinaryPath
	}

	return &Executor{
		dboxBinary: filepath.Join(dboxBinary, "distrobox"),
		engine:     engine,
		envVars:    make(map[string]string),
	}, nil
}

func getDboxBinary() (string, error) {
	location, err := exec.LookPath("distrobox")
	if err != nil {
		return "", err
	}

	return location, nil
}

func (p *Executor) AddEnvVar(key, value string) {
	p.envVars[key] = value
}

func (p *Executor) NewCmd(useEngineFlags bool, useEngineEnv bool, args ...string) (*exec.Cmd, error) {
	allArgs := append([]string{}, args...)
	allArgs = append(allArgs, "--dry-run")

	if len(p.engine.AdditionalFlags) > 0 && useEngineFlags {
		allArgs = append(allArgs, "--additional-flags")
		allArgs = append(allArgs, strings.Join(p.engine.AdditionalFlags, " "))
	}

	cmd := exec.Command(p.dboxBinary, allArgs...)

	if len(p.engine.EnvVars) > 0 && useEngineEnv {
		for key, value := range p.engine.EnvVars {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
	}

	for key, value := range p.envVars {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	cmd.Env = append(cmd.Env, "DBX_CONTAINER_MANAGER="+p.engine.Name)

	cmd.Env = append(cmd.Env, os.Environ()...)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(string(output))
		return nil, fmt.Errorf("error building command: %s", err)
	}

	finalCmd, err := shellwords.Parse(string(output))
	if err != nil {
		return nil, fmt.Errorf("error parsing command: %s", err)
	}

	cmd = exec.Command(finalCmd[0], finalCmd[1:]...)
	return cmd, nil
}
