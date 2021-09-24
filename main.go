package main

import (
	"encoding/json"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/alecthomas/kingpin"
	foundation "github.com/estafette/estafette-foundation"
	"github.com/rs/zerolog/log"
)

var (
	appgroup  string
	app       string
	version   string
	branch    string
	revision  string
	buildDate string
	goVersion = runtime.Version()
)

var (
	// flags
	gitRepoSource        = kingpin.Flag("git-repo-source", "The source of the git repository, bitbucket.org in this case.").Envar("ESTAFETTE_GIT_SOURCE").Required().String()
	gitRepoFullname      = kingpin.Flag("git-repo-fullname", "The owner and repo name of the Bitbucket repository.").Envar("ESTAFETTE_GIT_FULLNAME").Required().String()
	gitRevision          = kingpin.Flag("git-revision", "The hash of the revision to set build status for.").Envar("ESTAFETTE_GIT_REVISION").Required().String()
	estafetteBuildStatus = kingpin.Flag("estafette-build-status", "The current build status of the Estafette pipeline.").Envar("ESTAFETTE_BUILD_STATUS").Required().String()
	statusOverride       = kingpin.Flag("status-override", "Allow status property in manifest to override the actual build status.").Envar("ESTAFETTE_EXTENSION_STATUS").String()
	ciBaseURL            = kingpin.Flag("estafette-ci-server-base-url", "The base url of the ci server.").Envar("ESTAFETTE_CI_SERVER_BASE_URL").Required().String()
	estafetteBuildID     = kingpin.Flag("estafette-build-id", "The build id of this particular build.").Envar("ESTAFETTE_BUILD_ID").Required().String()

	estafetteBuildVersion = kingpin.Flag("estafette-build-version", "The current build version of the Estafette pipeline.").Envar("ESTAFETTE_BUILD_VERSION").Required().String()
	releaseName           = kingpin.Flag("release-name", "Name of the release section, automatically set by Estafette CI.").Envar("ESTAFETTE_RELEASE_NAME").String()
	releaseAction         = kingpin.Flag("release-action", "Name of the release action, automatically set by Estafette CI.").Envar("ESTAFETTE_RELEASE_ACTION").String()

	apiTokenPath = kingpin.Flag("credentials-path", "Path to file with Bitbucket api token credentials configured at the CI server, passed in to this trusted extension.").Default("/credentials/bitbucket_api_token.json").String()
)

func main() {

	// parse command line parameters
	kingpin.Parse()

	// init log format from envvar ESTAFETTE_LOG_FORMAT
	foundation.InitLoggingFromEnv(appgroup, app, version, branch, revision, buildDate)

	// check if there's a status override
	status := *estafetteBuildStatus
	if *statusOverride != "" {
		status = *statusOverride
	}

	// make sure ciBaseURL ends with a slash
	if !strings.HasSuffix(*ciBaseURL, "/") {
		*ciBaseURL = *ciBaseURL + "/"
	}

	// get api token from injected credentials
	var credentials []APITokenCredentials

	// use mounted credential file if present instead of relying on an envvar
	if runtime.GOOS == "windows" {
		*apiTokenPath = "C:" + *apiTokenPath
	}
	if foundation.FileExists(*apiTokenPath) {
		log.Info().Msgf("Reading credentials from file at path %v...", *apiTokenPath)
		credentialsFileContent, err := ioutil.ReadFile(*apiTokenPath)
		if err != nil {
			log.Fatal().Msgf("Failed reading credential file at path %v.", *apiTokenPath)
		}
		err = json.Unmarshal(credentialsFileContent, &credentials)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed unmarshalling injected credentials")
		}
		if len(credentials) == 0 {
			log.Warn().Str("data", string(credentialsFileContent)).Msgf("Found 0 credentials in file %v", *apiTokenPath)
		}
		log.Debug().Msgf("Read %v credentials", len(credentials))
	}
	if len(credentials) == 0 {
		log.Fatal().Msg("No credentials have been injected")
	}

	// set build status
	bitbucketAPIClient := newBitbucketAPIClient()
	err := bitbucketAPIClient.SetBuildStatus(credentials[0].AdditionalProperties.Token, *gitRepoFullname, *gitRevision, status, *estafetteBuildVersion, *releaseName, *releaseAction)
	if err != nil {
		log.Fatal().Err(err).Msg("Updating Bitbucket build status failed")
	}

	log.Info().Msg("Finished estafette-extension-bitbucket-status...")
}
