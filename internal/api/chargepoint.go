package api

import (
	"github.com/go-fuego/fuego"
	"github.com/vpmv/chargepoint-api/internal/dto"
)

func (api *API) ListChargePoints(c fuego.ContextNoBody) (*dto.ChargePoints, error) {
	page := c.QueryParamInt(ParamPage)
	pageSize := c.QueryParamInt(ParamPageSize)

	return api.store.GetChargePoints(page, pageSize)
}

func (api *API) GetChargePointByID(c fuego.ContextNoBody) (*dto.ChargePoint, error) {
	id := c.PathParam("id")
	return &dto.ChargePoint{ID: id}, nil
}

func (api *API) CreateChargePoints(c fuego.ContextWithBody[dto.ChargePoints]) (*dto.ChargePoints, error) {
	chargePoints, err := c.Body()
	if err != nil {
		return nil, err
	}

	return api.store.CreateChargePoints(&chargePoints)
}
