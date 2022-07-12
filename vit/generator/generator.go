// Generates go code that defines a component described by a vit file.
package generator

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/parse"
)

const vitPackage = "github.com/omniskop/vitrum/vit"
const parsePackage = "github.com/omniskop/vitrum/vit/parse"
const stdPackage = "github.com/omniskop/vitrum/vit/std"

// Indicates that file locations in the source file should not be brought into the generated code
const hideSourceFiles = true

// All generator specific tags
const (
	onChangeTag    = "gen-onchange"
	internalTag    = "gen-internal"
	typeTag        = "gen-type"
	initializerTag = "gen-initializer"
	privateTag     = "gen-private"
	optionalTag    = "gen-optional"
)

// GenerateFromFileAndSave takes the path to a vit file and generates a go file at the destination path.
// The package name defines to which package the go file should belong.
func GenerateFromFileAndSave(srcPath string, packageName string, dstPath string) error {
	doc, err := parseVit(srcPath)
	if err != nil {
		return fmt.Errorf("unable to parse: %v", err)
	}

	// We first store the generated output in a buffer and then write it into the file afterwards
	// to make sure that we don't overwrite an existing file if the generation failed.
	var output bytes.Buffer

	err = GenerateFromDocument(doc, packageName, &output)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("unable to create destination file: %w", err)
	}
	_, err = dstFile.Write(output.Bytes())
	if err != nil {
		return fmt.Errorf("unable to write to destination file: %w", err)
	}

	return nil
}

// Generate takes a reader for a vit file and outputs generated go code into the destination writer.
// It uses the source path to extract the component name and to provide useful error messages.
// The package name defines to which package the go file should belong.
func Generate(src io.Reader, srcPath string, packageName string, dst io.Writer) error {
	lexer := parse.NewLexer(src, srcPath)

	doc, err := parse.Parse(parse.NewTokenBuffer(lexer.Lex))
	if err != nil {
		return err
	}
	doc.Name = getComponentName(srcPath)

	return GenerateFromDocument(doc, packageName, dst)
}

// GenerateFromFile takes the path of a vit file and outputs generated go code into the destination writer.
// The package name defines to which package the go file should belong.
func GenerateFromFile(srcFile string, packageName string, dst io.Writer) error {
	doc, err := parseVit(srcFile)
	if err != nil {
		return err
	}

	return GenerateFromDocument(doc, packageName, dst)
}

