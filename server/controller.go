package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func yank(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Error().Msgf("error reading body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	yr := &Clip{}
	err = json.Unmarshal(body, yr)
	if err != nil {
		logger.Error().Msgf("invalid json: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if yr.Reg == "" {
		yr.Reg = defaultRegister
	}

	memMut.Lock()
	mem[yr.Reg] = yr.Content
	memMut.Unlock()
	w.WriteHeader(http.StatusOK)
}

func paste(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	reg := p.ByName("reg")
	if reg == "" {
		reg = defaultRegister
	}

	var content string

	memMut.RLock()
	if v, ok := mem[reg]; ok {
		content = v
	}
	memMut.RUnlock()

	if content == "" {
		logger.Debug().Msgf("clip not found with reg %q", reg)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	yr := &Clip{
		Reg:     reg,
		Content: content,
	}

	b, err := json.Marshal(yr)
	if err != nil {
		logger.Error().Msgf("error marshaling request %v: %s", yr, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
