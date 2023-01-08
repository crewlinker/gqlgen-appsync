package genappsync

import (
	"embed"

	"github.com/99designs/gqlgen/codegen"
	"github.com/99designs/gqlgen/codegen/templates"
)

//go:embed *.gotpl
var tmpls embed.FS

// Plugin for gqlgen
type Plugin struct {
	fname   string
	pkgname string
}

// New inits the plugin
func New(fname, pkgname string) *Plugin {
	return &Plugin{fname: fname, pkgname: pkgname}
}

// Name names the plugin
func (Plugin) Name() string { return "clgqlgen" }

// GenerateCode performs the actual code generation
func (p Plugin) GenerateCode(gendata *codegen.Data) error {
	data := &PluginData{Data: gendata, HasResolverArguments: false}
	for _, obj := range data.Objects {
		if !obj.HasResolvers() {
			continue
		}
		for _, fld := range obj.Fields {
			if !fld.IsResolver {
				continue
			}

			// if there is at least one resolver field with arguments we
			// need to add extra code. To prevent compilation errors
			if len(fld.Args) > 0 {
				data.HasResolverArguments = true
			}
		}
	}

	return templates.Render(templates.Options{
		TemplateFS:      tmpls,
		PackageName:     p.pkgname,
		Filename:        p.fname,
		Data:            data,
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
	})
}

// PluginData is what is passed to the template
type PluginData struct {
	// adds a flag that is set to true if the schema has ANY arguments to unmarshal
	HasResolverArguments bool

	// Embedded codegen data
	*codegen.Data
}
