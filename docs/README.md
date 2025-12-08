# Azure Developer CLI Documentation

This directory contains comprehensive documentation for Azure Developer CLI (azd) features and functionality.

## Available Documentation

### ServicePrincipal and Microsoft Graph Application Features

Comprehensive guides to the ServicePrincipal and Microsoft Graph Application management capabilities in azd:

- **[service-principal-and-graph-application-features.md](./service-principal-and-graph-application-features.md)** (中文版)
  - 完整介绍 azd 中 ServicePrincipal 和 Microsoft Graph Application 的所有功能
  - 包含详细的 API 示例、使用场景和最佳实践
  
- **[service-principal-and-graph-application-features-en.md](./service-principal-and-graph-application-features-en.md)** (English)
  - Complete overview of ServicePrincipal and Microsoft Graph Application features in azd
  - Includes detailed API examples, usage scenarios, and best practices

## Documentation Contents

These documents cover:

1. **Core Components**
   - Microsoft Graph SDK integration
   - EntraId service layer architecture
   
2. **Features**
   - ServicePrincipal CRUD operations
   - Application management
   - Federated Identity Credentials (OIDC)
   - Password credential management
   - Azure RBAC role assignments

3. **Command-Line Integration**
   - `azd pipeline config` - CI/CD pipeline configuration
   - `azd auth login` - Service principal authentication

4. **Usage Scenarios**
   - GitHub Actions OIDC setup
   - Azure DevOps integration
   - Service Tree ID compliance
   - Multi-environment configurations
   - Credential rotation

5. **Best Practices**
   - Security recommendations
   - Least privilege principles
   - Credential management
   - Audit and monitoring

## Quick Links

### For CI/CD Setup
- See "命令行集成" / "Command-Line Integration" section
- Check "使用场景" / "Usage Scenarios" for platform-specific examples

### For API Integration
- See "Microsoft Graph SDK 集成" / "Microsoft Graph SDK Integration" section
- Check code examples in each feature section

### For Enterprise Compliance
- See "Service Tree ID" sections
- Review security best practices

## Contributing

If you find any issues or have suggestions for improving this documentation, please open an issue or submit a pull request.
