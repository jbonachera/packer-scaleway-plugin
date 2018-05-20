package volumesurrogate

import (
	"context"

	"github.com/hashicorp/packer/helper/multistep"
	"github.com/hashicorp/packer/packer"
	"github.com/scaleway/scaleway-cli/pkg/api"
)

type step_create_image struct {
	config Config
	api    *api.ScalewayAPI
}

func (s *step_create_image) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	id := state.Get("server-id").(string)
	scw := s.api
	config := s.config
	ui.Say("stopping server")
	scw.PostServerAction(id, "poweroff")
	_, err := api.WaitForServerStopped(scw, id)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	ui.Say("server stopped")
	ui.Say("creating snapshot")
	srv, err := scw.GetServer(id)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	rootVolume, ok := srv.Volumes["1"]
	if !ok {
		ui.Error("could not find surrogate volume")
		ui.Error("volumes: ")
		for key := range srv.Volumes {
			ui.Error("  * " + key)
		}
		return multistep.ActionHalt
	}
	snapId, err := scw.PostSnapshot(rootVolume.Identifier, config.VmName)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("snapshot-id", snapId)
	image, err := scw.PostImage(snapId, s.config.VmName, "", "x86_64")
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("image-id", image)
	return multistep.ActionContinue
}

func (s *step_create_image) Cleanup(state multistep.StateBag) {
	snapId := state.Get("snapshot-id").(string)
	if snapId != "" {
		s.api.DeleteSnapshot(snapId)
	}
}
