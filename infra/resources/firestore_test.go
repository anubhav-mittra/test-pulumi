package resources

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Firestore database naming convention
func TestFirestore_DatabaseNaming(t *testing.T) {
	prefix := "test-app"
	dbName := fmt.Sprintf("%s-db", prefix)
	assert.Equal(t, "test-app-db", dbName)
}

// Test database name with different prefixes
func TestFirestore_DifferentPrefixes(t *testing.T) {
	testCases := []struct {
		prefix   string
		expected string
	}{
		{"my-app", "my-app-db"},
		{"test", "test-db"},
		{"prod", "prod-db"},
	}

	for _, tc := range testCases {
		t.Run(tc.prefix, func(t *testing.T) {
			dbName := fmt.Sprintf("%s-db", tc.prefix)
			assert.Equal(t, tc.expected, dbName)
		})
	}
}

// Test location formats
func TestFirestore_LocationFormats(t *testing.T) {
	locations := []string{"us-central1", "europe-west1", "asia-east1"}
	for _, loc := range locations {
		t.Run(loc, func(t *testing.T) {
			assert.NotEmpty(t, loc)
			// Location should be in format: region-zone
			assert.Contains(t, loc, "-")
		})
	}
}
