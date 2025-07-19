package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"unicode/utf8"

	"github.com/englandrecoil/go-marketplace-service/internal/auth"
	"github.com/englandrecoil/go-marketplace-service/internal/constants"
	"github.com/englandrecoil/go-marketplace-service/internal/database"
	"github.com/englandrecoil/go-marketplace-service/internal/dto"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidLengthTitle       = errors.New("invalid length of title")
	ErrInvalidLengthDescription = errors.New("invalid length of description")
	ErrInvalidPrice             = errors.New("incorrect price")
)

// HandlerCreateAd godoc
//
//	@Summary		Создать новое объявление
//	@Description	Создаёт новое объявление с заданными параметрами
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			Authorization	header		string					true	"Bearer токен"	example(Bearer J2bc3Cd0F...)
//	@Param			body			body		dto.CreateAdsRequest	true	"Параметры объявления"
//	@Success		201				{object}	dto.CreateAdsResponse	"Успешное создание объявления"
//	@Failure		400				{object}	dto.ErrorResponse		"Неверный формат запроса"
//	@Failure		401				{object}	dto.ErrorResponse		"Невалидный или просроченный токен-доступа"
//	@Failure		500				{object}	dto.ErrorResponse		"Внутренняя ошибка сервера"
//	@Router			/api/ads [post]
func (cfg *ApiConfig) HandlerCreateAd(c *gin.Context) {
	token, err := auth.GetBearerToken(c)
	if err != nil {
		dto.ResponseWithError(c, http.StatusUnauthorized, err.Error(), err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.Secret)
	if err != nil {
		dto.ResponseWithError(c, http.StatusUnauthorized, "invalid or expired access token", err)
		return
	}

	inputAdParams := dto.CreateAdsRequest{}
	if err := c.BindJSON(&inputAdParams); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid request body format", err)
		return
	}

	// validation part
	err = validateAdParams(
		inputAdParams.Title,
		inputAdParams.Description,
		inputAdParams.ImageAddress,
		inputAdParams.Price,
	)
	if err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	// create new record of ad in db
	ad, err := cfg.DB.CreateAdvertisement(
		c.Request.Context(),
		database.CreateAdvertisementParams{
			Title:        inputAdParams.Title,
			Description:  inputAdParams.Description,
			ImageAddress: inputAdParams.ImageAddress,
			Price:        int32(inputAdParams.Price),
			CreatedAt:    time.Now().UTC(),
			UpdatedAt:    time.Now().UTC(),
			UserID:       userID,
		},
	)
	if err != nil {
		dto.ResponseWithError(c, http.StatusInternalServerError, "internal server error", err)
		return
	}

	c.JSON(
		http.StatusCreated,
		dto.CreateAdsResponse{
			ID:           ad.ID.String(),
			Title:        ad.Title,
			Description:  ad.Description,
			ImageAddress: ad.ImageAddress,
			Price:        int(ad.Price),
			CreatedAt:    ad.CreatedAt,
		},
	)
}

func validateAdParams(title, description, imageUrl string, price int) error {
	if utf8.RuneCountInString(title) < constants.MinTitleLength || utf8.RuneCountInString(title) > constants.MaxTitleLength {
		return ErrInvalidLengthTitle
	}
	if utf8.RuneCountInString(description) < constants.MinDescLength || utf8.RuneCountInString(description) > constants.MaxDescLength {
		return ErrInvalidLengthDescription
	}
	if price < constants.MinPrice || price > constants.MaxPrice {
		return ErrInvalidPrice
	}
	if err := validateImage(imageUrl); err != nil {
		return err
	}
	return nil
}

func validateImage(imageUrl string) error {
	if _, err := url.ParseRequestURI(imageUrl); err != nil {
		return err
	}

	req, err := http.NewRequest("HEAD", imageUrl, nil)
	if err != nil {
		return fmt.Errorf("couldn't get image metadata: %s", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("can't send request to server: %s", err)
	}
	defer res.Body.Close()

	imageType := res.Header.Get("content-type")
	if imageType != "image/jpeg" && imageType != "image/jpg" && imageType != "image/png" {
		return errors.New("invalid image format")
	}
	if res.ContentLength > constants.MaxImageSize {
		return errors.New("image size is too big")
	}
	return nil
}
