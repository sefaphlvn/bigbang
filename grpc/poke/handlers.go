package poke

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (p *Poke) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/ping":
		p.handlePing(w, r)
	case "/poke":
		p.handlePoke(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (p *Poke) handlePing(w http.ResponseWriter, r *http.Request) {
	p.logger.Info(fmt.Fprint(w, "OK"))
}

func (p *Poke) handlePoke(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()

	serviceValue := queryParams.Get("service")

	if serviceValue == "" {
		http.Error(w, "Service query parameter is required", http.StatusBadRequest)
		return // Query parametre boşsa, burada fonksiyonu sonlandır
	}

	allResources, err := p.getAllResourcesFromListener(serviceValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.logger.Error(err)
		return // Bir hata oluşursa, burada fonksiyonu sonlandır
	}

	err = p.ctx.SetSnapshot(allResources, p.logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.logger.Error(err)
		return // Bir hata oluşursa, burada fonksiyonu sonlandır
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // İşlemler başarılıysa HTTP 200 durum kodunu ayarla
	err = json.NewEncoder(w).Encode(allResources)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		p.logger.Error(err)
		return // JSON encode işlemi sırasında bir hata oluşursa, burada fonksiyonu sonlandır
	}
}
