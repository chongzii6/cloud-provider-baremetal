package backup

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
)

// ResourceVaultConfigsClient is the open API 2.0 Specs for Azure RecoveryServices Backup service
type ResourceVaultConfigsClient struct {
	BaseClient
}

// NewResourceVaultConfigsClient creates an instance of the ResourceVaultConfigsClient client.
func NewResourceVaultConfigsClient(subscriptionID string) ResourceVaultConfigsClient {
	return NewResourceVaultConfigsClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewResourceVaultConfigsClientWithBaseURI creates an instance of the ResourceVaultConfigsClient client.
func NewResourceVaultConfigsClientWithBaseURI(baseURI string, subscriptionID string) ResourceVaultConfigsClient {
	return ResourceVaultConfigsClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get fetches resource vault config.
//
// vaultName is the name of the recovery services vault. resourceGroupName is the name of the resource group where
// the recovery services vault is present.
func (client ResourceVaultConfigsClient) Get(ctx context.Context, vaultName string, resourceGroupName string) (result ResourceVaultConfigResource, err error) {
	req, err := client.GetPreparer(ctx, vaultName, resourceGroupName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ResourceVaultConfigsClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "backup.ResourceVaultConfigsClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ResourceVaultConfigsClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client ResourceVaultConfigsClient) GetPreparer(ctx context.Context, vaultName string, resourceGroupName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2016-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupconfig/vaultconfig", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client ResourceVaultConfigsClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client ResourceVaultConfigsClient) GetResponder(resp *http.Response) (result ResourceVaultConfigResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// Update updates vault security config.
//
// vaultName is the name of the recovery services vault. resourceGroupName is the name of the resource group where
// the recovery services vault is present. parameters is resource config request
func (client ResourceVaultConfigsClient) Update(ctx context.Context, vaultName string, resourceGroupName string, parameters ResourceVaultConfigResource) (result ResourceVaultConfigResource, err error) {
	req, err := client.UpdatePreparer(ctx, vaultName, resourceGroupName, parameters)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ResourceVaultConfigsClient", "Update", nil, "Failure preparing request")
		return
	}

	resp, err := client.UpdateSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "backup.ResourceVaultConfigsClient", "Update", resp, "Failure sending request")
		return
	}

	result, err = client.UpdateResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "backup.ResourceVaultConfigsClient", "Update", resp, "Failure responding to request")
	}

	return
}

// UpdatePreparer prepares the Update request.
func (client ResourceVaultConfigsClient) UpdatePreparer(ctx context.Context, vaultName string, resourceGroupName string, parameters ResourceVaultConfigResource) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"vaultName":         autorest.Encode("path", vaultName),
	}

	const APIVersion = "2016-12-01"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsJSON(),
		autorest.AsPatch(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/Subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupconfig/vaultconfig", pathParameters),
		autorest.WithJSON(parameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// UpdateSender sends the Update request. The method will close the
// http.Response Body if it receives an error.
func (client ResourceVaultConfigsClient) UpdateSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// UpdateResponder handles the response to the Update request. The method always
// closes the http.Response Body.
func (client ResourceVaultConfigsClient) UpdateResponder(resp *http.Response) (result ResourceVaultConfigResource, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}
