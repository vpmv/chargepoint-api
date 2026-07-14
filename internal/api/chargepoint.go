package api

import (
	"github.com/go-fuego/fuego"
	"github.com/vpmv/chargepoint-api/internal/dto"
)

func (api *API) ListChargePoints(c fuego.ContextNoBody) ([]*dto.ChargePoint, error) {
	//page := c.QueryParamInt(ParamPage)
	//pageSize := c.QueryParamInt(ParamPageSize)

	return nil, nil
}

func (api *API) GetChargePointByID(c fuego.ContextNoBody) (*dto.ChargePoint, error) {
	id := c.PathParam("id")
	return &dto.ChargePoint{ID: id}, nil
}

func (api *API) CreateChargePoint(c fuego.ContextWithBody[[]*dto.ChargePoint]) (*dto.ChargePoint, error) {
	return nil, nil
}
