package bindings

import (
	"fmt"
	"strings"

	"github.com/linux-immutability-tools/DistroboxLib/core"
	"github.com/linux-immutability-tools/DistroboxLib/types"
)

type Dbox struct {
	Executor *core.Executor
}

func NewDbox(dboxBinaryPath string) (*Dbox, error) {
	engine, err := core.DetectEngine()
	if err != nil {
		return nil, err
	}

	executor, err := core.NewExecutor(engine, dboxBinaryPath)
	if err != nil {
		return nil, err
	}

	return &Dbox{
		Executor: executor,
	}, nil
}

func NewDboxWithOpts(dboxBinaryPath string, engine *core.Engine) (*Dbox, error) {
	executor, err := core.NewExecutor(engine, dboxBinaryPath)
	if err != nil {
		return nil, err
	}

	return &Dbox{
		Executor: executor,
	}, nil
}

func newContainerFromLine(line string) (*types.Container, error) {
	container := types.Container{}

	fields := strings.Split(line, "|")
	if len(fields) != 4 {
		return nil, fmt.Errorf("invalid line: %s", line)
	}

	container.ID = strings.TrimSpace(fields[0])
	container.Name = strings.TrimSpace(fields[1])

	status := strings.TrimSpace(fields[2])
	if strings.HasPrefix(status, "Up") {
		container.Status = types.ContainerStatusRunning
	} else {
		container.Status = types.ContainerStatusStopped
	}

	container.Image = fields[3]

	return &container, nil
}

func (p *Dbox) ListContainers() ([]types.Container, error) {
	cmd, err := p.Executor.NewCmd(false, true, "list")
	if err != nil {
		return nil, err
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, err
	}

	containers := []types.Container{}
	for i, line := range strings.Split(string(output), "\n") {
		if i == 0 {
			continue
		}

		if line == "" {
			continue
		}

		container, err := newContainerFromLine(line)
		if err != nil {
			return nil, err
		}

		containers = append(containers, *container)
	}

	return containers, nil
}

func (p *Dbox) CreateContainer(image string, name string, additionalFlags []string, additionalPackages []string, initHooks []string, preInitHooks []string, init bool, nvidia bool, home string, volumes []string) error {
	args := []string{"create"}

	if image != "" {
		args = append(args, "--image", image)
	}

	if name != "" {
		args = append(args, "--name", name)
	}

	if len(additionalFlags) > 0 {
		args = append(args, "--additional-flags", strings.Join(additionalFlags, " "))
	}

	if len(additionalPackages) > 0 {
		args = append(args, "--additional-packages", strings.Join(additionalPackages, " "))
	}

	if len(initHooks) > 0 {
		args = append(args, "--init-hooks", strings.Join(initHooks, " "))
	}

	if len(preInitHooks) > 0 {
		args = append(args, "--pre-init-hooks", strings.Join(preInitHooks, " "))
	}

	if init {
		args = append(args, "--init")
	}

	if nvidia {
		args = append(args, "--nvidia")
	}

	if home != "" {
		args = append(args, "--home", home)
	}

	if len(volumes) > 0 {
		args = append(args, "--volume", strings.Join(volumes, " "))
	}

	args = append(args, "--yes")

	cmd, err := p.Executor.NewCmd(true, true, args...)
	if err != nil {
		return err
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error creating container: %s", string(output))
	}

	return nil
}
