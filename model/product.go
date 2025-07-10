package model

type Product struct {
	BaseModel
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Price       int64  `gorm:"not null" json:"price"`
	Description string `gorm:"type:text" json:"description"`
	Quantity    int    `gorm:"type:int;not null" json:"quantity"`
}
