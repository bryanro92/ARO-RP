package api

import (
	hivev1 "github.com/openshift/hive/apis/hive/v1"
)

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

// ClusterManagerConfiguration represents the configuration from OpenShift Cluster Manager (OCM)
type ClusterManagerConfiguration struct {
	MissingFields

	ID       string `json:"id,omitempty"`
	Deleting bool   `json:"deleting,omitempty"` // https://docs.microsoft.com/en-us/azure/cosmos-db/change-feed-design-patterns#deletes

	ClusterDeployment    hivev1.ClusterDeployment    `json:"clusterDeployment,omitempty"`
	SyncIdentityProvider hivev1.SyncIdentityProvider `json:"syncIdentityProvider,omitempty"`
	// consider a single record
	SyncSet     hivev1.SyncSet     `json:"syncSets,omitempty"`
	MachinePool hivev1.MachinePool `json:"machinePool,omitempty"`
}

// single item/resource per record/row
// look at if we should use the hive types directly elad
