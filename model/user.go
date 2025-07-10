package model

type User struct {
	BaseModel
	Username string `gorm:"size:191;uniqueIndex" json:"username"`
	Password string `json:"password"`
	FullName string `gorm:"type:varchar(100)" json:"full_name"`
}
