package registry

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ginger-core/errors"
)

type gitlabReader struct {
	*base
}

func newGitlabReader(base *base) reader {
	return &gitlabReader{base: base}
}

type remoteResponse struct {
	Content string `json:"content"`
}

func (c *gitlabReader) read() errors.Error {
	url := os.Getenv("CONFIG_URL")
	if url == "" {
		panic("invalid config remote base url")
	}
	ref := os.Getenv("CONFIG_REF")
	if ref == "" {
		ref = "master"
	}
	token := os.Getenv("CONFIG_TOKEN")
	if token == "" {
		panic("invalid config remote token")
	}
	//
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?ref=%s", url, ref), nil)
	if err != nil {
		return errors.New(err)
	}
	req.Header.Set("PRIVATE-TOKEN", token)
	client := &http.Client{Timeout: time.Second * 5}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.New(err)
	}

	if resp.StatusCode != 200 {
		return errors.New().
			WithMessage(fmt.Sprintf("config server returned unexpected status %s. "+
				"err %v", resp.Status, string(data)))
	}

	respBody := new(remoteResponse)
	if err := json.Unmarshal(data, respBody); err != nil {
		return errors.New(err)
	}

	data, err = base64.StdEncoding.DecodeString(respBody.Content)
	if err != nil {
		return errors.New(err)
	}

	err = c.viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return errors.New(err)
	}
	return nil
}
