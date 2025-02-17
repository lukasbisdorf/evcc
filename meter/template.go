package meter

import (
	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/templates"
	"gopkg.in/yaml.v3"
)

func init() {
	registry.Add("template", NewMeterFromTemplateConfig)
}

func NewMeterFromTemplateConfig(other map[string]interface{}) (api.Meter, error) {
	cc := struct {
		Template string
		Other    map[string]interface{} `mapstructure:",remain"`
	}{}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	tmpl, err := templates.ByTemplate(cc.Template, templates.Meter)
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
