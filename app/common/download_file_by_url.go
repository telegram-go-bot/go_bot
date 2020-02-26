package cmn

import (
	"errors"
	"io"
	"log"
	"net/http"
)

// DownloadFileByURL you know. Returns file contents array
func DownloadFileByURL(url string) (buf []byte, err error) {
	return downloadFile(url, httpClient())
}

func downloadFile(fullURLFile string, client *http.Client) (buf []byte, err error) {
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return nil, err
	}

	if resp.ContentLength <= 0 {
		return nil, errors.New("resp.ContentLength <= 0")
	}

	buf = make([]byte, resp.ContentLength)
	bytesRead, err := io.ReadFull(resp.Body, buf)
	if err != nil {
		return nil, err
	}

	log.Printf("Just Downloaded a file %s with size %d", fullURLFile, bytesRead)
	return buf, nil
}

func httpClient() *http.Client {
	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	return &client
}
