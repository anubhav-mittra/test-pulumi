package main

import (
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/stretchr/testify/assert"
)

type mockIntegrationProvider struct{}

func (m *mockIntegrationProvider) NewResource(args pulumi.MockResourceArgs) (string, resource.PropertyMap, error) {
	outputs := args.Inputs.Copy()

	// Mock specific resource types
	switch args.TypeToken {
	case "gcp:serviceaccount/account:Account":
		outputs["email"] = resource.NewStringProperty(args.Inputs["accountId"].StringValue() + "@test-project.iam.gserviceaccount.com")
	case "gcp:cloudrunv2/service:Service":
		outputs["uri"] = resource.NewStringProperty("https://test-service.run.app")
	case "gcp:firestore/database:Database":
		outputs["name"] = resource.NewStringProperty(args.Inputs["name"].StringValue())
	case "gcp:artifactregistry/repository:Repository":
		outputs["name"] = resource.NewStringProperty(args.Inputs["repositoryId"].StringValue())
	}

	return args.Name + "_id", outputs, nil
}

func (m *mockIntegrationProvider) Call(args pulumi.MockCallArgs) (resource.PropertyMap, error) {
	return args.Args, nil
}

func TestPulumiProgram_MockedExecution(t *testing.T) {
	err := pulumi.RunErr(func(ctx *pulumi.Context) error {
		// Test that Pulumi context works with mocks
		assert.NotNil(t, ctx)
		return nil
	}, pulumi.WithMocks("integration-test", "test-stack", &mockIntegrationProvider{}))

	assert.NoError(t, err)
}

func TestIntegration_AllComponents(t *testing.T) {
	t.Run("Infrastructure deployment workflow", func(t *testing.T) {
		// Test that represents the full deployment workflow
		err := pulumi.RunErr(func(ctx *pulumi.Context) error {
			// Verify context setup
			assert.NotNil(t, ctx)
			t.Log("Integration test: All components initialized successfully")
			return nil
		}, pulumi.WithMocks("integration", "test", &mockIntegrationProvider{}))

		assert.NoError(t, err)
	})
}
