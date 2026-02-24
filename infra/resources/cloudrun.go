package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/cloudrunv2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// CloudRunOutputs contains the outputs from Cloud Run resources
type CloudRunOutputs struct {
	Service    *cloudrunv2.Service
	ServiceURL pulumi.StringOutput
}

// CreateCloudRunService creates a Cloud Run service with the specified configuration
func CreateCloudRunService(
	ctx *pulumi.Context,
	project string,
	region string,
	prefix string,
	imageName pulumi.StringOutput,
	serviceAccountEmail pulumi.StringOutput,
	firestoreDatabaseID pulumi.StringOutput,
	allowUnauthenticated bool,
	provider pulumi.ProviderResource,
	imageDependency pulumi.Resource,
) (*CloudRunOutputs, error) {
	serviceName := fmt.Sprintf("%s-service", prefix)

	// Create Cloud Run service
	service, err := cloudrunv2.NewService(ctx, serviceName, &cloudrunv2.ServiceArgs{
		Name:     pulumi.String(serviceName),
		Location: pulumi.String(region),
		Project:  pulumi.String(project),
		Template: &cloudrunv2.ServiceTemplateArgs{
			ServiceAccount: serviceAccountEmail,
			Containers: cloudrunv2.ServiceTemplateContainerArray{
				&cloudrunv2.ServiceTemplateContainerArgs{
					Image: imageName,
					Envs: cloudrunv2.ServiceTemplateContainerEnvArray{
						&cloudrunv2.ServiceTemplateContainerEnvArgs{
							Name:  pulumi.String("FIRESTORE_DATABASE"),
							Value: firestoreDatabaseID,
						},
						&cloudrunv2.ServiceTemplateContainerEnvArgs{
							Name:  pulumi.String("GCP_PROJECT"),
							Value: pulumi.String(project),
						},
					},
					Ports: cloudrunv2.ServiceTemplateContainerPortArray{
						&cloudrunv2.ServiceTemplateContainerPortArgs{
							ContainerPort: pulumi.Int(8080),
						},
					},
					Resources: &cloudrunv2.ServiceTemplateContainerResourcesArgs{
						Limits: pulumi.StringMap{
							"cpu":    pulumi.String("1"),
							"memory": pulumi.String("512Mi"),
						},
					},
				},
			},
			Scaling: &cloudrunv2.ServiceTemplateScalingArgs{
				MinInstanceCount: pulumi.Int(0),
				MaxInstanceCount: pulumi.Int(10),
			},
		},
	}, pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{imageDependency}))
	if err != nil {
		return nil, err
	}

	// Configure IAM policy for public access if enabled
	if allowUnauthenticated {
		_, err = cloudrunv2.NewServiceIamMember(ctx, fmt.Sprintf("%s-invoker", serviceName), &cloudrunv2.ServiceIamMemberArgs{
			Project:  service.Project,
			Location: service.Location,
			Name:     service.Name,
			Role:     pulumi.String("roles/run.invoker"),
			Member:   pulumi.String("allUsers"),
		}, pulumi.Provider(provider))
		if err != nil {
			return nil, err
		}
	}

	return &CloudRunOutputs{
		Service:    service,
		ServiceURL: service.Uri,
	}, nil
}
