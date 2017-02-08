package packages

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

const (
	indexAddr  = "https://packages.solus-project.com/shannon/eopkg-index.xml"
	verifyAddr = "https://packages.solus-project.com/shannon/eopkg-index.xml.sha1sum"
)

func Download(indexAddr, verifyAddr, cacheDir string) error {
	// create cache directory if it doesn't exist
	err := os.MkdirAll(cacheDir, 0700)
	if err != nil {
		return err
	}

	indexPath := filepath.Join(cacheDir, "eopkg-index.xml")
	tempPath := indexPath + ".tmp"
	defer os.Remove(tempPath)

	// create/truncate temp file
	tempFile, err := os.Create(tempPath)
	defer tempFile.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, indexAddr, nil)
	if err != nil {
		return err
	}

	info, err := os.Stat(indexPath)
	if err == nil {
		// Update file only if changed
		time := info.ModTime().Format(http.TimeFormat)
		req.Header.Add("If-Modified-Since", time)
	}

	// download the package list
	status, err := consumeReq(req, tempFile)
	if err != nil || status != http.StatusOK {
		return err
	}

	// check the package list integrity
	verified, err := verifyIndex(verifyAddr, tempFile)
	if err != nil {
		return err
	}
	if !verified {
		return errors.New("index failed verification")
	}

	err = tempFile.Close()
	if err != nil {
		return err
	}

	return os.Rename(tempPath, indexPath)
}

func consumeReq(req *http.Request, dest io.Writer) (statusCode int, err error) {
	resp, err := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	if err != nil || statusCode != http.StatusOK {
		return
	}
	_, err = io.Copy(dest, resp.Body)
	return
}

func verifyIndex(url string, file io.ReadSeeker) (bool, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return false, err
	}

	resp, err := http.Get(verifyAddr)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	checksum, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return false, err
	}

	decoded := make([]byte, 20)
	hex.Decode(decoded, checksum)

	return bytes.Equal(h.Sum(nil), decoded), nil
}
