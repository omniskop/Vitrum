package std

import (
	"fmt"

	vit "github.com/omniskop/vitrum/vit"
)

type Rectangle struct {
	Item
	id string

	color  vit.ColorValue
	radius vit.FloatValue
	border vit.GroupValue
}

func NewRectangle(id string, scope vit.ComponentContainer) *Rectangle {
	return &Rectangle{
		Item:   *NewItem(id, scope),
		id:     id,
		color:  *vit.NewEmptyColorValue(),
		radius: *vit.NewEmptyFloatValue(),
		border: *vit.NewEmptyGroupValue(map[string]vit.Value{
			"color": vit.NewEmptyColorValue(),
			"width": vit.NewEmptyIntValue(),
		}),
	}
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle(%s)", r.id)
}

func (r *Rectangle) Property(key string) (vit.Value, bool) {
	switch key {
	case "color":
		return &r.color, true
	case "radius":
		return &r.radius, true
	case "border":
		return &r.border, true
	default:
		return r.Item.Property(key)
	}
}

func (r *Rectangle) MustProperty(key string) vit.Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Rectangle) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "color":
		err = r.color.SetValue(value)
	case "radius":
		err = r.radius.SetValue(value)
	case "border":
		r.border.SetValue(value)
	default:
		return r.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Rectangle", key, r.id, err)
	}
	return nil
}

func (r *Rectangle) SetPropertyExpression(key string, code string, pos *vit.PositionRange) error {
	switch key {
	case "color":
		r.color.SetExpression(code, pos)
	case "radius":
		r.radius.SetExpression(code, pos)
	case "border":
		r.border.SetExpression(code, pos)
	default:
		return r.Item.SetPropertyExpression(key, code, pos)
	}
	return nil
}

func (r *Rectangle) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case r.id:
		return r, true
	case "color":
		return &r.color, true
	case "radius":
		return &r.radius, true
	case "border":
		return &r.border, true
	default:
		return r.Item.ResolveVariable(key)
	}
}

func (r *Rectangle) AddChild(child vit.Component) {
	child.SetParent(r)
	r.AddChildButKeepParent(child)
}

func (r *Rectangle) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	var targetType vit.Component = afterThis

	for ind, child := range r.Children() {
		if child.As(&targetType) {
			addThis.SetParent(r)
			r.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	r.AddChild(addThis)
}

func (r *Rectangle) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if changed, err := r.color.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rectangle", "color", r.id, err))
		}
	}
	if changed, err := r.radius.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rectangle", "radius", r.id, err))
		}
	}
	if changed, err := r.border.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Rectangle", "border", r.id, err))
		}
	}

	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	n, err := r.UpdatePropertiesInContext(r)
	sum += n
	errs.AddGroup(err)
	n, err = r.Item.UpdateExpressions()
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (r *Rectangle) As(target *vit.Component) bool {
	if _, ok := (*target).(*Rectangle); ok {
		*target = r
		return true
	}
	return r.Item.As(target)
}

func (r *Rectangle) ID() string {
	return r.id
}

func (r *Rectangle) Finish() error {
	return r.RootC().FinishInContext(r)
}
