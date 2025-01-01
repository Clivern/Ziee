// Copyright 2026 Actx0. All rights reserved.
// License can be found in the LICENSE file.

package qdrant

import (
	"context"
	"fmt"

	qdrantsdk "github.com/qdrant/go-client/qdrant"
	"github.com/spf13/viper"
)

// Client wraps the Qdrant SDK for storing and querying embeddings.
type Client struct {
	api *qdrantsdk.Client
}

// New returns a Qdrant client loaded from app.ai.embed.qdb config.
func New() (*Client, error) {
	api, err := qdrantsdk.NewClient(&qdrantsdk.Config{
		Host:     viper.GetString("app.ai.embed.qdb.host"),
		Port:     viper.GetInt("app.ai.embed.qdb.port"),
		APIKey:   viper.GetString("app.ai.embed.qdb.api_key"),
		UseTLS:   viper.GetBool("app.ai.embed.qdb.use_tls"),
		PoolSize: uint(viper.GetInt("app.ai.embed.qdb.pool_size")),
	})

	if err != nil {
		return nil, err
	}

	return &Client{api: api}, nil
}

// EnsureCollection creates the collection when it does not exist.
func (c *Client) EnsureCollection(ctx context.Context, collection string, indexes ...Index) error {
	exists, err := c.api.CollectionExists(ctx, collection)

	if err != nil {
		return fmt.Errorf("qdrant collection exists: %w", err)
	}

	if exists {
		return nil
	}

	return c.CreateCollection(ctx, collection, indexes...)
}

// CreateCollection creates a collection and its payload indexes.
func (c *Client) CreateCollection(ctx context.Context, collection string, indexes ...Index) error {
	onDisk := true
	err := c.api.CreateCollection(ctx, &qdrantsdk.CreateCollection{
		CollectionName: collection,
		OnDiskPayload:  &onDisk,
		VectorsConfig: qdrantsdk.NewVectorsConfig(&qdrantsdk.VectorParams{
			Size:     uint64(viper.GetInt("app.ai.embed.qdb.vector_size")),
			Distance: qdrantsdk.Distance_Cosine,
			OnDisk:   &onDisk,
			HnswConfig: &qdrantsdk.HnswConfigDiff{
				OnDisk: &onDisk,
			},
		}),
	})

	if err != nil {
		return fmt.Errorf("qdrant create collection: %w", err)
	}

	for _, index := range indexes {
		_, err := c.api.CreateFieldIndex(ctx, &qdrantsdk.CreateFieldIndexCollection{
			CollectionName: collection,
			FieldName:      index.Field,
			FieldType:      ParseFieldType(index.Type),
			Wait:           new(true),
		})

		if err != nil && !IsIndexExists(err) {
			return fmt.Errorf("qdrant create field index %q: %w", index.Field, err)
		}
	}

	return nil
}

// Upsert stores or updates vector points in the given collection.
func (c *Client) Upsert(ctx context.Context, collection string, points []Point) error {
	qpoints := make([]*qdrantsdk.PointStruct, 0, len(points))

	for _, point := range points {
		if len(point.Vector) == 0 {
			return ErrEmptyVector
		}

		payload := make(map[string]*qdrantsdk.Value, len(point.Payload))
		for key, value := range point.Payload {
			payload[key] = qdrantsdk.NewValueString(value)
		}

		qpoints = append(qpoints, &qdrantsdk.PointStruct{
			Id:      qdrantsdk.NewID(point.Id),
			Vectors: qdrantsdk.NewVectors(point.Vector...),
			Payload: payload,
		})
	}

	_, err := c.api.Upsert(ctx, &qdrantsdk.UpsertPoints{
		CollectionName: collection,
		Wait:           new(true),
		Points:         qpoints,
	})

	if err != nil {
		return fmt.Errorf("qdrant upsert: %w", err)
	}

	return nil
}

// Search finds the nearest vectors for the given query in a collection.
func (c *Client) Search(ctx context.Context, collection string, query Query) ([]Result, error) {
	request := &qdrantsdk.QueryPoints{
		CollectionName: collection,
		Query:          qdrantsdk.NewQuery(query.Vector...),
		Limit:          &query.Limit,
		WithPayload:    qdrantsdk.NewWithPayload(true),
	}

	if len(query.Filters) > 0 {
		conditions := make([]*qdrantsdk.Condition, 0, len(query.Filters))
		for field, value := range query.Filters {
			conditions = append(conditions, qdrantsdk.NewMatchKeyword(field, value))
		}
		request.Filter = &qdrantsdk.Filter{Must: conditions}
	}

	points, err := c.api.Query(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("qdrant search: %w", err)
	}

	results := make([]Result, 0, len(points))
	for _, point := range points {
		payload := make(map[string]string, len(point.Payload))
		for key, value := range point.Payload {
			if v, ok := value.Kind.(*qdrantsdk.Value_StringValue); ok {
				payload[key] = v.StringValue
			}
		}

		results = append(results, Result{
			Id:      PointIdString(point.Id),
			Score:   point.Score,
			Payload: payload,
		})
	}

	return results, nil
}

// Delete removes points from the given collection by Id.
func (c *Client) Delete(ctx context.Context, collection string, ids []string) error {
	pointIds := make([]*qdrantsdk.PointId, 0, len(ids))

	for _, id := range ids {
		pointIds = append(pointIds, qdrantsdk.NewID(id))
	}

	_, err := c.api.Delete(ctx, &qdrantsdk.DeletePoints{
		CollectionName: collection,
		Wait:           new(true),
		Points:         qdrantsdk.NewPointsSelector(pointIds...),
	})

	if err != nil {
		return fmt.Errorf("qdrant delete: %w", err)
	}

	return nil
}

// DeleteByFilter removes points matching all payload filters in the given collection.
func (c *Client) DeleteByFilter(ctx context.Context, collection string, filters map[string]string) error {
	if len(filters) == 0 {
		return nil
	}

	conditions := make([]*qdrantsdk.Condition, 0, len(filters))
	for field, value := range filters {
		conditions = append(conditions, qdrantsdk.NewMatchKeyword(field, value))
	}

	_, err := c.api.Delete(ctx, &qdrantsdk.DeletePoints{
		CollectionName: collection,
		Wait:           new(true),
		Points: qdrantsdk.NewPointsSelectorFilter(&qdrantsdk.Filter{
			Must: conditions,
		}),
	})

	if err != nil {
		return fmt.Errorf("qdrant delete by filter: %w", err)
	}

	return nil
}

// Close closes the underlying Qdrant connection pool.
func (c *Client) Close() error {
	c.api.Close()
	return nil
}
