package privateendpointconnection

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type PrivateLinkServiceConnectionState struct {
	ActionsRequired *string `json:"actionsRequired,omitempty"`
	Description     *string `json:"description,omitempty"`
	Status          *Status `json:"status,omitempty"`
}
