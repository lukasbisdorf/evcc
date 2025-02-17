package charger

import (
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/templates"
	"gopkg.in/yaml.v3"
)

func init() {
	registry.Add("template", NewChargerFromTemplateConfig)
}

func NewChargerFromTemplateConfig(other map[string]interface{}) (api.Charger, error) {
	cc := struct {
		Template string
		Other    map[string]interface{} `mapstructure:",remain"`
	}{}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	tmpl, err := templates.ByTemplate(cc.Template, templates.Charger)
	if err != nil {
		return nil, err
	}

	b, _, err := tmpl.RenderResult(false, other)
	if err != nil {
		return nil, err
	}

	var instance struct {
		Type  string
		Other map[string]interface{} `yaml:",inline"`
	}

	if err := yaml.Unmarshal(b, &instance); err != nil {
		return nil, err
	}

	return NewFromConfig(instance.Type, instance.Other)
}
