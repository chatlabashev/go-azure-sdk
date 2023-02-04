package mastersites

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PatchSiteOperationResponse struct {
	HttpResponse *http.Response
	Model        *MasterSite
}

// PatchSite ...
func (c MasterSitesClient) PatchSite(ctx context.Context, id MasterSiteId, input MasterSite) (result PatchSiteOperationResponse, err error) {
	req, err := c.preparerForPatchSite(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mastersites.MasterSitesClient", "PatchSite", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "mastersites.MasterSitesClient", "PatchSite", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForPatchSite(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "mastersites.MasterSitesClient", "PatchSite", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForPatchSite prepares the PatchSite request.
func (c MasterSitesClient) preparerForPatchSite(ctx context.Context, id MasterSiteId, input MasterSite) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPatch(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForPatchSite handles the response to the PatchSite request. The method always
// closes the http.Response Body.
func (c MasterSitesClient) responderForPatchSite(resp *http.Response) (result PatchSiteOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusAccepted, http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
