package main

import (
	"fmt"
	"github.com/JKallhoff/gofig"
	"github.com/JKallhoff/risk-analyzer-web/riskEngine"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

var (
	session    mgo.Session
	collection *mgo.Collection
)

type battleRepository interface {
	SaveBattleResult(*riskEngine.BattleResult) error
	FetchBattleResult(attackingArmies, defendingArmies int) (*riskEngine.BattleResult, error)
	Close()
}

type mongoRepository struct {
}

func (*mongoRepository) Close() {
	session.Close()
}

func (*mongoRepository) SaveBattleResult(result *riskEngine.BattleResult) (err error) {
	err = collection.Insert(result)
	return
}

func (*mongoRepository) FetchBattleResult(attackingArmies, defendingArmies int) (*riskEngine.BattleResult, error) {
	result := riskEngine.BattleResult{}
	var err error

	if err = collection.Find(bson.M{"AttackingArmies": attackingArmies, "DefendingArmies": defendingArmies}).One(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func init() {
	session, err := mgo.Dial(fmt.Sprintf("%s:%s@localhost:%s", gofig.Str("mongoUsername"), gofig.Str("mongoPassword"), gofig.Str("mongoPort")))
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	collection = session.DB("risk").C("battleResults")
}
