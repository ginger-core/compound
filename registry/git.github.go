package registry

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ginger-core/errors"
)

type githubReader struct {
	*base
}

func newGithubReader(base *base) reader {
	return &githubReader{base: base}
}

func (c *githubReader) read() errors.Error {
	url := os.Getenv("CONFIG_URL")
	if url == "" {
		panic("invalid config remote base url")
	}
	token := os.Getenv("CONFIG_TOKEN")
	if token == "" {
		panic("invalid config remote token")
	}
	//
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.New(err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
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
	err = c.viper.ReadConfig(bytes.NewBuffer(data))
	if err != nil {
		return errors.New(err)
	}
	return nil
}
