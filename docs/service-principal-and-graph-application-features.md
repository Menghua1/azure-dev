# Azure Developer CLI (azd) - ServicePrincipal 和 Microsoft Graph Application 功能文档

本文档详细介绍了 Azure Developer CLI (azd) 中关于 ServicePrincipal（服务主体）和 Microsoft Graph Application（Microsoft Graph 应用程序）的现有功能。

## 目录

1. [概述](#概述)
2. [核心组件](#核心组件)
3. [Microsoft Graph SDK 集成](#microsoft-graph-sdk-集成)
4. [ServicePrincipal 功能](#serviceprincipal-功能)
5. [Application 功能](#application-功能)
6. [Federated Identity Credentials 功能](#federated-identity-credentials-功能)
7. [EntraId 服务层](#entraid-服务层)
8. [命令行集成](#命令行集成)
9. [使用场景](#使用场景)

## 概述

azd 提供了完整的 ServicePrincipal 和 Microsoft Graph Application 管理功能，主要用于：

- **CI/CD 流水线配置**：自动创建和配置服务主体用于 GitHub Actions 和 Azure DevOps
- **应用程序身份管理**：管理 Microsoft Entra ID (Azure AD) 应用程序和服务主体
- **联合身份凭证**：支持 OpenID Connect (OIDC) 联合身份验证
- **密码凭证管理**：管理应用程序的客户端密钥
- **角色分配**：自动配置 Azure RBAC 角色

## 核心组件

### 1. Microsoft Graph SDK (`cli/azd/pkg/graphsdk/`)

azd 实现了自定义的 Microsoft Graph SDK 客户端，包括：

#### 核心文件：
- `graph_client.go` - Graph API 客户端入口
- `service_principal_request_builders.go` - ServicePrincipal API 构建器
- `service_principal_models.go` - ServicePrincipal 数据模型
- `application_request_builders.go` - Application API 构建器
- `applications_models.go` - Application 数据模型
- `fic_request_builders.go` - Federated Identity Credential API 构建器
- `fic_models.go` - Federated Identity Credential 数据模型

### 2. EntraId 服务层 (`cli/azd/pkg/entraid/`)

提供高级业务逻辑封装：

- `entraid.go` - EntraId 服务实现
- `entraid_test.go` - 单元测试

## Microsoft Graph SDK 集成

### GraphClient 结构

```go
type GraphClient struct {
    pipeline runtime.Pipeline
    host     string
}
```

### 主要方法

```go
// Applications
func (c *GraphClient) Applications() *ApplicationListRequestBuilder
func (c *GraphClient) ApplicationById(id string) *ApplicationItemRequestBuilder

// ServicePrincipals
func (c *GraphClient) ServicePrincipals() *ServicePrincipalListRequestBuilder
func (c *GraphClient) ServicePrincipalById(id string) *ServicePrincipalItemRequestBuilder

// Me (当前用户)
func (c *GraphClient) Me() *MeItemRequestBuilder
```

## ServicePrincipal 功能

### 1. 数据模型

```go
type ServicePrincipal struct {
    Id                     *string  // 服务主体 ID
    DisplayName            string   // 显示名称
    AppId                  string   // 应用程序（客户端）ID
    AppOwnerOrganizationId *string  // 租户 ID
    AppDisplayName         *string  // 应用程序显示名称
    Description            *string  // 描述
    Type                   *string  // 服务主体类型
}
```

### 2. 支持的操作

#### 列出服务主体
```go
servicePrincipals, err := client.ServicePrincipals().Get(ctx)
```

**功能**：获取当前用户有权访问的所有服务主体列表。

#### 获取特定服务主体
```go
servicePrincipal, err := client.ServicePrincipalById(id).Get(ctx)
```

**功能**：根据 ID 获取特定的服务主体详细信息。

#### 创建服务主体
```go
newSpn := &graphsdk.ServicePrincipal{
    AppId:       "application-app-id",
    DisplayName: "My Service Principal",
}
servicePrincipal, err := client.ServicePrincipals().Post(ctx, newSpn)
```

**功能**：基于现有应用程序创建新的服务主体。

#### 删除服务主体
```go
err := client.ServicePrincipalById(id).Delete(ctx)
```

**功能**：删除指定的服务主体。

#### 筛选查询
```go
servicePrincipals, err := client.ServicePrincipals().
    Filter("displayName eq 'MyApp'").
    Get(ctx)
```

**功能**：使用 OData 过滤器查询服务主体。

## Application 功能

### 1. 数据模型

```go
type Application struct {
    Id                         *string                          // 应用程序对象 ID
    AppId                      *string                          // 应用程序（客户端）ID
    DisplayName                string                           // 显示名称
    Description                *string                          // 描述
    PasswordCredentials        []*ApplicationPasswordCredential // 密码凭证列表
    ServiceManagementReference *string                          // 服务管理引用 (Service Tree ID)
}
```

### 2. 支持的操作

#### 列出应用程序
```go
applications, err := client.Applications().Get(ctx)
```

**功能**：获取当前用户有权访问的所有应用程序。

#### 获取特定应用程序
```go
// 根据对象 ID 获取
application, err := client.ApplicationById(objectId).Get(ctx)

// 根据应用程序 ID (AppId/ClientId) 获取
application, err := client.ApplicationById(appId).GetByAppId(ctx)
```

**功能**：根据对象 ID 或应用程序 ID 获取应用程序详细信息。

#### 创建应用程序
```go
newApp := &graphsdk.Application{
    DisplayName:                "My Application",
    Description:                to.Ptr("Created by Azure Developer CLI"),
    ServiceManagementReference: to.Ptr("uuid-service-tree-id"),
}
application, err := client.Applications().Post(ctx, newApp)
```

**功能**：创建新的 Microsoft Entra ID 应用程序。支持设置 Service Tree ID 以满足企业合规性要求。

#### 删除应用程序
```go
err := client.ApplicationById(id).Delete(ctx)
```

**功能**：删除指定的应用程序。

#### 筛选查询
```go
applications, err := client.Applications().
    Filter("startswith(displayName, 'MyApp')").
    Get(ctx)
```

**功能**：使用 OData 过滤器查询应用程序。

### 3. 密码凭证管理

#### 密码凭证数据模型
```go
type ApplicationPasswordCredential struct {
    KeyId               *string    // 凭证 ID
    CustomKeyIdentifier *string    // 自定义标识符
    DisplayName         *string    // 显示名称
    StartDateTime       *time.Time // 开始时间
    EndDateTime         *time.Time // 结束时间（过期时间）
    SecretText          *string    // 密钥文本（仅在创建时返回）
    Hint                *string    // 密钥提示
}
```

#### 添加密码凭证
```go
credential, err := client.ApplicationById(appId).AddPassword(ctx)
```

**功能**：
- 为应用程序创建新的客户端密钥
- 自动设置显示名称为 "Azure Developer CLI"
- 默认有效期为 180 天
- 返回包含密钥文本的凭证对象（仅此时可获取）

#### 删除密码凭证
```go
err := client.ApplicationById(appId).RemovePassword(ctx, keyId)
```

**功能**：删除应用程序的指定密码凭证。

## Federated Identity Credentials 功能

联合身份凭证支持使用 OpenID Connect (OIDC) 进行无密码身份验证，特别适用于 GitHub Actions 和其他 CI/CD 系统。

### 1. 数据模型

```go
type FederatedIdentityCredential struct {
    Id          *string   // 凭证 ID
    Name        string    // 凭证名称
    Issuer      string    // 令牌颁发者 URL
    Subject     string    // 主题标识符
    Description *string   // 描述
    Audiences   []string  // 受众列表
}
```

### 2. 支持的操作

#### 列出联合身份凭证
```go
credentials, err := client.ApplicationById(appId).
    FederatedIdentityCredentials().
    Get(ctx)
```

**功能**：获取应用程序的所有联合身份凭证。

#### 获取特定联合身份凭证
```go
credential, err := client.ApplicationById(appId).
    FederatedIdentityCredentialById(credId).
    Get(ctx)
```

#### 创建联合身份凭证
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

**功能**：
- 为 GitHub Actions、Azure DevOps 等创建 OIDC 联合凭证
- 无需管理客户端密钥
- 支持基于存储库、分支、环境等的细粒度访问控制

#### 更新联合身份凭证
```go
updatedCred := &graphsdk.FederatedIdentityCredential{
    Subject: "repo:owner/repo:ref:refs/heads/develop",
}
err := client.ApplicationById(appId).
    FederatedIdentityCredentialById(credId).
    Update(ctx, updatedCred)
```

#### 删除联合身份凭证
```go
err := client.ApplicationById(appId).
    FederatedIdentityCredentialById(credId).
    Delete(ctx)
```

### 3. 常见配置

#### GitHub Actions
```go
FederatedIdentityCredential{
    Issuer:    "https://token.actions.githubusercontent.com",
    Subject:   "repo:org/repo:environment:production", // 特定环境
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

## EntraId 服务层

`EntraIdService` 提供了高级业务逻辑封装，简化了常见操作。

### 1. 接口定义

```go
type EntraIdService interface {
    // 获取服务主体
    GetServicePrincipal(ctx context.Context, subscriptionId string, appIdOrName string) 
        (*graphsdk.ServicePrincipal, error)
    
    // 创建或更新服务主体
    CreateOrUpdateServicePrincipal(ctx context.Context, subscriptionId string, 
        appIdOrName string, options CreateOrUpdateServicePrincipalOptions) 
        (*graphsdk.ServicePrincipal, error)
    
    // 重置密码凭证
    ResetPasswordCredentials(ctx context.Context, subscriptionId string, appId string) 
        (*AzureCredentials, error)
    
    // 应用联合身份凭证
    ApplyFederatedCredentials(ctx context.Context, subscriptionId string, 
        clientId string, federatedCredentials []*graphsdk.FederatedIdentityCredential) 
        ([]*graphsdk.FederatedIdentityCredential, error)
    
    // 创建 RBAC 角色分配
    CreateRbac(ctx context.Context, subscriptionId string, 
        scope, roleId, principalId string) error
    
    // 确保角色分配
    EnsureRoleAssignments(ctx context.Context, subscriptionId string, 
        roleNames []string, servicePrincipal *graphsdk.ServicePrincipal, 
        options *EnsureRoleAssignmentsOptions) error
}
```

### 2. 主要功能

#### 创建或更新服务主体

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

**功能**：
- 自动查找或创建应用程序
- 自动创建关联的服务主体
- 自动分配 Azure RBAC 角色
- 支持 Service Tree ID 合规性要求
- 处理 Azure AD 传播延迟的重试逻辑

**特殊错误处理**：
- `ServiceTreeNullValueError` - 当企业策略要求 Service Tree ID 时
- `ServiceTreeInvalidError` - 当提供的 Service Tree ID 无效时

#### 重置密码凭证

```go
credentials, err := entraIdService.ResetPasswordCredentials(
    ctx,
    subscriptionId,
    "client-id",
)
```

**功能**：
- 删除所有现有密码凭证
- 创建新的密码凭证（180天有效期）
- 返回完整的 Azure 凭证信息

**返回的凭证结构**：
```go
type AzureCredentials struct {
    ClientId       string  // 客户端 ID
    ClientSecret   string  // 客户端密钥
    SubscriptionId string  // 订阅 ID
    TenantId       string  // 租户 ID
}
```

#### 应用联合身份凭证

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

**功能**：
- 智能检测现有凭证，避免重复创建
- 支持批量创建多个联合凭证
- 自动处理凭证命名冲突

#### 角色分配管理

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

**功能**：
- 支持订阅级别和资源组级别的角色分配
- 自动重试机制处理 Azure AD 传播延迟
- 幂等操作，不会重复分配相同角色
- 自动检测权限不足并提供建议

## 命令行集成

### 1. `azd pipeline config` 命令

用于配置 CI/CD 流水线的服务主体和凭证。

#### 主要参数：

```bash
azd pipeline config \
  --principal-name "my-app-pipeline" \
  --principal-role "Contributor" \
  --principal-role "User Access Administrator" \
  --auth-type "federated" \
  --provider "github" \
  --applicationServiceManagementReference "uuid"
```

**参数说明**：
- `--principal-id` - 使用现有服务主体的客户端 ID
- `--principal-name` - 创建或查找指定名称的服务主体
- `--principal-role` - 分配给服务主体的 Azure 角色（可多次使用）
- `--auth-type` - 认证类型：`federated`（OIDC）或 `client-credentials`（密钥）
- `--provider` - CI/CD 提供商：`github` 或 `azdo`
- `--applicationServiceManagementReference` / `-m` - Service Tree ID

#### 支持的身份验证类型：

**1. Federated (OIDC) - 推荐**
```bash
azd pipeline config --auth-type federated --provider github
```
- 无密钥管理
- 更高安全性
- 自动创建联合身份凭证
- 支持 GitHub Actions 和 Azure DevOps

**2. Client Credentials (密钥)**
```bash
azd pipeline config --auth-type client-credentials --provider github
```
- 创建客户端密钥
- 需要在 CI/CD 中安全存储密钥
- 180 天有效期

#### 工作流程：

1. **创建或查找应用程序**
   - 根据 `--principal-name` 或 `--principal-id` 查找现有应用
   - 如不存在则创建新应用
   - 处理 Service Tree ID 要求

2. **创建或更新服务主体**
   - 基于应用程序创建服务主体
   - 如已存在则复用

3. **配置角色分配**
   - 默认分配 "Contributor" 和 "User Access Administrator" 角色
   - 可通过 `--principal-role` 自定义

4. **配置凭证**
   - **Federated**：创建 OIDC 联合凭证
   - **Client Credentials**：重置并创建新密钥

5. **保存到环境**
   - 将客户端 ID 保存到 `.azure/<environment>/.env`
   - 密钥或联合凭证配置到 CI/CD 系统

### 2. `azd auth login` 命令

支持使用服务主体进行身份验证。

#### 服务主体登录方式：

**1. 使用客户端密钥**
```bash
azd auth login \
  --client-id "app-id" \
  --client-secret "secret" \
  --tenant-id "tenant-id"
```

**2. 使用客户端证书**
```bash
azd auth login \
  --client-id "app-id" \
  --client-certificate "/path/to/cert.pem" \
  --tenant-id "tenant-id"
```

**3. 使用联合凭证（OIDC）**
```bash
azd auth login \
  --client-id "app-id" \
  --federated-credential-provider "github" \
  --tenant-id "tenant-id"
```

**支持的联合凭证提供商**：
- `github` - GitHub Actions
- `azure-pipelines` - Azure DevOps

## 使用场景

### 场景 1：为 GitHub Actions 配置 OIDC 联合身份验证

**命令**：
```bash
azd pipeline config --auth-type federated --provider github
```

**自动执行的操作**：
1. 创建 Microsoft Entra ID 应用程序（如不存在）
2. 创建关联的服务主体
3. 分配 Azure 订阅的 Contributor 和 User Access Administrator 角色
4. 创建联合身份凭证，配置为信任 GitHub Actions
5. 在 GitHub 仓库中配置必要的 secrets

**生成的联合凭证配置**：
```go
FederatedIdentityCredential{
    Name:      "github-{org}-{repo}-{branch}",
    Issuer:    "https://token.actions.githubusercontent.com",
    Subject:   "repo:org/repo:ref:refs/heads/main",
    Audiences: []string{"api://AzureADTokenExchange"},
}
```

### 场景 2：为 Azure DevOps 配置密钥身份验证

**命令**：
```bash
azd pipeline config --auth-type client-credentials --provider azdo
```

**自动执行的操作**：
1. 创建应用程序和服务主体
2. 分配必要的 Azure 角色
3. 生成新的客户端密钥（180天有效期）
4. 在 Azure DevOps 中创建服务连接

### 场景 3：使用 Service Tree ID 满足合规性要求

**命令**：
```bash
azd pipeline config \
  --principal-name "my-app" \
  --applicationServiceManagementReference "12345678-1234-1234-1234-123456789abc"
```

**功能**：
- 在创建应用程序时附加 Service Tree ID
- 满足企业 Azure AD 策略要求
- 自动处理 `ServiceTreeNullValueProvided` 错误并提示输入

### 场景 4：自定义角色分配

**命令**：
```bash
azd pipeline config \
  --principal-name "my-app" \
  --principal-role "Contributor" \
  --principal-role "Key Vault Secrets User" \
  --principal-role "Storage Blob Data Contributor"
```

**功能**：
- 分配多个自定义 Azure 角色
- 支持内置角色和自定义角色
- 自动处理角色分配的传播延迟

### 场景 5：在 CI/CD 中使用服务主体

**GitHub Actions 工作流示例**：
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

**Azure DevOps 管道示例**：
```yaml
- task: AzureCLI@2
  inputs:
    azureSubscription: 'my-service-connection'
    scriptType: 'bash'
    scriptLocation: 'inlineScript'
    inlineScript: |
      azd deploy
```

### 场景 6：轮换服务主体凭证

**使用 EntraIdService**：
```go
// 重置密码凭证（删除旧的，创建新的）
newCredentials, err := entraIdService.ResetPasswordCredentials(
    ctx,
    subscriptionId,
    clientId,
)

// 更新 CI/CD 系统中的密钥
// newCredentials.ClientSecret 包含新密钥
```

### 场景 7：为多个环境配置联合凭证

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

## 安全最佳实践

### 1. 优先使用 OIDC 联合身份验证
- 无需管理长期凭证
- 自动轮换令牌
- 更细粒度的访问控制

### 2. 最小权限原则
- 仅分配所需的 Azure 角色
- 使用资源组级别的角色分配（而非订阅级别）
- 定期审查服务主体权限

### 3. Service Tree ID 使用
- 在企业环境中始终使用 Service Tree ID
- 便于跟踪和审计应用程序

### 4. 凭证轮换
- 定期轮换客户端密钥（如使用）
- azd 默认设置 180 天有效期
- 使用 `ResetPasswordCredentials` 简化轮换流程

### 5. 审计和监控
- 监控服务主体的使用情况
- 定期检查和清理未使用的应用程序和服务主体
- 使用 Azure AD 审计日志

## 技术实现细节

### 1. Azure AD 传播延迟处理

azd 实现了自动重试机制处理 Azure AD 的传播延迟：

```go
// 最多重试 10 次，每次间隔 5 秒
retry.Do(ctx, retry.WithMaxRetries(10, retry.NewConstant(time.Second*5)), 
    func(ctx context.Context) error {
        _, err = roleAssignmentsClient.Create(ctx, scope, roleAssignmentId, parameters, nil)
        // ...
        return retry.RetryableError(err)
    })
```

### 2. 幂等操作

所有操作都设计为幂等：
- 重复创建应用程序会返回现有应用程序
- 重复分配角色会被忽略（409 冲突）
- 联合凭证按 Subject 去重

### 3. 错误处理

特殊错误类型：
- `ServiceTreeNullValueError` - 缺少必需的 Service Tree ID
- `ServiceTreeInvalidError` - Service Tree ID 格式无效
- `ErrorWithSuggestion` - 权限不足时提供可操作建议

### 4. 客户端缓存

GraphClient 实例按订阅 ID 缓存：
```go
clientCache map[string]*graphsdk.GraphClient
```

## 相关代码文件

### Microsoft Graph SDK
- `cli/azd/pkg/graphsdk/graph_client.go`
- `cli/azd/pkg/graphsdk/service_principal_request_builders.go`
- `cli/azd/pkg/graphsdk/service_principal_models.go`
- `cli/azd/pkg/graphsdk/application_request_builders.go`
- `cli/azd/pkg/graphsdk/applications_models.go`
- `cli/azd/pkg/graphsdk/fic_request_builders.go`
- `cli/azd/pkg/graphsdk/fic_models.go`

### EntraId 服务
- `cli/azd/pkg/entraid/entraid.go`
- `cli/azd/pkg/entraid/entraid_test.go`

### 命令集成
- `cli/azd/cmd/pipeline.go` - Pipeline 命令
- `cli/azd/cmd/auth_login.go` - 身份验证命令
- `cli/azd/pkg/pipeline/pipeline_manager.go` - Pipeline 管理器

### 测试
- `cli/azd/pkg/graphsdk/*_test.go` - Graph SDK 单元测试
- `cli/azd/pkg/entraid/entraid_test.go` - EntraId 服务测试
- `cli/azd/test/mocks/mockgraphsdk/mocks.go` - Mock 实现

## 总结

Azure Developer CLI 提供了全面的 ServicePrincipal 和 Microsoft Graph Application 管理功能，主要特性包括：

### 核心功能
✅ **完整的 Microsoft Graph API 集成**
- ServicePrincipal CRUD 操作
- Application CRUD 操作
- Federated Identity Credentials 管理
- 密码凭证管理

✅ **企业级功能**
- Service Tree ID 支持
- 自动角色分配
- Azure AD 传播延迟处理
- 详细的错误处理和建议

✅ **CI/CD 集成**
- GitHub Actions 支持（OIDC 和密钥）
- Azure DevOps 支持（OIDC 和密钥）
- 自动化流水线配置

✅ **安全性**
- OIDC 联合身份验证（推荐）
- 自动凭证轮换
- 最小权限原则支持

### 使用建议
1. **新项目**：使用 `azd pipeline config --auth-type federated` 配置 OIDC
2. **现有应用**：使用 `--principal-id` 复用现有服务主体
3. **企业环境**：始终使用 `--applicationServiceManagementReference` 设置 Service Tree ID
4. **多环境**：为每个环境创建单独的联合凭证

这些功能使 azd 成为管理 Azure 应用程序身份和 CI/CD 流水线配置的强大工具。
