package api

import (
	"bitcoin-app-golang/config"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type API struct{}

func NewAPI() *API {
	return &API{}
}

func (api *API) Do(method string, reqModel, resModel any, url string, headerMap map[string]any) error {
	reqJson, err := marshalJson(reqModel)
	if err != nil {
		return err
	}

	res, err := request(method, url, reqJson, convertToStringMap(headerMap))
	if err != nil {
		return err
	}

	resJson, err := readResponse(res)
	if err != nil {
		return err
	}

	if resModel == nil {
		return nil
	}

	return json.Unmarshal(resJson, resModel)
}

func request(method, url string, body []byte, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	return http.DefaultClient.Do(req)
}

func readResponse(resp *http.Response) ([]byte, error) {
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		fmt.Printf("status code: %d, body: %s", resp.StatusCode, string(body))
		return nil, errors.New(string(body))
	}

	return body, nil
}

func marshalJson(v any) ([]byte, error) {
	if v == nil || v == "" {
		return []byte{}, nil
	}
	return json.Marshal(v)
}

func convertToStringMap(m map[string]any) map[string]string {
	if m == nil {
		return nil
	}
	stringMap := make(map[string]string, len(m))
	for k, v := range m {
		switch val := v.(type) {
		case string:
			stringMap[k] = val
		case config.Credential:
			stringMap[k] = string(val)
		}
	}
	return stringMap
}
