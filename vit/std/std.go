package std

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

//go:generate go build -o gencmd github.com/omniskop/vitrum/vit/generator/gencmd
//go:generate ./gencmd -i Rectangle.vit -o rectangle_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Repeater.vit -o repeater_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Row.vit -o row_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Column.vit -o column_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Grid.vit -o grid_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i Text.vit -o text_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate ./gencmd -i MouseArea.vit -o mouseArea_gen.go -p github.com/omniskop/vitrum/vit/std
//go:generate rm ./gencmd

func init() {
	parse.RegisterLibrary("Vit", StdLib{})
}

type StdLib struct {
}

func (l StdLib) ComponentNames() []string {
	return []string{"Item", "Rectangle", "Repeater", "Container", "Row", "Column", "Grid", "Text", "MouseArea"}
}

func (l StdLib) NewComponent(name string, id string, globalCtx *vit.GlobalContext) (vit.Component, bool) {
	var comp vit.Component
	var err error
	switch name {
	case "Item":
		var fileCtx = vit.NewFileContext(globalCtx)
		return NewItem(id, fileCtx), true
	case "Rectangle":
		comp, err = newRectangleInGlobal(id, globalCtx)
	case "Repeater":
		comp, err = newRepeaterInGlobal(id, globalCtx)
	case "Container":
		var fileCtx = vit.NewFileContext(globalCtx)
		return NewContainer(id, fileCtx), true
	case "Row":
		comp, err = newRowInGlobal(id, globalCtx)
	case "Column":
		comp, err = newColumnInGlobal(id, globalCtx)
	case "Grid":
		comp, err = newGridInGlobal(id, globalCtx)
	case "Text":
		comp, err = newTextInGlobal(id, globalCtx)
	case "MouseArea":
		comp, err = newMouseAreaInGlobal(id, globalCtx)
	default:
		return nil, false
	}
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return comp, true
}

func (l StdLib) StaticAttribute(componentName string, attributeName string) (interface{}, bool) {
	switch componentName {
	case "Grid":
		return (*Grid)(nil).staticAttribute(attributeName)
	case "Text":
		return (*Text)(nil).staticAttribute(attributeName)
	case "MouseArea":
		return (*MouseArea)(nil).staticAttribute(attributeName)
	}
	return nil, false
}
