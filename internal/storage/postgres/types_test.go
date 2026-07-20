package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChargePoint_Status(t *testing.T) {
	var cp *ChargePoint

	// default statuses
	cp = new(ChargePoint)
	cp.SetStatus(`available`)
	assert.Equal(t, statusIdle.String(), cp.Status.String())

	cp = new(ChargePoint)
	cp.SetStatus(`offline`)
	assert.Equal(t, statusOffline.String(), cp.Status.String())

	cp = new(ChargePoint)
	cp.SetStatus(`occupied`)
	assert.Equal(t, statusOccupied.String(), cp.Status.String())

	// invalid statuses
	cp = new(ChargePoint)
	cp.SetStatus(`invalid`)
	assert.Equal(t, statusIdle.String(), cp.Status.String())

	cp = new(ChargePoint)
	cp.SetStatus(``)
	assert.Equal(t, statusIdle.String(), cp.Status.String())

	// manual assign
	cp.Status = statusOccupied
	assert.Equal(t, statusOccupied.String(), cp.Status.String())

	// invalid manual assign
	cp.Status = 6
	assert.Equal(t, ``, cp.Status.String())
}
