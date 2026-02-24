package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/firestore"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// FirestoreOutputs contains the outputs from Firestore resources
type FirestoreOutputs struct {
	Database   *firestore.Database
	DatabaseID pulumi.StringOutput
}

// CreateFirestore creates a Firestore database in Native mode
func CreateFirestore(ctx *pulumi.Context, project string, location string, prefix string, provider pulumi.ProviderResource) (*FirestoreOutputs, error) {
	databaseID := fmt.Sprintf("%s-db", prefix)

	// Create Firestore database in Native mode
	database, err := firestore.NewDatabase(ctx, databaseID, &firestore.DatabaseArgs{
		Project:                       pulumi.String(project),
		Name:                          pulumi.String(databaseID),
		LocationId:                    pulumi.String(location),
		Type:                          pulumi.String("FIRESTORE_NATIVE"),
		ConcurrencyMode:               pulumi.String("OPTIMISTIC"),
		AppEngineIntegrationMode:      pulumi.String("DISABLED"),
		PointInTimeRecoveryEnablement: pulumi.String("POINT_IN_TIME_RECOVERY_DISABLED"),
		DeleteProtectionState:         pulumi.String("DELETE_PROTECTION_DISABLED"),
		DeletionPolicy:                pulumi.String("DELETE"),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	return &FirestoreOutputs{
		Database:   database,
		DatabaseID: database.Name,
	}, nil
}
