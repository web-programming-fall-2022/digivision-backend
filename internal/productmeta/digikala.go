package productmeta

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/arimanius/digivision-backend/pkg/api/v1"
	"github.com/go-resty/resty/v2"
	"strconv"
	"time"
)

type empty interface{}
type semaphore chan empty

type DigikalaFetcher struct {
	baseUrl    string
	apiBaseUrl string
	client     *resty.Client
	maxRetry   int
	sem        semaphore
}

func NewDigikalaFetcher(
	baseUrl string, apiBaseUrl string, client *resty.Client, maxRetry int, concurrencyFactor int,
) DigikalaFetcher {
	return DigikalaFetcher{
		baseUrl:    baseUrl,
		apiBaseUrl: apiBaseUrl,
		client:     client,
		maxRetry:   maxRetry,
		sem:        make(semaphore, concurrencyFactor),
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
		Price:      p.DefaultVariant.Price.SellingPrice,
	}, nil
}

func (f DigikalaFetcher) AsyncFetch(ctx context.Context, productIds []string, count int) (chan *v1.Product, chan error) {
	resp := make(chan *v1.Product)
	err := make(chan error)

	go func() {
		defer close(resp)
		defer close(err)
		var emp empty
		c := 0
		for _, productId := range productIds {
			f.sem <- emp
			p, e := f.Fetch(ctx, productId)
			retryCount := 0
			for e != nil {
				err <- e
				time.Sleep(1 * time.Second)
				if retryCount >= f.maxRetry {
					break
				}
				p, e = f.Fetch(ctx, productId)
				retryCount++
			}
			<-f.sem
			resp <- p
			c++
			if c >= count {
				break
			}
		}
	}()

	return resp, err
}
