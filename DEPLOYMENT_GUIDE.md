# Deployment Guide - GCP Cloud Run Infrastructure

This guide documents the Pulumi infrastructure setup for deploying a containerized application to Google Cloud Run with Artifact Registry and Firestore.

## 📋 What Was Built

A complete GCP infrastructure stack that includes:

1. **Artifact Registry** - Docker repository for container images
2. **Docker Image Build & Push** - Automated using Pulumi Command resources
3. **Firestore Database** - Native mode NoSQL database
4. **IAM Service Account** - With least-privilege `roles/datastore.user` permission
5. **Cloud Run Service** - Serverless container deployment with Firestore integration
6. **Public Access** - Cloud Run service accessible without authentication

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                     GCP Project                         │
│                                                         │
│  ┌──────────────────┐      ┌────────────────────┐     │
│  │ Artifact Registry│      │  Firestore Native  │     │
│  │   (Docker repo)  │      │    (Database)      │     │
│  └──────────────────┘      └────────────────────┘     │
│           │                          │                 │
│           │ pulls image              │ read/write      │
│           ▼                          ▼                 │
│  ┌──────────────────────────────────────────────┐     │
│  │         Cloud Run Service                    │     │
│  │  - Container: app:latest                     │     │
│  │  - Service Account: *-cloudrun-sa            │     │
│  │  - Public Access: allUsers                   │     │
│  │  - Env: FIRESTORE_DATABASE, GCP_PROJECT      │     │
│  └──────────────────────────────────────────────┘     │
│                      │                                 │
│                      │ HTTPS                           │
│                      ▼                                 │
│              Public Internet                           │
└─────────────────────────────────────────────────────────┘
```

## 🎯 Implementation Details

### Key Configuration Files

#### `Pulumi.dev.yaml` (Current Configuration)
```yaml
config:
  gcp:project: "dev-proj-480923"
  gcp:region: europe-north1
  hello-pulumi:resourcePrefix: test-pulumi-trial
  hello-pulumi:cloudRunAllowUnauthenticated: "true"
  hello-pulumi:firestoreLocation: europe-north1
```

### Resources Created

| Resource Type | Name | Purpose |
|--------------|------|---------|
| Artifact Registry | `test-pulumi-trial-registry` | Stores Docker images |
| Firestore Database | `test-pulumi-trial-db` | NoSQL data persistence |
| Service Account | `test-pulumi-trial-cloudrun-sa` | Cloud Run identity |
| Cloud Run Service | `test-pulumi-trial-service` | Runs containerized app |
| IAM Binding | `test-pulumi-trial-firestore-user` | Grants Firestore access |
| IAM Binding | `test-pulumi-trial-service-invoker` | Allows public access |

### Docker Image Build Strategy

**Initial Problem:** Pulumi Docker provider had BuildKit cache authentication issues (403 errors)

**Solution:** Used Pulumi Command provider to run native Docker commands:
1. `docker build` - Builds image locally
2. `docker push` - Pushes to Artifact Registry using gcloud credentials
3. Explicit dependency ensures Cloud Run waits for image availability

### IAM Permissions (Least Privilege)

```
Service Account: test-pulumi-trial-cloudrun-sa@dev-proj-480923.iam.gserviceaccount.com
├─ roles/datastore.user (Project-level)
│  └─ Allows: Firestore read/write operations
│
Cloud Run Service IAM:
└─ roles/run.invoker → allUsers
   └─ Allows: Unauthenticated public access
```

## 🚀 Deployment Steps

### Prerequisites Checklist

- [x] GCP Project: `dev-proj-480923`
- [x] Pulumi CLI installed
- [x] Go 1.25+ installed
- [x] Docker installed and running
- [x] gcloud CLI authenticated
- [x] GCS bucket for state: `scratch-test-pulumi-state-dev`

### Step-by-Step Deployment

#### 1. Authenticate with GCP

```bash
# Login with your GCP account
gcloud auth application-default login

# Set your project
gcloud config set project dev-proj-480923

# Configure Docker for Artifact Registry
gcloud auth configure-docker europe-north1-docker.pkg.dev --quiet
```

#### 2. Enable Required GCP APIs

```bash
gcloud services enable \
  artifactregistry.googleapis.com \
  run.googleapis.com \
  firestore.googleapis.com \
  iamcredentials.googleapis.com \
  --project=dev-proj-480923
