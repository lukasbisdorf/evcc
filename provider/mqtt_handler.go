package provider

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"

	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/util/jq"
	"github.com/itchyny/gojq"
)

type msgHandler struct {
	mux     *util.Waiter
	scale   float64
	topic   string
	payload string
	re      *regexp.Regexp
	jq      *gojq.Query
}

func (h *msgHandler) receive(payload string) {
	h.mux.Lock()
	defer h.mux.Unlock()

	h.payload = payload
	h.mux.Update()
}

// hasValue returned the received and processed payload as string
func (h *msgHandler) hasValue() (string, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	if late := h.mux.Overdue(); late > 0 {
		return "", fmt.Errorf("%s outdated: %v", h.topic, late.Truncate(time.Second))
	}

	var err error
	payload := h.payload

	if h.re != nil {
		m := h.re.FindStringSubmatch(payload)
		if len(m) > 1 {
			payload = m[1] // first submatch
		}
	}

	if h.jq != nil {
		var val interface{}
		if val, err = jq.Query(h.jq, []byte(payload)); err == nil {
			payload = fmt.Sprintf("%v", val)
		}
	}

	return payload, err
}

func (h *msgHandler) floatGetter() (float64, error) {
	v, err := h.hasValue()
	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, fmt.Errorf("%s invalid: '%s'", h.topic, v)
	}

	return f * h.scale, nil
}

func (h *msgHandler) intGetter() (int64, error) {
	f, err := h.floatGetter()
	return int64(math.Round(f)), err
}

func (h *msgHandler) stringGetter() (string, error) {
	v, err := h.hasValue()
	if err != nil {
		return "", err
	}

	return string(v), nil
}

func (h *msgHandler) boolGetter() (bool, error) {
	v, err := h.hasValue()
	if err != nil {
		return false, err
	}

	return util.Truish(v), nil
}
