package controllers

import (
	"fmt"
	"gcom-backend/models"
	"gcom-backend/responses"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

func UploadImage(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	defer src.Close()

	match, err := regexp.MatchString("(\\d{10})", file.Filename)
	if err != nil || !match {
		return c.JSON(http.StatusBadRequest, "Invalid image name, use UNIX timestamp")
	}
	dst, err := os.Create("imgs/" + file.Filename)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error saving image")
	}
	defer dst.Close()

	timestamp, _ := strconv.Atoi(file.Filename[:10])
	fmt.Println(file.Filename[:13])
	fmt.Println(timestamp)
	image := &models.Image{
		Timestamp: int64(timestamp),
		Filename:  file.Filename,
	}

	db, _ := c.Get("db").(*gorm.DB)
	if createErr := db.Create(&image).Error; createErr != nil {
		return c.JSON(http.StatusInternalServerError, createErr.Error())
	}

	return c.JSON(http.StatusAccepted, "Upload sucessful")
}

func ListImages(c echo.Context) error {
	var images []models.Image
	db, _ := c.Get("db").(*gorm.DB)
	if err := db.Find(&images).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Message: "Error whilst querying waypoints!",
			Data:    err.Error()})
	}

	return c.JSON(http.StatusOK, images)
}

func GetImage(c echo.Context) error {
	return nil
}

func GetImageDetails(c echo.Context) error {
	return nil
}
