package command

import (
	"github.com/IkezawaYuki/lucky-strike/protobuf/protoc-gen-gogo/generator"
	plugin_go "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
)

func Read() *plugin_go.CodeGeneratorRequest {
	g := generator.New()
}
