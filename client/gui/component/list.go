package component

type List struct {
	Component
}

func NewList(children ...ComponentInterface) *List {
	list := &List{}

	list.Child(children...)

	return list
}

func (list *List) Render() ComponentInterface {
	list.rows = float64(len(list.children))
	list.CalculateBounds()

	return list
}
