package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Picture struct {
	Photos []struct {
		ID     int `json:"id"`
		Sol    int `json:"sol"`
		Camera struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			RoverID  int    `json:"rover_id"`
			FullName string `json:"full_name"`
		} `json:"camera"`
		ImgSrc    string `json:"img_src"`
		EarthDate string `json:"earth_date"`
		Rover     struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			LandingDate string `json:"landing_date"`
			LaunchDate  string `json:"launch_date"`
			Status      string `json:"status"`
		} `json:"rover"`
	} `json:"photos"`
}

func main() {
	router := gin.Default()
	router.GET("/pictures/:date", MakeHttpRequest)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Error while launching the server: %v", err)
	}
}

func MakeHttpRequest(c *gin.Context) {
	date := c.Param("date")
	url := "https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos?earth_date=" + date + "&api_key=qvFPrKOzjjBt2FBT6uZfRyTkwWXxFxKsnjNIgsNC"
	response, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error while making the GET request: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Error while reading the response: %v", err)
	}
	defer response.Body.Close()

	var picture Picture
	err = json.Unmarshal([]byte(body), &picture)
	if err != nil {
		log.Fatalf("Error while unmarshalling the json: %v", err)
	}

	if len(picture.Photos) > 0 {
		err := DownloadPictures(picture, date)

		if err != nil {
			log.Fatalf("Error while downloading the pictures: %v", err)
		}

		c.JSON(http.StatusOK, "You downloaded "+strconv.Itoa(len(picture.Photos))+" pictures!")
	} else {
		c.JSON(http.StatusOK, "There were no new pictures to download")
	}
}
func DownloadPictures(picture Picture, date string) error {
	err := os.Mkdir(date, os.ModePerm)
	if err != nil {
		log.Printf("Error while creating the folder: %v", err)
	}
	err = os.Chdir(date)
	if err != nil {
		log.Printf("Error while changing the folder: %v", err)
		return err
	}

	for _, a := range picture.Photos {
		resp, err := http.Get(a.ImgSrc)

		if err != nil {
			log.Printf("Error while getting the image: %v", err)
			return err
		}
		defer resp.Body.Close()

		out, err := os.Create(strconv.Itoa(a.ID) + ".png")
		if err != nil {
			log.Printf("Error while creating the .png file: %v", err)
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Printf("Error while copying the contents into the file: %v", err)
			return err
		}
	}
	err = os.Chdir("..")
	if err != nil {
		log.Printf("Error while changing the folder: %v", err)
	}
	return nil
}
