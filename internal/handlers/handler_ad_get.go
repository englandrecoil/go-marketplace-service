package handlers

import (
	"errors"
	"net/http"

	"github.com/englandrecoil/go-marketplace-service/internal/auth"
	"github.com/englandrecoil/go-marketplace-service/internal/constants"
	"github.com/englandrecoil/go-marketplace-service/internal/database"
	"github.com/englandrecoil/go-marketplace-service/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var (
	ErrInvalidMinPrice       = errors.New("invalid min price value")
	ErrInvalidMaxPrice       = errors.New("invalid max price value")
	ErrInvalidFormatOfPrice  = errors.New("invalid format of price")
	ErrInvalidFormatOfLimit  = errors.New("invalid format of limit")
	ErrInvalidFormatOfOffset = errors.New("invalid format of offset")
)

// HandlerGetAds godoc
//
//	@Summary		Получить объявления
//	@Description	Позволяет получить объявления пользователей. Авторизованным пользователям доступно получение параметра `is_owner`.
//	@Produce		json
//	@Param			Authorization	header		string				false	"Bearer токен"							example(Bearer J2bc3Cd0F...)
//	@Param			page			query		int					false	"Номер страницы"						default(1)	minimum(1)
//	@Param			page_size		query		int					false	"Количество возвращаемых объявлений"	default(25)	minimum(1)	maximum(100)
//	@Param			min_price		query		int					false	"Минимальная цена"						minimum(0)
//	@Param			max_price		query		int					false	"Максимальная цена"						maximum(99999999)
//	@Param			sort_by			query		string				false	"Поле для сортировки"					default(created_at)	Enums(price, created_at)
//	@Param			order			query		string				false	"Направление сортировки"				default(desc)		Enums(asc, desc)
//	@Success		200				{array}		dto.GetAdsResponse	"Успешный ответ"
//	@Failure		400				{object}	dto.ErrorResponse	"Неверные параметры запроса"
//	@Failure		401				{object}	dto.ErrorResponse	"Невалидный или просроченный токен-доступа"
//	@Failure		500				{object}	dto.ErrorResponse	"Внутренняя ошибка сервера"
//	@Router			/api/ads [get]
func (cfg *ApiConfig) HandlerGetAds(c *gin.Context) {
	// if no auth header provided - simply exclude `is_owner` field from response in the future and do nothing lol
	// if auth header is provided - validate it as usual
	var userID uuid.UUID
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		token, err := auth.GetBearerToken(c)
		if err != nil {
			dto.ResponseWithError(c, http.StatusUnauthorized, err.Error(), err)
			return
		}
		userID, err = auth.ValidateJWT(token, cfg.Secret)
		if err != nil {
			dto.ResponseWithError(c, http.StatusUnauthorized, "invalid or expired access token", err)
			return
		}
	}

	query := dto.GetAdsQueryParamsRequest{}
	if err := c.BindQuery(&query); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid query parameters", err)
		return
	}

	// validate query params
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 25
	}
	if query.SortBy != "price" && query.SortBy != "created_at" {
		query.SortBy = "created_at"
	}
	if query.Order != "asc" && query.Order != "desc" {
		query.Order = "desc"
	}
	if query.MinPrice == nil {
		defaultMinPrice := constants.MinPrice
		query.MinPrice = &defaultMinPrice
	}
	if query.MaxPrice == nil {
		defaultMaxPrice := constants.MaxPrice
		query.MaxPrice = &defaultMaxPrice
	}
	if *query.MinPrice > *query.MaxPrice {
		dto.ResponseWithError(c, http.StatusBadRequest, "min_price cannot be greater than max_price", nil)
		return
	}
	offset := (query.Page - 1) * query.PageSize

	// get ads from db
	dbAds, err := cfg.DB.GetAdvertisements(
		c.Request.Context(),
		database.GetAdvertisementsParams{
			Limit:    int32(query.PageSize),
			Offset:   int32(offset),
			MinPrice: int32(*query.MinPrice),
			MaxPrice: int32(*query.MaxPrice),
			OrderBy:  query.SortBy,
			OrderDir: query.Order,
		},
	)
	if err != nil {
		dto.ResponseWithError(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	// aggregate ads from db to custom responseAds struct
	responseAds := make([]dto.GetAdsResponse, len(dbAds))
	for index, ad := range dbAds {
		var isOwner *bool
		if userID != uuid.Nil {
			isOwnerVal := ad.UserID == userID
			isOwner = &isOwnerVal
		}

		responseAd := dto.GetAdsResponse{
			Title:        ad.Title,
			Description:  ad.Description,
			ImageAddress: ad.ImageAddress,
			AuthorLogin:  ad.AuthorLogin,
			Price:        int(ad.Price),
			IsOwner:      isOwner,
		}
		responseAds[index] = responseAd
	}
	c.JSON(http.StatusOK, responseAds)
}
