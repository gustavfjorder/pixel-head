package framework

import "errors"

type ControllerInterface interface {
	Init()
	Run()
	Update()
}

type Controller struct {
	ControllerInterface

	App       *Application
	Container *Container
}

func (c *Controller) Init() {
	//panic(errors.New("NOT YET IMPLEMENTED"))
}

func (c *Controller) Run() {
	panic(errors.New("NOT YET IMPLEMENTED"))
}

func (c *Controller) Update() {
	panic(errors.New("NOT YET IMPLEMENTED"))
}