package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func newfileUploadRequest(path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/photos", *simpicUrl), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}

func upload(client http.Client, source <-chan string, results chan<- bool) {
	for {
		select {
		case file, more := <-source:
			if more {
				request, err := newfileUploadRequest(file)
				if err != nil {
					log.Printf("Unable to create request to upload '%s': %v\n", file, err)
					results <- false
					continue
				}

				res, err := client.Do(request)
				if err != nil {
					log.Printf("Unable to upload '%s': %v\n", file, err)
					results <- false
					continue
				}

				if res.StatusCode != http.StatusOK {
					log.Printf("Upload failed for '%s'. Server responded: %s\n", file, res.Status)
					results <- false
					continue
				}

				_ = res.Body.Close()

				results <- true
			} else {
				return
			}
		}
	}
}
