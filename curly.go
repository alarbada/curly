package curly

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

type Curly struct {
	req *http.Request
	err error
}

func New(method, path string) *Curly {
	req, err := http.NewRequest(method, path, nil)

	return &Curly{req, err}
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

	rawOutput := make(map[string]any)
	if err := json.Unmarshal(body, &rawOutput); err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		// Prevents mistakes with typos in struct tags / field names instead of
		// just putting a useless 0 value
		ErrorUnset: true,
		Result:     output,
		TagName:    "json",
	})
	if err != nil {
		return err
	}

	return decoder.Decode(rawOutput)
}
