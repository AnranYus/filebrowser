package http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

var downloadFileToRemoteHandle = withUser(func(w http.ResponseWriter, r *http.Request, d *data) (int, error) {

	url := r.URL.Query().Get("targetUrl")

	if len(url) < 0 {
		err := errors.New("can not download file without target url")
		return errToStatus(err), err
	}

	err := d.RunHook(func() error {
		response, err := http.Get(url)
		if err != nil {
			return err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(response.Body)

		lastIndex := strings.LastIndex(url, "/")
		fileName := url[lastIndex+1:]

		outFile, err := os.Create(fileName)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer func(outFile *os.File) {
			err := outFile.Close()
			if err != nil {
				fmt.Println(err)
			}
		}(outFile)

		_, err = io.Copy(outFile, response.Body)
		if err != nil {
			return err
		}

		return nil
	}, "download", r.URL.Path, url, d.user)

	return errToStatus(err), err
})
