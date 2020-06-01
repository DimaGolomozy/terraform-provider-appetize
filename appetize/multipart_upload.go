package appetize

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func newFileUploadRequest(uri string, path string, params map[string]*string) (*http.Request, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(abs)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)

	for key, val := range params {
		if val != nil {
			_ = writer.WriteField(key, *val)
		}
	}

	writer.Close()

	req, _ := http.NewRequest("POST", uri, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	return req, nil
}

func uploadFile(uri string, path string, params map[string]*string) (*http.Response, error) {
	request, err := newFileUploadRequest(uri, path, params)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
