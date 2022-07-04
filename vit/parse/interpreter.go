package parse

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

type componentError struct {
	component vit.AbstractComponent
	err       error
}

func (e componentError) Error() string {
	var cErr componentError
	if errors.As(e.err, &cErr) {
		return fmt.Sprintf("%s > %s", e.component.Name(), e.err) // looks nicer
	}
	return fmt.Sprintf("%s: %s", e.component.Name(), e.err)
}

func (e componentError) Is(target error) bool {
	_, ok := target.(componentError)
	return ok
}

func (e componentError) Unwrap() error {
	return e.err
}

type unknownComponentError struct {
	name string
}

func (e unknownComponentError) Error() string {
	return fmt.Sprintf("unknown component %q", e.name)
}

func (e unknownComponentError) Is(target error) bool {
	_, ok := target.(unknownComponentError)
	return ok
}

type genericError struct {
	position vit.PositionRange
	err      error
}

func genericErrorf(position vit.PositionRange, format string, args ...interface{}) error {
	return genericError{position, fmt.Errorf(format, args...)}
}

func (e genericError) Error() string {
	if e.err == nil {
		return "no error"
	}
	return e.err.Error()
}

func (e genericError) Is(target error) bool {
	_, ok := target.(genericError)
	return ok
}

func (e genericError) Unwrap() error {
	return e.err
}

func init() {
	vit.InstantiateComponent = instantiateComponent
}

// parseFile parsed a given file into a document with the given component name.
func parseFile(fileName string, componentName string) (*VitDocument, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lexer := NewLexer(file, fileName)

	doc, err := Parse(NewTokenBuffer(lexer.Lex))
	if err != nil {
		return nil, err
	}
	doc.Name = componentName

	return doc, nil
}

// interpret takes the parsed document and creates the appropriate component tree.
func interpret(document VitDocument, id string, context vit.ComponentContext) ([]vit.Component, error) {
	for _, imp := range document.Imports {
		if len(imp.file) != 0 {
			// file import
			return nil, genericErrorf(imp.position, "not yet implemented")
		} else if len(imp.namespace) != 0 {
			// namespace import
			lib, err := resolveLibraryImport(imp.namespace)
			if err != nil {
				return nil, ParseError{imp.position, err}
			}
			for _, name := range lib.ComponentNames() {
				context.KnownComponents.Set(name, &LibraryInstantiator{lib, name})
			}
		} else {
			return nil, genericErrorf(imp.position, "incomplete namespace")
		}
	}

	var instances []vit.Component
	for _, comp := range document.Components {
		instance, err := instantiateCustomComponent(comp, id, document.Name, context)
		if err != nil {
			return nil, err
		}
		instances = append(instances, instance)
	}

	return instances, nil
}

// instantiateCustomComponent creates a component described by a componentDefinition and wraps it in a Custom component with the given id.
func instantiateCustomComponent(def *vit.ComponentDefinition, id string, name string, context vit.ComponentContext) (vit.Component, error) {
	comp, err := instantiateComponent(def, context)
	if err != nil {
		return nil, err
	}

	cst := vit.NewCustom(id, name, comp)

	return cst, nil
}

// instantiateComponent creates a component described by a componentDefinition.
func instantiateComponent(def *vit.ComponentDefinition, context vit.ComponentContext) (vit.Component, error) {
	src, ok := context.KnownComponents.Get(def.BaseName)
	if !ok {
		// TODO: improve context for error; either here or upstream
		return nil, unknownComponentError{def.BaseName}
	}
	instance, err := src.Instantiate(def.ID, context)
	if err != nil {
		return nil, componentError{src, err}
	}

	err = populateComponent(instance, def, context)
	if err != nil {
		return instance, componentError{src, err}
	}

	context.Environment.RegisterComponent(instance)
	// TODO: figure out where the components will be unregistered again

	return instance, nil
}

// populateComponent takes a fresh component instance as well as it's definition and populates all attributes and children with their correct values.
func populateComponent(instance vit.Component, def *vit.ComponentDefinition, context vit.ComponentContext) error {
	for _, enum := range def.Enumerations {
		if !instance.DefineEnum(enum) {
			return genericErrorf(*enum.Position, "enum %q already defined", enum.Name)
		}
	}

	for _, prop := range def.Properties {
		if prop.VitType != "" {
			// this defines a new property
			if err := instance.DefineProperty(prop); err != nil {
				return err
			}
			// instance.SetProperty(prop.identifier[0], prop.expression)
		} else if len(prop.Identifier) == 1 {
			// simple property assignment
			var err error
			if len(prop.Components) == 0 {
				err = instance.SetPropertyExpression(prop.Identifier[0], prop.Expression, &prop.Pos)
			} else if len(prop.Components) == 1 {
				err = instance.SetProperty(prop.Identifier[0], prop.Components[0])
			} else {
				err = instance.SetProperty(prop.Identifier[0], prop.Components)
			}
			if err != nil {
				return genericError{prop.Pos, err}
			}
		} else {
			// assign property with qualifier
			// TODO: make this universal?
			if prop.Identifier[0] == "anchors" {
				v, ok := instance.Property(prop.Identifier[0])
				if !ok {
					return genericErrorf(prop.Pos, "unknown property %q of component %q", prop.Identifier[0], def.BaseName)
				}
				anchors, ok := v.(*vit.AnchorsValue)
				if !ok {
					return genericErrorf(prop.Pos, "cannot assign to non group-property %q of component %q", prop.Identifier[0], def.BaseName)
				}

				ok = anchors.SetPropertyExpression(prop.Identifier[1], prop.Expression, &prop.Pos)
				if !ok {
					return genericErrorf(prop.Pos, "unknown property %q of component %q", strings.Join(prop.Identifier, "."), def.BaseName)
				}
			} else {
				v, ok := instance.Property(prop.Identifier[0])
				if !ok {
					return genericErrorf(prop.Pos, "unknown property %q of component %q", prop.Identifier[0], def.BaseName)
				}
				group, ok := v.(*vit.GroupValue)
				if !ok {
					return genericErrorf(prop.Pos, "cannot assign to non group-property %q of component %q", prop.Identifier[0], def.BaseName)
				}

				err := group.SetExpressionOf(prop.Identifier[1], prop.Expression, &prop.Pos)
				if err != nil {
					return genericErrorf(prop.Pos, "group-property %q of component %q: %w", prop.Identifier[0], def.BaseName, err)
				}
			}
		}
	}

	for _, childDef := range def.Children {
		childInstance, err := instantiateComponent(childDef, context)
		if err != nil {
			return err
		}
		instance.AddChild(childInstance)
	}

	return nil
}
