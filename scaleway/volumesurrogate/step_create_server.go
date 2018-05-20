package volumesurrogate

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

type step_create_server struct {
	config Config
	api    *api.ScalewayAPI
}

func (s *step_create_server) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	scw := s.api
	config := s.config
	trueVal := true
	ui.Say("Creating surrogate volume")
	volumeId, err := scw.PostVolume(api.ScalewayVolumeDefinition{
		Name:         "packer-builder",
		Size:         uint64(s.config.Size),
		Organization: scw.Organization,
		Type:         "l_ssd",
	})
	ui.Say("Surrogate volume created with id " + volumeId)
	state.Put("volume-id", volumeId)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}

	id, err := scw.PostServer(api.ScalewayServerDefinition{
		Name:              config.VmName,
		Image:             &config.SourceImage,
		Tags:              []string{"packer-builder"},
		CommercialType:    config.InstanceType,
		DynamicIPRequired: &trueVal,
		Volumes: map[string]string{
			"1": volumeId,
		},
	})
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("server-id", id)

	ui.Say("Starting server")
	err = scw.PostServerAction(id, "poweron")
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("Waiting for server to be ready")
	srv, err := api.WaitForServerReady(scw, id, "")
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ip := srv.PublicAddress.IP
	ui.Say("Server started on " + ip)
	state.Put("public-ip", ip)
	return multistep.ActionContinue
}

func (s *step_create_server) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	volumeId := state.Get("volume-id").(string)
	id := state.Get("server-id").(string)
	if id != "" {
		srv, err := s.api.GetServer(id)
		if srv.State == "running" {
			ui.Say("powering off server")
			s.api.PostServerAction(id, "poweroff")
			if err != nil {
				return
			}
			api.WaitForServerStopped(s.api, id)
		}
		ui.Say("Deleting server")
		err = s.api.DeleteServer(id)
		if err != nil {
			return
		}
		ui.Say("Deleting server volumes:")
		for _, vol := range srv.Volumes {
			ui.Say("  * " + vol.Name + " (" + vol.Identifier + ")")
			s.api.DeleteVolume(vol.Identifier)
		}
	}
	if volumeId != "" {
		s.api.DeleteVolume(volumeId)
	}
}
