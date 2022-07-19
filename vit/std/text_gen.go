// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	canvas "github.com/tdewolff/canvas"
)

func newFileContextForText(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	return vit.NewFileContext(globalCtx), nil
}

type Text_HorizontalAlignment uint

const (
	Text_HorizontalAlignment_AlignLeft    Text_HorizontalAlignment = 0
	Text_HorizontalAlignment_AlignHCenter Text_HorizontalAlignment = 1
	Text_HorizontalAlignment_AlignRight   Text_HorizontalAlignment = 2
)

func (enum Text_HorizontalAlignment) String() string {
	switch enum {
	case Text_HorizontalAlignment_AlignLeft:
		return "AlignLeft"
	case Text_HorizontalAlignment_AlignHCenter:
		return "AlignHCenter"
	case Text_HorizontalAlignment_AlignRight:
		return "AlignRight"
	default:
		return "<unknownHorizontalAlignment>"
	}
}

type Text_VerticalAlignment uint

const (
	Text_VerticalAlignment_AlignTop     Text_VerticalAlignment = 0
	Text_VerticalAlignment_AlignVCenter Text_VerticalAlignment = 1
	Text_VerticalAlignment_AlignBottom  Text_VerticalAlignment = 2
)

func (enum Text_VerticalAlignment) String() string {
	switch enum {
	case Text_VerticalAlignment_AlignTop:
		return "AlignTop"
	case Text_VerticalAlignment_AlignVCenter:
		return "AlignVCenter"
	case Text_VerticalAlignment_AlignBottom:
		return "AlignBottom"
	default:
		return "<unknownVerticalAlignment>"
	}
}

type Text_FontWeight uint

const (
	Text_FontWeight_Thin       Text_FontWeight = 100
	Text_FontWeight_ExtraLight Text_FontWeight = 200
	Text_FontWeight_UltraLight Text_FontWeight = 200
	Text_FontWeight_Light      Text_FontWeight = 300
	Text_FontWeight_Normal     Text_FontWeight = 400
	Text_FontWeight_Regular    Text_FontWeight = 400
	Text_FontWeight_Medium     Text_FontWeight = 500
	Text_FontWeight_DemiBold   Text_FontWeight = 600
	Text_FontWeight_SemiBold   Text_FontWeight = 600
	Text_FontWeight_Bold       Text_FontWeight = 700
	Text_FontWeight_ExtraBold  Text_FontWeight = 800
	Text_FontWeight_UltraBold  Text_FontWeight = 800
	Text_FontWeight_Black      Text_FontWeight = 900
	Text_FontWeight_Heavy      Text_FontWeight = 900
)

func (enum Text_FontWeight) String() string {
	switch enum {
	case Text_FontWeight_Thin:
		return "Thin"
	case Text_FontWeight_ExtraLight:
		return "ExtraLight"
	// key UltraLight is omitted as it's the same as ExtraLight
	case Text_FontWeight_Light:
		return "Light"
	case Text_FontWeight_Normal:
		return "Normal"
	// key Regular is omitted as it's the same as Normal
	case Text_FontWeight_Medium:
		return "Medium"
	case Text_FontWeight_DemiBold:
		return "DemiBold"
	// key SemiBold is omitted as it's the same as DemiBold
	case Text_FontWeight_Bold:
		return "Bold"
	case Text_FontWeight_ExtraBold:
		return "ExtraBold"
	// key UltraBold is omitted as it's the same as ExtraBold
	case Text_FontWeight_Black:
		return "Black"
	// key Heavy is omitted as it's the same as Black
	default:
		return "<unknownFontWeight>"
	}
}

type Text_Elide uint

const (
	Text_Elide_ElideNone   Text_Elide = 0
	Text_Elide_ElideLeft   Text_Elide = 1
	Text_Elide_ElideMiddle Text_Elide = 2
	Text_Elide_ElideRight  Text_Elide = 3
)

func (enum Text_Elide) String() string {
	switch enum {
	case Text_Elide_ElideNone:
		return "ElideNone"
	case Text_Elide_ElideLeft:
		return "ElideLeft"
	case Text_Elide_ElideMiddle:
		return "ElideMiddle"
	case Text_Elide_ElideRight:
		return "ElideRight"
	default:
		return "<unknownElide>"
	}
}

