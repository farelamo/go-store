package controller

import (
	"errors"
	"store/constant"
	"store/repository"
	"store/utils"
	"strconv"

	schema "store/schema/product"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductController interface {
	GetProducts(c *gin.Context)
	GetProductByID(c *gin.Context)
	CreateProduct(c *gin.Context)
}

type productController struct {
	productRepository repository.ProductRepository
}

func NewProductController(productRepo repository.ProductRepository) ProductController {
	return &productController{
		productRepository: productRepo,
	}
}

func (p *productController) GetProducts(c *gin.Context) {
	var req schema.GetProductReq
	if err := c.ShouldBindQuery(&req); err != nil {
		errMsg := utils.FormatValidationError(err)
		utils.Response(c, 400, nil, nil, nil, "Invalid request parameters: "+errMsg, nil, false)
		return
	}

	req.Page, req.PageSize = utils.Paginate(req.Page, req.PageSize)
	var sort string

	if req.Sort != "" {
		_sort, err := utils.SortChecker(constant.ProductSort, req.Sort)
		if err != nil {
			utils.Response(c, 400, nil, nil, nil, err.Error(), nil, false)
			return
		}
		sort = _sort
	}

	products, total, err := p.productRepository.GetProducts(req.Page, req.PageSize, sort, req.Search)
	if err != nil {
		utils.Response(c, 500, nil, nil, nil, "Failed to retrieve products: "+err.Error(), nil, false)
		return
	}

	utils.Response(c, 200, &req.Page, &req.PageSize, &total, "Product retrieved successfully", products, false)
}

func (p *productController) GetProductByID(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		utils.Response(c, 400, nil, nil, nil, "Product ID is required", nil, false)
		return
	}

	id, err := strconv.Atoi(productID)
	if err != nil {
		utils.Response(c, 400, nil, nil, nil, "Invalid Product ID: "+err.Error(), nil, false)
		return
	}

	product, err := p.productRepository.GetProductByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Response(c, 404, nil, nil, nil, "data not found", nil, false)
			return
		}
		utils.Response(c, 500, nil, nil, nil, "Failed to retrieve product: "+err.Error(), nil, false)
		return
	}

	utils.Response(c, 200, nil, nil, nil,  "Product retrieved successfully", product, false)
}

func (p *productController) CreateProduct(c *gin.Context) {
	var product schema.CreateProductReq
	if err := c.ShouldBindJSON(&product); err != nil {
		errMsg := utils.FormatValidationError(err)
		utils.Response(c, 400, nil, nil, nil, "Invalid product data: "+errMsg, nil, false)
		return
	}

	if err := p.productRepository.CreateProduct(&product); err != nil {
		utils.Response(c, 500, nil, nil, nil, "Failed to create product: "+err.Error(), nil, false)
		return
	}

	utils.Response(c, 201, nil, nil, nil,  "Product created successfully", nil, false)
}