```

#### 3. Initialize Pulumi

```bash
cd infra

# Login to GCS backend
pulumi login gs://scratch-test-pulumi-state-dev

# Initialize stack (if not already done)
pulumi stack init dev
```

#### 4. Configure Stack (Optional - Already Set)

Your configuration is already set in `Pulumi.dev.yaml`. To modify:

```bash
# Change project
pulumi config set gcp:project YOUR_PROJECT_ID

# Change region
pulumi config set gcp:region YOUR_REGION

# Change resource prefix
pulumi config set resourcePrefix YOUR_PREFIX

# Toggle public access
pulumi config set cloudRunAllowUnauthenticated true/false
```

#### 5. Deploy Infrastructure

```bash
# Preview changes
pulumi preview

# Deploy (with confirmation)
pulumi up

# Or deploy without confirmation
pulumi up --yes
```

**Deployment Time:** ~30-45 seconds

### What Happens During Deployment

1. ✅ Creates Artifact Registry repository
2. ✅ Builds Docker image from `../app/` directory
3. ✅ Pushes image to Artifact Registry
4. ✅ Creates Firestore Native database
5. ✅ Creates IAM service account
6. ✅ Grants Firestore permissions to service account
7. ✅ Deploys Cloud Run service with environment variables
8. ✅ Configures public access IAM policy

## 🧪 Testing Your Deployment

### Get Service URL

```bash
# From Pulumi
pulumi stack output cloudRunURL

# From gcloud
gcloud run services describe test-pulumi-trial-service \
  --region=europe-north1 \
  --project=dev-proj-480923 \
  --format="value(status.url)"
```

### Test Endpoints

```bash
# Set the URL
URL=$(pulumi stack output cloudRunURL)

# Health check
curl $URL/health
# Expected: OK

# Root endpoint
curl $URL/
# Expected: Hello from Cloud Run! 🚀
#           Version: 1.0.0
#           Firestore: Connected

# Write to Firestore
curl $URL/write
# Expected: Successfully wrote to Firestore!

# Read from Firestore
curl $URL/read
# Expected: Data from Firestore: map[message:Hello from Cloud Run timestamp:...]
```

## 🔄 Updating Your Application

### Modify Application Code

1. Edit `app/main.go` with your changes
2. Run `pulumi up` - this will:
   - Rebuild the Docker image
   - Push new image to Artifact Registry
   - Update Cloud Run service with new image

### Versioning Strategy

Current setup uses `latest` tag. For production, consider:

```go
// In infra/resources/image.go, modify to use git commit or timestamp
imageName := pulumi.Sprintf("%s/%s-app:%s", registryURL, prefix, version)
```

## 📊 Stack Outputs

After deployment, view all outputs:

```bash
pulumi stack output
```

**Available Outputs:**
- `artifactRegistryURL` - Docker registry URL
- `cloudRunURL` - Public service URL
- `dockerImageName` - Full image name with tag
- `firestoreDatabase` - Database ID
- `serviceAccountEmail` - Service account email

## 🐛 Troubleshooting

### Issue: Docker Push 403 Forbidden

**Symptom:** `error reading build output: failed to load cache key: invalid response status 403`

**Solution:**
```bash
gcloud auth configure-docker europe-north1-docker.pkg.dev --quiet
gcloud auth application-default login
```

### Issue: Cloud Run Can't Find Image

**Symptom:** `Error code 5, message: Image '...' not found`

**Solution:** Ensure Docker push completes before Cloud Run deployment. The code now has explicit dependencies to prevent this.

### Issue: Service Already Exists (409)

**Symptom:** `Error 409: Resource 'test-pulumi-trial-service' already exists`

**Solution:**
```bash
# Delete the service manually
gcloud run services delete test-pulumi-trial-service \
  --region=europe-north1 \
  --project=dev-proj-480923 \
  --quiet

