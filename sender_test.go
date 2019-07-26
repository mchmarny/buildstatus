package main

import (
	"testing"
)

func TestSlackSend(t *testing.T) {

	msg := &CloudBuildNotification{
		Status: "SUCCESS",
		LogURL: "https://console.cloud.google.com/gcr/builds/7e1135b1-edd3-45f8-a66b-a94b8f323e16?project=546930041843",
		Source: CloudBuildSource{
			RepoSource: RepoSource{
				RepoName: "github_mchmarny_knative-gitops-using-cloud-build",
				TagName:  "release-v0.3.5",
			},
		},
	}

	err := send(msg)
	if err != nil {
		t.Errorf("Error on send: %v", err)
	}

}
