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

func GetTypeNameAndPrefix(importFile pgs.File, pgsEntity pgs.Entity, typeName string,  f *FileGenerator) (string, string, error) {
	packageFile := pgsEntity.File()
	if importFile != nil {
		packageFile = importFile
	}
	typeName, err := GetMessageName(typeName, packageFile.Package())
	prefix := ""
	if importFile != nil && importFile != pgsEntity.File() {
		prefix, err = f.RegisterImport(typeName, importFile)
		if err != nil {
			return "", "", err
		}
		prefix = prefix + "."
	}
	return prefix , typeName, nil
}

func GetTypeNameAndPrefixForField(pgsField pgs.Field, f *FileGenerator) (string, string, error) {
	imports := pgsField.Imports()
	if len(imports) > 1 {
		return "", "", fmt.Errorf("too many imports for an entity %s", pgsField.Name().String())
	}
	var importFile pgs.File
	if imports != nil {
		importFile = imports[0]
	}
	return GetTypeNameAndPrefix(importFile, pgsField, pgsField.Descriptor().GetTypeName(), f)
}