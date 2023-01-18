package btfs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

type service struct {
	apiHost string
	addPath string
	catPath string
}

func newService(apiHost, addPath, catPath string) *service {
	return &service{
		apiHost: apiHost,
		addPath: addPath,
		catPath: catPath,
	}
}

func (s *service) parseAddUrl(options *AddOptions) (uploadUrl string, err error) {
	uri, err := url.Parse(s.apiHost)
	if err != nil {
		return
	}
	uri.Path = s.addPath
	queries := uri.Query()
	if options.Pin {
		queries.Set("pin", "true")
	}
	uri.RawQuery = queries.Encode()
	uploadUrl = uri.String()
	return
}

func (s *service) Add(ctx context.Context, content []byte, name string, options *AddOptions) (result *AddResult, err error) {
	addUri, err := s.parseAddUrl(options)
	if err != nil {
		return
	}

	body := &bytes.Buffer{}
	multi := multipart.NewWriter(body)

	part, err := multi.CreateFormFile("file", name)
	if err != nil {
		return
	}

	_, err = io.Copy(part, bytes.NewBuffer(content))
	if err != nil {
		return
	}

	err = multi.Close()
	if err != nil {
		return
	}

	results, err := s.doAddRequest(
		ctx, addUri,
		multi.FormDataContentType(),
		body,
	)

	if err != nil {
		return
	}
	if len(results) < 1 {
		err = errors.New("no result")
		return
	}

	result = results[0]
	return
}

func (s *service) doAddRequest(ctx context.Context, addUri string, contentType string, body io.Reader) (results []*AddResult, err error) {
	req, err := http.NewRequest("POST", addUri, body)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", contentType)

	timeOut := 10 * time.Minute
	if deadline, ok := ctx.Deadline(); ok {
		timeOut = time.Until(deadline)
	}
	if timeOut <= 0 {
		err = errors.New("request timeout")
		return
	}

	client := &http.Client{
		Timeout: timeOut,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var result AddResult
		err = decoder.Decode(&result)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		results = append(results, &result)
	}
	return
}

func (s *service) parseCatUrl(hash string) (uploadUrl string, err error) {
	uri, err := url.Parse(s.apiHost)
	if err != nil {
		return
	}
	uri.Path = s.catPath
	queries := uri.Query()
	queries.Set("arg", hash)
	uri.RawQuery = queries.Encode()
	uploadUrl = uri.String()
	return
}

func (s *service) Cat(ctx context.Context, hash string) (content []byte, err error) {
	catUri, err := s.parseCatUrl(hash)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", catUri, nil)
	if err != nil {
		return
	}

	timeOut := 10 * time.Minute
	if deadline, ok := ctx.Deadline(); ok {
		timeOut = time.Until(deadline)
	}
	if timeOut <= 0 {
		err = errors.New("request timeout")
		return
	}

	client := &http.Client{
		Timeout: timeOut,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	content, err = io.ReadAll(resp.Body)
	return
}
