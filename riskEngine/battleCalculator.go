package riskEngine

import (
	"math/rand"
	"sort"
	"time"
)

type BattleCalculator interface {
	CalculateBattleResults() *BattleResult
}

type BattleResult struct {
	AttackingArmies              int     "AttackingArmies"
	DefendingArmies              int     "DefendingArmies"
	AverageNumberOfAttackersLeft int     "AverageNumberOfAttackersLeft"
	PercentageThatWereWins       float32 "PercentageThatWereWins"
}

type BattleRequest struct {
	AttackingArmies, DefendingArmies int
	NumberOfBattles                  int
	singleBattleResolver             func(*BattleRequest) *singleBattleResult
	diceRoller                       func(numberOfDiceToRoll int) []int
}

type singleBattleResult struct {
	AttackingArmiesLeft int
	AttackerWon         bool
}

func (request *BattleRequest) CalculateBattleResults() (result *BattleResult) {

	rand.Seed(time.Now().UnixNano())
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
	result.AttackingArmies = request.AttackingArmies
	result.DefendingArmies = request.DefendingArmies
	return
}

type diceResultSorter struct {
	sort.Interface
}

func (sorter diceResultSorter) Less(i, j int) bool {
	return sorter.Interface.Less(j, i)
}

func rollTheDice(numberOfDiceToRoll int) (results []int) {
	results = make([]int, numberOfDiceToRoll)
	for index, _ := range results {
		results[index] = rand.Intn(6) + 1
	}
	sort.Sort(diceResultSorter{sort.IntSlice(results)})
	return
}

func (request *BattleRequest) calculateBattleResult(battles []*singleBattleResult) (result *BattleResult) {
	result = new(BattleResult)
	sum := 0
	numberOfBattles := 0
	numberOfWins := 0

	for _, battleResult := range battles {
		sum = sum + battleResult.AttackingArmiesLeft
		numberOfBattles++
		if battleResult.AttackerWon {
			numberOfWins++
		}
	}

	result.AverageNumberOfAttackersLeft = int(sum / numberOfBattles)
	result.PercentageThatWereWins = (float32(numberOfWins) / float32(numberOfBattles)) * 100
	return
}

func determineSingleBattle(request *BattleRequest) (result *singleBattleResult) {
	var (
		remainingAttackers = request.AttackingArmies
		remainingDefenders = request.DefendingArmies
		numberOfDice       int
	)

	for remainingAttackers > 1 && remainingDefenders > 0 {
		if remainingAttackers > 4 {
			numberOfDice = 3
		} else {
			numberOfDice = remainingAttackers - 1
		}

		attackingDice := request.diceRoller(numberOfDice)

		if remainingDefenders >= 2 {
			numberOfDice = 2
		} else {
			numberOfDice = 1
		}

		defendingDice := request.diceRoller(numberOfDice)

		if attackingDice[0] > defendingDice[0] {
			remainingDefenders--
		} else {
			remainingAttackers--
		}

		if len(attackingDice) > 1 && len(defendingDice) > 1 {
			if attackingDice[1] > defendingDice[1] {
				remainingDefenders--
			} else {
				remainingAttackers--
			}
		}
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
