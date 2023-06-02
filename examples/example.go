package main

import (
	"fmt"

	"github.com/linux-immutability-tools/DistroboxLib/bindings"
)

func main() {
	// engine, err := core.NewEngine("podman", "storage", "overlay")
	// if err != nil {
	// 	panic(err)
	// }

	//dbox, err := bindings.NewDboxWithOpts("/usr/lib/apx/distrobox", engine)
	dbox, err := bindings.NewDbox("/usr/lib/apx/distrobox")
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
