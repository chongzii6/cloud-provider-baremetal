package logic

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

// WorkflowTriggersClient is the REST API for Azure Logic Apps.
type WorkflowTriggersClient struct {
	BaseClient
}

// NewWorkflowTriggersClient creates an instance of the WorkflowTriggersClient client.
func NewWorkflowTriggersClient(subscriptionID string) WorkflowTriggersClient {
	return NewWorkflowTriggersClientWithBaseURI(DefaultBaseURI, subscriptionID)
}

// NewWorkflowTriggersClientWithBaseURI creates an instance of the WorkflowTriggersClient client.
func NewWorkflowTriggersClientWithBaseURI(baseURI string, subscriptionID string) WorkflowTriggersClient {
	return WorkflowTriggersClient{NewWithBaseURI(baseURI, subscriptionID)}
}

// Get gets a workflow trigger.
//
// resourceGroupName is the resource group name. workflowName is the workflow name. triggerName is the workflow
// trigger name.
func (client WorkflowTriggersClient) Get(ctx context.Context, resourceGroupName string, workflowName string, triggerName string) (result WorkflowTrigger, err error) {
	req, err := client.GetPreparer(ctx, resourceGroupName, workflowName, triggerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "Get", nil, "Failure preparing request")
		return
	}

	resp, err := client.GetSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "Get", resp, "Failure sending request")
		return
	}

	result, err = client.GetResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "Get", resp, "Failure responding to request")
	}

	return
}

// GetPreparer prepares the Get request.
func (client WorkflowTriggersClient) GetPreparer(ctx context.Context, resourceGroupName string, workflowName string, triggerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"triggerName":       autorest.Encode("path", triggerName),
		"workflowName":      autorest.Encode("path", workflowName),
	}

	const APIVersion = "2015-02-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/triggers/{triggerName}", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// GetSender sends the Get request. The method will close the
// http.Response Body if it receives an error.
func (client WorkflowTriggersClient) GetSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// GetResponder handles the response to the Get request. The method always
// closes the http.Response Body.
func (client WorkflowTriggersClient) GetResponder(resp *http.Response) (result WorkflowTrigger, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// List gets a list of workflow triggers.
//
// resourceGroupName is the resource group name. workflowName is the workflow name. top is the number of items to
// be included in the result. filter is the filter to apply on the operation.
func (client WorkflowTriggersClient) List(ctx context.Context, resourceGroupName string, workflowName string, top *int32, filter string) (result WorkflowTriggerListResultPage, err error) {
	result.fn = client.listNextResults
	req, err := client.ListPreparer(ctx, resourceGroupName, workflowName, top, filter)
	if err != nil {
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "List", nil, "Failure preparing request")
		return
	}

	resp, err := client.ListSender(req)
	if err != nil {
		result.wtlr.Response = autorest.Response{Response: resp}
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "List", resp, "Failure sending request")
		return
	}

	result.wtlr, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "List", resp, "Failure responding to request")
	}

	return
}

// ListPreparer prepares the List request.
func (client WorkflowTriggersClient) ListPreparer(ctx context.Context, resourceGroupName string, workflowName string, top *int32, filter string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"workflowName":      autorest.Encode("path", workflowName),
	}

	const APIVersion = "2015-02-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}
	if top != nil {
		queryParameters["$top"] = autorest.Encode("query", *top)
	}
	if len(filter) > 0 {
		queryParameters["$filter"] = autorest.Encode("query", filter)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsGet(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/triggers/", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// ListSender sends the List request. The method will close the
// http.Response Body if it receives an error.
func (client WorkflowTriggersClient) ListSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// ListResponder handles the response to the List request. The method always
// closes the http.Response Body.
func (client WorkflowTriggersClient) ListResponder(resp *http.Response) (result WorkflowTriggerListResult, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result),
		autorest.ByClosing())
	result.Response = autorest.Response{Response: resp}
	return
}

// listNextResults retrieves the next set of results, if any.
func (client WorkflowTriggersClient) listNextResults(lastResults WorkflowTriggerListResult) (result WorkflowTriggerListResult, err error) {
	req, err := lastResults.workflowTriggerListResultPreparer()
	if err != nil {
		return result, autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "listNextResults", nil, "Failure preparing next results request")
	}
	if req == nil {
		return
	}
	resp, err := client.ListSender(req)
	if err != nil {
		result.Response = autorest.Response{Response: resp}
		return result, autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "listNextResults", resp, "Failure sending next results request")
	}
	result, err = client.ListResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "listNextResults", resp, "Failure responding to next results request")
	}
	return
}

// ListComplete enumerates all values, automatically crossing page boundaries as required.
func (client WorkflowTriggersClient) ListComplete(ctx context.Context, resourceGroupName string, workflowName string, top *int32, filter string) (result WorkflowTriggerListResultIterator, err error) {
	result.page, err = client.List(ctx, resourceGroupName, workflowName, top, filter)
	return
}

// Run runs a workflow trigger.
//
// resourceGroupName is the resource group name. workflowName is the workflow name. triggerName is the workflow
// trigger name.
func (client WorkflowTriggersClient) Run(ctx context.Context, resourceGroupName string, workflowName string, triggerName string) (result autorest.Response, err error) {
	req, err := client.RunPreparer(ctx, resourceGroupName, workflowName, triggerName)
	if err != nil {
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "Run", nil, "Failure preparing request")
		return
	}

	resp, err := client.RunSender(req)
	if err != nil {
		result.Response = resp
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "Run", resp, "Failure sending request")
		return
	}

	result, err = client.RunResponder(resp)
	if err != nil {
		err = autorest.NewErrorWithError(err, "logic.WorkflowTriggersClient", "Run", resp, "Failure responding to request")
	}

	return
}

// RunPreparer prepares the Run request.
func (client WorkflowTriggersClient) RunPreparer(ctx context.Context, resourceGroupName string, workflowName string, triggerName string) (*http.Request, error) {
	pathParameters := map[string]interface{}{
		"resourceGroupName": autorest.Encode("path", resourceGroupName),
		"subscriptionId":    autorest.Encode("path", client.SubscriptionID),
		"triggerName":       autorest.Encode("path", triggerName),
		"workflowName":      autorest.Encode("path", workflowName),
	}

	const APIVersion = "2015-02-01-preview"
	queryParameters := map[string]interface{}{
		"api-version": APIVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsPost(),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters("/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/triggers/{triggerName}/run", pathParameters),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// RunSender sends the Run request. The method will close the
// http.Response Body if it receives an error.
func (client WorkflowTriggersClient) RunSender(req *http.Request) (*http.Response, error) {
	return autorest.SendWithSender(client, req,
		azure.DoRetryWithRegistration(client.Client))
}

// RunResponder handles the response to the Run request. The method always
// closes the http.Response Body.
func (client WorkflowTriggersClient) RunResponder(resp *http.Response) (result autorest.Response, err error) {
	err = autorest.Respond(
		resp,
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.Response = resp
	return
}
