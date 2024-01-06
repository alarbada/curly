package curly

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type Curly struct {
	req *http.Request
	err error
}

func New(method, path string) *Curly {
	req, err := http.NewRequest(method, path, nil)
	return &Curly{req, err}
}

func Base(path string) *Curly {
	url := &url.URL{Path: path}
	return &Curly{req: &http.Request{URL: url}}
}

func (this *Curly) Header(key, value string) *Curly {
	if err := this.err; err != nil {
		return this
	}
	this.req.Header.Add(key, value)
	return this
}

func (this *Curly) Body(input any) *Curly {
	bs, err := json.Marshal(input)
	if err != nil {
		this.err = err
		return this
	}

	this.req.Body = io.NopCloser(bytes.NewReader(bs))

	return this
}

func (this *Curly) Do(output any) error {
	if err := this.err; err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(this.req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, output)
}
