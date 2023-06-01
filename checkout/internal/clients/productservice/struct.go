package productservice

type BaseRequest struct {
	Token string `json:"token"`
}

type ListSKUsRequest struct {
	BaseRequest
	StartAfterSku int64 `json:"startAfterSku"`
	Count         int64 `json:"count"`
}

type ListSKUsResponse struct {
	SKUs []int64 `json:"skus"`
}

type GetProductRequest struct {
	BaseRequest
	SKU int64 `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price int64  `json:"price"`
}
