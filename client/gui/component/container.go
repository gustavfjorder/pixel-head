package component

type Container struct {
	Component
}

func NewContainer(children ...ComponentInterface) *Container {
	container := &Container{}

	container.Child(children...)
	container.Center()

	return container
}

func (container *Container) Render() ComponentInterface {
	return container
}
