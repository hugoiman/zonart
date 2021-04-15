package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
)

// Cloudinary is class
type Cloudinary struct {
	baseURL   string
	apiKey    string
	apiSecret string
}

// setVariable is setter
func (c *Cloudinary) setVariable() {
	c.baseURL = "*"
	c.apiKey = "*"
	c.apiSecret = "*"
}

// DeleteImages is func
func (c Cloudinary) DeleteImages(images []string) error {
	fmt.Println(images)
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
