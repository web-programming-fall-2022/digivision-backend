package productmeta

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/arimanius/digivision-backend/pkg/api/v1"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type DigikalaFetcher struct {
	baseUrl    string
	apiBaseUrl string
	client     *resty.Client
}

func NewDigikalaFetcher(baseUrl string, apiBaseUrl string, client *resty.Client) DigikalaFetcher {
	return DigikalaFetcher{
		baseUrl:    baseUrl,
		apiBaseUrl: apiBaseUrl,
		client:     client,
	}
}

func (f DigikalaFetcher) Fetch(ctx context.Context, productId string) (*v1.Product, error) {
	pid, err := strconv.Atoi(productId)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/%d/", f.apiBaseUrl, pid)
	resp, err := f.client.R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}
	if resp.Status() != "200 OK" {
		return nil, fmt.Errorf("failed to fetch product %s. status: %s", productId, resp.Status())
	}
	product := DigikalaProduct{}
	err = json.Unmarshal(resp.Body(), &product)
	if err != nil {
		return nil, err
	}
	p := product.Data.Product
	price := int64(0)
	switch variant := p.DefaultVariant.(type) {
	case Variant:
		price = variant.SellingPrice
	}
	return &v1.Product{
		Id:       int32(pid),
		Title:    p.TitleFa,
		Url:      fmt.Sprintf("%s%s", f.baseUrl, p.Url),
		Status:   p.Status,
		ImageUrl: p.Images.Main.Url[0],
		Rate: &v1.Rating{
			Rate:  p.Rating.Rate,
			Count: p.Rating.Count,
		},
		Categories: ToCategories(f.baseUrl, p.Breadcrumb),
		Price:      price,
	}, nil
}
