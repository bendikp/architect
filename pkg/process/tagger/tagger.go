package tagger

import (
	"github.com/Sirupsen/logrus"
	"github.com/bendikp/architect/pkg/config"
	"github.com/bendikp/architect/pkg/config/runtime"
	"github.com/bendikp/architect/pkg/docker"
	"github.com/pkg/errors"
)

type TagResolver interface {
	ResolveTags(appVersion *runtime.AuroraVersion, pushExtratags config.PushExtraTags) ([]string, error)
}

type SingleTagTagResolver struct {
	Registry   string
	Repository string
	Tag        string
}

func (m *SingleTagTagResolver) ResolveTags(appVersion *runtime.AuroraVersion, pushExtratags config.PushExtraTags) ([]string, error) {
	return docker.CreateImageNameFromSpecAndTags([]string{m.Tag}, m.Registry, m.Repository), nil
}

type NormalTagResolver struct {
	Registry   string
	Repository string
	Overwrite  bool
	Provider   docker.ImageInfoProvider
}

func (m *NormalTagResolver) ResolveTags(appVersion *runtime.AuroraVersion, pushExtratags config.PushExtraTags) ([]string, error) {
	tags, err := findCandidateTags(appVersion, m.Overwrite, m.Repository, pushExtratags, m.Provider)
	if err != nil {
		return nil, err
	}
	tags = append(tags)
	return docker.CreateImageNameFromSpecAndTags(tags, m.Registry, m.Repository), nil
}

func findCandidateTags(appVersion *runtime.AuroraVersion, tagOverwrite bool, outputRepository string,
	pushExtraTags config.PushExtraTags, provider docker.ImageInfoProvider) ([]string, error) {
	var repositoryTags []string
	if !tagOverwrite {

		repositoryTags, err := provider.GetTags(outputRepository)
		logrus.Debug("Tags in repository ", repositoryTags)

		if err != nil {
			return nil, errors.Wrapf(err, "Error in GetTags, repository=%s", outputRepository)
		}

	}
	versionTags, err := appVersion.GetApplicationVersionTagsToPush(repositoryTags, pushExtraTags)
	if err != nil {
		return nil, errors.Wrapf(err, "Error in FilterVersionTags, app_version=%s, "+
			"versionTags=%v, repositoryTags=%v",
			appVersion, versionTags, repositoryTags)
	}
	return versionTags, nil
}
