package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapSlice(t *testing.T) {
	type exposedStruct struct {
		ExternalField string
	}
	type internalStruct struct {
		InternalField string
	}

	extType := []*exposedStruct{
		{
			ExternalField: "Test1",
		},
	}
	intType := []*internalStruct{
		{
			InternalField: "Test1",
		},
	}

	assert.Equal(t, extType, MapSlice(intType, func(s *internalStruct) *exposedStruct {
		return &exposedStruct{
			ExternalField: s.InternalField,
		}
	}))

	assert.Equal(t, intType, MapSlice(extType, func(s *exposedStruct) *internalStruct {
		return &internalStruct{
			InternalField: s.ExternalField,
		}
	}))
}
