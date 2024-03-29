package controllers

import (
	"encoding/json"
	"github.com/solarlabsteam/dvpn-openwrt/services/keys"
	"io/ioutil"
	"net/http"
)

func ListKeys(w http.ResponseWriter, r *http.Request) {
	keys, err := keys.Wallet.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(keys)
}

func AddRecoverKeys(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	addKeys, err := keys.ValidateAndUnmarshalRecovery(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = keys.Wallet.AddRecover(addKeys)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func DeleteKeys(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	deleteKeys, err := keys.ValidateAndUnmarshalDeletion(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = keys.Wallet.Delete(deleteKeys)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
