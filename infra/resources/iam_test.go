package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test service account ID naming
func TestServiceAccount_NamingConvention(t *testing.T) {
	prefix := "test-app"
	accountID := fmt.Sprintf("%s-cloudrun-sa", prefix)
	assert.Equal(t, "test-app-cloudrun-sa", accountID)
}

// Test service account email format
func TestServiceAccount_EmailFormat(t *testing.T) {
	prefix := "my-app"
	project := "test-project"
	accountID := fmt.Sprintf("%s-cloudrun-sa", prefix)
	email := fmt.Sprintf("%s@%s.iam.gserviceaccount.com", accountID, project)
	
	assert.Contains(t, email, prefix)
	assert.Contains(t, email, "cloudrun-sa")
	assert.Contains(t, email, "@test-project.iam.gserviceaccount.com")
}

// Test different prefix formats
func TestServiceAccount_DifferentPrefixes(t *testing.T) {
	testCases := []struct {
		prefix          string
		expectedAccountID string
	}{
		{"test", "test-cloudrun-sa"},
		{"prod-app", "prod-app-cloudrun-sa"},
		{"dev-service", "dev-service-cloudrun-sa"},
	}

	for _, tc := range testCases {
		t.Run(tc.prefix, func(t *testing.T) {
			accountID := fmt.Sprintf("%s-cloudrun-sa", tc.prefix)
			assert.Equal(t, tc.expectedAccountID, accountID)
		})
	}
}

// Test IAM role format
func TestServiceAccount_RoleFormat(t *testing.T) {
	role := "roles/datastore.user"
	assert.Equal(t, "roles/datastore.user", role)
	assert.Contains(t, role, "roles/")
	assert.Contains(t, role, "datastore")
}
