package main

import (
	stdlog "log"
	"os"
	"runtime"

	"github.com/alecthomas/kingpin"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	// flags
	bitbucketAPIAccessToken = kingpin.Flag("bitbucket-api-token", "The time-limited access token to access the Bitbucket api.").Envar("ESTAFETTE_BITBUCKET_API_TOKEN").Required().String()
	gitRepoFullname         = kingpin.Flag("git-repo-fullname", "The owner and repo name of the Bitbucket repository.").Envar("ESTAFETTE_GIT_NAME").Required().String()
	gitRevision             = kingpin.Flag("git-revision", "The hash of the revision to set build status for.").Envar("ESTAFETTE_GIT_REVISION").Required().String()
	estafetteBuildStatus    = kingpin.Flag("estafette-build-status", "The current build status of the Estafette pipeline.").Envar("ESTAFETTE_BUILD_STATUS").Required().String()
	statusOverride          = kingpin.Flag("status-override", "Allow status property in manifest to override the actual build status.").Envar("ESTAFETTE_EXTENSION_STATUS").String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// pretty print to make build logs more readable
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().
		Timestamp().
		Logger()

	// use zerolog for any logs sent via standard log library
	stdlog.SetFlags(0)
	stdlog.SetOutput(log.Logger)

	// log startup message
	log.Info().
		Str("branch", branch).
		Str("revision", revision).
		Str("buildDate", buildDate).
		Str("goVersion", goVersion).
		Msg("Starting estafette-extension-bitbucket-status...")

	// check if there's a status override
	status := *estafetteBuildStatus
	if *statusOverride != "" {
		status = *statusOverride
	}

	// set build status
	bitbucketAPIClient := newBitbucketAPIClient()
	err := bitbucketAPIClient.SetBuildStatus(*bitbucketAPIAccessToken, *gitRepoFullname, *gitRevision, status)
	if err != nil {
		log.Fatal().Err(err).Msg("Updating Bitbucket build status failed")
	}

	log.Info().
		Msg("Finished estafette-extension-bitbucket-status...")
}
