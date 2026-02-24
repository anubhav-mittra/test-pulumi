package resources

import (
	"fmt"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ImageOutputs contains the outputs from Docker image resources
type ImageOutputs struct {
	BuildCommand *local.Command
	PushCommand  *local.Command
	ImageName    pulumi.StringOutput
}

// BuildAndPushImage builds a Docker image and pushes it to Artifact Registry using docker commands
func BuildAndPushImage(ctx *pulumi.Context, registryURL pulumi.StringOutput, prefix string, location string, registry pulumi.Resource) (*ImageOutputs, error) {
	// Construct image name with tag
	imageName := pulumi.Sprintf("%s/%s-app:latest", registryURL, prefix)

	// Build Docker image locally
	buildCommand, err := local.NewCommand(ctx, fmt.Sprintf("%s-docker-build", prefix), &local.CommandArgs{
		Create: pulumi.Sprintf("docker build -t %s --platform linux/amd64 ../app", imageName),
		Update: pulumi.Sprintf("docker build -t %s --platform linux/amd64 ../app", imageName),
	}, pulumi.DependsOn([]pulumi.Resource{registry}))
	if err != nil {
		return nil, err
	}

	// Push Docker image to Artifact Registry
	pushCommand, err := local.NewCommand(ctx, fmt.Sprintf("%s-docker-push", prefix), &local.CommandArgs{
		Create: pulumi.Sprintf("docker push %s", imageName),
		Update: pulumi.Sprintf("docker push %s", imageName),
	}, pulumi.DependsOn([]pulumi.Resource{buildCommand}))
	if err != nil {
		return nil, err
	}

	return &ImageOutputs{
		BuildCommand: buildCommand,
		PushCommand:  pushCommand,
		ImageName:    imageName,
	}, nil
}
