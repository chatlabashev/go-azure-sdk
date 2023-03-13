package managedcassandras

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DataCenterResourceProperties struct {
	AvailabilityZone                   *bool                              `json:"availabilityZone,omitempty"`
	BackupStorageCustomerKeyUri        *string                            `json:"backupStorageCustomerKeyUri,omitempty"`
	Base64EncodedCassandraYamlFragment *string                            `json:"base64EncodedCassandraYamlFragment,omitempty"`
	DataCenterLocation                 *string                            `json:"dataCenterLocation,omitempty"`
	DelegatedSubnetId                  *string                            `json:"delegatedSubnetId,omitempty"`
	DiskCapacity                       *int64                             `json:"diskCapacity,omitempty"`
	DiskSku                            *string                            `json:"diskSku,omitempty"`
	ManagedDiskCustomerKeyUri          *string                            `json:"managedDiskCustomerKeyUri,omitempty"`
	NodeCount                          *int64                             `json:"nodeCount,omitempty"`
	ProvisioningState                  *ManagedCassandraProvisioningState `json:"provisioningState,omitempty"`
	SeedNodes                          *[]SeedNode                        `json:"seedNodes,omitempty"`
	Sku                                *string                            `json:"sku,omitempty"`
}
