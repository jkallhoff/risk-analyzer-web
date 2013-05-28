package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jkallhoff/gofig"
	"github.com/jkallhoff/risk-analyzer-web/riskEngine"
	"log"
	"net/http"
	"strconv"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/BattleRequest", battleRequestHandler).Methods("POST")

	http.Handle("/", router)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r)
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	panic(http.ListenAndServe(gofig.Str("webPort"), nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/views/home.html")
	return
}

func battleRequestHandler(w http.ResponseWriter, r *http.Request) {
	var result *riskEngine.BattleResult
	attackingArmies, _ := strconv.Atoi(r.FormValue("attackingArmies")) //r.FormValue("body")
	defendingArmies, _ := strconv.Atoi(r.FormValue("defendingArmies"))

	repo := new(mongoRepository)
	defer repo.Close()

	result = repo.FetchBattleResult(attackingArmies, defendingArmies)

	if result == nil {
		battleRequest := &riskEngine.BattleRequest{AttackingArmies: attackingArmies, DefendingArmies: defendingArmies, NumberOfBattles: 10000}
		result = battleRequest.CalculateBattleResults()
		repo.SaveBattleResult(result)
	}

	returnData, _ := json.Marshal(result)
	w.Write(returnData)
	return
}
