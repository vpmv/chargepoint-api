package api

import (
	"fmt"

	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
	"github.com/vpmv/chargepoint-api/internal/dto"
	"github.com/vpmv/chargepoint-api/internal/helpers"
)

type LocationParams struct {
	Latitude  float64 `query:"lat" validate:"numeric,min=-90,max=90,ne=0,required"`
	Longitude float64 `query:"lon" validate:"numeric,min=-180,max=180,ne=0,required"`
	Radius    int     `query:"radius" validate:"min=5,max=100,required"`
}

func (api *API) CreateChargePoints(c fuego.ContextWithBody[dto.ChargePoints]) (*dto.ChargePoints, error) {
	var (
		chargePoints dto.ChargePoints
		err          error
		errs         map[string]validator.ValidationErrors
	)

	chargePoints, err = c.Body()
	if err != nil {
		return nil, err
	}

	chargePoints, errs = chargePoints.Validate(true)
	if len(errs) > 0 {
		api.log.Warnf("Ignoring validation errors: %s", parseValidationErrors(errs))
	}
	if len(chargePoints) == 0 {
		return nil, fuego.BadRequestError{Detail: `no (valid) charge points provided`}
	}

	chargePoints, err = api.store.CreateChargePoints(chargePoints)
	if err != nil {
		return nil, err
	}

	return &chargePoints, nil
}

func (api *API) DeleteChargePoint(c fuego.ContextNoBody) (string, error) {
	id := c.PathParam("id")

	err := api.store.DeleteChargePoint(id)
	if err != nil {
		return "", err
	}

	return "ok", nil
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
	if err := helpers.ValidateEntity(params); err != nil {
		return nil, fuego.BadRequestError{Err: fmt.Errorf(`Invalid field "%s" - did not match rule "%s"`, err[0].StructField(), err[0].Tag())}
	}

	cps, err := api.store.GetChargePointByLocation(params.Latitude, params.Longitude, params.Radius)
	if err != nil {
		return nil, fuego.BadRequestError{Err: err}
	}

	return &cps, nil

}

func parseValidationErrors(errs map[string]validator.ValidationErrors) string {
	var err string
	for id, validationErrors := range errs {
		for _, ve := range validationErrors {
			err += fmt.Sprintf("@ID: \"%s\" - Invalid field \"%s\" - did not match rule \"%s\" with value \"%v\"\n", id, ve.Field(), ve.Tag(), ve.Value())
		}
	}

	return err
}
