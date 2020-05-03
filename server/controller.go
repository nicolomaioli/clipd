package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickmn/go-cache"
	"github.com/rs/zerolog"
)

// ClipdController contains the handler for the "/clipd" routes
type ClipdController struct {
	logger *zerolog.Logger
	cache  *cache.Cache
}

// NewClipdController instantiates a new clipdController
func NewClipdController(l *zerolog.Logger, c *cache.Cache) *ClipdController {
	return &ClipdController{
		logger: l,
		cache:  c,
	}
}

// Yank POST /clipd
func (c *ClipdController) Yank(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.logger.Error().Msgf("error reading body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	yr := &Clip{}
	err = json.Unmarshal(body, yr)
	if err != nil {
		c.logger.Error().Msgf("invalid json: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if yr.Reg == "" {
		yr.Reg = DefaultRegister
	}

	c.cache.Set(yr.Reg, yr.Content, 0)
	w.WriteHeader(http.StatusOK)
}

// Paste GET /clipd
func (c *ClipdController) Paste(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reg := p.ByName("reg")
	if reg == "" {
		reg = DefaultRegister
	}

	var content string

	if v, ok := c.cache.Get(reg); ok {
		content = v.(string)
	}

	if content == "" {
		c.logger.Debug().Msgf("clip not found with reg %q", reg)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	yr := &Clip{
		Reg:     reg,
		Content: content,
	}

	b, err := json.Marshal(yr)
	if err != nil {
		c.logger.Error().Msgf("error marshaling request %v: %s", yr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
