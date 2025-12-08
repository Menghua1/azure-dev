// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package graphsdk

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// DirectoryDeletedItemsRequestBuilder provides methods to interact with deleted items in the directory
type DirectoryDeletedItemsRequestBuilder struct {
	client *GraphClient
}

// NewDirectoryDeletedItemsRequestBuilder creates a new DirectoryDeletedItemsRequestBuilder
func NewDirectoryDeletedItemsRequestBuilder(client *GraphClient) *DirectoryDeletedItemsRequestBuilder {
	return &DirectoryDeletedItemsRequestBuilder{
		client: client,
	}
}

// Get retrieves a deleted item from the directory
func (b *DirectoryDeletedItemsRequestBuilder) Get(ctx context.Context, id string) error {
	req, err := runtime.NewRequest(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/directory/deleteditems/%s", b.client.host, id),
	)
	if err != nil {
		return fmt.Errorf("failed creating request: %w", err)
	}

	res, err := b.client.pipeline.Do(req)
	if err != nil {
		return err
	}

	if !runtime.HasStatusCode(res, http.StatusOK) {
		return runtime.NewResponseError(res)
	}

	return nil
}

// Delete permanently deletes an item from the directory's deleted items
func (b *DirectoryDeletedItemsRequestBuilder) Delete(ctx context.Context, id string) error {
	req, err := runtime.NewRequest(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/directory/deleteditems/%s", b.client.host, id),
	)
	if err != nil {
		return fmt.Errorf("failed creating request: %w", err)
	}

	res, err := b.client.pipeline.Do(req)
	if err != nil {
		return err
	}

	if !runtime.HasStatusCode(res, http.StatusNoContent) {
		return runtime.NewResponseError(res)
	}

	return nil
}
