// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
)

type Row struct {
	Item
	id string

	topPadding    vit.OptionalValue[*vit.FloatValue]
	rightPadding  vit.OptionalValue[*vit.FloatValue]
	bottomPadding vit.OptionalValue[*vit.FloatValue]
	leftPadding   vit.OptionalValue[*vit.FloatValue]
	padding       vit.FloatValue
	spacing       vit.FloatValue
	childLayouts  vit.LayoutList
}

func NewRow(id string, scope vit.ComponentContainer) *Row {
	r := &Row{
		Item:          *NewItem(id, scope),
		id:            id,
		topPadding:    *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		rightPadding:  *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		bottomPadding: *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		leftPadding:   *vit.NewOptionalValue(vit.NewFloatValueFromExpression("0", nil)),
		padding:       *vit.NewFloatValueFromExpression("0", nil),
		spacing:       *vit.NewFloatValueFromExpression("0", nil),
		childLayouts:  make(vit.LayoutList),
	}
	return r
}

func (r *Row) String() string {
	return fmt.Sprintf("Row(%s)", r.id)
}

func (r *Row) Property(key string) (vit.Value, bool) {
	switch key {
	case "topPadding":
		return &r.topPadding, true
	case "rightPadding":
		return &r.rightPadding, true
	case "bottomPadding":
		return &r.bottomPadding, true
	case "leftPadding":
		return &r.leftPadding, true
	case "padding":
		return &r.padding, true
	case "spacing":
		return &r.spacing, true
	default:
		return r.Item.Property(key)
	}
}

func (r *Row) MustProperty(key string) vit.Value {
	v, ok := r.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (r *Row) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "topPadding":
		err = r.topPadding.SetValue(value)
	case "rightPadding":
		err = r.rightPadding.SetValue(value)
	case "bottomPadding":
		err = r.bottomPadding.SetValue(value)
	case "leftPadding":
		err = r.leftPadding.SetValue(value)
	case "padding":
		err = r.padding.SetValue(value)
	case "spacing":
		err = r.spacing.SetValue(value)
	default:
		return r.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Row", key, r.id, err)
	}
	return nil
}

func (r *Row) SetPropertyExpression(key string, code string, pos *vit.PositionRange) error {
	switch key {
	case "topPadding":
		r.topPadding.SetExpression(code, pos)
	case "rightPadding":
		r.rightPadding.SetExpression(code, pos)
	case "bottomPadding":
		r.bottomPadding.SetExpression(code, pos)
	case "leftPadding":
		r.leftPadding.SetExpression(code, pos)
	case "padding":
		r.padding.SetExpression(code, pos)
	case "spacing":
		r.spacing.SetExpression(code, pos)
	default:
		return r.Item.SetPropertyExpression(key, code, pos)
	}
	return nil
}

func (r *Row) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case r.id:
		return r, true
	case "topPadding":
		return &r.topPadding, true
	case "rightPadding":
		return &r.rightPadding, true
	case "bottomPadding":
		return &r.bottomPadding, true
	case "leftPadding":
		return &r.leftPadding, true
	case "padding":
		return &r.padding, true
	case "spacing":
		return &r.spacing, true
	default:
		return r.Item.ResolveVariable(key)
	}
}

func (r *Row) AddChild(child vit.Component) {
	defer r.childWasAdded(child)
	child.SetParent(r)
	r.AddChildButKeepParent(child)
}

func (r *Row) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	defer r.childWasAdded(addThis)
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

func (r *Row) UpdateExpressions() (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if changed, err := r.topPadding.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "topPadding", r.id, err))
		}
		r.recalculateLayout(r.topPadding)
	}
	if changed, err := r.rightPadding.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "rightPadding", r.id, err))
		}
		r.recalculateLayout(r.rightPadding)
	}
	if changed, err := r.bottomPadding.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "bottomPadding", r.id, err))
		}
		r.recalculateLayout(r.bottomPadding)
	}
	if changed, err := r.leftPadding.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "leftPadding", r.id, err))
		}
		r.recalculateLayout(r.leftPadding)
	}
	if changed, err := r.padding.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "padding", r.id, err))
		}
		r.recalculateLayout(r.padding)
	}
	if changed, err := r.spacing.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "spacing", r.id, err))
		}
		r.recalculateLayout(r.spacing)
	}
	if changed, err := r.childLayouts.Update(r); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Row", "childLayouts", r.id, err))
		}
		r.recalculateLayout(r.childLayouts)
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

func (r *Row) As(target *vit.Component) bool {
	if _, ok := (*target).(*Row); ok {
		*target = r
		return true
	}
	return r.Item.As(target)
}

func (r *Row) ID() string {
	return r.id
}

func (r *Row) Finish() error {
	return r.RootC().FinishInContext(r)
}