package generator

import (
	"sort"

	"github.com/dave/jennifer/jen"
	"github.com/omniskop/vitrum/vit"
)

// This file contains various helper functions that make generating code easier.

// mapProperties calls function 'f' for every property and combines their generated codes into one block.
// It skips properties that are considered to be internal.
// 'f' will be called with it's property definition and identifier of the property.
// It should return the corresponding code for the property. It may return nil to indicate that no code should be added.
func mapProperties(props []vit.PropertyDefinition, f func(vit.PropertyDefinition, string) jen.Code) []jen.Code {
	var result []jen.Code
	for _, prop := range props {
		if isInternalProperty(prop) {
			continue
		}
		// NOTE: properties with multiple identifiers are currently not supported
		result = append(result, f(prop, prop.Identifier[0]))
	}
	return result
}

func prepend(first jen.Code, rest []jen.Code, tail ...jen.Code) []jen.Code {
	return append([]jen.Code{first}, append(rest, tail...)...)
}

// addMultiple adds all code to the group as separate statements
func addMultiple(g *jen.Group, code []jen.Code) {
	for _, c := range code {
		if c != nil {
			g.Add(c)
		}
	}
}

// generateCallbackForAddedChild checks if the component provides a callback for the event that a child has been added to the component.
// If that is the case it returns a defer statement that should be added to the top of all methods that add a child.
// If no callback is provided it returns nil.
func generateCallbackForAddedChild(comp *vit.ComponentDefinition, receiverName, parameterName string) jen.Code {
	if childrenProp, ok := getProperty(comp, "children"); ok && childrenProp.HasTag(onChangeTag) {
		// if the children property is explicitly provided and has the 'onchange' tag we will call the provided method with the added child.
		return jen.Defer().Id(receiverName).Dot(childrenProp.Tags[onChangeTag]).Call(jen.Id(parameterName))
	}
	return nil
}

// unpointer removes the first '*' from statements.
// TODO: currently this works by just always removing the first code without checking if that's actually a '*'.
func unpointer(code *jen.Statement) *jen.Statement {
	statement := (*code)[1:]
	return &statement
}

// isReadable return true if the property can be read from outside of the component.
func isReadable(prop vit.PropertyDefinition) bool {
	return !prop.HasTag(privateTag)
}

// isWritable returns true if the property can be written to from outside of the component.
func isWritable(prop vit.PropertyDefinition) bool {
	return !prop.HasTag(privateTag) && !prop.ReadOnly
}

// isInternalProperty returns true if the property is internal and no code should be generated for it.
func isInternalProperty(prop vit.PropertyDefinition) bool {
	return prop.Identifier[0] == "children"
}

// getProperty returns the property with the given name from the component.
// The book indicated that a property was found.
func getProperty(prop *vit.ComponentDefinition, identifier string) (*vit.PropertyDefinition, bool) {
	for _, p := range prop.Properties {
		if p.Identifier[0] == identifier {
			return &p, true
		}
	}
	return nil, false
}

// orderEnumValues orders an enum map by it's values (and by it's name if values are equal).
func orderEnumValues(values map[string]int) []enumValue {
	var list = make(enumValueList, 0, len(values))
	for k, v := range values {
		list = append(list, enumValue{k, v})
	}
	sort.Sort(list)
	return list
}

type enumValue struct {
	name  string
	value int
}

type enumValueList []enumValue

func (v enumValueList) Len() int {
	return len(v)
}

func (v enumValueList) Less(i, j int) bool {
	if v[i].value == v[j].value {
		// if values are equal use name as fallback
		return v[i].name < v[j].name
	}
	return v[i].value < v[j].value
}

func (v enumValueList) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
