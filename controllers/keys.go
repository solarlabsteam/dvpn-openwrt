package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/audi70r/dvpn-openwrt/services/keys"
	"io/ioutil"
	"net/http"
)

func ListKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	keys, err := keys.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(keys)
}

func AddRecoverKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	addKeys, err := keys.ValidateAndUnmarshal(body)

	fmt.Println(addKeys)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = keys.AddRecover(addKeys)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
