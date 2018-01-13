package framework

import (
	"reflect"
	"fmt"
	"errors"
)

type Application struct {
	container         *Container
	CurrentController ControllerInterface

	controllers map[string]reflect.Type
}

func NewApplication(container *Container) *Application {
	app := &Application{
		container:   container,
		controllers: make(map[string]reflect.Type),
	}

	return app
}

func (app *Application) AddController(name string, c ControllerInterface) {
	app.controllers[name] = reflect.Indirect(reflect.ValueOf(c)).Type()
}

func (app *Application) Run() {
	app.CurrentController.Run()
}

func (app *Application) Update() {
	app.CurrentController.Update()
}

func (app *Application) SetController(name string) {
	app.CurrentController = app.prepareController(name)
}

func (app *Application) ChangeTo(name string) {
	app.SetController(name)
	app.CurrentController.Run()
}

func (app *Application) prepareController(name string) ControllerInterface {
	controllerType, found := app.controllers[name]
	if ! found {
		panic(errors.New(fmt.Sprintf("Controller '%s' is not found", name)))
	}

	controller := reflect.New(controllerType)

	// Set application
	f := controller.Elem().FieldByName("App")
	if f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(app))
	}

	// Set container
	f = controller.Elem().FieldByName("Container")
	if f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(app.container))
	}

	// Initialize controller
	controller.MethodByName("Init").Call([]reflect.Value{})

	return controller.Interface().(ControllerInterface)
}