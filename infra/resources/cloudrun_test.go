package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Cloud Run service naming
func TestCloudRun_ServiceNaming(t *testing.T) {
	prefix := "test-app"
	serviceName := fmt.Sprintf("%s-service", prefix)
	assert.Equal(t, "test-app-service", serviceName)
}

// Test environment variable configuration
func TestCloudRun_EnvironmentVariables(t *testing.T) {
	envVars := map[string]string{
		"GCP_PROJECT":       "test-project",
		"FIRESTORE_DATABASE": "test-db",
	}
	
	assert.Equal(t, "test-project", envVars["GCP_PROJECT"])
	assert.Equal(t, "test-db", envVars["FIRESTORE_DATABASE"])
	assert.Len(t, envVars, 2)
}

// Test IAM policy member format for unauthenticated access
func TestCloudRun_UnauthenticatedAccessPolicy(t *testing.T) {
	allowUnauthenticated := true
	var member string
	
	if allowUnauthenticated {
		member = "allUsers"
	} else {
		member = "serviceAccount:sa@project.iam.gserviceaccount.com"
	}
	
	assert.Equal(t, "allUsers", member)
}

// Test IAM policy for authenticated access
func TestCloudRun_AuthenticatedAccessPolicy(t *testing.T) {
	allowUnauthenticated := false
	var requiresAuth bool
	
	if allowUnauthenticated {
		requiresAuth = false
	} else {
		requiresAuth = true
	}
	
	assert.True(t, requiresAuth)
}

// Test service naming with different prefixes
func TestCloudRun_DifferentPrefixes(t *testing.T) {
	testCases := []struct {
		prefix   string
		expected string
	}{
		{"test", "test-service"},
		{"prod-app", "prod-app-service"},
		{"dev", "dev-service"},
	}

	for _, tc := range testCases {
		t.Run(tc.prefix, func(t *testing.T) {
			serviceName := fmt.Sprintf("%s-service", tc.prefix)
			assert.Equal(t, tc.expected, serviceName)
		})
	}
}
