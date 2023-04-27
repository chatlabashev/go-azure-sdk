package invitation

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

type ListByShareOperationResponse struct {
	HttpResponse *http.Response
	Model        *[]Invitation

	nextLink     *string
	nextPageFunc func(ctx context.Context, nextLink string) (ListByShareOperationResponse, error)
}

type ListByShareCompleteResult struct {
	Items []Invitation
}

func (r ListByShareOperationResponse) HasMore() bool {
	return r.nextLink != nil
}

func (r ListByShareOperationResponse) LoadMore(ctx context.Context) (resp ListByShareOperationResponse, err error) {
	if !r.HasMore() {
		err = fmt.Errorf("no more pages returned")
		return
	}
	return r.nextPageFunc(ctx, *r.nextLink)
}

type ListByShareOperationOptions struct {
	Filter  *string
	Orderby *string
}

func DefaultListByShareOperationOptions() ListByShareOperationOptions {
	return ListByShareOperationOptions{}
}

func (o ListByShareOperationOptions) toHeaders() map[string]interface{} {
	out := make(map[string]interface{})

	return out
}

func (o ListByShareOperationOptions) toQueryString() map[string]interface{} {
	out := make(map[string]interface{})

	if o.Filter != nil {
		out["$filter"] = *o.Filter
	}

	if o.Orderby != nil {
		out["$orderby"] = *o.Orderby
	}

	return out
}

// ListByShare ...
func (c InvitationClient) ListByShare(ctx context.Context, id ShareId, options ListByShareOperationOptions) (resp ListByShareOperationResponse, err error) {
	req, err := c.preparerForListByShare(ctx, id, options)
	if err != nil {
		err = autorest.NewErrorWithError(err, "invitation.InvitationClient", "ListByShare", nil, "Failure preparing request")
		return
	}

	resp.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "invitation.InvitationClient", "ListByShare", resp.HttpResponse, "Failure sending request")
		return
	}

	resp, err = c.responderForListByShare(resp.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "invitation.InvitationClient", "ListByShare", resp.HttpResponse, "Failure responding to request")
		return
	}
	return
}

// preparerForListByShare prepares the ListByShare request.
func (c InvitationClient) preparerForListByShare(ctx context.Context, id ShareId, options ListByShareOperationOptions) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	for k, v := range options.toQueryString() {
		queryParameters[k] = autorest.Encode("query", v)
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithHeaders(options.toHeaders()),
		autorest.WithPath(fmt.Sprintf("%s/invitations", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// preparerForListByShareWithNextLink prepares the ListByShare request with the given nextLink token.
func (c InvitationClient) preparerForListByShareWithNextLink(ctx context.Context, nextLink string) (*http.Request, error) {
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

// responderForListByShare handles the response to the ListByShare request. The method always
// closes the http.Response Body.
func (c InvitationClient) responderForListByShare(resp *http.Response) (result ListByShareOperationResponse, err error) {
	type page struct {
		Values   []Invitation `json:"value"`
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
		result.nextPageFunc = func(ctx context.Context, nextLink string) (result ListByShareOperationResponse, err error) {
			req, err := c.preparerForListByShareWithNextLink(ctx, nextLink)
			if err != nil {
				err = autorest.NewErrorWithError(err, "invitation.InvitationClient", "ListByShare", nil, "Failure preparing request")
				return
			}

			result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
			if err != nil {
				err = autorest.NewErrorWithError(err, "invitation.InvitationClient", "ListByShare", result.HttpResponse, "Failure sending request")
				return
			}

			result, err = c.responderForListByShare(result.HttpResponse)
			if err != nil {
				err = autorest.NewErrorWithError(err, "invitation.InvitationClient", "ListByShare", result.HttpResponse, "Failure responding to request")
				return
			}

			return
		}
	}
	return
}

// ListByShareComplete retrieves all of the results into a single object
func (c InvitationClient) ListByShareComplete(ctx context.Context, id ShareId, options ListByShareOperationOptions) (ListByShareCompleteResult, error) {
	return c.ListByShareCompleteMatchingPredicate(ctx, id, options, InvitationOperationPredicate{})
}

// ListByShareCompleteMatchingPredicate retrieves all of the results and then applied the predicate
func (c InvitationClient) ListByShareCompleteMatchingPredicate(ctx context.Context, id ShareId, options ListByShareOperationOptions, predicate InvitationOperationPredicate) (resp ListByShareCompleteResult, err error) {
	items := make([]Invitation, 0)

	page, err := c.ListByShare(ctx, id, options)
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

	out := ListByShareCompleteResult{
		Items: items,
	}
	return out, nil
}
