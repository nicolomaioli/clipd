package server

import (
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
func (c *ClipdController) Yank(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.logger.Error().Msgf("error reading body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reg := p.ByName("reg")
	if reg == "" {
		reg = DefaultRegister
	}

	c.cache.Set(reg, content, 0)
	w.WriteHeader(http.StatusOK)
}

// Paste GET /clipd
func (c *ClipdController) Paste(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reg := p.ByName("reg")
	if reg == "" {
		reg = DefaultRegister
	}

	var content []byte

	if v, ok := c.cache.Get(reg); ok {
		content = v.([]byte)
	} else {
		c.logger.Debug().Msgf("clip not found in reg %q", reg)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}
