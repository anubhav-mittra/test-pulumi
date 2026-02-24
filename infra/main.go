package main

import (
	"github.com/anubhav-mittra/hello-pulumi/infra/config"
	"github.com/anubhav-mittra/hello-pulumi/infra/resources"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Load configuration
		cfg, err := config.LoadConfig(ctx)
		if err != nil {
			return err
		}

		// Create GCP provider
		provider, err := config.GetProvider(ctx, cfg)
		if err != nil {
			return err
		}

		// 1. Create Artifact Registry repository
		registry, err := resources.CreateArtifactRegistry(
			ctx,
			cfg.GCPProject,
			cfg.GCPRegion,
			cfg.ResourcePrefix,
			provider,
		)
		if err != nil {
			return err
		}

		// 2. Build and push Docker image to Artifact Registry
		image, err := resources.BuildAndPushImage(
			ctx,
			registry.RepositoryURL,
			cfg.ResourcePrefix,
			cfg.GCPRegion,
			registry.Repository,
		)
		if err != nil {
			return err
		}

		// 3. Create Firestore database
		firestoreDB, err := resources.CreateFirestore(
			ctx,
			cfg.GCPProject,
			cfg.FirestoreLocation,
			cfg.ResourcePrefix,
			provider,
		)
		if err != nil {
			return err
		}

		// 4. Create service account with Firestore permissions
		iam, err := resources.CreateServiceAccount(
			ctx,
			cfg.GCPProject,
			cfg.ResourcePrefix,
			firestoreDB.Database,
			provider,
		)
		if err != nil {
			return err
		}

		// 5. Deploy Cloud Run service
		cloudRun, err := resources.CreateCloudRunService(
			ctx,
			cfg.GCPProject,
			cfg.GCPRegion,
			cfg.ResourcePrefix,
			image.ImageName,
			iam.ServiceAccountEmail,
			firestoreDB.DatabaseID,
			cfg.CloudRunAllowUnauthenticated,
			provider,
			image.PushCommand,
		)
		if err != nil {
			return err
		}

		// Export stack outputs
		ctx.Export("artifactRegistryURL", registry.RepositoryURL)
		ctx.Export("dockerImageName", image.ImageName)
		ctx.Export("firestoreDatabase", firestoreDB.DatabaseID)
		ctx.Export("serviceAccountEmail", iam.ServiceAccountEmail)
		ctx.Export("cloudRunURL", cloudRun.ServiceURL)

		return nil
	})
}
