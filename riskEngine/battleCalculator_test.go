package riskEngine

import (
	//"log"
	"testing"
)

func TestCustomSingleBattleResolver(t *testing.T) {
	//Arrange
	var mockResolver = func(*BattleRequest) (result *singleBattleResult) {
		result = new(singleBattleResult)
		result.AttackingArmiesLeft = 20 //arbitrary amount
		return
	}

	//Act
	request := BattleRequest{singleBattleResolver: mockResolver}

	//Assert
	if result := request.singleBattleResolver(&request); result.AttackingArmiesLeft != 20 {
		t.Error("We were unable to override the single battle resolver")
	}
}

func TestCustomDiceRoller(t *testing.T) {
	//Arrange
	var mockRoller = func(numberOfDiceToRoll int) (results []int) {
		results = make([]int, 1)
		results[0] = 100
		return
	}

	//Act
	request := BattleRequest{diceRoller: mockRoller}

	//Assert
	if result := request.diceRoller(1); result[0] != 100 {
		t.Error("We were unable to override the dice roller")
	}
}

//test rollTheDice(int)
func TestDefaultRollTheDiceLengthOfReturn(t *testing.T) {
	//Arrange,Act,Assert
	results := rollTheDice(3)
	if length := len(results); length != 3 {
		t.Errorf("The number of dice rolled should have been 3 but instead was: %v", length)
	}
}

func TestDefaultRollTheDiceSortsResultOrder(t *testing.T) {
	//Arrange,Act,Assert
	for i := 0; i < 20; i++ {
		results := rollTheDice(10)

		previous := results[0]
		for _, digit := range results {
			if digit > previous {
				t.Errorf("%v Values are unsorted", results)
				break
			}
			previous = digit
		}
	}
}

//test determineSingleBattle(*BattleRequest)
func TestSingleBattleResultsWhereAttackerWins(t *testing.T) {
	//Arrange
	requestOne := FetchBattleRequest(4, 2, 1)
	requestTwo := FetchBattleRequest(4, 8, 1)
	requestThree := FetchBattleRequest(2, 8, 1)

	//Act,Assert
	RunSingleBattleTest(t, requestOne, 6, 5, 4)
	RunSingleBattleTest(t, requestTwo, 6, 5, 4)
	RunSingleBattleTest(t, requestThree, 6, 5, 2)

}

func TestSingleBattleResultsWhereDefenderTies(t *testing.T) {
	//Arrange
	requestOne := FetchBattleRequest(4, 2, 1)
	requestTwo := FetchBattleRequest(3, 8, 1)
	requestThree := FetchBattleRequest(2, 8, 1)

	//Act,Assert
	RunSingleBattleTest(t, requestOne, 5, 5, 1)
	RunSingleBattleTest(t, requestTwo, 5, 5, 1)
	RunSingleBattleTest(t, requestThree, 5, 5, 1)

}

//Helpers
func RunSingleBattleTest(t *testing.T, request *BattleRequest, attackerDice, defenderDice, expectedNumberOfAttackersLeft int) {
	state := 1
	var mockRoller = func(numberOfDiceToRoll int) (results []int) {
		results = make([]int, numberOfDiceToRoll)
		if state == 1 {
			for i := range results {
				results[i] = attackerDice
			}
			state = 2
		} else {
			for i := range results {
				results[i] = defenderDice
			}
			state = 1
		}
		return
	}
	request.diceRoller = mockRoller
	if singleResult := determineSingleBattle(request); singleResult.AttackingArmiesLeft != expectedNumberOfAttackersLeft {
		t.Errorf("Expected %v attackers left but instead found %v", expectedNumberOfAttackersLeft, singleResult.AttackingArmiesLeft)
	}
}

func FetchBattleRequest(attackingArmies, defendingArmies, numberOfBattles int) (request *BattleRequest) {
	request = new(BattleRequest)
	request.AttackingArmies = attackingArmies
	request.DefendingArmies = defendingArmies
	request.NumberOfBattles = numberOfBattles
	return
}
