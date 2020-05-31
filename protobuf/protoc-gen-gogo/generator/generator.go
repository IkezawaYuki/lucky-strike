package generator

import (
	"bytes"
	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	plugin_go "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
	"strconv"
)

type GoImportPath string

func (p GoImportPath) String() string {
	return strconv.Quote(string(p))
}

type GoPackageName string

type common struct {
	file *FileDescriptor
}

func (c *common) GoImportPath() GoImportPath {
	return c.file.importPath
}

func (c *common) File() *FileDescriptor {
	return c.file
}

func fileIsProto3(file *descriptor.FileDescriptorProto) bool {
	return file.GetSyntax() == "proto3"
}

func (c *common) proto3() bool {
	return fileIsProto3(c.file.FileDescriptorProto)
}

type Descriptor struct {
	common
	*descriptor.DescriptorProto
	parent *Descriptor
	nested []*Descriptor
	enums  []*EnumDescriptor
}

type EnumDescriptor struct {
	common
	*descriptor.EnumDescriptorProto
	parent   *Descriptor
	typeName []string
	index    int
	path     string
}

func (e *EnumDescriptor) TypeName() (s []string) {
	if e.typeName != nil {
		return e.typeName
	}
	name := e.GetName()
	if e.parent == nil {
		s = make([]string, 1)
	} else {
		pname := e.parent.TypeName()
		s = make([]string, len(pname)+1)
		copy(s, pname)
	}
	s[len(s)-1] = name
	e.typename = s
	return s
}

func (e *EnumDescriptor) alias() (s []string) {
	s = e.TypeName()
	if gogoproto.IsEnumCustomName(e.EnumDescriptorProto) {
		s[len(s)-1] = gogoproto.GetEnumCustomName(e.EnumDescriptorProto)
	}
	return
}

func (e *EnumDescriptor) prefix() string {
	typeName := e.alias()
	if e.parent == nil {
		return CamelCase(typeName[len(typeName)-1]) + "_"
	}
}

type FileDescriptor struct {
	*descriptor.FileDescriptorProto
	desc []*Descriptor
}

type Generator struct {
	*bytes.Buffer

	Request  *plugin_go.CodeGeneratorRequest
	Response *plugin_go.CodeGeneratorResponse

	Param             map[string]string
	PackageImportPath string
	ImportPrefix      string
	ImportMap         map[string]string

	Pkg map[string]string

	outputImportPath GoImportPath
	allFiles         []*FileDescriptor
}

func New() *Generator {
	g := new(Generator)

}
