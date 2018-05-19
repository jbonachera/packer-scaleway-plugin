package volumesurrogate

import (
	"crypto/rand"
	"crypto/rsa"
	"strings"

	"github.com/hashicorp/packer/packer"
	"github.com/mitchellh/multistep"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"golang.org/x/crypto/ssh"
)

type step_create_keypair struct {
	config Config
	api    *api.ScalewayAPI
}

func (s *step_create_keypair) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)
	scw := s.api
	u, err := scw.GetUser()
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	keys := u.SSHPublicKeys
	reader := rand.Reader
	newRSAKey, err := rsa.GenerateKey(reader, 2048)
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	pubkey, err := ssh.NewPublicKey(newRSAKey.Public())
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	key := api.ScalewayKeyDefinition{
		Key: strings.Trim(string(ssh.MarshalAuthorizedKey(pubkey)), "\n"),
	}
	for i := range keys {
		keys[i].Fingerprint = ""
	}
	ui.Say("uploading ssh key")
	err = scw.PatchUserSSHKey(u.ID, api.ScalewayUserPatchSSHKeyDefinition{
		SSHPublicKeys: append(keys, key),
	})
	if err != nil {
		ui.Error(err.Error())
		return multistep.ActionHalt
	}
	state.Put("private-key", newRSAKey)
	state.Put("public-key", key.Key)
	return multistep.ActionContinue
}

func (s *step_create_keypair) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)
	wantedKey := state.Get("public-key").(string)
	scw := s.api
	u, err := scw.GetUser()
	if err != nil {
		return
	}
	keys := append([]api.ScalewayKeyDefinition(nil), u.SSHPublicKeys...)
	for i, key := range u.SSHPublicKeys {
		if key.Key == wantedKey {
			keys = append(keys[:i], keys[i+1:]...)
		} else {
			keys[i].Fingerprint = ""
		}
	}
	ui.Say("removing temporary ssh key")
	err = scw.PatchUserSSHKey(u.ID, api.ScalewayUserPatchSSHKeyDefinition{
		SSHPublicKeys: keys,
	})
}
