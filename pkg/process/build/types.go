package process

import (
	"github.com/bendikp/architect/pkg/config"
	"github.com/bendikp/architect/pkg/config/runtime"
	"github.com/bendikp/architect/pkg/docker"
	"github.com/bendikp/architect/pkg/nexus"
)

// Prepper is a fuction used to prepare a docker image. It is called within the context of
// The
type Prepper func(
	cfg *config.Config,
	auroraVersion *runtime.AuroraVersion,
	deliverable nexus.Deliverable,
	baseImage runtime.DockerImage) ([]docker.DockerBuildConfig, error)
