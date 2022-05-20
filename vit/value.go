package vit

import (
	"fmt"
	"strconv"
	"strings"
)

type Value interface {
	SetFromProperty(PropertyDefinition)
	Update(context Component) error
	GetValue() interface{}
	MakeDirty([]*Expression)
	GetExpression() *Expression
	AddDependent(*Expression)
	RemoveDependent(*Expression)
	ShouldEvaluate() bool
	Err() error
}

func ValueConstructorForType(vitType string, value interface{}, position *PositionRange) (Value, error) {
	switch vitType {
	case "string":
		return NewStringValue(value.(string), position), nil
	case "int":
		return NewIntValue(value.(string), position), nil
	case "float":
		return NewFloatValue(value.(string), position), nil
	case "var":
		return NewAnyValue(value.(string), position), nil
	case "component":
		return NewComponentValue(value.(*ComponentDefinition), position), nil
	}
	return nil, UnknownTypeError{vitType}
}

// ========================================= List Value ============================================

type ListValue[ElementType Value] struct {
	Expression
	Value []ElementType
}

func NewListValue[ElementType Value](expression string, position *PositionRange) *ListValue[ElementType] {
	v := new(ListValue[ElementType])
	if expression == "" {
		v.Expression = *NewExpression("[]", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (c *ListValue[ElementType]) SetFromProperty(prop PropertyDefinition) {
	c.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *ListValue[ElementType]) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	castVal, ok := castList[ElementType](val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type list but %T instead", val)
	}
	v.Value = castVal
	return nil
}

func (v *ListValue[ElementType]) GetValue() interface{} {
	var out []interface{}
	for _, element := range v.Value {
		out = append(out, element.GetValue())
	}
	return out
}

// ========================================== Int Value ============================================

type IntValue struct {
	Expression
	Value int
}

func NewIntValue(expression string, position *PositionRange) *IntValue {
	v := new(IntValue)
	if expression == "" {
		v.Expression = *NewExpression("0", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *IntValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *IntValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	castVal, ok := castInt(val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type int but %T instead", val)
	}
	v.Value = int(castVal)
	return nil
}

func (v *IntValue) GetValue() interface{} {
	return v.Value
}

// ========================================= Float Value ===========================================

type FloatValue struct {
	Expression
	Value float64
}

func NewFloatValue(expression string, position *PositionRange) *FloatValue {
	v := new(FloatValue)
	if expression == "" {
		v.Expression = *NewExpression("0", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *FloatValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *FloatValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	castVal, ok := castFloat(val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type int but %T instead", val)
	}
	v.Value = float64(castVal)
	return nil
}

func (c *FloatValue) GetValue() interface{} {
	return c.Value
}

// ======================================== String Value ===========================================

type StringValue struct {
	Expression
	Value string
}

func NewStringValue(expression string, position *PositionRange) *StringValue {
	v := new(StringValue)
	if expression == "" {
		v.Expression = *NewExpression(`""`, position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *StringValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *StringValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	var ok bool
	v.Value, ok = convertJSValueToString(val)
	if !ok {
		return fmt.Errorf("did not evaluate to expected type string but %T instead", val)
	}
	return nil
}

func (c *StringValue) GetValue() interface{} {
	return c.Value
}

func convertJSValueToString(value interface{}) (string, bool) {
	switch actual := value.(type) {
	case string:
		return actual, true
	case int64:
		return strconv.FormatInt(actual, 10), true
	case float64:
		return strconv.FormatFloat(actual, 'f', 64, 10), true
	}
	return "", false
}

// ========================================= Alias Value ===========================================

// Alias value points to a different value, potentially of another component
type AliasValue struct {
	Expression string
	Position   *PositionRange
	other      Value
}

func NewAliasValue(expression string, position *PositionRange) *AliasValue {
	v := new(AliasValue)
	v.Position = position
	v.Expression = expression
	return v
}

func (v *AliasValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression = prop.Expression
}

func (v *AliasValue) Update(context Component) error {
	if v.other != nil {
		return nil
	}
	if v.Expression == "" {
		// TODO: add position to error
		return fmt.Errorf("alias reference is empty")
	}
	parts := strings.Split(v.Expression, ".")
	var currentComponent Component = context
	var currentProperty Value
	// find component using the id's listed in the expression
	for {
		part := parts[0]

		if strings.Contains(part, " ") {
			return fmt.Errorf("invalid alias reference: %q", v.Expression)
		}

		if currentComponent.ID() == part {
			parts = parts[1:]
			continue // no change
		}
		if childComp, ok := currentComponent.ResolveID(part); ok {
			currentComponent = childComp
			parts = parts[1:]
			continue
		}
		break
	}
	// find property using the remaining parts
	for _, part := range parts {
		val, ok := currentComponent.Property(part)
		if !ok {
			return fmt.Errorf("unable to resolve alias reference: %q", v.Expression)
		}
		currentProperty = val
	}

	// nothing found
	if currentProperty == nil {
		return fmt.Errorf("unable to resolve alias reference: %q", v.Expression)
	}
	// referenced itself
	if currentProperty == v {
		return fmt.Errorf("alias cannot reference itself: %q", v.Expression)
	}

	v.other = currentProperty // saving this also marks the alias as updated, preventing an infinite loop in the next check

	// if we referenced another alias we need will update that as well and make sure there are no circular references
	if otherAlias, ok := currentProperty.(*AliasValue); ok {
		err := otherAlias.Update(currentComponent)
		if err != nil {
			return fmt.Errorf("error in nested alias update: %w", err)
		}
		if yes, chain := isAliasRecursive(v, nil); yes {
			return fmt.Errorf("alias contains circular reference: %v", formatAliasChain(chain))
		}
	}

	return nil
}

func isAliasRecursive(alias *AliasValue, chain []*AliasValue) (bool, []*AliasValue) {
	if subAlias, ok := alias.other.(*AliasValue); ok {
		for _, a := range chain {
			if a == subAlias {
				return true, append(chain, alias, subAlias)
			}
		}
		return isAliasRecursive(subAlias, append(chain, alias))
	}

	return false, nil
}

func formatAliasChain(chain []*AliasValue) string {
	var steps []string
	for _, a := range chain {
		steps = append(steps, fmt.Sprintf("%q", a.Expression))

	}
	return strings.Join(steps, " -> ")
}

func (v *AliasValue) GetValue() interface{} {
	if v.other == nil {
		return nil
	}
	return v.other.GetValue()
}

func (v *AliasValue) MakeDirty(stack []*Expression) {}

func (v *AliasValue) GetExpression() *Expression {
	return NewExpression(v.Expression, v.Position)
}

func (v *AliasValue) AddDependent(exp *Expression) {
	if v.other == nil {
		v.other.AddDependent(exp)
	}
}

func (v *AliasValue) RemoveDependent(exp *Expression) {
	if v.other == nil {
		v.other.RemoveDependent(exp)
	}
}

func (v *AliasValue) ShouldEvaluate() bool {
	return v.other == nil || v.other.ShouldEvaluate()
}

func (v *AliasValue) Err() error {
	return nil
}

// ========================================== Any Value ============================================

type AnyValue struct {
	Expression
	Value interface{}
}

func NewAnyValue(expression string, position *PositionRange) *AnyValue {
	v := new(AnyValue)
	if expression == "" {
		v.Expression = *NewExpression("null", position)
	} else {
		v.Expression = *NewExpression(expression, position)
	}
	return v
}

func (v *AnyValue) SetFromProperty(prop PropertyDefinition) {
	v.Expression.ChangeCode(prop.Expression, &prop.Pos)
}

func (v *AnyValue) Update(context Component) error {
	val, err := v.Expression.Evaluate(context)
	if err != nil {
		if err == unsettledDependenciesError {
			return nil
		}
		return err
	}
	v.Value = val
	return nil
}

func (c *AnyValue) GetValue() interface{} {
	return c.Value
}

// ======================================= Component Value =========================================

type ComponentValue struct {
	Value   *ComponentDefinition
	Changed bool
	err     error
}

func NewComponentValue(component *ComponentDefinition, position *PositionRange) *ComponentValue {
	return &ComponentValue{
		Value:   component,
		Changed: true,
	}
}

func (v *ComponentValue) ChangeComponent(component *ComponentDefinition) {
	v.Value = component
	v.Changed = true
	v.err = nil
}

func (v *ComponentValue) SetFromProperty(prop PropertyDefinition) {
	if len(prop.Components) == 0 {
		v.Value = nil
		v.err = nil
	} else if len(prop.Components) == 1 {
		v.Value = prop.Components[0]
		v.err = nil
	} else {
		v.Value = prop.Components[0]
		v.err = fmt.Errorf("cannot assign multiple components to a single component value at %s", prop.Pos.String())
	}
	v.Changed = true
}

func (v *ComponentValue) Update(context Component) error {
	v.Changed = false
	return v.err
}

func (v *ComponentValue) GetValue() interface{} {
	return v.Value
}

func (v *ComponentValue) MakeDirty(stack []*Expression) {
	v.Changed = true
}

func (v *ComponentValue) GetExpression() *Expression {
	return NewExpression("", nil)
}

func (v *ComponentValue) AddDependent(exp *Expression) {}

func (v *ComponentValue) RemoveDependent(exp *Expression) {}

func (v *ComponentValue) ShouldEvaluate() bool {
	return v.Changed
}

func (v *ComponentValue) Err() error {
	return nil
}

// ========================================= Static List ===========================================

type StaticBaseValue struct {
	Changed bool
}

func (v *StaticBaseValue) SetFromProperty(prop PropertyDefinition) {
	v.Changed = true
}

func (v *StaticBaseValue) Update(context Component) error {
	v.Changed = false
	return nil
}

func (v *StaticBaseValue) MakeDirty(stack []*Expression) {
	v.Changed = true
}

func (v *StaticBaseValue) GetExpression() *Expression {
	return NewExpression("", nil)
}

func (v *StaticBaseValue) AddDependent(exp *Expression) {}

func (v *StaticBaseValue) RemoveDependent(exp *Expression) {}

func (v *StaticBaseValue) ShouldEvaluate() bool {
	return v.Changed
}

func (v *StaticBaseValue) Err() error {
	return nil
}

type StaticListValue[ElementType Value] struct {
	StaticBaseValue
	Items []ElementType
}

func NewStaticListValue[ElementType Value](items []ElementType, position *PositionRange) *StaticListValue[ElementType] {
	return &StaticListValue[ElementType]{
		StaticBaseValue: StaticBaseValue{true},
		Items:           items,
	}
}

func (v *StaticListValue[ElementType]) GetValue() interface{} {
	return v.Items
}

func (v *StaticListValue[ElementType]) Set(value []ElementType) {
	v.Items = value
	v.Changed = true
}
