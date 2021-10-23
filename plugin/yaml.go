package plugin

import (
	"fmt"
	yaml "gopkg.in/yaml.v3"
)

type YamlNodeDecoder struct {
	Node *yaml.Node
}

func (d *YamlNodeDecoder) Decode(cfg interface{}) error {
	{
		if d.Node == nil {
			return fmt.Errorf("yaml node empty")
		}
		return d.Node.Decode(cfg)
	}
}
