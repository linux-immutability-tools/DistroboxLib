package main

import (
	"fmt"

	"github.com/linux-immutability-tools/DistroboxLib/bindings"
	"github.com/linux-immutability-tools/DistroboxLib/core"
)

func main() {
	engine, err := core.NewPodmanEngine("storage", "overlay")
	if err != nil {
		panic(err)
	}

	dbox, err := bindings.NewDboxWithOpts("/home/mirko/Downloads/89luca89-distrobox-1.4.2-1-g849131c/89luca89-distrobox-849131c/", engine)
	//dbox, err := bindings.NewDbox("/usr/lib/apx/distrobox")
	if err != nil {
		panic(err)
	}

	err = dbox.CreateContainer("alpine:latest", "test", []string{}, []string{}, []string{}, []string{}, false, false, "", []string{})
	if err != nil {
		panic(err)
	}

	containers, err := dbox.ListContainers()
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Printf("Name: %s\n", container.Name)
		fmt.Printf("ID: %s\n", container.ID)
		fmt.Printf("Status: %s\n", container.Status)
		fmt.Printf("Image: %s\n\n", container.Image)
	}
}
