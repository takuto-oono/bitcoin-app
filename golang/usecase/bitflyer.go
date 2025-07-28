package usecase

import (
	"errors"
	"fmt"
	"net/http"

	"bitcoin-app-golang/api"
	"bitcoin-app-golang/config"
	"bitcoin-app-golang/consts"
)

type IBitFlyerUsecase interface {
	GetTicker(productCode string) (api.TickerFromBitFlyer, int, error)
	BuyOrder(dto BuyOrderDTO) (api.SendChildOrderResponse, int, error)
	SellOrder(dto SellOrderDTO) (api.SendChildOrderResponse, int, error)
}

type BuyOrderDTO struct {
	ProductCode    ProductCode    `json:"product_code"`
	ChildOrderType ChildOrderType `json:"child_order_type"`
	Price          float64        `json:"price"`
	Size           float64        `json:"size"`
	MinuteToExpire MinuteToExpire `json:"minute_to_expire"`
	TimeInForce    TimeInForce    `json:"time_in_force"`
	IsDry          bool           `json:"is_dry"`
}

type SellOrderDTO struct {
	ProductCode    ProductCode    `json:"product_code"`
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
	pc, err := NewProductCode(productCode)
	if err != nil {
		return api.TickerFromBitFlyer{}, http.StatusBadRequest, err
	}

	res, err := b.BitFlyerAPI.GetTicker(string(pc))
	if err != nil {
		return api.TickerFromBitFlyer{}, http.StatusInternalServerError, err
	}

	return res, http.StatusOK, nil
}

func (b *BitFlyerUsecase) BuyOrder(dto BuyOrderDTO) (api.SendChildOrderResponse, int, error) {
	if err := validateBuyOrSellOrder(dto); err != nil {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, err
	}

	args := api.SendChildOrderRequest{
		ProductCode:    string(dto.ProductCode),
		ChildOrderType: string(dto.ChildOrderType),
		Side:           consts.SideBuy,
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
	if err := validateBuyOrSellOrder(dto); err != nil {
		return api.SendChildOrderResponse{}, http.StatusBadRequest, err
	}

	args := api.SendChildOrderRequest{
		ProductCode:    string(dto.ProductCode),
		ChildOrderType: string(dto.ChildOrderType),
		Side:           consts.SideSell,
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

func validateBuyOrSellOrder(dto any) error {
	switch v := dto.(type) {
	case BuyOrderDTO:
		if err := v.ProductCode.validate(); err != nil {
			return err
		}
		if err := v.ChildOrderType.validate(); err != nil {
			return err
		}
		if err := v.TimeInForce.validate(); err != nil {
			return err
		}
		if err := v.MinuteToExpire.validate(); err != nil {
			return err
		}
		if v.ChildOrderType == consts.ChildOrderTypeLimit && v.Price <= 0 {
			return errors.New("price must be greater than 0 for LIMIT orders")
		}
	case SellOrderDTO:
		if err := v.ProductCode.validate(); err != nil {
			return err
		}
		if err := v.ChildOrderType.validate(); err != nil {
			return err
		}
		if err := v.TimeInForce.validate(); err != nil {
			return err
		}
		if err := v.MinuteToExpire.validate(); err != nil {
			return err
		}
		if v.ChildOrderType == consts.ChildOrderTypeLimit && v.Price <= 0 {
			return errors.New("price must be greater than 0 for LIMIT orders")
		}
	default:
		return fmt.Errorf("unsupported order type: %T", dto)
	}
	return nil
}

type ProductCode string

func (p ProductCode) validate() error {
	switch p {
	case consts.ProductCodeBTCJPY, consts.ProductCodeXRPJPY, consts.ProductCodeETHJPY, consts.ProductCodeXLMJPY, consts.ProductCodeMONAJPY,
		consts.ProductCodeETHBTC, consts.ProductCodeBCHBTC, consts.ProductCodeFXBTCJPY:
		return nil
	default:
		return fmt.Errorf("invalid product code: %s", p)
	}
}

func NewProductCode(code string) (ProductCode, error) {
	pc := ProductCode(code)
	if err := pc.validate(); err != nil {
		return "", err
	}
	return pc, nil
}

type ChildOrderType string

func (c ChildOrderType) validate() error {
	switch c {
	case consts.ChildOrderTypeLimit, consts.ChildOrderTypeMarket:
		return nil
	default:
		return errors.New("invalid child order type")
	}
}

type TimeInForce string

func (t TimeInForce) validate() error {
	switch t {
	case consts.TimeInForceGTC, consts.TimeInForceIOC, consts.TimeInForceFOK:
		return nil
	default:
		return errors.New("invalid time in force")
	}
}

type MinuteToExpire int

func (m MinuteToExpire) validate() error {
	if m < consts.MinMinuteToExpire || m > consts.MaxMinuteToExpire {
		return fmt.Errorf("minute to expire must be between %d and %d", consts.MinMinuteToExpire, consts.MaxMinuteToExpire)
	}
	return nil
}