// parseVit takes the path to a vit file. It reads, lexes and parses the file.
func parseVit(srcFile string) (*parse.VitDocument, error) {
	file, err := os.Open(srcFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lexer := parse.NewLexer(file, srcFile)

	doc, err := parse.Parse(parse.NewTokenBuffer(lexer.Lex))
	if err != nil {
		return nil, err
	}
	doc.Name = getComponentName(srcFile)

	return doc, nil
}

// GenerateFromDocument takes a vit document and outputs generated go code into the destination writer.
// The package name defines to which package the go file should belong.
func GenerateFromDocument(doc *parse.VitDocument, packageName string, dst io.Writer) error {
	f := jen.NewFilePath(packageName)
	f.HeaderComment("Code generated by vitrum gencmd. DO NOT EDIT.")

	for _, comp := range doc.Components {
		err := generateComponent(f, doc.Name, comp)
		if err != nil {
			return err
		}
	}

	return f.Render(dst)
}

// getComponentName takes the path of a vit file and extracts the appropriate component name from it.
// For example: "path/to/file.vit" will result in the component name "File".
func getComponentName(fileName string) string {
	fileName = filepath.Base(fileName)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

// generateComponent generates the code for a full component definition under the given component name.
// The resulting code is added to the given jen file.
func generateComponent(f *jen.File, compName string, comp *vit.ComponentDefinition) error {
	// TODO: Figure out from which package the base component should be imported from. Currently this is hardcoded to be the std package.

	f.Add(generateComponentEnums(compName, comp))

	properties := []jen.Code{
		jen.Op("*").Qual(stdPackage, comp.BaseName),
		jen.Id("id").String(),
		jen.Line(),
	}

	// All property instantiations
	// we could use jen.Dict here but I wan't to preserve the property order
	propertyInstantiations := []jen.Code{
		jen.Line().Id(comp.BaseName).Op(":").Qual(stdPackage, fmt.Sprintf("New%s", comp.BaseName)).Op("(").Id("id").Op(",").Id("context").Op(")"),
		jen.Line().Id("id").Op(":").Id("id"),
	}

	// name of the variable that will hold the receiver of the components methods
	receiverName := strings.ToLower(string(compName[0]))

	// setup all properties for the struct definition as well as the instantiations
	for _, prop := range comp.Properties {
		if isInternalProperty(prop) || !prop.IsNewDefinition() {
			continue
		}
		propType, propConstructor, err := vitTypeInfo(comp, prop)
		if err != nil {
			return err
		}

		properties = append(properties, jen.Id(prop.Identifier[0]).Add(propType))

		// property instantiation
		propertyInstantiations = append(propertyInstantiations, jen.Line().Id(prop.Identifier[0]).Op(":").Add(propConstructor))
	}
	properties = append(properties, jen.Line())

	// setup all event attributes
	for _, ev := range comp.Events {
		var eventType *jen.Statement
		if len(ev.Parameters) == 0 {
			return fmt.Errorf("event %q has no parameters", ev.Name)
		}
		if typeName, ok := ev.Parameters[0].Tags[typeTag]; ok {
			eventType = generateCustomType(typeName)
		} else {
			var err error
			eventType, _, err = vitTypeInfo(comp, ev.Parameters[0])
			if err != nil {
				return err
			}
		}
		properties = append(properties, jen.Id(ev.Name).Qual(vitPackage, "EventAttribute").Types(eventType))
		propertyInstantiations = append(propertyInstantiations, jen.Line().Id(ev.Name).Op(":").Op("*").Qual(vitPackage, "NewEventAttribute").Types(eventType).Call())
	}
	properties = append(properties, jen.Line())

	// setup all method attributes
	for _, m := range comp.Methods {
		properties = append(properties, jen.Id(m.Name).Qual(vitPackage, "Method"))
		propertyInstantiations = append(propertyInstantiations, jen.Line().Id(m.Name).Op(":").Qual(vitPackage, "NewMethod").Call(jen.Lit(m.Name), jen.Lit(m.Code()), generatePositionRange(*m.Position)))
	}

	propertyInstantiations = append(propertyInstantiations, jen.Line())

	f.Type().Id(compName).Struct(properties...)

	// constructor
	f.Func().
		Id(fmt.Sprintf("New%s", compName)).
		Params(jen.Id("id").String(), jen.Id("context").Qual(vitPackage, "ComponentContext")).
		Params(jen.Op("*").Id(compName)).
		BlockFunc(func(g *jen.Group) {
			g.Id(receiverName).Op(":=").Op("&").Id(compName).Values(propertyInstantiations...)
			// property assignments
			g.Comment("property assignments on embedded components")
			for _, prop := range comp.Properties {
				if isInternalProperty(prop) || prop.IsNewDefinition() || len(prop.Identifier) > 1 {
					continue
				}
				// TODO: implement group properties
				g.Id(receiverName).Dot(comp.BaseName).Dot("SetPropertyExpression").Call(jen.Lit(prop.Identifier[0]), jen.Lit(prop.Expression), generatePositionRange(*prop.ValuePos))
			}
			// property change listeners
			g.Comment("register listeners for when a property changes")
			addMultiple(g, mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propId string) jen.Code {
				if tag, ok := prop.Tags[onChangeTag]; ok {
					return jen.Id(receiverName).Dot(propId).Dot("AddDependent").Call(jen.Id("vit").Dot("FuncDep").Call(jen.Id(receiverName).Dot(tag)))
				}
				return nil
			}))
			// event listeners
			g.Comment("register event listeners")
			firstEventListener := true
			for _, prop := range comp.Properties {
				if isInternalProperty(prop) || len(prop.Identifier) <= 1 {
					continue
				}
				if firstEventListener {
					firstEventListener = false
					g.Var().Id("event").Qual(vitPackage, "Listenable")
					g.Var().Id("listener").Qual(vitPackage, "Evaluater")
				}
				if prop.Identifier[0] == "Root" {
					g.Id("event").Op(",").Id("_").Op("=").Id(receiverName).Dot("Root").Dot("Event").Call(jen.Lit(prop.Identifier[1]))
					g.Id("listener").Op("=").Id("event").Dot("CreateListener").Call(jen.Lit(prop.Expression), generatePositionRange(prop.Pos))
					g.Id(receiverName).Dot("RootC").Call().Dot("AddListenerFunction").Call(jen.Id("listener"))
				} else {
					fmt.Fprintf(os.Stderr, "setting a property like %q is not currently supported\r\n", strings.Join(prop.Identifier, "."))
				}
			}
			// enumerations
			g.Comment("register enumerations")
			for _, enum := range comp.Enumerations {
				g.Id(receiverName).Dot("DefineEnum").Call(generateEnumeration(enum))
			}
			// children
			g.Comment("add child components")
			if len(comp.Children) > 0 {
				g.Var().Id("child").Qual(vitPackage, "Component")
				for _, child := range comp.Children {
					g.List(jen.Id("child"), jen.Op("_")).Op("=").Qual(parsePackage, "InstantiateComponent").Call(generateComponentDefinition(child), jen.Id("context"))
					g.Id(receiverName).Dot("AddChild").Call(jen.Id("child"))
				}
			}
			g.Line()
			g.Id("context").Dot("Environment").Dot("RegisterComponent").Call(jen.Id(receiverName))
			g.Line()
			g.Return(jen.Id(receiverName))
		})

	f.Line()

	// .String() string
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("String").
		Params().
		Params(jen.String()).
		Block(
			jen.Return(jen.Qual("fmt", "Sprintf").Call(jen.Lit(compName+"(%s)"), jen.Id(receiverName).Dot("id"))),
		).
		Line()

	// .Property(key string) (Value, bool)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("Property").
		Params(jen.Id("key").String()).
		Params(jen.Qual(vitPackage, "Value"), jen.Bool()).
		Block(
			jen.Switch(jen.Id("key")).Block(
				append(mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propId string) jen.Code {
					if !isReadable(prop) {
						return nil // don't add unreadable properties
					}
					return jen.Case(jen.Lit(propId)).Block(
						jen.Return(jen.Op("&").Id(receiverName).Dot(propId), jen.True()),
					)
				}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("Property").Call(jen.Id("key"))),
					),
				)...,
			),
		).
		Line()

	// .MustProperty(key string) Value
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("MustProperty").
		Params(jen.Id("key").String()).
		Params(jen.Qual(vitPackage, "Value")).
		Block(
			jen.List(jen.Id("v"), jen.Id("ok")).Op(":=").Id(receiverName).Dot("Property").Call(jen.Id("key")),
			jen.If(jen.Op("!").Id("ok")).Block(
				jen.Panic(jen.Qual("fmt", "Errorf").Call(jen.Lit("MustProperty called with unknown key %q"), jen.Id("key"))),
			),
			jen.Return(jen.Id("v")),
		).
		Line()

	// .SetProperty(key string, value interface{}) error
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("SetProperty").
		Params(jen.Id("key").String(), jen.Id("value").Interface()).
		Params(jen.Error()).
		Block(
			jen.Var().Id("err").Error(),
			jen.Switch(jen.Id("key")).Block(
				append(mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propId string) jen.Code {
					if !isWritable(prop) {
						return nil // don't add unwritable properties
					}
					return jen.Case(jen.Lit(propId)).Block(
						jen.Id("err").Op("=").Id(receiverName).Dot(propId).Op(".").Id("SetValue").Call(jen.Id("value")),
					)
				}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("SetProperty").Call(jen.Id("key"), jen.Id("value"))),
					),
				)...,
			),
			jen.If(jen.Id("err").Op("!=").Nil()).Block(
				jen.Return().Qual(vitPackage, "NewPropertyError").Call(jen.Lit(compName), jen.Id("key"), jen.Id(receiverName).Dot("id"), jen.Id("err")),
			),
			jen.Return(jen.Nil()),
		).
		Line()

	// .SetPropertyExpression(key string, code string, pos *vit.PositionRange) error
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("SetPropertyExpression").
		Params(jen.Id("key").String(), jen.Id("code").String(), jen.Id("pos").Op("*").Qual(vitPackage, "PositionRange")).
		Params(jen.Error()).
		Block(
			jen.Switch(jen.Id("key")).Block(
				append(mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propId string) jen.Code {
					if !isWritable(prop) {
						return nil // don't add unwritable properties
					}
					return jen.Case(jen.Lit(propId)).Block(
						jen.Id(receiverName).Dot(propId).Op(".").Id("SetExpression").Call(jen.Id("code"), jen.Id("pos")),
					)
				}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("SetPropertyExpression").Call(jen.Id("key"), jen.Id("code"), jen.Id("pos"))),
					),
				)...,
			),
			jen.Return().Nil(),
		).
		Line()

	// .ResolveVariable(key string) (interface{}, bool)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("ResolveVariable").
		Params(jen.Id("key").String()).
		Params(jen.Interface(), jen.Bool()).
		Block(
			jen.Switch(jen.Id("key")).BlockFunc(func(g *jen.Group) {
				g.Case(jen.Id(receiverName).Dot("id")).Block(
					jen.Return(jen.Id(receiverName), jen.True()),
				)
				for _, prop := range comp.Properties {
					if isInternalProperty(prop) || !isReadable(prop) || !prop.IsNewDefinition() {
						continue
					}
					g.Case(jen.Lit(prop.Identifier[0])).Block(
						jen.Return(jen.Op("&").Id(receiverName).Dot(prop.Identifier[0]), jen.True()),
					)
				}
				for _, method := range comp.Methods {
					g.Case(jen.Lit(method.Name)).Block(
						jen.Return(jen.Op("&").Id(receiverName).Dot(method.Name), jen.True()),
					)
				}
				for _, event := range comp.Events {
					g.Case(jen.Lit(event.Name)).Block(
						jen.Return(jen.Op("&").Id(receiverName).Dot(event.Name), jen.True()),
					)
				}
				g.Default().Block(
					jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("ResolveVariable").Call(jen.Id("key"))),
				)
			}),
		).
		Line()

	// .AddChild(child Component)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("AddChild").
		Params(jen.Id("child").Qual(vitPackage, "Component")).
		Block(
			generateCallbackForAddedChild(comp, receiverName, "child"),
			jen.Id("child").Dot("SetParent").Call(jen.Id(receiverName)),
			jen.Id(receiverName).Dot("AddChildButKeepParent").Call(jen.Id("child")),
		).
		Line()

	// .AddChildAfter(afterThis, addThis Component)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("AddChildAfter").
		Params(jen.Id("afterThis").Qual(vitPackage, "Component"), jen.Id("addThis").Qual(vitPackage, "Component")).
		Block(
			generateCallbackForAddedChild(comp, receiverName, "addThis"),
			jen.Var().Id("targetType").Qual(vitPackage, "Component").Op("=").Id("afterThis"),
			jen.Line(),
			jen.For(jen.List(jen.Id("ind"), jen.Id("child")).Op(":=").Range().Id(receiverName).Dot("Children").Call()).Block(
				jen.If(jen.Id("child").Dot("As").Call(jen.Op("&").Id("targetType"))).Block(
					jen.Id("addThis").Dot("SetParent").Call(jen.Id(receiverName)),
					jen.Id(receiverName).Dot("AddChildAtButKeepParent").Call(jen.Id("addThis"), jen.Id("ind").Op("+").Lit(1)),
					jen.Return(),
				),
			),
			jen.Id(receiverName).Dot("AddChild").Call(jen.Id("addThis")),
		).
		Line()

	// .UpdateExpressions() (int, ErrorGroup)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("UpdateExpressions").
		Params().
		Params(jen.Int(), jen.Qual(vitPackage, "ErrorGroup")).
		BlockFunc(func(g *jen.Group) {
			// initialize 'sum' and 'errs' variables
			g.Var().Id("sum").Int()
			g.Var().Id("errs").Qual(vitPackage, "ErrorGroup")
			g.Line()
			// now handle changes for all necessary properties
			addMultiple(g, mapProperties(comp.Properties, func(prop vit.PropertyDefinition, propID string) jen.Code {
				if prop.HasTag(typeTag) {
					// We will not handle changes for properties with custom types.
					return nil
				}
				return jen.If(jen.List(jen.Id("changed"), jen.Id("err")).Op(":=").Id(receiverName).Dot(propID).Dot("Update").Call(jen.Id(receiverName)).Op(";").Id("changed").Op("||").Id("err").Op("!=").Nil()).Block(
					jen.Id("sum").Op("++"),
					jen.If(jen.Id("err").Op("!=").Nil()).Block(
						jen.Id("errs").Dot("Add").Call(jen.Qual(vitPackage, "NewPropertyError").Call(
							jen.Lit(compName),
							jen.Lit(propID),
							jen.Id(receiverName).Dot("id"),
							jen.Id("err"),
						)),
					),
				)
			}))
			g.Line()
			g.Comment("this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables")
			g.Id("n").Op(",").Id("err").Op(":=").Id(receiverName).Dot("UpdatePropertiesInContext").Call(jen.Id(receiverName)) // n, err := receiver.UpdatePropertiesInContext(receiver) // just approximate code, names will vary
			g.Id("sum").Op("+=").Id("n")                                                                                      // sum += n
			g.Id("errs").Dot("AddGroup").Call(jen.Id("err"))                                                                  // errs.AddGroup(err)
			g.Id("n").Op(",").Id("err").Op("=").Id(receiverName).Dot(comp.BaseName).Dot("UpdateExpressions").Call()           // n, err = receiver.BaseComponent.UpdateExpressions()
			g.Id("sum").Op("+=").Id("n")                                                                                      // sum += n
			g.Id("errs").Dot("AddGroup").Call(jen.Id("err"))                                                                  // errs.AddGroup(err)
			g.Return(jen.Id("sum"), jen.Id("errs"))                                                                           // return sum, errs
		}).
		Line()

	// .As(*Component) (bool)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("As").
		Params(jen.Id("target").Op("*").Qual(vitPackage, "Component")).
		Params(jen.Bool()).
		Block(
			jen.If(jen.List(jen.Id("_"), jen.Id("ok")).Op(":=").Parens(jen.Op("*").Id("target")).Op(".").Parens(jen.Op("*").Id(compName)).Op(";").Id("ok")).Block(
				jen.Op("*").Id("target").Op("=").Id(receiverName),
				jen.Return(jen.True()),
			),
			jen.Return(jen.Id(receiverName).Dot("Item").Dot("As").Call(jen.Id("target"))),
		).
		Line()

	// ID() string
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("ID").
		Params().
		Params(jen.String()).
		Block(jen.Return(jen.Id(receiverName).Dot("id"))).
		Line()

	// Finish() error
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("Finish").
		Params().
		Params(jen.Error()).
		Block(
			jen.Return(jen.Id(receiverName).Dot("RootC").Call().Dot("FinishInContext").Call(jen.Id(receiverName))),
		).
		Line()

	f.Add(generateStaticAttributeMethod(receiverName, compName, comp))

	return nil
}

func generateStaticAttributeMethod(receiverName, compName string, comp *vit.ComponentDefinition) jen.Code {
	var didGenerateSomething bool

	// staticAttribute(name string) (interface{}, bool)
	code := jen.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("staticAttribute").
		Params(jen.Id("name").String()).
		Params(jen.Interface(), jen.Bool()).
		Block(
			jen.Switch(jen.Id("name")).BlockFunc(func(g *jen.Group) {
				// list all values of embedded enums
				for _, enum := range comp.Enumerations {
					if !enum.Embedded {
						continue
					}
					didGenerateSomething = true
					for _, value := range orderEnumValues(enum.Values) {
						g.Case(jen.Lit(value.name)).Block(
							jen.Return().List(
								// the type needs to be converted to uint to be usable in an expression
								jen.Uint().Call(jen.Id(fmt.Sprintf("%s_%s_%s", compName, enum.Name, value.name))),
								jen.True(),
							),
						)
					}
				}
				g.Default().Block(
					jen.Return().List(jen.Nil(), jen.False()),
				)
			}),
		).
		Line()

	if !didGenerateSomething {
		return jen.Null()
	}
	return code
}
