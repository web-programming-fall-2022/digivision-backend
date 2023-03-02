package productmeta

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	v1 "github.com/web-programming-fall-2022/digivision-backend/pkg/api/v1"
)

type Breadcrumb struct {
	Title string `json:"title"`
	Url   struct {
		Uri string `json:"uri"`
	}
}

type Variant struct {
	Price struct {
		SellingPrice int64 `json:"selling_price"`
	} `json:"price"`
}

func (v *Variant) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" || (data[0] == '[' && data[len(data)-1] == ']') {
		return nil
	}
	if data[0] == '{' && data[len(data)-1] == '}' { // object?
		type TempVariant struct {
			Price struct {
				SellingPrice int64 `json:"selling_price"`
			} `json:"price"`
		}
		var tempVariant TempVariant
		err := json.Unmarshal(data, &tempVariant)
		if err != nil {
			return err
		}
		v.Price.SellingPrice = tempVariant.Price.SellingPrice
		return nil
	}
	return errors.New("invalid variant")
}

type DigikalaProduct struct {
	Data struct {
		Product struct {
			IsInactive bool   `json:"is_inactive"`
			TitleFa    string `json:"title_fa"`
			Url        struct {
				Uri string `json:"uri"`
			} `json:"url"`
			Status string `json:"status"`
			Images struct {
				Main struct {
					Url []string `json:"url"`
				}
			} `json:"images"`
			Rating struct {
				Rate  int32 `json:"rate"`
				Count int32 `json:"count"`
			} `json:"rating"`
			Breadcrumb     []Breadcrumb `json:"breadcrumb"`
			DefaultVariant Variant      `json:"default_variant"`
		}
	}
}

func (b Breadcrumb) ToCategory(baseUrl string) *v1.Category {
	return &v1.Category{
		Title: b.Title,
		Url:   fmt.Sprintf("%s%s", baseUrl, b.Url.Uri),
	}
}

func ToCategories(baseUrl string, breadcrumb []Breadcrumb) []*v1.Category {
	categories := make([]*v1.Category, len(breadcrumb)-1)
	for i, b := range breadcrumb {
		if i == len(breadcrumb)-1 {
			break
		}
		categories[i] = b.ToCategory(baseUrl)
	}
	return categories
}
