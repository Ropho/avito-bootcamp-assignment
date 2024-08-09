package testintegration

import (
	"os"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestServiceStart(t *testing.T) {
	suite.RunSuite(t, new(ServiceStartSuite))
}

func TestFlactCreate(t *testing.T) {
	suite.RunSuite(t, new(FlatCreateSuite))
}
