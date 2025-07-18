package dto

type CredentialsRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AdvertisementRequest struct {
	Title        string `json:"title" binding:"required"`
	Description  string `json:"description" binding:"required"`
	ImageAddress string `json:"image_address" binding:"required"`
	Price        int    `json:"price" binding:"required"`
}
