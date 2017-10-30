package volumesurrogate

import (
	"crypto/rsa"
	"fmt"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/communicator"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	"github.com/scaleway/scaleway-cli/pkg/api"
	gossh "golang.org/x/crypto/ssh"
	"log"
	"time"
)

const (
	BuilderId = "jbonachera.scalewayvolumesurrogate"
)

type Builder struct {
	Config   Config
	api      *api.ScalewayAPI
	serverId string
	runner   multistep.Runner
}

func (b *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	state := &multistep.BasicStateBag{}
	state.Put("api", b.api)
	state.Put("ui", ui)
	state.Put("hook", hook)
	state.Put("cache", cache)
	steps := []multistep.Step{
		&step_create_keypair{
			config: b.Config,
			api:    b.api,
		},
		&step_create_server{
			config: b.Config,
			api:    b.api,
		},
		&communicator.StepConnect{
			SSHConfig: func(multistep.StateBag) (*gossh.ClientConfig, error) {
				privKey := state.Get("private-key").(*rsa.PrivateKey)
				signer, err := gossh.NewSignerFromKey(privKey)
				if err != nil {
					return nil, fmt.Errorf("Error setting up SSH config: %s", err)
				}
				return &gossh.ClientConfig{
					User:            b.Config.Comm.SSHUsername,
					HostKeyCallback: gossh.InsecureIgnoreHostKey(),
					Auth: []gossh.AuthMethod{
						gossh.PublicKeys(signer),
					},
				}, nil
			},
			Host: func(state multistep.StateBag) (string, error) {
				return state.Get("public-ip").(string), nil
			},
			SSHPort: func(multistep.StateBag) (int, error) {
				return b.Config.Comm.SSHPort, nil
			},
			Config: &b.Config.Comm,
		},
		&common.StepProvision{},
		&step_create_image{
			config: b.Config,
			api:    b.api,
		},
	}

	b.runner = common.NewRunner(steps, b.Config.PackerConfig, ui)
	b.runner.Run(state)
	// If there was an error, return that
	if rawErr, ok := state.GetOk("error"); ok {
		return nil, rawErr.(error)
	}

	img := state.Get("image-id").(string)
	if img == "" {
		return nil, nil
	}

	artifact := &Artifact{
		Image:   img,
		Builder: BuilderId,
		Api:     b.api,
	}

	return artifact, nil
}
func (b *Builder) Prepare(keys ...interface{}) ([]string, error) {
	err := config.Decode(&b.Config, &config.DecodeOpts{Interpolate: false}, keys...)
	if err != nil {
		return nil, err
	}
	log.Println(b.Config.Dump())
	b.api, err = api.NewScalewayAPI(b.Config.AccessKey, b.Config.Token, "a", b.Config.Region)
	if err != nil {
		return nil, err
	}
	if b.Config.Comm.SSHUsername == "" {
		b.Config.Comm.SSHUsername = "root"
	}
	if b.Config.Comm.Type == "" {
		b.Config.Comm.Type = "ssh"
	}
	if b.Config.Comm.SSHPort == 0 {
		b.Config.Comm.SSHPort = 22
	}
	if b.Config.Comm.SSHTimeout == 0 {
		b.Config.Comm.SSHTimeout = 30 * time.Second
	}

	return []string{}, nil
}

func (b *Builder) Cancel() {
	b.runner.Cancel()
}
