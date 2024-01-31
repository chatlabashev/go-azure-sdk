package deletedservers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByLocationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeletedServer
}

type ListByLocationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeletedServer
}

// ListByLocation ...
func (c DeletedServersClient) ListByLocation(ctx context.Context, id LocationId) (result ListByLocationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/deletedServers", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]DeletedServer `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByLocationComplete retrieves all the results into a single object
func (c DeletedServersClient) ListByLocationComplete(ctx context.Context, id LocationId) (ListByLocationCompleteResult, error) {
	return c.ListByLocationCompleteMatchingPredicate(ctx, id, DeletedServerOperationPredicate{})
}

// ListByLocationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeletedServersClient) ListByLocationCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate DeletedServerOperationPredicate) (result ListByLocationCompleteResult, err error) {
	items := make([]DeletedServer, 0)

	resp, err := c.ListByLocation(ctx, id)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ListByLocationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
