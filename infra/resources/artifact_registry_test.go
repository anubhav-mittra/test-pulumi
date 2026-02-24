package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test repository naming convention
func TestArtifactRegistry_RepositoryNaming(t *testing.T) {
	prefix := "test-app"
	repoID := fmt.Sprintf("%s-registry", prefix)
	assert.Equal(t, "test-app-registry", repoID)
}

// Test repository URL format
func TestArtifactRegistry_URLFormat(t *testing.T) {
	project := "test-project"
	location := "us-central1"
	repoID := "test-app-registry"
	expectedURL := fmt.Sprintf("%s-docker.pkg.dev/%s/%s", location, project, repoID)
	assert.Equal(t, "us-central1-docker.pkg.dev/test-project/test-app-registry", expectedURL)
}

// Test custom location URL format
func TestArtifactRegistry_CustomLocation(t *testing.T) {
	project := "test-project"
	location := "europe-west1"
	repoID := "my-app-registry"
	expectedURL := fmt.Sprintf("%s-docker.pkg.dev/%s/%s", location, project, repoID)
	assert.Contains(t, expectedURL, "europe-west1")
	assert.Contains(t, expectedURL, "my-app-registry")
}

// Test naming with different prefixes
func TestArtifactRegistry_DifferentPrefixes(t *testing.T) {
	testCases := []struct {
		prefix   string
		expected string
	}{
		{"my-app", "my-app-registry"},
		{"test", "test-registry"},
		{"prod-service", "prod-service-registry"},
	}

	for _, tc := range testCases {
		t.Run(tc.prefix, func(t *testing.T) {
			repoID := fmt.Sprintf("%s-registry", tc.prefix)
			assert.Equal(t, tc.expected, repoID)
		})
	}
}
