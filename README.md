# **🚀 Go Application Deployment on Kubernetes (AKS)**

This repository contains a full-stack CI/CD solution for deploying a Go application to Azure Kubernetes Service (AKS). The pipeline automates code quality checks, container security scanning, and reliable deployment to separate **Development** and **Production** environments using Helm charts, Nginx Ingress, and Cert-Manager for automated TLS/SSL.

## **🌟 Architecture Overview**

The solution follows a modern GitOps-ready deployment pattern:

| Component | Technology | Role |
| :---- | :---- | :---- |
| **Application** | Go | The core microservice. |
| **CI/CD Pipeline** | Azure DevOps (azure-pipelines.yml) | Automates build, test, security scanning (Trivy), and deployment. |
| **Deployment** | Helm (my-go-app-chart) | Templates Kubernetes manifests and manages application releases across environments. |
| **Containerization** | Docker (docker-compose.ci.yml) | Defines the application image build process. |
| **Routing** | Nginx Ingress Controller | Manages external traffic routing based on hostnames. |
| **Security** | Cert-Manager | Automatically provisions and renews TLS certificates (Let's Encrypt) for Ingress. |
| **Environments** | Kubernetes Namespaces | Strict separation of dev and prod environments. |

## **📂 Repository Structure**

| Folder/File | Description |
| :---- | :---- |
| app/ | Contains the source code for the Go application. |
| my-go-app-chart/ | The Helm chart used to deploy the application to Kubernetes. Includes values.dev.yaml and values.prod.yaml. |
| local-deployment-docker/ | Configuration files (e.g., docker-compose.yaml) for running the application locally using Docker Compose. |
| local-deployment-kubernetes/ | Local Kubernetes manifests (e.g., Minikube/kind) for testing manifests before deployment. |
| azure-pipelines.yml | The complete CI/CD pipeline definition for Azure DevOps. |
| docker-compose.ci.yml | Used by the CI pipeline to build and push the final production-ready Docker image to ACR. |

## **⚙️ Local Development and Testing**

### **Running with Docker Compose**

To quickly run the application and any required services (like MySQL) locally:

1. Navigate to the local-deployment-docker directory.  
2. Run:  
   docker compose up \-d

3. The application should be accessible at http://localhost:8080 (or the port defined in your Compose file).

### **Running with Local Kubernetes**

If you need to test the full Helm chart locally (e.g., using Minikube or kind):

1. Ensure your local Kubernetes cluster is running.  
2. Use the Helm chart to install the application:  
   helm upgrade \--install my-app-local ./my-go-app-chart \-f my-go-app-chart/values.dev.yaml \--namespace default

   *(Note: This assumes your local cluster is configured with an Ingress controller for testing).*

## **☁️ CI/CD Pipeline (azure-pipelines.yml)**

The pipeline defines three main stages:

### **1\. Build & Test (CI)**

This stage ensures the application is robust and secure before deployment.

* **Go Testing:** Runs unit tests and publishes results to the Azure DevOps **Tests** tab.  
* **Code Coverage:** Measures test coverage and publishes results to the **Code Coverage** tab.  
* **SonarCloud:** Integrates static analysis, quality gates, and code smells detection.  
* **Docker Build & Push:** Builds the application image and pushes it to Azure Container Registry (ACR) tagged with $(Build.BuildId).  
* **Trivy Security Scan:** Runs a vulnerability scan on the Docker image, failing the pipeline if Critical or High vulnerabilities are found (published to the **Tests** tab).

### **2\. Deploy to Development (Deploy\_to\_dev)**

This stage deploys to the isolated dev namespace in AKS.

* Uses my-go-app-chart/values.dev.yaml.  
* Ingress Hostname: **dev.bigfirm.online**  
* **Automated TLS:** Cert-Manager watches the Ingress object and automatically provisions and renews the TLS certificate (bigfirm-online-tls) for the Dev environment.

### **3\. Deploy to Production (Deploy\_to\_production)**

This stage deploys the validated build artifact to the critical prod namespace.

* Uses my-go-app-chart/values.prod.yaml.  
* Ingress Hostname: **prod.bigfirm.online** (or www.bigfirm.online)  
* **Security:** Uses separate, dedicated environment variables for production secrets.  
* **Automated TLS:** Cert-Manager provisions a unique, separate TLS certificate (bigfirm-prod-tls) for the Production environment.

## **🔑 Configuration and Secrets Management**

All environment-specific configurations are handled in two ways:

1. **Helm Value Overrides:**  
   * **my-go-app-chart/values.dev.yaml**: Contains settings like the Dev hostname, image pull policy (often set to Always for CI), and Dev-specific resource limits.  
   * **my-go-app-chart/values.prod.yaml**: Contains settings like the Prod hostname, stricter resource limits, and specific environment labels.  
2. **Kubernetes Secrets (app-config):**  
   * Sensitive credentials (e.g., MYSQL\_USER, MYSQL\_PASSWORD) are passed into the pipeline as secret variables.  
   * The pipeline creates a Kubernetes Secret named app-config in both the dev and prod namespaces, which the application consumes.  
   * ***Best Practice:*** Production secrets are stored in dedicated Azure DevOps variable groups/Key Vaults, completely isolated from development secrets.

*Built with ❤️ and automation.*
