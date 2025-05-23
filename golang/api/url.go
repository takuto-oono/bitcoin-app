package api

import (
	"errors"
	"net/url"
	"path"
	"strings"
)

type BitFlyerURL string
type GolangServerURL string
type DRFServerURL string

func (b BitFlyerURL) GetTicker(productCode ProductCode) (string, error) {
	qVal := url.Values{}
	if productCode != "" {
		qVal.Set("product_code", string(productCode))
	}
	return createUrl(string(b), "v1/getticker", qVal)
}

func (g GolangServerURL) GetTicker(productCode ProductCode) (string, error) {
	qVal := url.Values{}
	if productCode != "" {
		qVal.Set("product_code", string(productCode))
	}
	return createUrl(string(g), "/bitflyer/ticker", qVal)
}

func (d DRFServerURL) PostTicker() (string, error) {
	qVal := url.Values{}
	return createUrl(string(d), "/api/tickers", qVal)
}

func createUrl(baseUrl, p string, qVal url.Values, el ...string) (string, error) {
	if baseUrl == "" {
		return "", errors.New("base url is empty")
	}

	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", err
	}

	pEls := append([]string{p}, el...)
	u.Path = withSuffixSlash(path.Join(pEls...))

	u.RawQuery = qVal.Encode()

	return u.String(), nil
}

func withSuffixSlash(s string) string {
	if strings.HasSuffix(s, "/") {
		return s
	}
	return s + "/"
}

func ExtractPort(urlString string) (string, error) {
	if urlString == "" {
		return "", errors.New("url is empty")
	}
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}
	return u.Port(), nil
}