type Text struct {
	*Item
	id string

	text                vit.StringValue
	color               vit.ColorValue
	horizontalAlignment vit.IntValue
	verticalAlignment   vit.IntValue
	font                vit.GroupValue
	elide               vit.IntValue
	fontData            *canvas.FontFamily
	fontFaceData        *canvas.FontFace
}

// newTextInGlobal creates an appropriate file context for the component and then returns a new Text instance.
// The returned error will only be set if a library import that is required by the component fails.
func newTextInGlobal(id string, globalCtx *vit.GlobalContext) (*Text, error) {
	fileCtx, err := newFileContextForText(globalCtx)
	if err != nil {
		return nil, err
	}
	return NewText(id, fileCtx), nil
}
func NewText(id string, context *vit.FileContext) *Text {
	t := &Text{
		Item:                NewItem("", context),
		id:                  id,
		text:                *vit.NewEmptyStringValue(),
		color:               *vit.NewColorValueFromCode(vit.Code{FileCtx: context, Code: "\"black\"", Position: nil}),
		horizontalAlignment: *vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "HorizontalAlignment.AlignLeft", Position: nil}),
		verticalAlignment:   *vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "VerticalAlignment.AlignTop", Position: nil}),
		font: *vit.NewEmptyGroupValue(map[string]vit.Value{
			"bold":      vit.NewBoolValueFromCode(vit.Code{FileCtx: context, Code: "false", Position: nil}),
			"italic":    vit.NewBoolValueFromCode(vit.Code{FileCtx: context, Code: "false", Position: nil}),
			"strikeout": vit.NewBoolValueFromCode(vit.Code{FileCtx: context, Code: "false", Position: nil}),
			"underline": vit.NewBoolValueFromCode(vit.Code{FileCtx: context, Code: "false", Position: nil}),
			"pixelSize": vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "12", Position: nil}),
			"pointSize": vit.NewFloatValueFromCode(vit.Code{FileCtx: context, Code: "12", Position: nil}),
			"family":    vit.NewStringValueFromCode(vit.Code{FileCtx: context, Code: "\"Arial\"", Position: nil}),
			"weight":    vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "FontWeight.Normal", Position: nil}),
		}),
		elide:        *vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "Elide.ElideNone", Position: nil}),
		fontData:     nil,
		fontFaceData: nil,
	}
	// property assignments on embedded components
	// register listeners for when a property changes
	t.font.AddDependent(vit.FuncDep(t.updateFont))
	// register event listeners
	// register enumerations
	t.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "HorizontalAlignment",
		Position: nil,
		Values:   map[string]int{"AlignLeft": 0, "AlignHCenter": 1, "AlignRight": 2},
	})
	t.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "VerticalAlignment",
		Position: nil,
		Values:   map[string]int{"AlignTop": 0, "AlignVCenter": 1, "AlignBottom": 2},
	})
	t.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "FontWeight",
		Position: nil,
		Values:   map[string]int{"Thin": 100, "ExtraLight": 200, "UltraLight": 200, "Light": 300, "Normal": 400, "Regular": 400, "Medium": 500, "DemiBold": 600, "SemiBold": 600, "Bold": 700, "ExtraBold": 800, "UltraBold": 800, "Black": 900, "Heavy": 900},
	})
	t.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "Elide",
		Position: nil,
		Values:   map[string]int{"ElideNone": 0, "ElideLeft": 1, "ElideMiddle": 2, "ElideRight": 3},
	})
	// add child components

	context.RegisterComponent("", t)

	return t
}

func (t *Text) String() string {
	return fmt.Sprintf("Text(%s)", t.id)
}

func (t *Text) Property(key string) (vit.Value, bool) {
	switch key {
	case "text":
		return &t.text, true
	case "color":
		return &t.color, true
	case "horizontalAlignment":
		return &t.horizontalAlignment, true
	case "verticalAlignment":
		return &t.verticalAlignment, true
	case "font":
		return &t.font, true
	case "elide":
		return &t.elide, true
	default:
		return t.Item.Property(key)
	}
}

func (t *Text) MustProperty(key string) vit.Value {
	v, ok := t.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (t *Text) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "text":
		err = t.text.SetValue(value)
	case "color":
		err = t.color.SetValue(value)
	case "horizontalAlignment":
		err = t.horizontalAlignment.SetValue(value)
	case "verticalAlignment":
		err = t.verticalAlignment.SetValue(value)
	case "font":
		err = t.font.SetValue(value)
	case "elide":
		err = t.elide.SetValue(value)
	default:
		return t.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Text", key, t.id, err)
	}
	return nil
}

