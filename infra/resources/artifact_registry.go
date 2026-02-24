package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ArtifactRegistryOutputs contains the outputs from Artifact Registry resources
type ArtifactRegistryOutputs struct {
	Repository    *artifactregistry.Repository
	RepositoryURL pulumi.StringOutput
}

// CreateArtifactRegistry creates an Artifact Registry repository for Docker images
func CreateArtifactRegistry(ctx *pulumi.Context, project string, location string, prefix string, provider pulumi.ProviderResource) (*ArtifactRegistryOutputs, error) {
	repoID := fmt.Sprintf("%s-registry", prefix)

	// Create Artifact Registry repository
	repo, err := artifactregistry.NewRepository(ctx, repoID, &artifactregistry.RepositoryArgs{
		Location:     pulumi.String(location),
		RepositoryId: pulumi.String(repoID),
		Description:  pulumi.String("Docker image repository for containerized application"),
		Format:       pulumi.String("DOCKER"),
		Project:      pulumi.String(project),
	}, pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}

	// Construct the repository URL for pushing images
	repoURL := pulumi.Sprintf("%s-docker.pkg.dev/%s/%s", location, project, repoID)

	return &ArtifactRegistryOutputs{
		Repository:    repo,
		RepositoryURL: repoURL,
	}, nil
}
