package vehicle

import (
	"fmt"
	"time"

	"github.com/evcc-io/evcc/api"
	"github.com/evcc-io/evcc/provider"
	"github.com/evcc-io/evcc/util"
)

//go:generate go run ../cmd/tools/decorate.go -f decorateVehicle -b api.Vehicle -t "api.ChargeState,Status,func() (api.ChargeStatus, error)" -t "api.VehicleRange,Range,func() (int64, error)" -t "api.VehicleOdometer,Odometer,func() (float64, error)"

// Vehicle is an api.Vehicle implementation with configurable getters and setters.
type Vehicle struct {
	*embed
	chargeG func() (float64, error)
	statusG func() (string, error)
}

func init() {
	registry.Add(api.Custom, NewConfigurableFromConfig)
}

// NewConfigurableFromConfig creates a new Vehicle
func NewConfigurableFromConfig(other map[string]interface{}) (api.Vehicle, error) {
	cc := struct {
		embed    `mapstructure:",squash"`
		Charge   provider.Config
		Status   *provider.Config
		Range    *provider.Config
		Odometer *provider.Config
		Cache    time.Duration
	}{
		Cache: interval,
	}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	getter, err := provider.NewFloatGetterFromConfig(cc.Charge)
	if err != nil {
		return nil, fmt.Errorf("charge: %w", err)
	}

	if cc.Cache > 0 {
		getter = provider.NewCached(getter, cc.Cache).FloatGetter()
	}

	v := &Vehicle{
		embed:   &cc.embed,
		chargeG: getter,
	}

	// decorate vehicle with Status
	var status func() (api.ChargeStatus, error)
	if cc.Status != nil {
		v.statusG, err = provider.NewStringGetterFromConfig(*cc.Status)
		if err != nil {
			return nil, fmt.Errorf("status: %w", err)
		}
		status = v.status
	}

	// decorate vehicle with Range
	var rng func() (int64, error)
	if cc.Range != nil {
		rangeG, err := provider.NewIntGetterFromConfig(*cc.Range)
		if err != nil {
			return nil, fmt.Errorf("range: %w", err)
		}
		rng = rangeG
	}

	// decorate vehicle with Range
	var odo func() (float64, error)
	if cc.Odometer != nil {
		odoG, err := provider.NewFloatGetterFromConfig(*cc.Odometer)
		if err != nil {
			return nil, fmt.Errorf("odometer: %w", err)
		}
		odo = odoG
	}

	res := decorateVehicle(v, status, rng, odo)

	return res, nil
}

// SoC implements the api.Vehicle interface
func (v *Vehicle) SoC() (float64, error) {
	return v.chargeG()
}

// SoC implements the api.ChargeState interface
func (v *Vehicle) status() (api.ChargeStatus, error) {
	status := api.StatusF

	statusS, err := v.statusG()
	if err == nil {
		status = api.ChargeStatus(statusS)
	}

	return status, err
}
