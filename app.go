package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jkallhoff/gofig"
	"github.com/jkallhoff/risk-analyzer-web/riskEngine"
	"net/http"
	"strconv"
)

type dependencyHandler func(w http.ResponseWriter, r *http.Request, repo battleRepository)

func (d dependencyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	repo := new(mongoRepository)
	d(w, r, repo)
	defer repo.Close()
	return
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.Handle("/BattleRequest", dependencyHandler(battleRequestHandler)).Methods("POST")

	http.Handle("/", router)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	panic(http.ListenAndServe(gofig.Str("webPort"), nil))
}

func battleRequestHandler(w http.ResponseWriter, r *http.Request, repo battleRepository) {
	var result *riskEngine.BattleResult
	var attackingArmies, defendingArmies int
	var err error

	if attackingArmies, err = strconv.Atoi(r.FormValue("attackingArmies")); err != nil {
		http.Error(w, "Invalid battle request detected", http.StatusInternalServerError)
		return
	}
	if defendingArmies, err = strconv.Atoi(r.FormValue("defendingArmies")); err != nil {
		http.Error(w, "Invalid battle request detected", http.StatusInternalServerError)
		return
	}

	if result, err = repo.FetchBattleResult(attackingArmies, defendingArmies); err != nil && err.Error() != "not found" {
		http.Error(w, "There was an error fetching the existing results", http.StatusInternalServerError)
		return
	}

	if result == nil {
		battleRequest := &riskEngine.BattleRequest{AttackingArmies: attackingArmies, DefendingArmies: defendingArmies, NumberOfBattles: 10000}
		result = battleRequest.CalculateBattleResults()
		if err = repo.SaveBattleResult(result); err != nil {
			http.Error(w, "There was an error saving the new results", http.StatusInternalServerError)
			return
		}
	}

	if returnData, err := json.Marshal(result); err != nil {
		http.Error(w, "There was an error with your request. Please try again", http.StatusInternalServerError)
		return
	} else {
		w.Write(returnData)
	}

	return
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/views/home.html")
	return
}
