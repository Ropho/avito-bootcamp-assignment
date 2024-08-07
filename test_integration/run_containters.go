package test_integration

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

var dockerComposePath = "test_integration/docker-compose.yaml"

func RunContainerDeps(t provider.T) {
	compose, err := tc.NewDockerCompose(dockerComposePath)
	assert.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		assert.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	assert.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")
}
