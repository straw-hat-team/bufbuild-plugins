package main

import (
	"github.com/bufbuild/protoplugin"
	"github.com/straw-hat-team/bufbuild-plugins/internal/protoschema"
	"github.com/straw-hat-team/bufbuild-plugins/internal/protoschema/plugin/pluginjsonschema"
)

func main() {
	protoplugin.Main(protoplugin.HandlerFunc(pluginjsonschema.Handle), protoplugin.WithVersion(protoschema.Version()))
}
