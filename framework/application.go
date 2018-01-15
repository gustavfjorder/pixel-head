package framework

import (
	"reflect"
	"fmt"
	"errors"
)

type Application struct {
	container         *Container
	CurrentController ControllerInterface

	controllers map[string]ControllerInterface
}

func NewApplication(container *Container) *Application {
	app := &Application{
		container:   container,
		controllers: make(map[string]ControllerInterface),
	}

	return app
}

func (app *Application) AddController(name string, c ControllerInterface) {
	app.controllers[name] = app.prepareController(name, reflect.Indirect(reflect.ValueOf(c)).Type())
}

func (app *Application) Run() {
	app.CurrentController.Run()
}

func (app *Application) Update() {
	app.CurrentController.Update()
}

func (app *Application) SetController(name string) {
	controller, found := app.controllers[name]
	if ! found {
		panic(errors.New(fmt.Sprintf("Controller '%s' is not found", name)))
	}

	app.CurrentController = controller
}

func (app *Application) ChangeTo(name string) {
	app.SetController(name)
	app.CurrentController.Run()
}

func (app *Application) prepareController(name string, controllerType reflect.Type) ControllerInterface {
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