package hypervsites

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetSiteHealthSummaryOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]SiteHealthSummary

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (GetSiteHealthSummaryOperationResponse, error)
}

type GetSiteHealthSummaryCompleteResult struct {
	Items []SiteHealthSummary
}

func (r GetSiteHealthSummaryOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r GetSiteHealthSummaryOperationResponse) LoadMore(ctx context.Context) (resp GetSiteHealthSummaryOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

// GetSiteHealthSummary ...
func (c HyperVSitesClient) GetSiteHealthSummary(ctx context.Context, id HyperVSiteId) (resp GetSiteHealthSummaryOperationResponse, err error) {
	req, err := c.preparerForGetSiteHealthSummary(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hypervsites.HyperVSitesClient", "GetSiteHealthSummary", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "hypervsites.HyperVSitesClient", "GetSiteHealthSummary", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForGetSiteHealthSummary(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "hypervsites.HyperVSitesClient", "GetSiteHealthSummary", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForGetSiteHealthSummary prepares the GetSiteHealthSummary request.
func (c HyperVSitesClient) preparerForGetSiteHealthSummary(ctx context.Context, id HyperVSiteId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/healthSummary", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForGetSiteHealthSummaryWithNextLink prepares the GetSiteHealthSummary request with the given nextLink token.
func (c HyperVSitesClient) preparerForGetSiteHealthSummaryWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(uri.Path),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetSiteHealthSummary handles the response to the GetSiteHealthSummary request. The method always
// closes the http.Response Body.
func (c HyperVSitesClient) responderForGetSiteHealthSummary(resp *http.Response) (result GetSiteHealthSummaryOperationResponse, err error) {
	type page struct {
		Values   []SiteHealthSummary `json:"value"`
		NextLink *string             `json:"nextLink"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result GetSiteHealthSummaryOperationResponse, err error) {
			req, err := c.preparerForGetSiteHealthSummaryWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "hypervsites.HyperVSitesClient", "GetSiteHealthSummary", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "hypervsites.HyperVSitesClient", "GetSiteHealthSummary", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForGetSiteHealthSummary(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "hypervsites.HyperVSitesClient", "GetSiteHealthSummary", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// GetSiteHealthSummaryComplete retrieves all of the results into a single object
func (c HyperVSitesClient) GetSiteHealthSummaryComplete(ctx context.Context, id HyperVSiteId) (GetSiteHealthSummaryCompleteResult, error) {
	return c.GetSiteHealthSummaryCompleteMatchingPredicate(ctx, id, SiteHealthSummaryOperationPredicate{})
}

// GetSiteHealthSummaryCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c HyperVSitesClient) GetSiteHealthSummaryCompleteMatchingPredicate(ctx context.Context, id HyperVSiteId, predicate SiteHealthSummaryOperationPredicate) (resp GetSiteHealthSummaryCompleteResult, err error) {
	items := make([]SiteHealthSummary, 0)

	page, err := c.GetSiteHealthSummary(ctx, id)
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

	out := GetSiteHealthSummaryCompleteResult{
		Items: items,
	}
	return out, nil
}
