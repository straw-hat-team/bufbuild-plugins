package wit

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type GeneratorOption func(*witGenerator)

// WithPackageName sets the package name for the generated WIT file.
func WithPackageName(name string) GeneratorOption {
	return func(g *witGenerator) {
		g.packageName = name
	}
}

type result struct {
	ID      string
	Content strings.Builder
}

// Generate generates a WASM WIT file for the given message descriptor.
func Generate(input protoreflect.MessageDescriptor, opts ...GeneratorOption) map[protoreflect.FullName]result {
	generator := &witGenerator{
		result:      make(map[protoreflect.FullName]result),
		sb:          strings.Builder{},
		packageName: "generated",
	}
	for _, opt := range opts {
		opt(generator)
	}

	generator.sb.WriteString(fmt.Sprintf("package %s;\n\n", generator.packageName))
	generator.generate(input)

	return generator.result
}

type witGenerator struct {
	result       map[protoreflect.FullName]result
	packageName  string
	useJSONNames bool
	sb           strings.Builder
}

func (g *witGenerator) generate(desc protoreflect.MessageDescriptor) {
	if g.isGenerated(desc) {
		return
	}

	var sb strings.Builder

	g.generateRecord(&sb, desc)

	g.result[desc.FullName()] = result{
		ID:      g.getID(desc),
		Content: sb,
	}

	for i := 0; i < desc.Fields().Len(); i++ {
		field := desc.Fields().Get(i)
		if field.Kind() == protoreflect.MessageKind || field.Kind() == protoreflect.GroupKind {
			g.generate(field.Message())
		}
	}
}

func (g *witGenerator) isGenerated(desc protoreflect.MessageDescriptor) bool {
	if _, ok := g.result[desc.FullName()]; ok {
		return true
	}
	return false
}

func (g *witGenerator) generateRecord(sb *strings.Builder, desc protoreflect.MessageDescriptor) {
	sb.WriteString(fmt.Sprintf("record %s {\n", desc.Name()))
	for i := 0; i < desc.Fields().Len(); i++ {
		field := desc.Fields().Get(i)
		fieldName := strings.ToLower(string(field.Name()))
		fieldType := g.getWITType(field)
		sb.WriteString(fmt.Sprintf("    %s: %s,\n", fieldName, fieldType))
	}
	sb.WriteString("}\n\n")
}

func (g *witGenerator) getWITType(field protoreflect.FieldDescriptor) string {
	switch field.Kind() {
	case protoreflect.BoolKind:
		return "bool"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		return "s32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		return "s64"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		return "u32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		return "u64"
	case protoreflect.FloatKind:
		return "float32"
	case protoreflect.DoubleKind:
		return "float64"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		return "list<u8>"
	case protoreflect.EnumKind:
		return g.generateEnum(field.Enum())
	case protoreflect.MessageKind, protoreflect.GroupKind:
		if field.IsMap() {
			keyType := g.getWITType(field.MapKey())
			valueType := g.getWITType(field.MapValue())
			return fmt.Sprintf("list<tuple<%s, %s>>", keyType, valueType)
		}
		return string(field.Message().Name())
	}
	return "unknown"
}

func (g *witGenerator) generateEnum(enum protoreflect.EnumDescriptor) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("enum %s {\n", enum.Name()))
	for i := 0; i < enum.Values().Len(); i++ {
		value := enum.Values().Get(i)
		sb.WriteString(fmt.Sprintf("    %s,\n", value.Name()))
	}
	sb.WriteString("}\n\n")
	g.result[enum.FullName()] = result{
		ID:      g.getID(enum),
		Content: sb,
	}
	return string(enum.Name())
}

func (p *witGenerator) getID(desc protoreflect.Descriptor) string {
	if p.useJSONNames {
		return string(desc.FullName()) + ".jsonschema.wit"
	}
	return string(desc.FullName()) + ".schema.wit"
}
