package domain

type Model struct {
	lomsService    ILomsService
	productService IProductService
}

func New(stockChecker ILomsService, productService IProductService) *Model {
	return &Model{
		lomsService:    stockChecker,
		productService: productService,
	}
}
