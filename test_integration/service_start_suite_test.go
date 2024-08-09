//go:build !unit
// +build !unit

package testintegration

import (
	"context"
	"time"

	"github.com/Ropho/avito-bootcamp-assignment/internal/boot"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/require"
)

type ServiceStartSuite struct {
	suite.Suite
	ctx context.Context
}

const (
	configPath = "../config"
	configName = "config"
)

func (s *ServiceStartSuite) BeforeAll(provider.T) {}

func (s *ServiceStartSuite) TestService_ServiceStart(t provider.T) {
	s.ctx = context.Background()

	go func(t provider.T) {

		bootErr := boot.App(configPath, configName)
		require.NoError(t, bootErr, "failed to init service")
	}(t)

}

func (s *ServiceStartSuite) AfterAll(t provider.T) {
	t.Log("waiting for server to start")
	time.Sleep(5 * time.Second)
}
