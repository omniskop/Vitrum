// Code generated by vitrum gencmd. DO NOT EDIT.

package controls

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
	std "github.com/omniskop/vitrum/vit/std"
)

func newFileContextForTextField(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	fileCtx := vit.NewFileContext(globalCtx)

	var lib parse.Library
	var err error
	lib, err = parse.ResolveLibrary([]string{"Vit"})
	if err != nil {
		// The file used to generate the "TextField" component imported a library called "Vit".
		// If this error occurs that imported failed. Probably because the library is not known.
		return nil, fmt.Errorf("unable to create file context for generated \"TextField\" component: %w", err)
	}
	parse.AddLibraryToContainer(lib, &fileCtx.KnownComponents)

	return fileCtx, nil
}

type TextField struct {
	*std.Item
	id string

	text    vit.StringValue
	focused vit.BoolValue
}

// newTextFieldInGlobal creates an appropriate file context for the component and then returns a new TextField instance.
// The returned error will only be set if a library import that is required by the component fails.
func newTextFieldInGlobal(id string, globalCtx *vit.GlobalContext) (*TextField, error) {
	fileCtx, err := newFileContextForTextField(globalCtx)
	if err != nil {
		return nil, err
	}
	return NewTextField(id, fileCtx), nil
}
func NewTextField(id string, context *vit.FileContext) *TextField {
	t := &TextField{
		Item:    std.NewItem("textField", context),
		id:      id,
		text:    *vit.NewEmptyStringValue(),
		focused: *vit.NewBoolValueFromCode(vit.Code{FileCtx: context, Code: "false", Position: nil}),
	}
	// property assignments on embedded components
	t.Item.SetPropertyCode("width", vit.Code{FileCtx: context, Code: "300", Position: nil})
	t.Item.SetPropertyCode("height", vit.Code{FileCtx: context, Code: "25", Position: nil})
	// register listeners for when a property changes
	// register event listeners
	var event vit.Listenable
	var listener vit.Evaluater
	event, _ = t.Root.Event("onCompleted")
	listener = event.CreateListener(vit.Code{FileCtx: context, Code: "function() {\n        // keyArea.onKeyDown.addEventListener()\n        // mouseArea.onClicked.AddEventListener(clicked)\n        textField.focused = true;\n    }", Position: nil})
	t.AddListenerFunction(listener)
	event.(*vit.EventAttribute[struct{}]).AddListener(vit.ListenerCB[struct{}](t.wasCompleted))
	// register enumerations
	// add child components
	var child vit.Component
	child, _ = parse.InstantiateComponent(&vit.ComponentDefinition{BaseName: "MouseArea", ID: "mouseArea", Properties: []vit.PropertyDefinition{vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 26, StartColumn: 9, EndLine: 26, EndColumn: 28}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 26, StartColumn: 23, EndLine: 26, EndColumn: 28}, Identifier: []string{"anchors", "fill"}, Expression: "parent", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 27, StartColumn: 9, EndLine: 27, EndColumn: 45}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 27, StartColumn: 26, EndLine: 27, EndColumn: 45}, Identifier: []string{"acceptedButtons"}, Expression: "MouseArea.leftButton", Tags: map[string]string{}}}}, context)
	t.AddChild(child)
	child, err := parse.InstantiateComponent(&vit.ComponentDefinition{BaseName: "KeyArea", ID: "keyArea", Properties: []vit.PropertyDefinition{vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 32, StartColumn: 9, EndLine: 32, EndColumn: 28}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 32, StartColumn: 23, EndLine: 32, EndColumn: 28}, Identifier: []string{"anchors", "fill"}, Expression: "parent", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 33, StartColumn: 9, EndLine: 33, EndColumn: 31}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 33, StartColumn: 18, EndLine: 33, EndColumn: 31}, Identifier: []string{"enabled"}, Expression: "parent.focused", Tags: map[string]string{}}}}, context)
	fmt.Println(err)
	t.AddChild(child)
	child, _ = parse.InstantiateComponent(&vit.ComponentDefinition{BaseName: "Rectangle", Properties: []vit.PropertyDefinition{vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 37, StartColumn: 9, EndLine: 37, EndColumn: 28}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 37, StartColumn: 23, EndLine: 37, EndColumn: 28}, Identifier: []string{"anchors", "fill"}, Expression: "parent", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 38, StartColumn: 9, EndLine: 38, EndColumn: 37}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 38, StartColumn: 16, EndLine: 38, EndColumn: 37}, Identifier: []string{"color"}, Expression: "Vit.rgb(230, 230, 230)", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 39, StartColumn: 9, EndLine: 39, EndColumn: 44}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 39, StartColumn: 23, EndLine: 39, EndColumn: 44}, Identifier: []string{"border", "color"}, Expression: "Vit.rgb(130, 130, 130)", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 40, StartColumn: 9, EndLine: 40, EndColumn: 44}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 40, StartColumn: 23, EndLine: 40, EndColumn: 44}, Identifier: []string{"border", "width"}, Expression: "parent.focused ? 2 : 0", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 41, StartColumn: 9, EndLine: 41, EndColumn: 17}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 41, StartColumn: 17, EndLine: 41, EndColumn: 17}, Identifier: []string{"radius"}, Expression: "5", Tags: map[string]string{}}}}, context)
	t.AddChild(child)
	child, _ = parse.InstantiateComponent(&vit.ComponentDefinition{BaseName: "Text", ID: "text", Properties: []vit.PropertyDefinition{vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 46, StartColumn: 9, EndLine: 46, EndColumn: 28}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 46, StartColumn: 23, EndLine: 46, EndColumn: 28}, Identifier: []string{"anchors", "fill"}, Expression: "parent", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 47, StartColumn: 9, EndLine: 47, EndColumn: 29}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 47, StartColumn: 29, EndLine: 47, EndColumn: 29}, Identifier: []string{"anchors", "leftMargin"}, Expression: "5", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 48, StartColumn: 9, EndLine: 48, EndColumn: 30}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 48, StartColumn: 30, EndLine: 48, EndColumn: 30}, Identifier: []string{"anchors", "rightMargin"}, Expression: "5", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 49, StartColumn: 9, EndLine: 49, EndColumn: 25}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 49, StartColumn: 15, EndLine: 49, EndColumn: 25}, Identifier: []string{"text"}, Expression: "parent.text", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 50, StartColumn: 9, EndLine: 50, EndColumn: 26}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 50, StartColumn: 25, EndLine: 50, EndColumn: 26}, Identifier: []string{"font", "pointSize"}, Expression: "40", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 51, StartColumn: 9, EndLine: 51, EndColumn: 33}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 51, StartColumn: 22, EndLine: 51, EndColumn: 33}, Identifier: []string{"font", "family"}, Expression: "\"Montserrat\"", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 52, StartColumn: 9, EndLine: 52, EndColumn: 32}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 52, StartColumn: 22, EndLine: 52, EndColumn: 32}, Identifier: []string{"font", "weight"}, Expression: "Text.Medium", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 53, StartColumn: 9, EndLine: 53, EndColumn: 44}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 53, StartColumn: 28, EndLine: 53, EndColumn: 44}, Identifier: []string{"verticalAlignment"}, Expression: "Text.AlignVCenter", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 54, StartColumn: 9, EndLine: 54, EndColumn: 43}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 54, StartColumn: 30, EndLine: 54, EndColumn: 43}, Identifier: []string{"horizontalAlignment"}, Expression: "Text.AlignLeft", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: "TextField.vit", StartLine: 55, StartColumn: 9, EndLine: 55, EndColumn: 29}, ValuePos: &vit.PositionRange{FilePath: "TextField.vit", StartLine: 55, StartColumn: 16, EndLine: 55, EndColumn: 29}, Identifier: []string{"elide"}, Expression: "Text.ElideLeft", Tags: map[string]string{}}}}, context)
	t.AddChild(child)

	context.RegisterComponent("textField", t)

	return t
}

