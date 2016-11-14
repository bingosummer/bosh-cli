package cmd

import (
	"github.com/cppforlife/go-patch/patch"

	boshtpl "github.com/cloudfoundry/bosh-cli/director/template"
	boshui "github.com/cloudfoundry/bosh-cli/ui"
)

type BuildManifestCmd struct {
	ui boshui.UI
}

func NewBuildManifestCmd(ui boshui.UI) BuildManifestCmd {
	return BuildManifestCmd{ui: ui}
}

func (c BuildManifestCmd) Run(opts BuildManifestOpts) error {
	tpl := boshtpl.NewTemplate(opts.Args.Manifest.Bytes)

	vars := opts.VarFlags.AsVariables()
	op := opts.OpsFlags.AsOp()
	evalOpts := boshtpl.EvaluateOpts{ExpectAllKeys: opts.VarErrors}

	if opts.Path.IsSet() {
		evalOpts.PostVarSubstitutionOp = patch.FindOp{Path: opts.Path}

		// Printing YAML indented multiline strings (eg SSH key) is not useful
		evalOpts.UnescapedMultiline = true
	}

	bytes, err := tpl.Evaluate(vars, op, evalOpts)
	if err != nil {
		return err
	}

	c.ui.PrintBlock(string(bytes))

	return nil
}
