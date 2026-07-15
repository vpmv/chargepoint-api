package api

import (
	"github.com/go-fuego/fuego"
	"github.com/vpmv/chargepoint-api/internal/dto"
	"github.com/vpmv/chargepoint-api/internal/helpers"
)

type LocationParams struct {
	Latitude  float64 `query:"lat" validate:"numeric,min=-90,max=90,required"`
	Longitude float64 `query:"lon" validate:"numeric,min=-180,max=180,required"`
	Radius    int     `query:"radius" validate:"min=5,max=100,required"`
}

func (api *API) CreateChargePoints(c fuego.ContextWithBody[dto.ChargePoints]) (*dto.ChargePoints, error) {
	chargePoints, err := c.Body()
	if err != nil {
		return nil, err
	}

	cps, err := api.store.CreateChargePoints(chargePoints)
	if err != nil {
		return nil, err
	}

	return &cps, nil
}

func (api *API) GetChargePointByID(c fuego.ContextNoBody) (*dto.ChargePoint, error) {
	id := c.PathParam("id")

	cp, err := api.store.GetChargePoint(id)
	if err != nil {
		return nil, err
	} else if cp == nil {
		return nil, fuego.NotFoundError{}
	}

	return cp, nil
}

func (api *API) ListChargePoints(c fuego.ContextNoBody) (*dto.ChargePoints, error) {
	page := c.QueryParamInt(ParamPage)
	pageSize := c.QueryParamInt(ParamPageSize)

	cps, err := api.store.GetChargePoints(page, pageSize)
	if err != nil {
		return nil, err
	}

	return &cps, nil
}

func (api *API) ListChargePointsByLocation(c fuego.ContextWithParams[LocationParams]) (*dto.ChargePoints, error) {
	params, err := c.Params()
	if err != nil {
		return nil, fuego.BadRequestError{Err: err}
	}
	if err := helpers.ValidateStruct(params); err != nil {
		return nil, fuego.BadRequestError{Err: err} // fixme
	}

	cps, err := api.store.GetChargePointByLocation(params.Latitude, params.Longitude, params.Radius)
	if err != nil {
		return nil, fuego.BadRequestError{Err: err}
	}

	return &cps, nil

}
