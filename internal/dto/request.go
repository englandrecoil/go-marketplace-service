package dto

type CredentialsRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateAdsRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description" binding:"required"`
	ImageAddress string `json:"image_address" binding:"required"`
	Price        int    `json:"price" binding:"required"`
}

type GetAdsQueryParamsRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	MinPrice *int   `form:"min_price"`
	MaxPrice *int   `form:"max_price"`
	SortBy   string `form:"sort_by"`
	Order    string `form:"order"`
}
