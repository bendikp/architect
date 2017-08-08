package architect

import (
	"github.com/Sirupsen/logrus"
	"github.com/skatteetaten/architect/pkg/config"
	"github.com/skatteetaten/architect/pkg/java"
	"github.com/skatteetaten/architect/pkg/java/nexus"
	"github.com/spf13/cobra"
)

var localRepo bool
var verbose bool

var JavaLeveransepakke = &cobra.Command{

	Use:   "build",
	Short: "Build Docker image from Zip",
	Long:  `TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		var configReader = config.NewInClusterConfigReader()
		var downloader nexus.Downloader
		if verbose {
			logrus.SetLevel(logrus.DebugLevel)
		} else {
			logrus.SetLevel(logrus.InfoLevel)
		}
		if len(cmd.Flag("fileconfig").Value.String()) != 0 {
			conf := cmd.Flag("fileconfig").Value.String()
			logrus.Debugf("Using config from %s", conf)
			configReader = config.NewFileConfigReader(conf)
		}
		if localRepo {
			logrus.Debugf("Using local maven repo")
			downloader = nexus.NewLocalDownloader()
		} else {
			mavenRepo := "http://aurora/nexus/service/local/artifact/maven/content"
			logrus.Debugf("Using Maven repo on %s", mavenRepo)
			downloader = nexus.NewNexusDownloader(mavenRepo)
		}

		RunArchitect(configReader, downloader)
	},
}

func init() {
	JavaLeveransepakke.Flags().StringP("fileconfig", "f", "", "Path to file config. If not set, the environment variable BUILD is read")
	JavaLeveransepakke.Flags().StringP("skippush", "s", "", "If set, Docker push will not be performed")
	JavaLeveransepakke.Flags().BoolVarP(&localRepo, "localrepo", "l", false, "If set, the Leveransepakke will be fetched from the local repo")
	JavaLeveransepakke.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logging")
}

func RunArchitect(configReader config.ConfigReader, downloader nexus.Downloader) {

	// Read build config
	cfg, err := configReader.ReadConfig()
	if err != nil {
		logrus.Fatalf("Could not read configuration: %s", err)
	}

	logrus.Debugf("Config %+v", cfg)

	err = configReader.AddRegistryCredentials(cfg)
	if err != nil {
		logrus.Fatalf("Could not read configuration: %s", err)
	}

	if cfg.DockerSpec.RetagWith != "" {
		logrus.Info("Perform retag")

		tagger, err := java.CreateTagger(*cfg)

		if err != nil {
			logrus.Fatalf("Failed to retag temporary image: %s", err)
		}

		err = tagger.RetagTemporary()

		if err != nil {
			logrus.Fatalf("Failed to retag temporary image: %s", err)
		}


	} else {
		logrus.Debugf("Download deliverable for GAV %-v", cfg.MavenGav)
		deliverable, err := downloader.DownloadArtifact(&cfg.MavenGav)

		if err != nil {
			logrus.Fatalf("Could not download deliverable %-v", cfg.MavenGav)
		}

		builder, err := java.CreateBuilder(*cfg, deliverable)

		logrus.Info("Perform build")

		err = builder.Build()

		if err != nil {
			logrus.Fatalf("Failed to build image: %s", err)
		}

		if cfg.DockerSpec.TagWith != "" {
			err = builder.TagTemporary()

			if err != nil {
				logrus.Fatalf("Failed to tag image: %s", err)
			}

		} else {
			err = builder.Tag()

			if err != nil {
				logrus.Fatalf("Failed to tag image: %s", err)
			}
		}
	}

}
