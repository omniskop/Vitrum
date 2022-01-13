package vit

type custom struct {
	Component
	RootComponent *Root
	id            string
	name          string
}

func NewCustom(id string, name string, parent Component) *custom {
	return &custom{
		Component:     parent,
		RootComponent: parent.root(),
		id:            id,
		name:          name,
	}
}

func (c *custom) ResolveVariable(key string) (interface{}, bool) {
	if key == c.id {
		return c, true
	}

	return c.Component.ResolveVariable(key)
}

func (c *custom) AddChild(child Component) {
	child.SetParent(c)
	c.RootComponent.children = append(c.RootComponent.children, child)
}

func (c *custom) UpdateExpressions() (int, ErrorGroup) {
	var errs ErrorGroup
	var sum int
	// this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables
	for name, prop := range c.RootComponent.properties {
		if prop.ShouldEvaluate() {
			sum++
			err := prop.Update(c)
			if err != nil {
				errs.Add(newExpressionError(c.name, name, c.id, *prop.GetExpression(), err))
			}
		}
	}

	s, e := c.Component.UpdateExpressions()
	sum += s
	errs.AddGroup(e)

	return sum, errs
}

func (c *custom) ID() string {
	return c.id
}
