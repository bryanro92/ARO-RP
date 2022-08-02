package database

// Copyright (c) Microsoft Corporation.
// Licensed under the Apache License 2.0.

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/ARO-RP/pkg/api"
	"github.com/Azure/ARO-RP/pkg/database/cosmosdb"
)

type clusterManagerConfiguration struct {
	c cosmosdb.ClusterManagerConfigurationDocumentClient
}

type ClusterManagerConfiguration interface {
	Create(context.Context, *api.OpenShiftClusterDocument) (*api.OpenShiftClusterDocument, error)
	Get(context.Context, string) (*api.OpenShiftClusterDocument, error)
	Patch(context.Context, string, func(*api.OpenShiftClusterDocument) error) (*api.OpenShiftClusterDocument, error)
	Delete(context.Context, *api.OpenShiftClusterDocument) error
	ChangeFeed() cosmosdb.OpenShiftClusterDocumentIterator
}

func NewClusterManagerConfiguration(ctx context.Context, isDevelopmentMode bool, dbc cosmosdb.DatabaseClient) (ClusterManagerConfiguration, error) {
	dbid, err := Name(isDevelopmentMode)
	if err != nil {
		return nil, err
	}

	collc := cosmosdb.NewCollectionClient(dbc, dbid)

	documentClient := cosmosdb.NewOpenShiftClusterDocumentClient(collc, collHiveResources)
	return NewClusterManagerConfigurationWithProvidedClient(documentClient), nil
}

func NewClusterManagerConfigurationWithProvidedClient(client cosmosdb.OpenShiftClusterDocumentClient) ClusterManagerConfiguration {
	return &clusterManagerConfiguration{c: client}
}

// Only used internally by Patch()
func (c *clusterManagerConfiguration) replace(ctx context.Context, doc *api.OpenShiftClusterDocument) (*api.OpenShiftClusterDocument, error) {
	if doc.ID != strings.ToLower(doc.ID) {
		return nil, fmt.Errorf("id %q is not lower case", doc.ID)
	}

	return c.c.Replace(ctx, doc.ID, doc, nil)
}

func (c *clusterManagerConfiguration) Create(ctx context.Context, doc *api.OpenShiftClusterDocument) (*api.OpenShiftClusterDocument, error) {
	if doc.ID != strings.ToLower(doc.ID) {
		return nil, fmt.Errorf("id %q is not lower case", doc.ID)
	}

	return c.c.Create(ctx, doc.ID, doc, nil)
}

func (c *clusterManagerConfiguration) Get(ctx context.Context, id string) (*api.OpenShiftClusterDocument, error) {
	if id != strings.ToLower(id) {
		return nil, fmt.Errorf("id %q is not lower case", id)
	}

	return c.c.Get(ctx, id, id, nil)
}

func (c *clusterManagerConfiguration) Patch(ctx context.Context, id string, callback func(*api.OpenShiftClusterDocument) error) (*api.OpenShiftClusterDocument, error) {
	var doc *api.OpenShiftClusterDocument

	err := cosmosdb.RetryOnPreconditionFailed(func() error {
		doc, err := c.Get(ctx, id)
		if err != nil {
			return err
		}

		err = callback(doc)
		if err != nil {
			return err
		}

		doc, err = c.replace(ctx, doc)
		return err
	})

	return doc, err
}

func (c *clusterManagerConfiguration) Delete(ctx context.Context, doc *api.OpenShiftClusterDocument) error {
	if doc.ID != strings.ToLower(doc.ID) {
		return fmt.Errorf("id %q is not lower case", doc.ID)
	}

	return c.c.Delete(ctx, doc.ID, doc, &cosmosdb.Options{NoETag: true})
}

func (c *clusterManagerConfiguration) ChangeFeed() cosmosdb.OpenShiftClusterDocumentIterator {
	return c.c.ChangeFeed(nil)
}
