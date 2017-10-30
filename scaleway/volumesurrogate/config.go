package volumesurrogate

import (
	"fmt"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	AccessKey           string              `mapstructure:"access_key"`
	Token               string              `mapstructure:"token"`
	InstanceType        string              `mapstructure:"instance_type"`
	SourceImage         string              `mapstructure:"source_image"`
	VmName              string              `mapstructure:"vm_name"`
	Region              string              `mapstructure:"region"`
	Comm                communicator.Config `mapstructure:",squash"`
	Size                int                 `mapstructure:"size"`
	ctx                 interpolate.Context
}

func (c *Config) Dump() string {
	return fmt.Sprintf(`loaded config:
  * access_key: [REDACTED]
  * token: [REDACTED]
  * Instance Type: "%s"
  * Source Image: "%s"
  * Region: "%s"
  * Size: "%d"
  * Vm Name: "%s"`,
		c.InstanceType, c.SourceImage, c.Region, c.Size, c.VmName)
}
