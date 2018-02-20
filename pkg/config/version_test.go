package config_test

import (
	"fmt"
	"github.com/bendikp/architect/pkg/config"
	"github.com/bendikp/architect/pkg/config/runtime"
	"github.com/bendikp/architect/pkg/docker"
	"github.com/docker/distribution/manifest/schema1"
	"testing"
)

const (
	CFG_APPLICATION_TYPE     = "Java"
	CFG_GAV_GROUP_ID         = "bar"
	CFG_GAV_ARTIFACT_ID      = "foo"
	CFG_GAV_VERSION          = "2.4.5"
	CFG_GAV_SNAPSHOT_VERSION = "2.4.6-SNAPSHOT"
	CFG_OUTPUT_REGISTRY      = "ouput.skead.no:5000"
	CFG_EXTERNAL_REGISTRY    = "ext.skead.no:5001"
	CFG_OUPUT_REPOSITORY     = "aurora/foo"
	CFG_PUSH_EXTRA_TAGS      = "major minor patch latest"
	CFG_BUILDER_VERSION      = "1.11.0"
	CFG_BASE_REPOSITORY      = "aurora/oracle8"
	CFG_BASE_VERSION         = "1"
	CFG_TAG_WITH             = "ab543b32de"
)

const (
	DELIVERABLE_PATH            = "/tmp/foo-JRA-100-Fix-20170329.115428-3-Leveransepakke.zip"
	INFERRED_BASE_IMAGE_VERSION = "1.2.3"
)

const (
	APP_VERSION     = "2.4.5"
	AURORA_VERSION  = "2.4.5-b1.11.0-oracle8-1.2.3"
	PUSH_EXTRA_TAGS = "major minor patch latest"
)

const (
	TAG_MAJOR    = "2"
	TAG_MINOR    = "2.4"
	TAG_PATCH    = "2.4.5"
	TAG_COMPLETE = "2.4.5-b1.11.0-oracle8-1.2.3"
	TAG_TEMP     = "ab543b32de"
)

const (
	SNAPSHOT_GIVEN_VERSION  = "branch_test-SNAPSHOT"
	SNAPSHOT_APP_VERSION    = "branch_test-201703929219"
	SNAPSHOT_AURORA_VERSION = "SNAPSHOT-201703929219-b1.11.0-oracle8-1.2.3"
	SNAPSHOT_TAG_COMPLETE   = "SNAPSHOT-201703929219-b1.11.0-oracle8-1.2.3"
)

type RegistryMock struct{}

var dockerSpec = config.DockerSpec{
	OutputRegistry:         CFG_OUTPUT_REGISTRY,
	PushExtraTags:          config.ParseExtraTags(CFG_PUSH_EXTRA_TAGS),
	OutputRepository:       CFG_OUPUT_REPOSITORY,
	ExternalDockerRegistry: CFG_EXTERNAL_REGISTRY,
}

var mavenGavRelease = config.ApplicationSpec{
	MavenGav: config.MavenGav{
		GroupId:    CFG_GAV_GROUP_ID,
		ArtifactId: CFG_GAV_ARTIFACT_ID,
		Version:    CFG_GAV_VERSION,
	},
	BaseImageSpec: config.DockerBaseImageSpec{
		BaseImage:   CFG_BASE_REPOSITORY,
		BaseVersion: CFG_BASE_VERSION,
	}}

var mavenGavSnapshot = config.ApplicationSpec{
	MavenGav: config.MavenGav{
		GroupId:    CFG_GAV_GROUP_ID,
		ArtifactId: CFG_GAV_ARTIFACT_ID,
		Version:    CFG_GAV_SNAPSHOT_VERSION,
	},
}

func TestTagInfoRelease(t *testing.T) {
	appVersion := runtime.NewAuroraVersion(APP_VERSION, false, APP_VERSION, runtime.CompleteVersion(AURORA_VERSION))
	tags, err := appVersion.GetApplicationVersionTagsToPush(make([]string, 0), config.ParseExtraTags(CFG_PUSH_EXTRA_TAGS))
	if err != nil {
		t.Fatalf("Failed to create target VersionInfo %v", err)
	}

	//TODO: Add the test for complete tag, but it should not be a part of appversion
	//expectedTags := []string{"latest", TAG_MAJOR, TAG_MINOR, TAG_PATCH, TAG_COMPLETE}
	expectedTags := []string{"latest", TAG_MAJOR, TAG_MINOR, TAG_PATCH, TAG_COMPLETE}

	verifyTagListContent(tags, expectedTags, t)
}

func TestTagInfoSnapshot(t *testing.T) {
	appVersion := runtime.NewAuroraVersion(SNAPSHOT_APP_VERSION, true, SNAPSHOT_GIVEN_VERSION, runtime.CompleteVersion(SNAPSHOT_AURORA_VERSION))
	tags, err := appVersion.GetApplicationVersionTagsToPush([]string{}, config.ParseExtraTags(CFG_PUSH_EXTRA_TAGS))
	if err != nil {
		t.Fatalf("Failed to create target VersionInfo %v", err)
	}

	//TODO: Add the test for complete tag, but it should not be a part of appversion
	verifyTagListContent(tags, []string{SNAPSHOT_GIVEN_VERSION, SNAPSHOT_TAG_COMPLETE}, t)
}

