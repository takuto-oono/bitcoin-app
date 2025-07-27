package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"bitcoin-app-golang/config"
)

const BitFlyerBaseURL = "https://api.bitflyer.com"

type IBitFlyerAPI interface {
	GetTicker(ProductCode) (TickerFromBitFlyer, error)
	SendChildOrder(SendChildOrderRequest, bool) (SendChildOrderResponse, error)
}

type BitFlyerAPI struct {
	Config config.Config
	API    *API
}

func NewBitFlyerAPI(cfg config.Config) IBitFlyerAPI {
	return &BitFlyerAPI{
		Config: cfg,
		API:    NewAPI(),
	}
}

func (b *BitFlyerAPI) GetTicker(productCode ProductCode) (TickerFromBitFlyer, error) {
	url, err := BitFlyerURL(BitFlyerBaseURL).GetTicker(productCode)
	if err != nil {
		return TickerFromBitFlyer{}, err
	}

	resModel := TickerFromBitFlyer{}
	if err := b.API.Do(http.MethodGet, nil, &resModel, url, nil); err != nil {
		return TickerFromBitFlyer{}, err
	}
	return resModel, nil
}

func (b *BitFlyerAPI) SendChildOrder(args SendChildOrderRequest, isDry bool) (SendChildOrderResponse, error) {
	url, err := BitFlyerURL(BitFlyerBaseURL).SendChildOrder()
	if err != nil {
		return SendChildOrderResponse{}, err
	}

	body, err := json.Marshal(args)
	if err != nil {
		return SendChildOrderResponse{}, err
	}

	authHeaders, err := b.privateRequestHeader(nowUnixTimestamp(), http.MethodPost, url, body)
	if err != nil {
		return SendChildOrderResponse{}, err
	}

	resModel := SendChildOrderResponse{}

	if isDry {
		log.Default().Println("Dry run: SendChildOrder is not executed")
	} else {
		if err := b.API.Do(http.MethodPost, args, &resModel, url, authHeaders); err != nil {
			return SendChildOrderResponse{}, err
		}
	}

	return resModel, nil
}

// https://lightning.bitflyer.com/docs#%E8%AA%8D%E8%A8%BC:~:text=%E4%BA%86%E6%89%BF%E3%81%8F%E3%81%A0%E3%81%95%E3%81%84%E3%80%82-,%E8%AA%8D%E8%A8%BC,-Private%20API%20%E3%81%AE
func (api *BitFlyerAPI) privateRequestHeader(timeStamp, method, url string, body []byte) (map[string]any, error) {
	path, err := extractPath(url)
	if err != nil {
		return nil, err
	}

	rawMessage := json.RawMessage(body)

	text := timeStamp + method + path + string(rawMessage)

	h := hmac.New(sha256.New, []byte(api.Config.BitFlyer.ApiSecret))
	h.Write([]byte(text))
	sign := hex.EncodeToString(h.Sum(nil))

	return map[string]any{
		"ACCESS-KEY":       api.Config.BitFlyer.ApiKey,
		"ACCESS-TIMESTAMP": timeStamp,
		"ACCESS-SIGN":      sign,
	}, nil
}

func extractPath(u string) (string, error) {
	uObj, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return uObj.Path, nil
}

func nowUnixTimestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}
