package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"unicode/utf8"

	"github.com/englandrecoil/go-marketplace-service/internal/auth"
	"github.com/englandrecoil/go-marketplace-service/internal/database"
	"github.com/englandrecoil/go-marketplace-service/internal/dto"
	"github.com/gin-gonic/gin"
)

const (
	minTitleLength = 1
	maxTitleLength = 50
	minDescLength  = 10
	maxDescLength  = 750
	minPrice       = 1
	maxPrice       = 99999999
	maxImageSize   = 10 * 1024 * 1024
)

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

	inputAdParams := dto.AdvertisementRequest{}
	if err := c.BindJSON(&inputAdParams); err != nil {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid request body format", err)
		return
	}

	// validation part
	if utf8.RuneCountInString(inputAdParams.Title) < minTitleLength || utf8.RuneCountInString(inputAdParams.Title) > maxTitleLength {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid length of title", err)
		return
	}
	if utf8.RuneCountInString(inputAdParams.Description) < minDescLength || utf8.RuneCountInString(inputAdParams.Description) > maxDescLength {
		dto.ResponseWithError(c, http.StatusBadRequest, "invalid length of description", err)
		return
	}
	if inputAdParams.Price < minPrice || inputAdParams.Price > maxPrice {
		dto.ResponseWithError(c, http.StatusBadRequest, "incorrect price", err)
		return
	}
	if err = validateImage(inputAdParams.ImageAddress); err != nil {
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
		dto.AdvertisementResponse{
			ID:           ad.ID.String(),
			Title:        ad.Title,
			Description:  ad.Description,
			ImageAddress: ad.ImageAddress,
			Price:        int(ad.Price),
			CreatedAt:    ad.CreatedAt,
		},
	)
}

func validateImage(imageUrl string) error {
	if _, err := url.Parse(imageUrl); err != nil {
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
	if res.ContentLength > maxImageSize {
		return errors.New("image size is too big")
	}
	return nil
}
