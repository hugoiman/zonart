package controllers

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

// Cloudinary is class
type Cloudinary struct {
	cloudName string
	baseURL   string
	apiKey    string
	apiSecret string
}

// setVariable is setter
func (c *Cloudinary) setVariable() {
	c.cloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
	c.baseURL = os.Getenv("CLOUDINARY_BASE_URL")
	c.apiKey = os.Getenv("CLOUDINARY_API_KEY")
	c.apiSecret = os.Getenv("CLOUDINARY_API_SECRET")
}

func (c Cloudinary) UploadImages(r *http.Request, maxSize int64, destinationFolder string) ([]string, error) {
	var cloud CloudTest
	var images = make([]string, 0)

	c.setVariable()
	cld, err := cloudinary.NewFromParams(c.cloudName, c.apiKey, c.apiSecret)
	if err != nil {
		return images, errors.New("Failed to intialize Cloudinary: " + err.Error())
	}

	var ctx = context.Background()
	const maxMemory = 32 << 20
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return images, err
	}

	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {
			if hdr.Size > maxSize {
				c.DeleteImages(images)
				return images, errors.New("Maksimal " + strconv.Itoa(int(maxSize/1024/1024)) + " MB per foto")
			}
			cloud.Gambar = append(cloud.Gambar, Gambar{Image: hdr.Filename})

			var infile multipart.File
			if infile, err = hdr.Open(); nil != err {
				fmt.Println("Cant open file, ", err)
				c.DeleteImages(images)
				return images, errors.New("Cant open file " + err.Error())
			}

			uploadResult, err := cld.Upload.Upload(ctx, infile, uploader.UploadParams{Folder: destinationFolder})
			if err != nil {
				fmt.Println("Failed to upload file, ", err)
				infile.Close()
				c.DeleteImages(images)
				return images, err
			}
			images = append(images, uploadResult.SecureURL)
			infile.Close()
		}
	}
	return images, nil
}

// stub UploadImages
// func (c Cloudinary) UploadImages(r *http.Request, maxSize int64, destinationFolder string) ([]string, error) {
// 	return []string{}, nil
//	return []string{"tes.jpg"}, errors.New("Terjadi error")
// }

// DeleteImages is func
func (c Cloudinary) DeleteImages(images []string) error {
	var data, filename, extension, public_ids string
	for _, url := range images {
		filename = url[strings.Index(url, "zonart"):]
		extension = filepath.Ext(filename)
		public_ids = filename[0 : len(filename)-len(extension)]
		data += "public_ids[]=" + public_ids + "&"
	}

	c.setVariable()
	url := "https://" + c.apiKey + ":" + c.apiSecret + "@" + strings.ReplaceAll(c.baseURL, "https://", "") + "/resources/image/upload"
	payload := strings.NewReader(data)
	req, _ := http.NewRequest("DELETE", url, payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	// body, _ := ioutil.ReadAll(res.Body)
	return err
}
