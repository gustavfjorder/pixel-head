package helper

import "reflect"

func TypeImplements(item interface{}, implementable interface{}) bool {
	interfaceType := reflect.Indirect(reflect.ValueOf(implementable)).Type()

	return reflected(item).Type().Implements(interfaceType)
}

func TypeHasMethod(item interface{}, method string) bool {
	return reflected(item).MethodByName(method).Kind() != reflect.Invalid
}

func TypeHasField(item interface{}, field string) bool {
	return reflected(item).Elem().FieldByName(field).Kind() != reflect.Invalid
}

func RealType(item interface{}) reflect.Type {
	return reflected(item).Type()
}

func CallMethod(item interface{}, method string) {
	reflected(item).MethodByName(method).Call([]reflect.Value{})
}

func reflected(item interface{}) reflect.Value {
	itemType := reflect.Indirect(reflect.ValueOf(item))

	return reflect.New(itemType.Type())
}