package search

import (
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/pkg/errors"
	"strconv"
)

const (
	VectorColumnName = "vector"
	TopKExpansion    = 10
)

// MilvusSearchHandler implements Handler interface{}
type MilvusSearchHandler struct {
	client         client.Client
	vectorDim      int
	metricType     entity.MetricType
	nProbe         int
	collectionName string
}

// NewMilvusSearchHandler returns a new MilvusSearchHandler
func NewMilvusSearchHandler(
	client client.Client,
	vectorDim int,
	metricType entity.MetricType,
	nProbe int,
	collectionName string) MilvusSearchHandler {
	return MilvusSearchHandler{
		client:         client,
		vectorDim:      vectorDim,
		metricType:     metricType,
		nProbe:         nProbe,
		collectionName: collectionName,
	}
}

// Search implements Handler interface{}
func (h MilvusSearchHandler) Search(ctx context.Context, query []float32, topK int) ([]ProductImage, error) {
	err := h.client.LoadCollection(ctx, h.collectionName, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load the partition")
	}
	sp, _ := entity.NewIndexIvfPQSearchParam(h.nProbe)

	searchResult, err := h.client.Search(
		ctx,
		h.collectionName,
		[]string{},
		"",
		[]string{"product_id"},
		[]entity.Vector{entity.FloatVector(query)},
		VectorColumnName,
		h.metricType,
		topK*TopKExpansion,
		sp,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to conduct search")
	}
	imageIds, ok := searchResult[0].IDs.(*entity.ColumnInt64)
	if !ok {
		return nil, errors.New("failed to convert pk to string column")
	}
	ids, ok := searchResult[0].Fields[0].(*entity.ColumnVarChar)
	if !ok {
		return nil, errors.New("failed to convert product id to string column")
	}
	distances := searchResult[0].Scores
	productImages := ids2ProductImages(ids.Data(), imageIds.Data(), distances)
	err = h.client.ReleaseCollection(ctx, h.collectionName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to release the partition")
	}
	return productImages, nil
}

func ids2ProductImages(ids []string, imageIds []int64, distances []float32) []ProductImage {
	products := make([]ProductImage, len(ids))
	for i := range ids {
		products[i] = ProductImage{
			ProductId: ids[i],
			ImageId:   strconv.FormatInt(imageIds[i], 10),
			Distance:  distances[i],
		}
	}
	return products
}
