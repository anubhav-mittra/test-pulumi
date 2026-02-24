# GCP Cloud Run Infrastructure with Pulumi

This repository contains Pulumi infrastructure code to deploy a containerized application on Google Cloud Platform using Cloud Run, Artifact Registry, and Firestore.

## Architecture

- **Artifact Registry**: Docker image repository for application containers
- **Cloud Run**: Serverless container deployment
- **Firestore**: NoSQL database (Native mode)
- **IAM**: Service account with least-privilege permissions (datastore.user role)

## Prerequisites

1. **GCP Account** with a project created
2. **Pulumi CLI** installed ([Installation guide](https://www.pulumi.com/docs/get-started/install/))
3. **Go 1.25+** installed
4. **Docker** installed and running
5. **GCP APIs enabled**:
   - Artifact Registry API (`artifactregistry.googleapis.com`)
   - Cloud Run API (`run.googleapis.com`)
   - Firestore API (`firestore.googleapis.com`)
   - IAM Service Account Credentials API (`iamcredentials.googleapis.com`)

## Configuration

The infrastructure uses the following configurable parameters:

| Parameter | Description | Default |
|-----------|-------------|---------|
| `gcp:project` | GCP Project ID | (required) |
| `gcp:region` | GCP region for resources | `us-central1` |
| `resourcePrefix` | Prefix for all resource names | `test-pulumi-trial` |
| `cloudRunAllowUnauthenticated` | Allow public access to Cloud Run | `true` |
| `firestoreLocation` | Firestore database location | Same as `gcp:region` |

## Setup

### 1. Initialize GCS Backend

The Pulumi state is stored in GCS bucket `scratch-test-pulumi-state-dev`. Log in to the backend:

```bash
cd infra
pulumi login gs://scratch-test-pulumi-state-dev
```

### 2. Create a Stack

Initialize a new stack (e.g., `dev`):

```bash
pulumi stack init dev
```

### 3. Configure GCP Project

Set your GCP project ID and region:

```bash
pulumi config set gcp:project YOUR_PROJECT_ID
pulumi config set gcp:region us-central1
```

### 4. (Optional) Customize Resource Names

Change the resource prefix if needed:

```bash
pulumi config set resourcePrefix myapp-dev
```

### 5. Authenticate with GCP

Ensure you're authenticated with GCP:

```bash
gcloud auth application-default login
```

## Deployment

### Deploy Infrastructure

Preview the changes:

```bash
pulumi preview
```

Deploy the infrastructure:

```bash
pulumi up
```

The deployment will:
1. Create an Artifact Registry repository
2. Build and push the Docker image
3. Create a Firestore Native database
4. Create a service account with Firestore access
5. Deploy the Cloud Run service

### Access the Service

After deployment, get the Cloud Run URL:

```bash
pulumi stack output cloudRunURL
```

Test the endpoints:

```bash
# Health check
curl $(pulumi stack output cloudRunURL)/health

# Root endpoint
curl $(pulumi stack output cloudRunURL)/

# Write to Firestore
curl $(pulumi stack output cloudRunURL)/write

# Read from Firestore
curl $(pulumi stack output cloudRunURL)/read
```

## Stack Outputs

The following outputs are exported:

- `artifactRegistryURL`: Artifact Registry repository URL
- `dockerImageName`: Full Docker image name with tag
- `firestoreDatabase`: Firestore database ID
- `serviceAccountEmail`: Service account email for Cloud Run
- `cloudRunURL`: Public URL of the Cloud Run service

## Updating the Application

To update the application code:

1. Modify the application code in `app/main.go`
2. Run `pulumi up` to rebuild and redeploy the image
3. Pulumi will automatically detect changes and update the Cloud Run service

## Resource Naming

With the default prefix `test-pulumi-trial`, resources are named:

- Artifact Registry: `test-pulumi-trial-registry`
- Firestore Database: `test-pulumi-trial-db`
- Service Account: `test-pulumi-trial-cloudrun-sa`
- Cloud Run Service: `test-pulumi-trial-service`
- Docker Image: `{location}-docker.pkg.dev/{project}/test-pulumi-trial-registry/test-pulumi-trial-app:latest`

## IAM Permissions

The infrastructure follows the principle of least privilege:

- **Service Account**: `test-pulumi-trial-cloudrun-sa@{project}.iam.gserviceaccount.com`
  - Role: `roles/datastore.user` (Firestore read/write access)
  - Scope: Project-level (required by GCP IAM structure)

## Cleanup

To destroy all resources:

```bash
pulumi destroy
```

To remove the stack:

```bash
pulumi stack rm dev
```

## Project Structure

```
.
├── app/                      # Application code
│   ├── main.go              # HTTP server with Firestore integration
│   ├── Dockerfile           # Multi-stage Docker build
│   └── go.mod               # Application dependencies
├── infra/                   # Infrastructure code
│   ├── main.go             # Pulumi program entry point
│   ├── config/
│   │   └── config.go       # Configuration loading
│   └── resources/
│       ├── artifact_registry.go  # Artifact Registry setup
│       ├── cloudrun.go          # Cloud Run service
│       ├── firestore.go         # Firestore database
│       ├── iam.go               # Service account & IAM
│       └── image.go             # Docker image build & push
├── Pulumi.yaml             # Pulumi project definition
└── Pulumi.dev.yaml         # Stack configuration (dev)
```

## Troubleshooting

### API Not Enabled

If you encounter errors about APIs not being enabled, run:

```bash
gcloud services enable artifactregistry.googleapis.com \
  run.googleapis.com \
  firestore.googleapis.com \
  iamcredentials.googleapis.com \
  --project=YOUR_PROJECT_ID
```

### Docker Build Fails

Ensure Docker is running:

```bash
docker info
```

### Firestore Access Denied

Verify the service account has the correct IAM role:

```bash
gcloud projects get-iam-policy YOUR_PROJECT_ID \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:*cloudrun-sa*"
```

## License

This project is provided as-is for educational and infrastructure provisioning purposes.
