package framework

import (
	"errors"
	"fmt"
	"sync"
)

type ContainerFunc func(container *Container) interface{}

type Container struct {
	values   map[string]ContainerFunc // contains the original closure to generate the service
	services map[string]interface{}   // contains the instantiated services
	mtx      *sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		services: make(map[string]interface{}),
		values:   make(map[string]ContainerFunc),
		mtx:      &sync.RWMutex{},
	}
}

func (container *Container) Set(name string, f ContainerFunc) {
	container.mtx.Lock()
	defer container.mtx.Unlock()

	if _, found := container.services[name]; found {
		panic(errors.New("CANNOT OVERWRITE INITIALIZED SERVICE"))
	}

	container.values[name] = f
}

func (container *Container) SetService(name string, service interface{}) {
	container.Set(name, func(container *Container) interface{} {
		return service
	})
}

func (container *Container) Get(name string) interface{} {
	container.mtx.RLock()
	_, found := container.services[name]
	container.mtx.RUnlock()
	if ! found {
		container.mtx.RLock()
		_, found = container.values[name]
		container.mtx.RUnlock()
		if ! found {
			panic(errors.New(fmt.Sprintf("SERVICE '%s' DOES NOT EXIST", name)))
		}

		v := container.values[name](container)

		container.mtx.Lock()
		container.services[name] = v
		container.mtx.Unlock()
	}

	container.mtx.RLock()
	defer container.mtx.RUnlock()

	return container.services[name]
}
