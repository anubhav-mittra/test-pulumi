package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/firestore"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// IAMOutputs contains the outputs from IAM resources
type IAMOutputs struct {
	ServiceAccount      *serviceaccount.Account
	ServiceAccountEmail pulumi.StringOutput
}

// CreateServiceAccount creates a service account for Cloud Run with minimal Firestore permissions
func CreateServiceAccount(ctx *pulumi.Context, project string, prefix string, firestoreDB *firestore.Database, provider pulumi.ProviderResource) (*IAMOutputs, error) {
	accountID := fmt.Sprintf("%s-cloudrun-sa", prefix)

	// Create service account
	sa, err := serviceaccount.NewAccount(ctx, accountID, &serviceaccount.AccountArgs{
		AccountId:   pulumi.String(accountID),
		DisplayName: pulumi.String(fmt.Sprintf("Cloud Run service account for %s", prefix)),
		Project:     pulumi.String(project),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	// Grant Firestore datastore user role to the service account
	// This follows least privilege - only access to Firestore, scoped to the database
	_, err = projects.NewIAMMember(ctx, fmt.Sprintf("%s-firestore-user", prefix), &projects.IAMMemberArgs{
		Project: pulumi.String(project),
		Role:    pulumi.String("roles/datastore.user"),
		Member:  pulumi.Sprintf("serviceAccount:%s", sa.Email),
	}, pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{sa, firestoreDB}))
	if err != nil {
		return nil, err
	}

	return &IAMOutputs{
		ServiceAccount:      sa,
		ServiceAccountEmail: sa.Email,
	}, nil
}
