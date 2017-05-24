package docker

import "fmt"

const (
	ENV_APP_VERSION     			= "APP_VERSION"
	ENV_AURORA_VERSION  			= "AURORA_VERSION"
	ENV_SNAPSHOT_TAG    			= "SNAPSHOT_TAG"
	ENV_PUSH_EXTRA_TAGS 			= "PUSH_EXTRA_TAGS"
	ENV_READINESS_CHECK_URL			= "READINESS_CHECK_URL"
	ENV_READINESS_ON_MANAGEMENT_PORT	= "READINESS_ON_MANAGEMENT_PORT"
)

type ImageName struct {
	Registry 	string
	Repository 	string
	Tag		string
}

func (n ImageName) String() string {
	if n.Registry == "" {
		return fmt.Sprintf("%s:%s", n.Repository, n.Tag)
	}

	return fmt.Sprintf("%s/%s:%s", n.Registry, n.Repository, n.Tag)
}
