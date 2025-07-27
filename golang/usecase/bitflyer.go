package usecase

import (
	"errors"
	"fmt"
	"net/http"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
)

const (
	ChildOrderTypeLimit  ChildOrderType = "LIMIT"
	ChildOrderTypeMarket ChildOrderType = "MARKET"

	SideBuy  = "BUY"
	SideSell = "SELL"

	TimeInForceGTC TimeInForce = "GTC"
	TimeInForceIOC TimeInForce = "IOC"
	TimeInForceFOK TimeInForce = "FOK"

	MinMinuteToExpire = 1
	MaxMinuteToExpire = 43200 // 30 days in minutes
)

type IBitFlyerUsecase interface {
	GetTicker(productCode string) (api.TickerFromBitFlyer, int, error)
	BuyOrder(dto BuyOrderDTO) (api.SendChildOrderResponse, int, error)
	SellOrder(dto SellOrderDTO) (api.SendChildOrderResponse, int, error)
}

type BuyOrderDTO struct {
	ProductCode    string         `json:"product_code"`
	ChildOrderType ChildOrderType `json:"child_order_type"`
	Price          float64        `json:"price"`
	Size           float64        `json:"size"`
	MinuteToExpire MinuteToExpire `json:"minute_to_expire"`
	TimeInForce    TimeInForce    `json:"time_in_force"`
	IsDry          bool           `json:"is_dry"`
}

type SellOrderDTO struct {
	ProductCode    string         `json:"product_code"`
	ChildOrderType ChildOrderType `json:"child_order_type"`
	Price          float64        `json:"price"`
	Size           float64        `json:"size"`
	MinuteToExpire MinuteToExpire `json:"minute_to_expire"`
	TimeInForce    TimeInForce    `json:"time_in_force"`
	IsDry          bool           `json:"is_dry"`
}

type BitFlyerUsecase struct {
	Config      config.Config
	BitFlyerAPI api.IBitFlyerAPI
}

func NewBitFlyerUsecase(cfg config.Config) IBitFlyerUsecase {
	return &BitFlyerUsecase{
		Config:      cfg,
		BitFlyerAPI: api.NewBitFlyerAPI(cfg),
	}
}

func (b *BitFlyerUsecase) GetTicker(productCode string) (api.TickerFromBitFlyer, int, error) {
	pc, err := api.NewProductCode(productCode)
	if err != nil {
		return api.TickerFromBitFlyer{}, http.StatusBadRequest, err
	}

	res, err := b.BitFlyerAPI.GetTicker(pc)
	if err != nil {
		return api.TickerFromBitFlyer{}, http.StatusInternalServerError, err
	}

	return res, http.StatusOK, nil
}

func (b *BitFlyerUsecase) BuyOrder(dto BuyOrderDTO) (api.SendChildOrderResponse, int, error) {
	if err := dto.ChildOrderType.validate(); err != nil {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, err
	}

	if err := dto.TimeInForce.validate(); err != nil {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, err
	}

	if err := dto.MinuteToExpire.validate(); err != nil {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, err
	}

	if dto.ChildOrderType == ChildOrderTypeLimit && dto.Price <= 0 {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, errors.New("price must be greater than 0 for LIMIT orders")
	}

	args := api.SendChildOrderRequest{
		ProductCode:    api.ProductCode(dto.ProductCode),
		ChildOrderType: string(dto.ChildOrderType),
		Side:           SideBuy,
		Price:          dto.Price,
		Size:           dto.Size,
		MinuteToExpire: int(dto.MinuteToExpire),
		TimeInForce:    string(dto.TimeInForce),
	}

	res, err := b.BitFlyerAPI.SendChildOrder(args, dto.IsDry)
	if err != nil {
		return api.SendChildOrderResponse{}, http.StatusInternalServerError, err
	}

	return res, http.StatusOK, nil
}

func (b *BitFlyerUsecase) SellOrder(dto SellOrderDTO) (api.SendChildOrderResponse, int, error) {
	if err := dto.ChildOrderType.validate(); err != nil {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, err
	}

	if err := dto.TimeInForce.validate(); err != nil {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, err
	}

	if dto.ChildOrderType == ChildOrderTypeLimit && dto.Price <= 0 {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, errors.New("price must be greater than 0 for LIMIT orders")
	}

	args := api.SendChildOrderRequest{
		ProductCode:    api.ProductCode(dto.ProductCode),
		ChildOrderType: string(dto.ChildOrderType),
		Side:           SideSell,
		Price:          dto.Price,
		Size:           dto.Size,
		MinuteToExpire: int(dto.MinuteToExpire),
		TimeInForce:    string(dto.TimeInForce),
	}

	res, err := b.BitFlyerAPI.SendChildOrder(args, dto.IsDry)
	if err != nil {
		return api.SendChildOrderResponse{}, http.StatusInternalServerError, err
	}

	return res, http.StatusOK, nil
}

type ChildOrderType string

func (c ChildOrderType) validate() error {
	switch c {
	case ChildOrderTypeLimit, ChildOrderTypeMarket:
		return nil
	default:
		return errors.New("invalid child order type")
	}
}

type TimeInForce string

func (t TimeInForce) validate() error {
	switch t {
	case TimeInForceGTC, TimeInForceIOC, TimeInForceFOK:
		return nil
	default:
		return errors.New("invalid time in force")
	}
}

type MinuteToExpire int

func (m MinuteToExpire) validate() error {
	if m < MinMinuteToExpire || m > MaxMinuteToExpire {
		return fmt.Errorf("minute to expire must be between %d and %d", MinMinuteToExpire, MaxMinuteToExpire)
	}
	return nil
}
