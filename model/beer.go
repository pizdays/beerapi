package model

import (
	"gorm.io/gorm"
)

type Beer struct {
	gorm.Model
	Name        string `json:"name" gorm:"type:VARCHAR(191) NOT NULL;index" example:"John Doe"`
	Category    string `json:"category" gorm:"type:VARCHAR(191) NULL"`
	Description string `json:"description" gorm:"type:VARCHAR(191) NULL"`
	Image       string `json:"image" gorm:"type:VARCHAR(191) NULL" example:"/images"`
}

type PagingQuery struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type RespPagingQuery struct {
	Limit  int   `form:"limit"`
	Offset int   `form:"offset"`
	Total  int64 `json:"total"`
}

type ErrorResponse struct {
	Message string
}

type DefaultResponse struct {
	Message string      `json:"message" example:"success"`
	Result  interface{} `json:"result"`
}

type GetsResponse struct {
	Message     string          `json:"message" example:"success"`
	Result      interface{}     `json:"result"`
	PagingQuery RespPagingQuery `json:"paging"`
}

type Logs struct {
	ClientIP   string `bson:"clientIP"`
	Method     string `bson:"method"`
	Path       string `bson:"path"`
	StatusCode int    `bson:"statusCode"`
	TimeStamp  string `bson:"timeStamp"`
}
