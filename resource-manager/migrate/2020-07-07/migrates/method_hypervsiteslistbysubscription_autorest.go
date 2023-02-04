package migrates

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HyperVSitesListBySubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]HyperVSite

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (HyperVSitesListBySubscriptionOperationResponse, error)
}

type HyperVSitesListBySubscriptionCompleteResult struct {
	Items []HyperVSite
}

func (r HyperVSitesListBySubscriptionOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r HyperVSitesListBySubscriptionOperationResponse) LoadMore(ctx context.Context) (resp HyperVSitesListBySubscriptionOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// HyperVSitesListBySubscription ...
func (c MigratesClient) HyperVSitesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (resp HyperVSitesListBySubscriptionOperationResponse, err error) {
	req, err := c.preparerForHyperVSitesListBySubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrates.MigratesClient", "HyperVSitesListBySubscription", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrates.MigratesClient", "HyperVSitesListBySubscription", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForHyperVSitesListBySubscription(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "migrates.MigratesClient", "HyperVSitesListBySubscription", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForHyperVSitesListBySubscription prepares the HyperVSitesListBySubscription request.
func (c MigratesClient) preparerForHyperVSitesListBySubscription(ctx context.Context, id commonids.SubscriptionId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.OffAzure/hyperVSites", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForHyperVSitesListBySubscriptionWithNextLink prepares the HyperVSitesListBySubscription request with the given nextLink token.
func (c MigratesClient) preparerForHyperVSitesListBySubscriptionWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
	uri, err := url.Parse(nextLink)
	if err != nil {
		return nil, fmt.Errorf("parsing nextLink %q: %+v", nextLink, err)
	}
	queryParameters := map[string]interface{}{}
	for k, v := range uri.Query() {
		if len(v) == 0 {
			continue
		}
		val := v[0]
		val = autorest.Encode("query", val)
		queryParameters[k] = val
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForHyperVSitesListBySubscription handles the response to the HyperVSitesListBySubscription request. The method always
// closes the http.Response Body.
func (c MigratesClient) responderForHyperVSitesListBySubscription(resp *http.Response) (result HyperVSitesListBySubscriptionOperationResponse, err error) {
	type page struct {
		Values   []HyperVSite `json:"value"`
		NextLink *string      `json:"nextLink"`
	}
	var respObj page
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&respObj),
		autorest.ByClosing())
	result.HttpResponse = resp
	result.Model = &respObj.Values
	result.nextLink = respObj.NextLink
	if respObj.NextLink != nil {
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result HyperVSitesListBySubscriptionOperationResponse, err error) {
			req, err := c.preparerForHyperVSitesListBySubscriptionWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "migrates.MigratesClient", "HyperVSitesListBySubscription", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "migrates.MigratesClient", "HyperVSitesListBySubscription", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForHyperVSitesListBySubscription(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "migrates.MigratesClient", "HyperVSitesListBySubscription", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// HyperVSitesListBySubscriptionComplete retrieves all of the results into a single object
func (c MigratesClient) HyperVSitesListBySubscriptionComplete(ctx context.Context, id commonids.SubscriptionId) (HyperVSitesListBySubscriptionCompleteResult, error) {
	return c.HyperVSitesListBySubscriptionCompleteMatchingPredicate(ctx, id, HyperVSiteOperationPredicate{})
}

// HyperVSitesListBySubscriptionCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c MigratesClient) HyperVSitesListBySubscriptionCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate HyperVSiteOperationPredicate) (resp HyperVSitesListBySubscriptionCompleteResult, err error) {
	items := make([]HyperVSite, 0)

	page, err := c.HyperVSitesListBySubscription(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading the initial page: %+v", err)
		return
	}
	if page.Model != nil {
		for _, v := range *page.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	for page.HasMore() {
		page, err = page.LoadMore(ctx)
		if err != nil {
			err = fmt.Errorf("loading the next page: %+v", err)
			return
		}

		if page.Model != nil {
			for _, v := range *page.Model {
				if predicate.Matches(v) {
					items = append(items, v)
				}
			}
		}
	}

	out := HyperVSitesListBySubscriptionCompleteResult{
		Items: items,
	}
	return out, nil
}
