package registry

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ginger-core/errors"
)

type git struct {
	*base
	baseUrl string
	path    string
	ref     string
	token   string
}

func newGit(ctx context.Context, format string, args ...interface{}) (Registry, errors.Error) {
	if len(args) < 4 {
		return nil, errors.Internal().
			WithContext(ctx).
			WithId("newGit.arg").
			WithMessage("invalid arguments")
	}

	c := &git{
		base:    _new(format),
		baseUrl: args[0].(string),
		path:    args[1].(string),
		ref:     args[2].(string),
		token:   args[3].(string),
	}
	if err := c.read(); err != nil {
		return nil, err
	}
	return c, nil
}

type remoteResponse struct {
	Content string `json:"content"`
}

func (c *git) read() errors.Error {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s?ref=%s", c.baseUrl, c.path, c.ref), nil)
	if err != nil {
		return errors.New(err)
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)
	client := &http.Client{Timeout: time.Second * 5}

	resp, err := client.Do(req)
	if err != nil {
		return errors.New(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
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
