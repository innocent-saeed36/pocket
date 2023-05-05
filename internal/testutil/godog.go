package testutil

import (
	"github.com/cucumber/godog"
	"testing"
)

type ScenarioInitilizer func(ctx *godog.ScenarioContext)

func RunGherkinFeature(t *testing.T, path string, initializer ScenarioInitilizer) {
	suite := godog.TestSuite{
		ScenarioInitializer: initializer,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{path},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
