package core

import (
	"os/exec"
	"strings"
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
		dboxBinary: dboxBinary,
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

func (p *Executor) NewCmd(dryRun bool, useAdditionalFlags bool, args ...string) *exec.Cmd {
	allArgs := append([]string{}, args...)

	if dryRun {
		allArgs = append(allArgs, "--dry-run")
	}

	if len(p.engine.AdditionalFlags) > 0 && useAdditionalFlags {
		allArgs = append(allArgs, "--additional-flags")
		allArgs = append(allArgs, strings.Join(p.engine.AdditionalFlags, " "))
	}

	cmd := exec.Command(p.dboxBinary, allArgs...)

	for key, value := range p.envVars {
		cmd.Env = append(cmd.Env, key+"="+value)
	}

	return cmd
}
