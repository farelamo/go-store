package product

type GetProductReq struct {
	Page     int    `json:"page" form:"page" binding:"required"`
	PageSize int    `json:"page_size" form:"page_size" binding:"required"`
	Sort     string `json:"sort" form:"sort"`
	Search   string `json:"search" form:"search"`
}

type CreateProductReq struct {
	Name        string `json:"name" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"required"`
}
