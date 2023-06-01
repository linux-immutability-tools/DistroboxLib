package types

type Container struct {
	ID     string
	Name   string
	Status ContainerStatus
	Image  string
}
