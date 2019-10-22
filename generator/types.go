package generator

import (
	"fmt"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	pgs "github.com/lyft/protoc-gen-star"
	"strings"
)

var ReadWriteNames = map[descriptor.FieldDescriptorProto_Type]string{
	descriptor.FieldDescriptorProto_TYPE_DOUBLE: "double",
	descriptor.FieldDescriptorProto_TYPE_FLOAT: "float",

	descriptor.FieldDescriptorProto_TYPE_INT64: "int64",
	descriptor.FieldDescriptorProto_TYPE_UINT64: "uint64",

	descriptor.FieldDescriptorProto_TYPE_INT32: "int32",
	descriptor.FieldDescriptorProto_TYPE_FIXED64: "fixed64",
	descriptor.FieldDescriptorProto_TYPE_FIXED32: "fixed32",
	descriptor.FieldDescriptorProto_TYPE_BOOL: "bool",
	descriptor.FieldDescriptorProto_TYPE_STRING: "string",

	descriptor.FieldDescriptorProto_TYPE_UINT32: "uint32",

	descriptor.FieldDescriptorProto_TYPE_SFIXED32: "sfixed32",
	descriptor.FieldDescriptorProto_TYPE_SFIXED64: "sfixed64",
	descriptor.FieldDescriptorProto_TYPE_SINT32: "sint32",
	descriptor.FieldDescriptorProto_TYPE_SINT64: "sint64",
	descriptor.FieldDescriptorProto_TYPE_BYTES: "bytes",
	descriptor.FieldDescriptorProto_TYPE_MESSAGE: "message",
	descriptor.FieldDescriptorProto_TYPE_ENUM: "enum",
}

func JsTypeForProtoType(protoType descriptor.FieldDescriptorProto_Type) string {
	switch protoType {
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE,
		descriptor.FieldDescriptorProto_TYPE_FLOAT,
		descriptor.FieldDescriptorProto_TYPE_INT64,
		descriptor.FieldDescriptorProto_TYPE_UINT64,
		descriptor.FieldDescriptorProto_TYPE_INT32,
		descriptor.FieldDescriptorProto_TYPE_FIXED64,
		descriptor.FieldDescriptorProto_TYPE_FIXED32,
		descriptor.FieldDescriptorProto_TYPE_UINT32,
		descriptor.FieldDescriptorProto_TYPE_SFIXED64,
		descriptor.FieldDescriptorProto_TYPE_SFIXED32,
		descriptor.FieldDescriptorProto_TYPE_SINT32,
		descriptor.FieldDescriptorProto_TYPE_SINT64:
		return "number"
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return "string"
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return "boolean"
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		return "Uint8Array"
	default:
		return ""
	}
}

func IsPacked(protoDesc *descriptor.FieldDescriptorProto) bool {
	if protoDesc.Options != nil && protoDesc.Options.Packed != nil {
		return *protoDesc.Options.Packed
	}
	switch *protoDesc.Type {
	case descriptor.FieldDescriptorProto_TYPE_STRING,
		 descriptor.FieldDescriptorProto_TYPE_GROUP,
		descriptor.FieldDescriptorProto_TYPE_MESSAGE,
		descriptor.FieldDescriptorProto_TYPE_BYTES:
		return false
	}
	return true
}

func GetMessageName(fqn string, pkg pgs.Package) (string, error) {
	packageName := "." + pkg.ProtoName().String()
	if strings.HasPrefix(fqn, packageName + ".") {
		return strings.TrimPrefix(fqn, packageName + "."), nil
	}
	return "", fmt.Errorf("cannot find type %s in package %s", fqn, packageName)
}

func GetTypeNameAndPrefix(pgsField pgs.Field, msg *MessageGenerator) (string, string, error) {
	imports := pgsField.Imports()
	if len(imports) > 1 {
		return "", "", fmt.Errorf("too many imports for a field %s", pgsField.Name().String())
	}
	packageFile := pgsField.File()
	if imports != nil {
		packageFile = imports[0].File()
	}
	typeName, err := GetMessageName(pgsField.Descriptor().GetTypeName(), packageFile.Package())
	prefix := ""
	if len(imports) == 1 {
		prefix, err = msg.File.RegisterImport(pgsField.Descriptor().GetTypeName(), imports[0])
		if err != nil {
			return "", "", err
		}
		prefix = prefix + "."
	}
	return prefix , typeName, nil
}