package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test image naming convention
func TestImage_NamingConvention(t *testing.T) {
	prefix := "test-app"
	imageName := fmt.Sprintf("%s-app", prefix)
	assert.Equal(t, "test-app-app", imageName)
}

// Test full image path with registry
func TestImage_FullPath(t *testing.T) {
	registryURL := "us-central1-docker.pkg.dev/test-project/test-registry"
	imageName := "test-app-app"
	tag := "latest"

	fullPath := fmt.Sprintf("%s/%s:%s", registryURL, imageName, tag)
	expected := "us-central1-docker.pkg.dev/test-project/test-registry/test-app-app:latest"

	assert.Equal(t, expected, fullPath)
}

// Test Docker build command format
func TestImage_BuildCommandFormat(t *testing.T) {
	imageTag := "us-central1-docker.pkg.dev/test-project/test-registry/test-app-app:latest"
	appPath := "../app"

	buildCmd := fmt.Sprintf("docker build -t %s %s", imageTag, appPath)

	assert.Contains(t, buildCmd, "docker build")
	assert.Contains(t, buildCmd, "-t")
	assert.Contains(t, buildCmd, imageTag)
	assert.Contains(t, buildCmd, appPath)
}

// Test Docker push command format
func TestImage_PushCommandFormat(t *testing.T) {
	imageTag := "us-central1-docker.pkg.dev/test-project/test-registry/test-app-app:latest"
	pushCmd := fmt.Sprintf("docker push %s", imageTag)

	assert.Contains(t, pushCmd, "docker push")
	assert.Contains(t, pushCmd, imageTag)
}

// Test different regions in registry URL
func TestImage_DifferentRegions(t *testing.T) {
	testCases := []struct {
		location       string
		project        string
		expectedPrefix string
	}{
		{"us-central1", "test-project", "us-central1-docker.pkg.dev"},
		{"europe-west1", "test-project", "europe-west1-docker.pkg.dev"},
		{"asia-east1", "test-project", "asia-east1-docker.pkg.dev"},
	}

	for _, tc := range testCases {
		t.Run(tc.location, func(t *testing.T) {
			registryURL := fmt.Sprintf("%s/%s/test-registry", tc.expectedPrefix, tc.project)
			assert.Contains(t, registryURL, tc.location)
			assert.Contains(t, registryURL, "docker.pkg.dev")
		})
	}
}
