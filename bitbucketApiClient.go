package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/sethgrid/pester"
)

// BitbucketAPIClient communicates with the Bitbucket api
type BitbucketAPIClient interface {
	SetBuildStatus(string, string, string, string) error
}

type bitbucketAPIClientImpl struct {
}

func newBitbucketAPIClient() BitbucketAPIClient {
	return &bitbucketAPIClientImpl{}
}

type buildStatusRequestBody struct {
	State       string `json:"state"`
	Key         string `json:"key"`
	Name        string `json:"name,omitempty"`
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

// SetBuildStatus sets the build status for a specific revision
func (gh *bitbucketAPIClientImpl) SetBuildStatus(accessToken, repoFullname, gitRevision, status string) (err error) {

	// https://confluence.atlassian.com/bitbucket/buildstatus-resource-779295267.html
	// estafette status: succeeded|failed|pending
	// bitbucket stat: INPROGRESS|SUCCESSFUL|FAILED|STOPPED

	state := "SUCCESSFUL"
	switch status {
	case "succeeded":
		state = "SUCCESSFUL"

	case "failed":
		state = "FAILED"

	case "pending":
		state = "INPROGRESS"
	}

	logsURL := fmt.Sprintf(
		"%vlogs/%v/%v/%v/%v",
		os.Getenv("ESTAFETTE_CI_SERVER_BASE_URL"),
		os.Getenv("ESTAFETTE_GIT_SOURCE"),
		os.Getenv("ESTAFETTE_GIT_NAME"),
		os.Getenv("ESTAFETTE_GIT_BRANCH"),
		os.Getenv("ESTAFETTE_GIT_REVISION"),
	)

	params := buildStatusRequestBody{
		State: state,
		Key:   "estafette",
		URL:   logsURL,
	}

	_, err = callBitbucketAPI("POST", fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%v/commit/%v/statuses/build", repoFullname, gitRevision), params, "Bearer", accessToken)

	return
}

func callBitbucketAPI(method, url string, params interface{}, authorizationType, token string) (body []byte, err error) {

	// convert params to json if they're present
	var requestBody io.Reader
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return body, err
		}
		requestBody = bytes.NewReader(data)
	}

	// create client, in order to add headers
	client := pester.New()
	client.MaxRetries = 3
	client.Backoff = pester.ExponentialJitterBackoff
	client.KeepLog = true
	request, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		return
	}

	// add headers
	request.Header.Add("Authorization", fmt.Sprintf("%v %v", authorizationType, token))
	request.Header.Add("Content-Type", "application/json")

	// perform actual request
	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	// unmarshal json body
	var b interface{}
	err = json.Unmarshal(body, &b)
	if err != nil {
		log.Error().Err(err).
			Str("url", url).
			Str("requestMethod", method).
			Interface("requestBody", params).
			Interface("requestHeaders", request.Header).
			Interface("responseHeaders", response.Header).
			Str("responseBody", string(body)).
			Msg("Deserializing response for '%v' Bitbucket api call failed")

		return
	}

	log.Debug().
		Str("url", url).
		Str("requestMethod", method).
		Interface("requestBody", params).
		Interface("requestHeaders", request.Header).
		Interface("responseHeaders", response.Header).
		Interface("responseBody", b).
		Msgf("Received response for '%v' Bitbucket api call...", url)

	return
}
