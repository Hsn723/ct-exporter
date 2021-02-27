package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/http-helper"
)

const (
	probeURL = "http://localhost:9809/probe"
)

func buildAndTestDockerImage(tag string, t *testing.T) {
	t.Helper()
	buildOptions := &docker.BuildOptions{
		Tags: []string{tag},
	}
	docker.Build(t, "../", buildOptions)
	opts := &docker.RunOptions{
		Detach:               true,
		Remove:               true,
		EnvironmentVariables: []string{"CERTSPOTTER_TOKEN=hoge"},
		OtherOptions:         []string{"-p", "9809:9809"},
	}
	id := docker.RunAndGetID(t, tag, opts)
	defer docker.Stop(t, []string{id}, &docker.StopOptions{})
	http_helper.HttpGetWithRetry(t, probeURL, nil, 400, "", 3, 10*time.Second)
}

func TestDocker(t *testing.T) {
	tag := "hsn723/ct-exporter-test"
	buildAndTestDockerImage(tag, t)
}
