# Azure Developer CLI (azd) - ServicePrincipal and Microsoft Graph Application Features

This document provides a comprehensive overview of the ServicePrincipal and Microsoft Graph Application functionality in the Azure Developer CLI (azd).

## Table of Contents

1. [Overview](#overview)
2. [Core Components](#core-components)
3. [Microsoft Graph SDK Integration](#microsoft-graph-sdk-integration)
4. [ServicePrincipal Features](#serviceprincipal-features)
5. [Application Features](#application-features)
6. [Federated Identity Credentials Features](#federated-identity-credentials-features)
7. [EntraId Service Layer](#entraid-service-layer)
8. [Command-Line Integration](#command-line-integration)
9. [Usage Scenarios](#usage-scenarios)

## Overview

azd provides comprehensive ServicePrincipal and Microsoft Graph Application management capabilities, primarily used for:

- **CI/CD Pipeline Configuration**: Automatically create and configure service principals for GitHub Actions and Azure DevOps
- **Application Identity Management**: Manage Microsoft Entra ID (Azure AD) applications and service principals
- **Federated Credentials**: Support OpenID Connect (OIDC) federated authentication
- **Password Credential Management**: Manage application client secrets
- **Role Assignment**: Automatically configure Azure RBAC roles

## Core Components

### 1. Microsoft Graph SDK (`cli/azd/pkg/graphsdk/`)

azd implements a custom Microsoft Graph SDK client, including:

#### Core Files:
- `graph_client.go` - Graph API client entry point
- `service_principal_request_builders.go` - ServicePrincipal API builders
- `service_principal_models.go` - ServicePrincipal data models
- `application_request_builders.go` - Application API builders
- `applications_models.go` - Application data models
- `fic_request_builders.go` - Federated Identity Credential API builders
- `fic_models.go` - Federated Identity Credential data models

### 2. EntraId Service Layer (`cli/azd/pkg/entraid/`)

Provides high-level business logic encapsulation:

- `entraid.go` - EntraId service implementation
- `entraid_test.go` - Unit tests

## Microsoft Graph SDK Integration

### GraphClient Structure

```go
type GraphClient struct {
    pipeline runtime.Pipeline
    host     string
}
```

### Main Methods

```go
// Applications
func (c *GraphClient) Applications() *ApplicationListRequestBuilder
func (c *GraphClient) ApplicationById(id string) *ApplicationItemRequestBuilder

// ServicePrincipals
func (c *GraphClient) ServicePrincipals() *ServicePrincipalListRequestBuilder
func (c *GraphClient) ServicePrincipalById(id string) *ServicePrincipalItemRequestBuilder

// Me (current user)
func (c *GraphClient) Me() *MeItemRequestBuilder
```

## ServicePrincipal Features

### 1. Data Model

```go
type ServicePrincipal struct {
    Id                     *string  // Service principal ID
    DisplayName            string   // Display name
    AppId                  string   // Application (client) ID
    AppOwnerOrganizationId *string  // Tenant ID
    AppDisplayName         *string  // Application display name
    Description            *string  // Description
    Type                   *string  // Service principal type
}
```

### 2. Supported Operations

#### List Service Principals
```go
servicePrincipals, err := client.ServicePrincipals().Get(ctx)
```

**Functionality**: Get all service principals the current user has access to.

#### Get Specific Service Principal
```go
servicePrincipal, err := client.ServicePrincipalById(id).Get(ctx)
```

**Functionality**: Get details of a specific service principal by ID.

#### Create Service Principal
```go
newSpn := &graphsdk.ServicePrincipal{
    AppId:       "application-app-id",
    DisplayName: "My Service Principal",
}
servicePrincipal, err := client.ServicePrincipals().Post(ctx, newSpn)
```

**Functionality**: Create a new service principal based on an existing application.

#### Delete Service Principal
```go
err := client.ServicePrincipalById(id).Delete(ctx)
```

**Functionality**: Delete a specified service principal.

#### Filter Queries
```go
servicePrincipals, err := client.ServicePrincipals().
    Filter("displayName eq 'MyApp'").
    Get(ctx)
```

**Functionality**: Query service principals using OData filters.

## Application Features

### 1. Data Model

```go
type Application struct {
    Id                         *string                          // Application object ID
    AppId                      *string                          // Application (client) ID
    DisplayName                string                           // Display name
    Description                *string                          // Description
    PasswordCredentials        []*ApplicationPasswordCredential // Password credentials list
    ServiceManagementReference *string                          // Service management reference (Service Tree ID)
}
```

### 2. Supported Operations

#### List Applications
```go
applications, err := client.Applications().Get(ctx)
```

**Functionality**: Get all applications the current user has access to.

#### Get Specific Application
```go
// Get by object ID
application, err := client.ApplicationById(objectId).Get(ctx)

// Get by application ID (AppId/ClientId)
application, err := client.ApplicationById(appId).GetByAppId(ctx)
```

**Functionality**: Get application details by object ID or application ID.

#### Create Application
```go
newApp := &graphsdk.Application{
    DisplayName:                "My Application",
    Description:                to.Ptr("Created by Azure Developer CLI"),
    ServiceManagementReference: to.Ptr("uuid-service-tree-id"),
}
application, err := client.Applications().Post(ctx, newApp)
```

**Functionality**: Create a new Microsoft Entra ID application. Supports setting Service Tree ID for enterprise compliance requirements.

#### Delete Application
```go
err := client.ApplicationById(id).Delete(ctx)
```

**Functionality**: Delete a specified application.

#### Filter Queries
```go
applications, err := client.Applications().
    Filter("startswith(displayName, 'MyApp')").
    Get(ctx)
```

**Functionality**: Query applications using OData filters.

### 3. Password Credential Management

#### Password Credential Data Model
```go
type ApplicationPasswordCredential struct {
    KeyId               *string    // Credential ID
    CustomKeyIdentifier *string    // Custom identifier
    DisplayName         *string    // Display name
    StartDateTime       *time.Time // Start time
    EndDateTime         *time.Time // End time (expiry)
    SecretText          *string    // Secret text (only returned on creation)
    Hint                *string    // Secret hint
}
```

#### Add Password Credential
```go
credential, err := client.ApplicationById(appId).AddPassword(ctx)
```

**Functionality**:
- Create a new client secret for the application
- Automatically sets display name to "Azure Developer CLI"
- Default validity period of 180 days
- Returns credential object with secret text (only available at creation)

#### Remove Password Credential
```go
err := client.ApplicationById(appId).RemovePassword(ctx, keyId)
```

**Functionality**: Delete a specified password credential from the application.

## Federated Identity Credentials Features

Federated identity credentials support passwordless authentication using OpenID Connect (OIDC), particularly suitable for GitHub Actions and other CI/CD systems.

### 1. Data Model

```go
type FederatedIdentityCredential struct {
    Id          *string   // Credential ID
    Name        string    // Credential name
    Issuer      string    // Token issuer URL
    Subject     string    // Subject identifier
    Description *string   // Description
    Audiences   []string  // Audience list
}
```

### 2. Supported Operations

#### List Federated Identity Credentials
```go
credentials, err := client.ApplicationById(appId).
    FederatedIdentityCredentials().
    Get(ctx)
```

**Functionality**: Get all federated identity credentials for an application.

#### Get Specific Federated Identity Credential
```go
credential, err := client.ApplicationById(appId).
    FederatedIdentityCredentialById(credId).
    Get(ctx)
```

#### Create Federated Identity Credential
```go
newCred := &graphsdk.FederatedIdentityCredential{
    Name:        "github-actions-credential",
    Issuer:      "https://token.actions.githubusercontent.com",
    Subject:     "repo:owner/repo:ref:refs/heads/main",
    Audiences:   []string{"api://AzureADTokenExchange"},
    Description: to.Ptr("Credential for GitHub Actions"),
}
credential, err := client.ApplicationById(appId).
    FederatedIdentityCredentials().
    Post(ctx, newCred)
```

**Functionality**:
- Create OIDC federated credentials for GitHub Actions, Azure DevOps, etc.
- No need to manage client secrets
- Support fine-grained access control based on repository, branch, environment, etc.

#### Update Federated Identity Credential
```go
updatedCred := &graphsdk.FederatedIdentityCredential{
    Subject: "repo:owner/repo:ref:refs/heads/develop",
}
err := client.ApplicationById(appId).
    FederatedIdentityCredentialById(credId).
    Update(ctx, updatedCred)
```

#### Delete Federated Identity Credential
```go
err := client.ApplicationById(appId).
    FederatedIdentityCredentialById(credId).
    Delete(ctx)
```

### 3. Common Configurations

#### GitHub Actions
```go
FederatedIdentityCredential{
    Issuer:    "https://token.actions.githubusercontent.com",
    Subject:   "repo:org/repo:environment:production", // Specific environment
    Audiences: []string{"api://AzureADTokenExchange"},
}
```

#### Azure DevOps
```go
FederatedIdentityCredential{
    Issuer:    "https://vstoken.dev.azure.com/{tenant-id}",
    Subject:   "sc://org/project/service-connection-name",
    Audiences: []string{"api://AzureADTokenExchange"},
}
```

## EntraId Service Layer

`EntraIdService` provides high-level business logic encapsulation, simplifying common operations.

### 1. Interface Definition

```go
type EntraIdService interface {
    // Get service principal
    GetServicePrincipal(ctx context.Context, subscriptionId string, appIdOrName string) 
        (*graphsdk.ServicePrincipal, error)
    
    // Create or update service principal
    CreateOrUpdateServicePrincipal(ctx context.Context, subscriptionId string, 
        appIdOrName string, options CreateOrUpdateServicePrincipalOptions) 
        (*graphsdk.ServicePrincipal, error)
    
    // Reset password credentials
    ResetPasswordCredentials(ctx context.Context, subscriptionId string, appId string) 
        (*AzureCredentials, error)
    
    // Apply federated credentials
    ApplyFederatedCredentials(ctx context.Context, subscriptionId string, 
        clientId string, federatedCredentials []*graphsdk.FederatedIdentityCredential) 
        ([]*graphsdk.FederatedIdentityCredential, error)
    
    // Create RBAC role assignment
    CreateRbac(ctx context.Context, subscriptionId string, 
        scope, roleId, principalId string) error
    
    // Ensure role assignments
    EnsureRoleAssignments(ctx context.Context, subscriptionId string, 
        roleNames []string, servicePrincipal *graphsdk.ServicePrincipal, 
        options *EnsureRoleAssignmentsOptions) error
}
```

### 2. Key Features

#### Create or Update Service Principal

```go
options := entraid.CreateOrUpdateServicePrincipalOptions{
    RolesToAssign:              []string{"Contributor", "User Access Administrator"},
    Description:                to.Ptr("Created by Azure Developer CLI"),
    ServiceManagementReference: to.Ptr("service-tree-uuid"),
}

servicePrincipal, err := entraIdService.CreateOrUpdateServicePrincipal(
    ctx,
    subscriptionId,
    "my-application",
    options,
)
```

**Functionality**:
- Automatically find or create application
- Automatically create associated service principal
- Automatically assign Azure RBAC roles
- Support Service Tree ID compliance requirements
- Handle Azure AD propagation delay with retry logic

**Special Error Handling**:
- `ServiceTreeNullValueError` - When enterprise policy requires Service Tree ID
- `ServiceTreeInvalidError` - When provided Service Tree ID is invalid

#### Reset Password Credentials

```go
credentials, err := entraIdService.ResetPasswordCredentials(
    ctx,
    subscriptionId,
    "client-id",
)
```

**Functionality**:
- Delete all existing password credentials
- Create new password credential (180-day validity)
- Return complete Azure credentials information

**Returned Credentials Structure**:
```go
type AzureCredentials struct {
    ClientId       string  // Client ID
    ClientSecret   string  // Client secret
    SubscriptionId string  // Subscription ID
    TenantId       string  // Tenant ID
}
```

#### Apply Federated Credentials

```go
federatedCreds := []*graphsdk.FederatedIdentityCredential{
    {
        Name:        "github-prod",
        Issuer:      "https://token.actions.githubusercontent.com",
        Subject:     "repo:org/repo:environment:production",
        Audiences:   []string{"api://AzureADTokenExchange"},
    },
}

createdCreds, err := entraIdService.ApplyFederatedCredentials(
    ctx,
    subscriptionId,
    clientId,
    federatedCreds,
)
```

**Functionality**:
- Intelligently detect existing credentials to avoid duplication
- Support batch creation of multiple federated credentials
- Automatically handle credential naming conflicts

#### Role Assignment Management

```go
err := entraIdService.EnsureRoleAssignments(
    ctx,
    subscriptionId,
    []string{"Contributor", "User Access Administrator"},
    servicePrincipal,
    &entraid.EnsureRoleAssignmentsOptions{
        Scope: to.Ptr("/subscriptions/sub-id/resourceGroups/rg-name"),
    },
)
```

**Functionality**:
- Support subscription-level and resource group-level role assignments
- Automatic retry mechanism for Azure AD propagation delay
- Idempotent operations, won't duplicate same role assignments
- Automatic permission detection with suggestions

## Command-Line Integration

### 1. `azd pipeline config` Command

Used to configure service principals and credentials for CI/CD pipelines.

#### Main Parameters:

```bash
azd pipeline config \
  --principal-name "my-app-pipeline" \
  --principal-role "Contributor" \
  --principal-role "User Access Administrator" \
  --auth-type "federated" \
  --provider "github" \
  --applicationServiceManagementReference "uuid"
```

**Parameter Description**:
- `--principal-id` - Use existing service principal's client ID
- `--principal-name` - Create or find service principal with specified name
- `--principal-role` - Azure role(s) to assign to service principal (can be used multiple times)
- `--auth-type` - Authentication type: `federated` (OIDC) or `client-credentials` (secret)
- `--provider` - CI/CD provider: `github` or `azdo`
- `--applicationServiceManagementReference` / `-m` - Service Tree ID

#### Supported Authentication Types:

**1. Federated (OIDC) - Recommended**
```bash
azd pipeline config --auth-type federated --provider github
```
- No secret management required
- Higher security
- Automatically creates federated identity credentials
- Supports GitHub Actions and Azure DevOps

**2. Client Credentials (Secret)**
```bash
azd pipeline config --auth-type client-credentials --provider github
```
- Creates client secret
- Requires secure storage in CI/CD
- 180-day validity period

#### Workflow:

1. **Create or Find Application**
   - Find existing app by `--principal-name` or `--principal-id`
   - Create new app if not exists
   - Handle Service Tree ID requirements

2. **Create or Update Service Principal**
   - Create service principal based on application
   - Reuse if already exists

3. **Configure Role Assignments**
   - Default assignment: "Contributor" and "User Access Administrator" roles
   - Customizable via `--principal-role`

4. **Configure Credentials**
   - **Federated**: Create OIDC federated credentials
   - **Client Credentials**: Reset and create new secret

5. **Save to Environment**
   - Save client ID to `.azure/<environment>/.env`
   - Configure secret or federated credentials in CI/CD system

### 2. `azd auth login` Command

Supports service principal authentication.

#### Service Principal Login Methods:

**1. Using Client Secret**
```bash
azd auth login \
  --client-id "app-id" \
  --client-secret "secret" \
  --tenant-id "tenant-id"
```

**2. Using Client Certificate**
```bash
azd auth login \
  --client-id "app-id" \
  --client-certificate "/path/to/cert.pem" \
  --tenant-id "tenant-id"
```

**3. Using Federated Credentials (OIDC)**
```bash
azd auth login \
  --client-id "app-id" \
  --federated-credential-provider "github" \
  --tenant-id "tenant-id"
```

**Supported Federated Credential Providers**:
- `github` - GitHub Actions
- `azure-pipelines` - Azure DevOps

## Usage Scenarios

### Scenario 1: Configure OIDC Federated Authentication for GitHub Actions

**Command**:
```bash
azd pipeline config --auth-type federated --provider github
```

**Automated Actions**:
1. Create Microsoft Entra ID application (if not exists)
2. Create associated service principal
3. Assign Contributor and User Access Administrator roles to Azure subscription
4. Create federated identity credential configured to trust GitHub Actions
5. Configure necessary secrets in GitHub repository

**Generated Federated Credential Configuration**:
```go
FederatedIdentityCredential{
    Name:      "github-{org}-{repo}-{branch}",
    Issuer:    "https://token.actions.githubusercontent.com",
    Subject:   "repo:org/repo:ref:refs/heads/main",
    Audiences: []string{"api://AzureADTokenExchange"},
}
```

### Scenario 2: Configure Secret Authentication for Azure DevOps

**Command**:
```bash
azd pipeline config --auth-type client-credentials --provider azdo
```

**Automated Actions**:
1. Create application and service principal
2. Assign necessary Azure roles
3. Generate new client secret (180-day validity)
4. Create service connection in Azure DevOps

### Scenario 3: Use Service Tree ID for Compliance

**Command**:
```bash
azd pipeline config \
  --principal-name "my-app" \
  --applicationServiceManagementReference "12345678-1234-1234-1234-123456789abc"
```

**Functionality**:
- Attach Service Tree ID when creating application
- Meet enterprise Azure AD policy requirements
- Automatically handle `ServiceTreeNullValueProvided` error with prompts

### Scenario 4: Custom Role Assignments

**Command**:
```bash
azd pipeline config \
  --principal-name "my-app" \
  --principal-role "Contributor" \
  --principal-role "Key Vault Secrets User" \
  --principal-role "Storage Blob Data Contributor"
```

**Functionality**:
- Assign multiple custom Azure roles
- Support built-in and custom roles
- Automatically handle role assignment propagation delays

### Scenario 5: Using Service Principal in CI/CD

**GitHub Actions Workflow Example**:
```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write
      contents: read
    steps:
      - uses: azure/login@v1
        with:
          client-id: ${{ secrets.AZURE_CLIENT_ID }}
          tenant-id: ${{ secrets.AZURE_TENANT_ID }}
          subscription-id: ${{ secrets.AZURE_SUBSCRIPTION_ID }}
      
      - name: Deploy with azd
        run: azd deploy
```

**Azure DevOps Pipeline Example**:
```yaml
- task: AzureCLI@2
  inputs:
    azureSubscription: 'my-service-connection'
    scriptType: 'bash'
    scriptLocation: 'inlineScript'
    inlineScript: |
      azd deploy
```

### Scenario 6: Rotate Service Principal Credentials

**Using EntraIdService**:
```go
// Reset password credentials (delete old, create new)
newCredentials, err := entraIdService.ResetPasswordCredentials(
    ctx,
    subscriptionId,
    clientId,
)

// Update secret in CI/CD system
// newCredentials.ClientSecret contains the new secret
```

### Scenario 7: Configure Federated Credentials for Multiple Environments

```go
environments := []string{"dev", "staging", "production"}
credentials := make([]*graphsdk.FederatedIdentityCredential, 0)

for _, env := range environments {
    credentials = append(credentials, &graphsdk.FederatedIdentityCredential{
        Name:      fmt.Sprintf("github-%s", env),
        Issuer:    "https://token.actions.githubusercontent.com",
        Subject:   fmt.Sprintf("repo:org/repo:environment:%s", env),
        Audiences: []string{"api://AzureADTokenExchange"},
    })
}

createdCreds, err := entraIdService.ApplyFederatedCredentials(
    ctx,
    subscriptionId,
    clientId,
    credentials,
)
```

## Security Best Practices

### 1. Prioritize OIDC Federated Authentication
- No need to manage long-term credentials
- Automatic token rotation
- More fine-grained access control

### 2. Principle of Least Privilege
- Assign only required Azure roles
- Use resource group-level role assignments (not subscription-level)
- Regularly review service principal permissions

### 3. Service Tree ID Usage
- Always use Service Tree ID in enterprise environments
- Facilitates tracking and auditing applications

### 4. Credential Rotation
- Regularly rotate client secrets (if used)
- azd defaults to 180-day validity
- Use `ResetPasswordCredentials` to simplify rotation process

### 5. Audit and Monitoring
- Monitor service principal usage
- Regularly review and clean up unused applications and service principals
- Use Azure AD audit logs

## Technical Implementation Details

### 1. Azure AD Propagation Delay Handling

azd implements automatic retry mechanism to handle Azure AD propagation delays:

```go
// Retry up to 10 times, 5-second intervals
retry.Do(ctx, retry.WithMaxRetries(10, retry.NewConstant(time.Second*5)), 
    func(ctx context.Context) error {
        _, err = roleAssignmentsClient.Create(ctx, scope, roleAssignmentId, parameters, nil)
        // ...
        return retry.RetryableError(err)
    })
```

### 2. Idempotent Operations

All operations are designed to be idempotent:
- Repeated application creation returns existing application
- Repeated role assignments are ignored (409 conflict)
- Federated credentials deduplicated by Subject

### 3. Error Handling

Special error types:
- `ServiceTreeNullValueError` - Missing required Service Tree ID
- `ServiceTreeInvalidError` - Invalid Service Tree ID format
- `ErrorWithSuggestion` - Provides actionable suggestions for permission issues

### 4. Client Caching

GraphClient instances are cached by subscription ID:
```go
clientCache map[string]*graphsdk.GraphClient
```

## Related Code Files

### Microsoft Graph SDK
- `cli/azd/pkg/graphsdk/graph_client.go`
- `cli/azd/pkg/graphsdk/service_principal_request_builders.go`
- `cli/azd/pkg/graphsdk/service_principal_models.go`
- `cli/azd/pkg/graphsdk/application_request_builders.go`
- `cli/azd/pkg/graphsdk/applications_models.go`
- `cli/azd/pkg/graphsdk/fic_request_builders.go`
- `cli/azd/pkg/graphsdk/fic_models.go`

### EntraId Service
- `cli/azd/pkg/entraid/entraid.go`
- `cli/azd/pkg/entraid/entraid_test.go`

### Command Integration
- `cli/azd/cmd/pipeline.go` - Pipeline command
- `cli/azd/cmd/auth_login.go` - Authentication command
- `cli/azd/pkg/pipeline/pipeline_manager.go` - Pipeline manager

### Tests
- `cli/azd/pkg/graphsdk/*_test.go` - Graph SDK unit tests
- `cli/azd/pkg/entraid/entraid_test.go` - EntraId service tests
- `cli/azd/test/mocks/mockgraphsdk/mocks.go` - Mock implementations

## Summary

Azure Developer CLI provides comprehensive ServicePrincipal and Microsoft Graph Application management capabilities, with key features including:

### Core Features
✅ **Complete Microsoft Graph API Integration**
- ServicePrincipal CRUD operations
- Application CRUD operations
- Federated Identity Credentials management
- Password credential management

✅ **Enterprise Features**
- Service Tree ID support
- Automatic role assignment
- Azure AD propagation delay handling
- Detailed error handling and suggestions

✅ **CI/CD Integration**
- GitHub Actions support (OIDC and secret)
- Azure DevOps support (OIDC and secret)
- Automated pipeline configuration

✅ **Security**
- OIDC federated authentication (recommended)
- Automatic credential rotation
- Least privilege principle support

### Usage Recommendations
1. **New Projects**: Use `azd pipeline config --auth-type federated` to configure OIDC
2. **Existing Apps**: Use `--principal-id` to reuse existing service principals
3. **Enterprise Environments**: Always use `--applicationServiceManagementReference` to set Service Tree ID
4. **Multi-Environment**: Create separate federated credentials for each environment

These features make azd a powerful tool for managing Azure application identities and CI/CD pipeline configurations.
