package charger

import (
	"errors"
	"fmt"
	"time"

	"github.com/evcc-io/evcc/api"
	goe "github.com/evcc-io/evcc/charger/go-e"
	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/sponsor"
)

// LICENSE

// Copyright (c) 2019-2021 andig

// This module is NOT covered by the MIT license. All rights reserved.

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// https://go-e.co/app/api.pdf
// https://github.com/goecharger/go-eCharger-API-v1/
// https://github.com/goecharger/go-eCharger-API-v2/

// GoE charger implementation
type GoE struct {
	api goe.API
}

func init() {
	registry.Add("go-e", NewGoEFromConfig)
}

// go:generate go run ../cmd/tools/decorate.go -f decorateGoE -b *GoE -r api.Charger -t "api.MeterEnergy,TotalEnergy,func() (float64, error)" -t "api.ChargePhases,Phases1p3p,func(int) (error)"

// NewGoEFromConfig creates a go-e charger from generic config
func NewGoEFromConfig(other map[string]interface{}) (api.Charger, error) {
	cc := struct {
		Token string
		URI   string
		Cache time.Duration
	}{}

	if err := util.DecodeOther(other, &cc); err != nil {
		return nil, err
	}

	if cc.URI != "" && cc.Token != "" {
		return nil, errors.New("should only have one of uri/token")
	}
	if cc.URI == "" && cc.Token == "" {
		return nil, errors.New("must have one of uri/token")
	}

	return NewGoE(cc.URI, cc.Token, cc.Cache)
}

// NewGoE creates GoE charger
func NewGoE(uri, token string, cache time.Duration) (api.Charger, error) {
	c := &GoE{}

	log := util.NewLogger("go-e").Redact(token)

	if token != "" {
		c.api = goe.NewCloud(log, token, cache)
	} else {
		c.api = goe.NewLocal(log, util.DefaultScheme(uri, "http"))
	}

	if c.api.IsV2() {
		var phases func(int) error
		if sponsor.IsAuthorized() {
			phases = c.phases1p3p
		} else {
			log.WARN.Println("automatic 1p3p phase switching requires sponsor token")
		}

		return decorateGoE(c, c.totalEnergy, phases), nil
	}

	return c, nil
}

// Status implements the api.Charger interface
func (c *GoE) Status() (api.ChargeStatus, error) {
	resp, err := c.api.Status()
	if err != nil {
		return api.StatusNone, err
	}

	switch car := resp.Status(); car {
	case 1:
		return api.StatusA, nil
	case 2:
		return api.StatusC, nil
	case 3, 4:
		return api.StatusB, nil
	default:
		return api.StatusNone, fmt.Errorf("car unknown result: %d", car)
	}
}

// Enabled implements the api.Charger interface
func (c *GoE) Enabled() (bool, error) {
	resp, err := c.api.Status()
	if err != nil {
		return false, err
	}

	return resp.Enabled(), nil
}

// Enable implements the api.Charger interface
func (c *GoE) Enable(enable bool) error {
	var b int
	if enable {
		b = 1
	}

	param := map[bool]string{false: "alw", true: "frc"}[c.api.IsV2()]
	if c.api.IsV2() {
		b += 1
	}

	return c.api.Update(fmt.Sprintf("%s=%d", param, b))
}

// MaxCurrent implements the api.Charger interface
func (c *GoE) MaxCurrent(current int64) error {
	param := map[bool]string{false: "amx", true: "amp"}[c.api.IsV2()]
	return c.api.Update(fmt.Sprintf("%s=%d", param, current))
}

var _ api.Meter = (*GoE)(nil)

// CurrentPower implements the api.Meter interface
func (c *GoE) CurrentPower() (float64, error) {
	resp, err := c.api.Status()
	if err != nil {
		return 0, err
	}

	return resp.CurrentPower(), err
}

var _ api.ChargeRater = (*GoE)(nil)

// ChargedEnergy implements the api.ChargeRater interface
func (c *GoE) ChargedEnergy() (float64, error) {
	resp, err := c.api.Status()
	if err != nil {
		return 0, err
	}

	return resp.ChargedEnergy(), err
}

var _ api.MeterCurrent = (*GoE)(nil)

// Currents implements the api.MeterCurrent interface
func (c *GoE) Currents() (float64, float64, float64, error) {
	resp, err := c.api.Status()
	if err != nil {
		return 0, 0, 0, err
	}

	i1, i2, i3 := resp.Currents()

	return i1, i2, i3, err
}

var _ api.Identifier = (*GoE)(nil)

// Identify implements the api.Identifier interface
func (c *GoE) Identify() (string, error) {
	resp, err := c.api.Status()
	return resp.Identify(), err
}

// totalEnergy implements the api.MeterEnergy interface - v2 only
func (c *GoE) totalEnergy() (float64, error) {
	resp, err := c.api.Status()
	if err != nil {
		return 0, err
	}

	var val float64
	if res, ok := resp.(*goe.StatusResponse2); ok {
		val = res.TotalEnergy()
	}

	return val, err
}

// phases1p3p implements the api.ChargePhases interface - v2 only
func (c *GoE) phases1p3p(phases int) error {
	if phases == 3 {
		phases = 2
	}

	return c.api.Update(fmt.Sprintf("psm=%d", phases))
}
