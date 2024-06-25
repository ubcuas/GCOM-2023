package controllers

import (
	"fmt"
	"gcom-backend/models"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var imgDirectory = "./imgs/"

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

	dst, err := os.Create(imgDirectory + file.Filename)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Error saving image")
	}
	defer dst.Close()

	timestamp, _ := strconv.Atoi(file.Filename[:14])
	fmt.Println(file.Filename[:14])
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
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, images)
}

func GetImage(c echo.Context) error {
	imgPath, err := os.Stat(imgDirectory + c.Param("filename"))
	fmt.Println(imgPath.Name())
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	} else {
		return c.Attachment(imgDirectory+c.Param("filename"), c.Param("filename"))
	}
}
