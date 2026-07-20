package helpers

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type _testCompany struct {
	Name      string `validate:"required"`
	Building  int    `validate:"number,gt=0"`
	SignedInt int    `validate:"number,min=-90,max=90"`
}

type _testEntity struct {
	Name    string `validate:"required,alpha"`
	Company _testCompany
}

func TestValidateEntity(t *testing.T) {
	var (
		entity *_testEntity
		errs   validator.ValidationErrors
	)

	entity = &_testEntity{
		Name: "John",
		Company: _testCompany{
			Name:      "Doe Corp.",
			Building:  1,
			SignedInt: 0,
		},
	}
	errs = ValidateEntity(entity)
	if errs != nil {
		t.Errorf("Did not expect ValidationError, got: %v", errs)
	}

	entity.Name = `123`
	errs = ValidateEntity(entity)
	if errs == nil {
		t.Error(`expected a ValidationError`)
		t.Fail()
	}

	assert.Equal(t, `Name`, errs[0].Field())
	assert.Equal(t, `123`, errs[0].Value())

	entity = &_testEntity{
		Name: "John",
		Company: _testCompany{
			Name:      "",
			Building:  1,
			SignedInt: 0,
		},
	}

	errs = ValidateEntity(entity)
	fmt.Printf(`%+v`, errs)
	if errs == nil {
		t.Error(`expected a ValidationError`)
		t.Fail()
	}

	assert.Equal(t, `Name`, errs[0].Field())
	assert.Equal(t, `required`, errs[0].Tag())
	assert.Equal(t, ``, errs[0].Value())
}
