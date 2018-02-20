package java

import (
	"github.com/Sirupsen/logrus"
	"github.com/bendikp/architect/pkg/config"
	"github.com/bendikp/architect/pkg/config/runtime"
	"github.com/bendikp/architect/pkg/docker"
	"github.com/bendikp/architect/pkg/java/prepare"
	"github.com/bendikp/architect/pkg/nexus"
	"github.com/bendikp/architect/pkg/process/build"
	"github.com/pkg/errors"
)

func Prepper() process.Prepper {
	return func(cfg *config.Config, auroraVersion *runtime.AuroraVersion, deliverable nexus.Deliverable,
		baseImage runtime.DockerImage) ([]docker.DockerBuildConfig, error) {

		logrus.Debug("Prepare output image")
		buildPath, err := prepare.Prepare(cfg.DockerSpec, auroraVersion, deliverable, baseImage)

		if err != nil {
			return nil, errors.Wrap(err, "Error prepare artifact")
		}

		buildConf := docker.DockerBuildConfig{
			AuroraVersion:    auroraVersion,
			BuildFolder:      buildPath,
			DockerRepository: cfg.DockerSpec.OutputRepository,
			Baseimage:        baseImage,
		}
		return []docker.DockerBuildConfig{buildConf}, nil
	}
}
