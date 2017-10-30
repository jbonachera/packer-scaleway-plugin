package volumesurrogate

import (
	"fmt"
	"github.com/scaleway/scaleway-cli/pkg/api"
	"log"
)

type Artifact struct {
	Image   string
	Builder string
	Api     *api.ScalewayAPI
}

func (a *Artifact) BuilderId() string {
	return a.Builder
}

func (*Artifact) Files() []string {
	return nil
}

func (a *Artifact) Id() string {
	return a.Image
}

func (a *Artifact) String() string {
	return fmt.Sprintf("Image was created: %s", a.Image)
}

func (a *Artifact) State(name string) interface{} {
	return nil
}

func (a *Artifact) Destroy() error {
	log.Printf("Deleting image ID (%s)", a.Image)
	img, err := a.Api.GetImage(a.Image)
	if err != nil {
		return err
	}
	err = a.Api.DeleteImage(a.Image)
	if err != nil {
		return err
	}
	err = a.Api.DeleteVolume(img.RootVolume.Identifier)
	if err != nil {
		return err
	}

	return nil
}
