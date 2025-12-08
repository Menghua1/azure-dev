// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package graphsdk_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/azure/azure-dev/cli/azd/pkg/graphsdk"
	"github.com/stretchr/testify/require"
)

func TestHasEnvTag(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		envName  string
		expected bool
	}{
		{
			name:     "Tag exists",
			tags:     []string{"azd-env-name:test-env", "other:tag"},
			envName:  "test-env",
			expected: true,
		},
		{
			name:     "Tag does not exist",
			tags:     []string{"other:tag", "another:tag"},
			envName:  "test-env",
			expected: false,
		},
		{
			name:     "Empty tags",
			tags:     []string{},
			envName:  "test-env",
			expected: false,
		},
		{
			name:     "Different environment name",
			tags:     []string{"azd-env-name:prod-env"},
			envName:  "test-env",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := graphsdk.HasEnvTag(tt.tags, tt.envName)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestApplicationTagsModel(t *testing.T) {
	// Test that the Application model includes the Tags field
	app := graphsdk.Application{
		Id:          to.Ptr("test-id"),
		AppId:       to.Ptr("app-id"),
		DisplayName: "Test Application",
		Tags:        []string{"azd-env-name:test-env", "azd-service-name:api"},
	}

	require.NotNil(t, app.Tags)
	require.Len(t, app.Tags, 2)
	require.Contains(t, app.Tags, "azd-env-name:test-env")
}

func TestServicePrincipalTagsModel(t *testing.T) {
	// Test that the ServicePrincipal model includes the Tags field
	sp := graphsdk.ServicePrincipal{
		Id:          to.Ptr("test-id"),
		AppId:       "app-id",
		DisplayName: "Test Service Principal",
		Tags:        []string{"azd-env-name:test-env"},
	}

	require.NotNil(t, sp.Tags)
	require.Len(t, sp.Tags, 1)
	require.Contains(t, sp.Tags, "azd-env-name:test-env")
}
