package nodes

import (
	"fmt"

	"github.com/hashicorp/go-azure-sdk/sdk/client/resourcemanager"
	sdkEnv "github.com/hashicorp/go-azure-sdk/sdk/environments"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NodesClient struct {
	Client *resourcemanager.Client
}

func NewNodesClientWithBaseURI(sdkApi sdkEnv.Api) (*NodesClient, error) {
	client, err := resourcemanager.NewResourceManagerClient(sdkApi, "nodes", defaultApiVersion)
	if err != nil {
		return nil, fmt.Errorf("instantiating NodesClient: %+v", err)
	}

	return &NodesClient{
		Client: client,
	}, nil
}
