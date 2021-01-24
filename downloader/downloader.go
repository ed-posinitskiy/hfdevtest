package downloader

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func CanHandle(src string) bool {
	res, _ := regexp.MatchString(`(?i)^https?://`, src)

	return res
}

func Download(url string) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "*")
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("non 200 code returned [%v], can't download the file", rsp.StatusCode)
	}

	_, err = io.Copy(tmpFile, rsp.Body)
	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}
