package config

import (
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

// AppConfig holds all configuration values for the infrastructure
type AppConfig struct {
	GCPProject                   string
	GCPRegion                    string
	ResourcePrefix               string
	CloudRunAllowUnauthenticated bool
	FirestoreLocation            string
}

// LoadConfig reads configuration from Pulumi config and returns AppConfig
func LoadConfig(ctx *pulumi.Context) (*AppConfig, error) {
	conf := config.New(ctx, "")
	gcpConfig := config.New(ctx, "gcp")

	// Get GCP project from provider config
	gcpProject := gcpConfig.Require("project")

	// Get GCP region with default
	gcpRegion := gcpConfig.Get("region")
	if gcpRegion == "" {
		gcpRegion = "us-central1"
	}

	// Get resource prefix with default
	resourcePrefix := conf.Get("resourcePrefix")
	if resourcePrefix == "" {
		resourcePrefix = "test-pulumi-trial"
	}

	// Get Cloud Run authentication setting with default
	cloudRunAllowUnauth := conf.GetBool("cloudRunAllowUnauthenticated")

	// Get Firestore location with fallback to GCP region
	firestoreLocation := conf.Get("firestoreLocation")
	if firestoreLocation == "" {
		firestoreLocation = gcpRegion
	}

	return &AppConfig{
		GCPProject:                   gcpProject,
		GCPRegion:                    gcpRegion,
		ResourcePrefix:               resourcePrefix,
		CloudRunAllowUnauthenticated: cloudRunAllowUnauth,
		FirestoreLocation:            firestoreLocation,
	}, nil
}

// GetProvider creates a GCP provider with the configured project and region
func GetProvider(ctx *pulumi.Context, cfg *AppConfig) (*gcp.Provider, error) {
	return gcp.NewProvider(ctx, "gcp-provider", &gcp.ProviderArgs{
		Project: pulumi.String(cfg.GCPProject),
		Region:  pulumi.String(cfg.GCPRegion),
	})
}
