package plugin

import (
	plugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin/host2plugin"
)

// PluginRoot is the root directory of the plugins
// This variable is exposed for testing.
var (
	PluginRoot      = "~/.tflint.d/plugins"
	localPluginRoot = "./.tflint.d/plugins"
)

// SDKVersionConstraints is the version constraint of the supported SDK version.
var SDKVersionConstraints = version.MustConstraints(version.NewConstraint(">= 0.16.0"))

// Plugin is an object handling plugins
// Basically, it is a wrapper for go-plugin and provides an API to handle them collectively.
type Plugin struct {
	RuleSets map[string]*host2plugin.Client

	clients map[string]*plugin.Client
}

// Clean is a helper for ending plugin processes
func (p *Plugin) Clean() {
	for _, client := range p.clients {
		client.Kill()
	}
}