func (t *TextField) String() string {
	return fmt.Sprintf("TextField(%s)", t.id)
}

func (t *TextField) Property(key string) (vit.Value, bool) {
	switch key {
	case "text":
		return &t.text, true
	case "focused":
		return &t.focused, true
	default:
		return t.Item.Property(key)
	}
}

func (t *TextField) MustProperty(key string) vit.Value {
	v, ok := t.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (t *TextField) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "text":
		err = t.text.SetValue(value)
	case "focused":
		err = t.focused.SetValue(value)
	default:
		return t.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("TextField", key, t.id, err)
	}
	return nil
}

func (t *TextField) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "text":
		t.text.SetCode(code)
	case "focused":
		t.focused.SetCode(code)
	default:
		return t.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (t *TextField) Event(name string) (vit.Listenable, bool) {
	switch name {
	default:
		return t.Item.Event(name)
	}
}

func (t *TextField) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "text":
		return &t.text, true
	case "focused":
		return &t.focused, true
	default:
		return t.Item.ResolveVariable(key)
	}
}

func (t *TextField) AddChild(child vit.Component) {
	child.SetParent(t)
	t.AddChildButKeepParent(child)
}

func (t *TextField) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	var targetType vit.Component = afterThis

	for ind, child := range t.Children() {
		if child.As(&targetType) {
			addThis.SetParent(t)
			t.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	t.AddChild(addThis)
}

func (t *TextField) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	// properties
	if changed, err := t.text.Update(t); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("TextField", "text", t.id, err))
		}
	}
	if changed, err := t.focused.Update(t); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("TextField", "focused", t.id, err))
		}
	}

	// methods

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	n, err := t.UpdatePropertiesInContext(t)
	sum += n
	errs.AddGroup(err)
	n, err = t.Item.UpdateExpressions()
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (t *TextField) As(target *vit.Component) bool {
	if _, ok := (*target).(*TextField); ok {
		*target = t
		return true
	}
	return t.Item.As(target)
}

func (t *TextField) ID() string {
	return t.id
}

func (t *TextField) Finish() error {
	return t.RootC().FinishInContext(t)
}