func (t *Text) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "text":
		t.text.SetCode(code)
	case "color":
		t.color.SetCode(code)
	case "horizontalAlignment":
		t.horizontalAlignment.SetCode(code)
	case "verticalAlignment":
		t.verticalAlignment.SetCode(code)
	case "font":
		t.font.SetCode(code)
	case "elide":
		t.elide.SetCode(code)
	default:
		return t.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (t *Text) Event(name string) (vit.Listenable, bool) {
	switch name {
	default:
		return t.Item.Event(name)
	}
}

func (t *Text) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "text":
		return &t.text, true
	case "color":
		return &t.color, true
	case "horizontalAlignment":
		return &t.horizontalAlignment, true
	case "verticalAlignment":
		return &t.verticalAlignment, true
	case "font":
		return &t.font, true
	case "elide":
		return &t.elide, true
	default:
		return t.Item.ResolveVariable(key)
	}
}

func (t *Text) AddChild(child vit.Component) {
	child.SetParent(t)
	t.AddChildButKeepParent(child)
}

func (t *Text) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
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

func (t *Text) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	// properties
	if changed, err := t.text.Update(t); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Text", "text", t.id, err))
		}
	}
	if changed, err := t.color.Update(t); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Text", "color", t.id, err))
		}
	}
	if changed, err := t.horizontalAlignment.Update(t); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Text", "horizontalAlignment", t.id, err))
		}
	}
	if changed, err := t.verticalAlignment.Update(t); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Text", "verticalAlignment", t.id, err))
		}
	}
	if changed, err := t.font.Update(t); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Text", "font", t.id, err))
		}
	}
	if changed, err := t.elide.Update(t); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Text", "elide", t.id, err))
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

func (t *Text) As(target *vit.Component) bool {
	if _, ok := (*target).(*Text); ok {
		*target = t
		return true
	}
	return t.Item.As(target)
}

func (t *Text) ID() string {
	return t.id
}

func (t *Text) Finish() error {
	return t.RootC().FinishInContext(t)
}

func (t *Text) staticAttribute(name string) (interface{}, bool) {
	switch name {
	case "AlignLeft":
		return uint(Text_HorizontalAlignment_AlignLeft), true
	case "AlignHCenter":
		return uint(Text_HorizontalAlignment_AlignHCenter), true
	case "AlignRight":
		return uint(Text_HorizontalAlignment_AlignRight), true
	case "AlignTop":
		return uint(Text_VerticalAlignment_AlignTop), true
	case "AlignVCenter":
		return uint(Text_VerticalAlignment_AlignVCenter), true
	case "AlignBottom":
		return uint(Text_VerticalAlignment_AlignBottom), true
	case "Thin":
		return uint(Text_FontWeight_Thin), true
	case "ExtraLight":
		return uint(Text_FontWeight_ExtraLight), true
	case "UltraLight":
		return uint(Text_FontWeight_UltraLight), true
	case "Light":
		return uint(Text_FontWeight_Light), true
	case "Normal":
		return uint(Text_FontWeight_Normal), true
	case "Regular":
		return uint(Text_FontWeight_Regular), true
	case "Medium":
		return uint(Text_FontWeight_Medium), true
	case "DemiBold":
		return uint(Text_FontWeight_DemiBold), true
	case "SemiBold":
		return uint(Text_FontWeight_SemiBold), true
	case "Bold":
		return uint(Text_FontWeight_Bold), true
	case "ExtraBold":
		return uint(Text_FontWeight_ExtraBold), true
	case "UltraBold":
		return uint(Text_FontWeight_UltraBold), true
	case "Black":
		return uint(Text_FontWeight_Black), true
	case "Heavy":
		return uint(Text_FontWeight_Heavy), true
	case "ElideNone":
		return uint(Text_Elide_ElideNone), true
	case "ElideLeft":
		return uint(Text_Elide_ElideLeft), true
	case "ElideMiddle":
		return uint(Text_Elide_ElideMiddle), true
	case "ElideRight":
		return uint(Text_Elide_ElideRight), true
	default:
		return nil, false
	}
}
