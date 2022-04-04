package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/omniskop/vitrum/vit/parse"
)

const vitPackage = "github.com/omniskop/vitrum/vit"

func GenerateFromFileAndSave(srcPath string, packageName string, dstPath string) error {
	doc, err := parseVit(srcPath)
	if err != nil {
		return fmt.Errorf("unable to parse: %v", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("unable to create destination file: %w", err)
	}
	defer dstFile.Close()

	return GenerateFromDocument(doc, packageName, dstFile)
}

func Generate(src io.Reader, srcPath string, packageName string, dst io.Writer) error {
	lexer := parse.NewLexer(src, srcPath)

	doc, err := parse.Parse(parse.NewTokenBuffer(lexer.Lex))
	if err != nil {
		return err
	}
	doc.Name = getComponentName(srcPath)

	return GenerateFromDocument(doc, packageName, dst)
}

func GenerateFromFile(srcFile string, packageName string, dst io.Writer) error {
	doc, err := parseVit(srcFile)
	if err != nil {
		return err
	}

	return GenerateFromDocument(doc, packageName, dst)
}

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

func GenerateFromDocument(doc *parse.VitDocument, packageName string, dst io.Writer) error {
	f := jen.NewFilePath(packageName)
	f.HeaderComment("Code generated by vitrum gencmd. DO NOT EDIT.")

	for _, comp := range doc.Components {
		generateComponent(f, doc.Name, comp)
	}

	return f.Render(dst)
}

func getComponentName(fileName string) string {
	fileName = filepath.Base(fileName)
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func generateComponent(f *jen.File, compName string, comp *parse.ComponentDefinition) {
	properties := []jen.Code{
		jen.Qual(vitPackage, comp.BaseName),
		jen.Id("id").String(),
		jen.Line(),
	}
	// we could use jen.Dict here but I wan't to preserve the property order
	propertyInstantiations := []jen.Code{
		jen.Line().Id(comp.BaseName).Op(":").Op("*").Qual(vitPackage, fmt.Sprintf("New%s", comp.BaseName)).Op("(").Id("id").Op(",").Id("scope").Op(")"),
		jen.Line().Id("id").Op(":").Id("id"),
	}

	receiverName := strings.ToLower(string(compName[0]))

	for _, prop := range comp.Properties {
		propType, propConstructor := vitTypeInfo(prop.VitType)
		properties = append(properties, jen.Id(prop.Identifier[0]).Add(propType))

		propertyInstantiations = append(propertyInstantiations, jen.Line().Id(prop.Identifier[0]).Op(":").Add(propConstructor.Call(jen.Lit(prop.Expression), jen.Nil())))
	}

	propertyInstantiations = append(propertyInstantiations, jen.Line())

	f.Type().Id(compName).Struct(properties...)

	// constructor
	f.Func().
		Id(fmt.Sprintf("New%s", compName)).
		Params(jen.Id("id").String(), jen.Id("scope").Qual(vitPackage, "ComponentContainer")).
		Params(jen.Op("*").Id(compName)).
		Block(
			jen.Return(jen.Op("&").Id(compName).Values(propertyInstantiations...)),
		)

	f.Line()

	// .String() string
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("String").
		Params().
		Params(jen.String()).
		Block(
			jen.Return(jen.Qual("fmt", "Sprintf").Call(jen.Lit("Rectangle{%s}"), jen.Id("r").Dot("id"))),
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
				append(mapProperties(comp.Properties, func(prop parse.Property, propId string) jen.Code {
					return jen.Case(jen.Lit(propId)).Block(
						jen.Return(jen.Id(receiverName).Dot(propId), jen.True()),
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

	// .SetProperty(key string, value interface{}, position *PositionRange) bool
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("SetProperty").
		Params(jen.Id("key").String(), jen.Id("value").Interface(), jen.Id("position").Op("*").Qual(vitPackage, "PositionRange")).
		Params(jen.Bool()).
		Block(
			jen.Switch(jen.Id("key")).Block(
				append(mapProperties(comp.Properties, func(prop parse.Property, propId string) jen.Code {
					return jen.Case(jen.Lit(propId)).Block(
						jen.Id(receiverName).Dot(propId).Dot("ChangeCode").Call(jen.Id("value").Assert(jen.String()), jen.Id("position")),
					)
				}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("SetProperty").Call(jen.Id("key"), jen.Id("value"), jen.Id("position"))),
					),
				)...,
			),
			jen.Return(jen.True()),
		).
		Line()

	// .ResolveVariable(key string) (interface{}, bool)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("ResolveVariable").
		Params(jen.Id("key").String()).
		Params(jen.Interface(), jen.Bool()).
		Block(
			jen.Switch(jen.Id("key")).Block(
				prepend(
					jen.Case(jen.Id(receiverName).Dot("id")).Block(
						jen.Return(jen.Id(receiverName), jen.True()),
					),
					mapProperties(comp.Properties, func(prop parse.Property, propId string) jen.Code {
						return jen.Case(jen.Lit(propId)).Block(
							jen.Return(jen.Id(receiverName).Dot(propId), jen.True()),
						)
					}),
					jen.Default().Block(
						jen.Return(jen.Id(receiverName).Dot(comp.BaseName).Dot("ResolveVariable").Call(jen.Id("key"))),
					),
				)...,
			),
		).
		Line()

	// .AddChild(child Component)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("AddChild").
		Params(jen.Id("child").Qual(vitPackage, "Component")).
		Block(
			jen.Id("child").Dot("SetParent").Call(jen.Id(receiverName)),
			jen.Id(receiverName).Dot("AddChildButKeepParent").Call(jen.Id("child")),
		).
		Line()

	// .UpdateExpression() (int, ErrorGroup)
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("UpdateExpression").
		Params().
		Params(jen.Int(), jen.Qual(vitPackage, "ErrorGroup")).
		BlockFunc(func(g *jen.Group) {
			g.Var().Id("sum").Int()
			g.Var().Id("errs").Qual(vitPackage, "ErrorGroup")
			g.Line()
			addMultiple(g, mapProperties(comp.Properties, func(prop parse.Property, propId string) jen.Code {
				return jen.If(jen.Id("r").Dot(propId).Dot("ShouldEvaluate").Call()).Block(
					jen.Id("sum").Op("++"),
					jen.Id("err").Op(":=").Id("r").Id(".").Id(propId).Dot("Update").Call(jen.Id(receiverName)),
					jen.If(jen.Id("err").Op("!=").Nil()).Block(
						jen.Id("errs").Dot("Add").Call(jen.Qual(vitPackage, "NewExpressionError").Call(
							jen.Lit(compName),
							jen.Lit(propId),
							jen.Id(receiverName).Dot("id"),
							jen.Id(receiverName).Dot(propId).Dot("Expression"),
							jen.Id("err"),
						)),
					),
				)
			}))
			g.Line()
			g.Comment("this needs to be done in every component and not just in root to give the expression the highest level component for resolving variables")
			g.Id("n").Op(",").Id("err").Op(":=").Id("r").Dot("UpdatePropertiesInContext").Call(jen.Id(receiverName))
			g.Id("sum").Op("+=").Id("n")
			g.Id("errs").Dot("AddGroup").Call(jen.Id("err"))
			g.Id("n").Op(",").Id("err").Op("=").Id("r").Dot("Item").Dot("UpdateExpressions").Call()
			g.Id("sum").Op("+=").Id("n")
			g.Id("errs").Dot("AddGroup").Call(jen.Id("err"))
			g.Return(jen.Id("sum"), jen.Id("errs"))
		}).
		Line()

	// ID() string
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("ID").
		Params().
		Params(jen.String()).
		Block(jen.Return(jen.Id(receiverName).Dot("id"))).
		Line()

	// finish() error
	f.Func().
		Params(jen.Id(receiverName).Op("*").Id(compName)).
		Id("finish").
		Params().
		Params(jen.Error()).
		Block(
			jen.Return(jen.Id(receiverName).Dot("RootC").Call().Dot("FinishInContext").Call(jen.Id(receiverName))),
		).
		Line()
}

func mapProperties(props []parse.Property, f func(parse.Property, string) jen.Code) []jen.Code {
	var result []jen.Code
	for _, prop := range props {
		// NOTE: properties with multiple identifiers are currently not supported
		result = append(result, f(prop, prop.Identifier[0]))
	}
	return result
}

func prepend(first jen.Code, rest []jen.Code, tail ...jen.Code) []jen.Code {
	return append([]jen.Code{first}, append(rest, tail...)...)
}

// add multiple adds all code to the group as separate statements
func addMultiple(g *jen.Group, code []jen.Code) {
	for _, c := range code {
		g.Add(c)
	}
}

func vitTypeInfo(vitType string) (*jen.Statement, *jen.Statement) {
	switch vitType {
	case "string":
		return jen.Op("*").Qual(vitPackage, "StringValue"), jen.Qual(vitPackage, "NewStringValue")
	case "int":
		return jen.Op("*").Qual(vitPackage, "IntValue"), jen.Qual(vitPackage, "NewIntValue")
	case "float":
		return jen.Op("*").Qual(vitPackage, "FloatValue"), jen.Qual(vitPackage, "NewFloatValue")
	case "bool":
		return jen.Op("*").Qual(vitPackage, "BoolValue"), jen.Qual(vitPackage, "NewBoolValue")
	case "color":
		return jen.Op("*").Qual(vitPackage, "ColorValue"), jen.Qual(vitPackage, "NewColorValue")
	default:
		panic(fmt.Errorf("unknown vit type %q", vitType))
	}
}
