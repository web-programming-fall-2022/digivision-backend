package productmeta

import (
	"fmt"
	v1 "github.com/arimanius/digivision-backend/pkg/api/v1"
)

type Breadcumb struct {
	Title string `json:"title"`
	Url   struct {
		Uri string `json:"uri"`
	}
}

type DigikalaProduct struct {
	Data struct {
		Product struct {
			TitleFa string `json:"title_fa"`
			Url     struct {
				Uri string `json:"uri"`
			}
			Status string `json:"status"`
			Images struct {
				Main struct {
					Url []string `json:"url"`
				}
			}
			Rating struct {
				Rate  int32 `json:"rate"`
				Count int32 `json:"count"`
			}
			Breadcrumb     []Breadcumb `json:"breadcrumb"`
			DefaultVariant struct {
				SellingPrice int64 `json:"selling_price"`
			} `json:"default_variant"`
		}
	}
}

func (b Breadcumb) ToCategory(baseUrl string) *v1.Category {
	return &v1.Category{
		Title: b.Title,
		Url:   fmt.Sprintf("%s%s", baseUrl, b.Url.Uri),
	}
}

func ToCategories(baseUrl string, breadcrumb []Breadcumb) []*v1.Category {
	categories := make([]*v1.Category, len(breadcrumb)-1)
	for i, b := range breadcrumb {
		if i == len(breadcrumb)-1 {
			break
		}
		categories[i] = b.ToCategory(baseUrl)
	}
	return categories
}
