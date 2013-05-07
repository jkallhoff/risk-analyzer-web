package riskEngine

import (
	"math/rand"
	"time"
)

//types
type BattleCalculator interface {
	CalculateBattleResults() *BattleResult
}

type BattleRequest struct {
	AttackingArmies, DefendingArmies int
	NumberOfBattles                  int
	singleBattleResolver             func(*BattleRequest) *singleBattleResult
	diceRoller                       func(numberOfDiceToRoll int) []int
}

type BattleResult struct {
	AverageNumberOfAttackersLeft int
	PercentageThatWereWins       float32
}

type singleBattleResult struct {
	AttackingArmiesLeft int
	AttackerWon         bool
}

//public methods
func (request *BattleRequest) CalculateBattleResults() (result *BattleResult) {

	battles := make([]*singleBattleResult, request.NumberOfBattles, request.NumberOfBattles)

	if request.singleBattleResolver == nil {
		request.singleBattleResolver = determineSingleBattle
	}

	if request.diceRoller == nil {
		request.diceRoller = rollTheDice
	}

	for index, _ := range battles {
		battles[index] = request.singleBattleResolver(request)
	}

	result = request.calculateBattleResult(battles)
	return
}

func (request *BattleRequest) calculateBattleResult(battles []*singleBattleResult) (result *BattleResult) {
	result = new(BattleResult)
	result.AverageNumberOfAttackersLeft = 2
	result.PercentageThatWereWins = 45.0
	return
}

func determineSingleBattle(request *BattleRequest) (result *singleBattleResult) {
	var (
		remainingAttackers = request.AttackingArmies
		remainingDefenders = request.DefendingArmies
	)

	for remainingAttackers > 1 && remainingDefenders > 0 {
		//var result = diceRoller.Roll()
	}

	result = new(singleBattleResult)
	result.AttackingArmiesLeft = remainingAttackers

	if remainingDefenders > 0 {
		result.AttackerWon = false
	} else {
		result.AttackerWon = true
	}
	return
}

func rollTheDice(numberOfDiceToRoll int) (results []int) {
	results = make([]int, numberOfDiceToRoll)
	rand.Seed(time.Now().UnixNano())
	for index, _ := range results {
		results[index] = rand.Intn(5) + 1
	}
	return
}
