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
func (p Plugin) GenerateCode(data *codegen.Data) error {
	return templates.Render(templates.Options{
		TemplateFS:      tmpls,
		PackageName:     p.pkgname,
		Filename:        p.fname,
		Data:            &PluginData{Data: data},
		GeneratedHeader: true,
		Packages:        data.Config.Packages,
	})
}

// PluginData is what is passed to the template
type PluginData struct {
	*codegen.Data
}
