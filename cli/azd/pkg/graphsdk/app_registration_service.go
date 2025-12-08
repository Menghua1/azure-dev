// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package graphsdk

import (
	"context"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/azure/azure-dev/cli/azd/pkg/azure"
)

const (
	// Retry configuration for waiting for service principals to appear in deleted items
	maxDeleteRetries = 6
	deleteRetryDelay = 5 * time.Second
)

// AppRegistrationService provides methods to manage App Registrations and Service Principals
type AppRegistrationService struct {
	client *GraphClient
}

// NewAppRegistrationService creates a new AppRegistrationService
func NewAppRegistrationService(client *GraphClient) *AppRegistrationService {
	return &AppRegistrationService{
		client: client,
	}
}

// FindByEnvName finds all applications and service principals tagged with the specified environment name
func (s *AppRegistrationService) FindByEnvName(
	ctx context.Context,
	envName string,
) ([]Application, []ServicePrincipal, error) {
	// Build the filter expression to find resources with the azd-env-name tag
	tagFilter := fmt.Sprintf("azd-env-name:%s", envName)

	// Query applications with the tag
	applications, err := s.client.Applications().Filter(fmt.Sprintf("tags/any(t:t eq '%s')", tagFilter)).Get(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query applications: %w", err)
	}

	// Query service principals with the tag
	servicePrincipals, err := s.client.ServicePrincipals().Filter(fmt.Sprintf("tags/any(t:t eq '%s')", tagFilter)).Get(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to query service principals: %w", err)
	}

	return applications.Value, servicePrincipals.Value, nil
}

// DeleteServicePrincipal deletes a service principal and permanently removes it from deleted items
func (s *AppRegistrationService) DeleteServicePrincipal(ctx context.Context, sp ServicePrincipal) error {
	if sp.Id == nil {
		return fmt.Errorf("service principal ID is nil")
	}

	spId := *sp.Id
	log.Printf("Deleting service principal %s (AppId: %s)", spId, sp.AppId)

	// Step 1: Delete the service principal (soft delete)
	err := s.client.ServicePrincipalById(spId).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete service principal %s: %w", spId, err)
	}

	// Step 2: Wait for the service principal to appear in deleted items
	// We'll retry for up to 30 seconds
	var deletedItemFound bool
	for i := 0; i < maxDeleteRetries; i++ {
		log.Printf("Verifying service principal %s is in deleted items (attempt %d/%d)...", spId, i+1, maxDeleteRetries)
		err = s.client.DirectoryDeletedItems().Get(ctx, spId)
		if err == nil {
			deletedItemFound = true
			break
		}
		if i < maxDeleteRetries-1 {
			time.Sleep(deleteRetryDelay)
		}
	}

	if !deletedItemFound {
		log.Printf("Warning: service principal %s may not have been moved to deleted items", spId)
		// Continue anyway - the soft delete succeeded
	}

	// Step 3: Permanently delete the service principal
	log.Printf("Permanently deleting service principal %s", spId)
	err = s.client.DirectoryDeletedItems().Delete(ctx, spId)
	if err != nil {
		return fmt.Errorf("failed to permanently delete service principal %s: %w", spId, err)
	}

	return nil
}

// DeleteApplication deletes an application
func (s *AppRegistrationService) DeleteApplication(ctx context.Context, app Application) error {
	if app.Id == nil {
		return fmt.Errorf("application ID is nil")
	}

	appId := *app.Id
	appIdStr := "<unknown>"
	if app.AppId != nil {
		appIdStr = *app.AppId
	}
	log.Printf("Deleting application %s (AppId: %s, DisplayName: %s)", appId, appIdStr, app.DisplayName)

	err := s.client.ApplicationById(appId).Delete(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete application %s: %w", appId, err)
	}

	return nil
}

// DeleteByEnvName deletes all applications and service principals tagged with the specified environment name
func (s *AppRegistrationService) DeleteByEnvName(ctx context.Context, envName string) error {
	applications, servicePrincipals, err := s.FindByEnvName(ctx, envName)
	if err != nil {
		return fmt.Errorf("failed to find resources: %w", err)
	}

	log.Printf("Found %d applications and %d service principals with tag azd-env-name:%s",
		len(applications), len(servicePrincipals), envName)

	// Delete service principals first
	for _, sp := range servicePrincipals {
		if err := s.DeleteServicePrincipal(ctx, sp); err != nil {
			log.Printf("Warning: failed to delete service principal %s: %v", *sp.Id, err)
			// Continue with other deletions
		}
	}

	// Delete applications
	for _, app := range applications {
		if err := s.DeleteApplication(ctx, app); err != nil {
			log.Printf("Warning: failed to delete application %s: %v", *app.Id, err)
			// Continue with other deletions
		}
	}

	return nil
}

// HasEnvTag checks if a resource has the azd-env-name tag with the specified environment name
func HasEnvTag(tags []string, envName string) bool {
	tagValue := fmt.Sprintf("%s:%s", azure.TagKeyAzdEnvName, envName)
	return slices.Contains(tags, tagValue)
}