# Retry deployment
pulumi up
```

### Issue: Firestore Access Denied

**Symptom:** Cloud Run can write/read from Firestore

**Solution:** Verify IAM binding:
```bash
gcloud projects get-iam-policy dev-proj-480923 \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:*cloudrun-sa*"
```

## 🗑️ Cleanup

### Destroy All Resources

```bash
cd infra
pulumi destroy --yes
```

This removes:
- Cloud Run service
- Service account and IAM bindings
- Firestore database
- Artifact Registry repository (and all images)

### Remove Stack

```bash
pulumi stack rm dev --yes
```

## 📁 Project Structure

```
pulumi-trial/
├── app/                              # Application code
│   ├── Dockerfile                    # Multi-stage Docker build
│   ├── go.mod                        # App dependencies
│   └── main.go                       # HTTP server with Firestore
│
├── infra/                            # Infrastructure as Code
│   ├── config/
│   │   └── config.go                 # Configuration management
│   ├── resources/
│   │   ├── artifact_registry.go      # Registry setup
│   │   ├── cloudrun.go               # Cloud Run service
│   │   ├── firestore.go              # Firestore database
│   │   ├── iam.go                    # Service account & IAM
│   │   └── image.go                  # Docker build & push
│   ├── go.mod                        # Infra dependencies
│   ├── main.go                       # Pulumi entry point
│   ├── Pulumi.yaml                   # Project definition
│   ├── Pulumi.dev.yaml               # Stack configuration
│   └── README.md                     # Detailed documentation
│
├── .gitignore                        # Git ignore rules
├── go.mod                            # Root module
└── DEPLOYMENT_GUIDE.md              # This file
```

## 🔐 Security Best Practices

### Current Implementation

✅ **Service Account** - Dedicated identity per service  
✅ **Least Privilege** - Only `datastore.user` role granted  
✅ **Resource Isolation** - Separate service account per Cloud Run service  
✅ **Artifact Registry** - Controlled image repository  

### Recommended for Production

⚠️ **Authentication** - Set `cloudRunAllowUnauthenticated: false`  
⚠️ **Secret Management** - Use Google Secret Manager for sensitive data  
⚠️ **Network Security** - Consider VPC connector for private resources  
⚠️ **Image Scanning** - Enable binary authorization  
⚠️ **Monitoring** - Set up Cloud Monitoring alerts  

## 🎓 Key Learnings

1. **Pulumi Docker Provider Limitations** - BuildKit cache can cause auth issues; native Docker commands are more reliable
2. **Resource Dependencies** - Explicit `pulumi.DependsOn` prevents race conditions
3. **GCP Authentication** - `gcloud auth configure-docker` is essential for Artifact Registry
4. **State Management** - GCS backend enables team collaboration
5. **Configuration Management** - Direct YAML editing works same as `pulumi config set`

## 📚 Next Steps

### To Use This Infrastructure

1. **Modify the App** - Replace `app/main.go` with your application
2. **Update Dockerfile** - Adjust for your app's requirements
3. **Configure Secrets** - Add Secret Manager integration if needed
4. **Set Up CI/CD** - Automate deployments with GitHub Actions
5. **Add Monitoring** - Integrate Cloud Monitoring/Logging
6. **Domain Mapping** - Add custom domain to Cloud Run

### To Extend This Infrastructure

- Add Cloud SQL database
- Implement Cloud Tasks for async jobs
- Add Cloud Storage buckets
- Set up Cloud CDN
- Configure Cloud Armor for DDoS protection
- Add multiple environments (staging, prod)

## 🆘 Support

For issues specific to:
- **Pulumi**: Check [Pulumi docs](https://www.pulumi.com/docs/)
- **GCP**: See [Google Cloud docs](https://cloud.google.com/docs)
- **Cloud Run**: Review [Cloud Run docs](https://cloud.google.com/run/docs)

## ✅ Deployment Checklist

Before deploying to production:

- [ ] Review and update resource names in `Pulumi.dev.yaml`
- [ ] Set `cloudRunAllowUnauthenticated: false` for authenticated access
- [ ] Enable Cloud Armor for DDoS protection
- [ ] Set up Cloud Monitoring alerts
- [ ] Configure backup strategy for Firestore
- [ ] Review IAM permissions and service accounts
- [ ] Set up multiple environments (dev/staging/prod)
- [ ] Configure custom domain and SSL
- [ ] Enable container image scanning
- [ ] Document runbook for incident response

---

**Status:** ✅ Fully deployed and tested on February 24, 2026

**Deployed URL:** `https://test-pulumi-trial-service-ytenlbc2aq-lz.a.run.app`

**Last Updated:** February 24, 2026
