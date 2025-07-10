package repository

import (
	"log"
	"store/model"
	reqSchema "store/schema/product"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetProducts(page, pageSize int, sort, search string) ([]model.Product, int64, error)
	GetProductByID(id uint) (*model.Product, error)
	CreateProduct(req *reqSchema.CreateProductReq) error
}

type productRepository struct {
	MySql *gorm.DB
}

func NewProductRepository(mysql *gorm.DB) ProductRepository {
	return &productRepository{
		MySql: mysql,
	}
}

func (p *productRepository) GetProducts(page, pageSize int, sort, search string) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := p.MySql.Model(&model.Product{})

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if sort != "" {
		query = query.Order(sort)
	} else {
		query = query.Order("created_at DESC")
	}

	log.Println("Querying products with pagination:", page, "pageSize:", pageSize, "sort:", sort, "search:", search)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (p *productRepository) GetProductByID(id uint) (*model.Product, error) {
	var product model.Product
	if err := p.MySql.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (p *productRepository) CreateProduct(req *reqSchema.CreateProductReq) error {
	var product model.Product
	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Quantity = req.Quantity

	if err := p.MySql.Create(&product).Error; err != nil {
		return err
	}
	return nil
}