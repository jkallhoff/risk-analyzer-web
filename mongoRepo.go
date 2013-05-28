package main

import (
	"fmt"
	"github.com/jkallhoff/gofig"
	"github.com/jkallhoff/risk-analyzer-web/riskEngine"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	session    mgo.Session
	collection *mgo.Collection
)

type battleRepository interface {
	SaveBattleResult(*riskEngine.BattleResult)
	FetchBattleResult(attackingArmies, defendingArmies int) *riskEngine.BattleResult
	Close()
}

type mongoRepository struct {
}

func (*mongoRepository) Close() {
	session.Close()
}

func (*mongoRepository) SaveBattleResult(result *riskEngine.BattleResult) {
	err := collection.Insert(result)
	if err != nil {
		panic(err)
	}
}

func (*mongoRepository) FetchBattleResult(attackingArmies, defendingArmies int) *riskEngine.BattleResult {
	result := riskEngine.BattleResult{}
	if err := collection.Find(bson.M{"AttackingArmies": attackingArmies, "DefendingArmies": defendingArmies}).One(&result); err != nil {
		return nil
	}
	return &result
}

func init() {
	session, err := mgo.Dial(fmt.Sprintf("%s:%s@localhost:%s", gofig.Str("mongoUsername"), gofig.Str("mongoPassword"), gofig.Str("mongoPort")))
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	collection = session.DB("risk").C("battleResults")
}
