package main

import (
	"testing"

	"github.com/gopcua/opcua/ua"
	prom "github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func getTestHandler() OpcValueHandler {
	var testGuage = prom.NewGauge(prom.GaugeOpts{Name: "foo"})
	return OpcValueHandler{testGuage}
}

func TestCoerceBooleanValues(t *testing.T) {
	handler := getTestHandler()

	variant, _ := ua.NewVariant(true)
	res, err := handler.FloatValue(*variant)
	assert.Nil(t, err)
	assert.Equal(t, 1.0, res)

	variant, _ = ua.NewVariant(false)
	res, err = handler.FloatValue(*variant)
	assert.Nil(t, err)
	assert.Equal(t, 0.0, res)
}

func TestCoerceNumericValues(t *testing.T) {
	handler := getTestHandler()

	type floatTest struct {
		input  interface{}
		output float64
	}

	testCases := []floatTest{
		floatTest{byte(0x03), 3.0},
		floatTest{int8(-4), -4.0},
		floatTest{int16(2), 2.0},
		floatTest{int32(33), 33.0},
		floatTest{int64(25), 25.0},
		floatTest{uint8(4), 4.0},
		floatTest{uint16(2), 2.0},
		floatTest{uint32(33), 33.0},
		floatTest{uint64(25), 25.0},
		floatTest{float32(8.8), float64(float32(8.8))}, // float32 --> float64 actually introduces rounding errors on the order of 1e-7
		floatTest{float64(238.4), 238.4},               // float32 --> float64 actually introduces rounding errors on the order of 1e-7
	}
	for _, testCase := range testCases {
		variant, e := ua.NewVariant(testCase.input)
		if e != nil {
			panic(e)
		}
		result, err := handler.FloatValue(*variant)
		assert.Nil(t, err)
		assert.Equal(t, testCase.output, result)
	}

}

func TestValueHandlerErrors(t *testing.T) {
	handler := getTestHandler()
	errorValues := []interface{}{
		"not a number",
	}

	for _, v := range errorValues {
		variant, vErr := ua.NewVariant(v)
		assert.Nil(t, vErr)
		_, err := handler.FloatValue(*variant)
		assert.Error(t, err)
	}
}