//TODO: This test should be in Java-buidler
/*
func TestBuildInfoReleaset(t *testing.T) {

	cfg := config.Config{ApplicationType: CFG_APPLICATION_TYPE, DockerSpec: dockerSpec, MavenGav: &mavenGavRelease,
		BuilderSpec: config.BuilderSpec{CFG_BUILDER_VERSION}}

	actual, err := config.NewBuildInfo(cfg, config.Deliverable{DELIVERABLE_PATH}, RegistryMock{})

	if err != nil {
		t.Fatalf("Failed to create target BuildInfo %v", err)
	}

	expectedEnv := map[string]string{
		docker.ENV_APP_VERSION:     APP_VERSION,
		docker.ENV_AURORA_VERSION:  AURORA_VERSION,
		docker.ENV_PUSH_EXTRA_TAGS: PUSH_EXTRA_TAGS,
	}

	expectedTags := []string{"latest", TAG_MAJOR, TAG_MINOR, TAG_PATCH, TAG_COMPLETE}

	verifyEnvMapContent(actual.Env, expectedEnv, t)
	verifyTagListContent(actual.OutputImage.VersionTags, expectedTags, t)
	verifyEquals(actual.BaseImage.Version, INFERRED_BASE_IMAGE_VERSION, t)
	verifyEquals(actual.BaseImage.Repository, CFG_BASE_REPOSITORY, t)
	verifyEquals(actual.OutputImage.Repository, CFG_OUPUT_REPOSITORY, t)
}


func TestBuildInfoTemporary(t *testing.T) {
	dockerSpec.TagWith = CFG_TAG_WITH

	cfg := config.Config{ApplicationType: CFG_APPLICATION_TYPE, DockerSpec: dockerSpec, MavenGav: &mavenGavRelease,
		BuilderSpec: config.BuilderSpec{CFG_BUILDER_VERSION}}

	actual, err := config.NewBuildInfo(cfg, config.Deliverable{DELIVERABLE_PATH}, RegistryMock{})

	if err != nil {
		t.Fatalf("Failed to create BuildInfo %v", err)
	}

	expectedEnv := map[string]string{
		docker.ENV_APP_VERSION:     APP_VERSION,
		docker.ENV_AURORA_VERSION:  AURORA_VERSION,
		docker.ENV_PUSH_EXTRA_TAGS: PUSH_EXTRA_TAGS,
	}

	expectedTags := []string{TAG_TEMP}

	verifyEnvMapContent(actual.Env, expectedEnv, t)
	verifyTagListContent(actual.OutputImage.VersionTags, expectedTags, t)
	verifyEquals(actual.BaseImage.Version, INFERRED_BASE_IMAGE_VERSION, t)
	verifyEquals(actual.BaseImage.Repository, CFG_BASE_REPOSITORY, t)
	verifyEquals(actual.OutputImage.Repository, CFG_OUPUT_REPOSITORY, t)
}
*/

func TestFilterTags(t *testing.T) {

	r := repositoryTester{
		t: t,
		tagsFromRegistry: []string{"latest",
			"1.1.2", "1.1", "1", "1.2.1", "1.2", "1.3.0", "1.3",
			"2.0.0", "2.0", "2"},
	}
	r.testTagFiltering(
		"1.1.1",
		[]string{"latest", "1.1.1", "1.1", "1", "1.1.1-b0.0.0-oracle8-1.4.0", "someothertag"},
		[]string{"1.1.1", "1.1.1-b0.0.0-oracle8-1.4.0", "someothertag"})

	r.testTagFiltering(
		"1.2.2",
		[]string{"latest", "1.2.2", "1.2", "1"},
		[]string{"1.2.2", "1.2"})

	r.testTagFiltering(
		"1.3.1",
		[]string{"latest", "1.3.1", "1.3", "1"},
		[]string{"1.3.1", "1.3", "1"})

	r.testTagFiltering(
		"2.0.1",
		[]string{"latest", "2.0.1", "2.0", "2"},
		[]string{"2.0.1", "2.0", "2", "latest"})
}

type repositoryTester struct {
	t                *testing.T
	tagsFromRegistry []string
}

//TODO: We don't filter tags, we return the tags we need... Need to refactor test!
func (m repositoryTester) testTagFiltering(appVersion string, candidateTags []string, excpectedFilteringResult []string) {
	a := runtime.NewAuroraVersion(appVersion, false, appVersion, runtime.CompleteVersion(AURORA_VERSION))
	_, err := a.GetApplicationVersionTagsToPush(m.tagsFromRegistry, config.ParseExtraTags("lastest major minor patch"))

	if err != nil {
		m.t.Fatalf("Failed to call FilterTags %v", err)
	}

	//TODO: Refactory!
	//verifyTagListContent(myTags, excpectedFilteringResult, m.t)
}

func (registry RegistryMock) GetTags(repository string) (*docker.TagsAPIResponse, error) {
	tags := []string{"a", "b"}
	return &docker.TagsAPIResponse{Name: "jalla", Tags: tags}, nil // Do not need this
}

func (registry RegistryMock) GetManifest(repository string, tag string) (*schema1.SignedManifest, error) {
	return nil, nil // Do not need this
}

func (registry RegistryMock) GetManifestEnv(repository string, tag string, name string) (string, error) {
	if name == "BASE_IMAGE_VERSION" {
		return INFERRED_BASE_IMAGE_VERSION, nil
	} else {
		return "", fmt.Errorf("ENV variable not found")
	}
}

func verifyTagListContent(actualList []string, expectedList []string, t *testing.T) {
	if len(actualList) != len(expectedList) {
		t.Errorf("Expected %v tags, actual is %v", expectedList, actualList)
	}

	for _, e := range expectedList {
		verifyContainsTag(actualList, e, t)
	}
}

func verifyContainsTag(actual []string, expected string, t *testing.T) {
	if !contains(actual, expected) {
		t.Errorf("Expected tag %s does not exist", expected)
	}
}

func contains(target []string, value string) bool {
	for _, t := range target {
		if t == value {
			return true
		}
	}

	return false
}
