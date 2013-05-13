package riskEngine

import (
	"log"
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

//test calculateBattleResult
func TestSingleCalculatedBattleResult(t *testing.T) {
	//Arrange
	request := new(BattleRequest)
	singleBattleResults := []*singleBattleResult{&singleBattleResult{AttackingArmiesLeft: 5, AttackerWon: true}}

	//Assert,Act
	if battleResult := request.calculateBattleResult(singleBattleResults); battleResult.AverageNumberOfAttackersLeft != 5 {
		t.Errorf("Excepted average of 5, instead got %v", battleResult.AverageNumberOfAttackersLeft)
	}

}

func TestPercentageOfWinsFor100(t *testing.T) {
	//Arrange
	request := new(BattleRequest)
	singleBattleResults := []*singleBattleResult{&singleBattleResult{AttackingArmiesLeft: 5, AttackerWon: true}}

	//Assert,Act
	if battleResult := request.calculateBattleResult(singleBattleResults); battleResult.PercentageThatWereWins != 100.0 {
		t.Errorf("Excepted percentage of wins of 100.0, instead got %v", battleResult.PercentageThatWereWins)
	}
}

func TestMultipleCalculatedBattleResult(t *testing.T) {
	//Arrange
	request := new(BattleRequest)
	singleBattleResults := []*singleBattleResult{
		&singleBattleResult{AttackingArmiesLeft: 5, AttackerWon: true},
		&singleBattleResult{AttackingArmiesLeft: 5, AttackerWon: true}}

	//Assert,Act
	if battleResult := request.calculateBattleResult(singleBattleResults); battleResult.AverageNumberOfAttackersLeft != 5 {
		t.Errorf("Excepted average of 5, instead got %v", battleResult.AverageNumberOfAttackersLeft)
	}

}

func TestMultipleCalculatedBattleResultWithVaryingAttackers(t *testing.T) {
	//Arrange
	request := new(BattleRequest)
	singleBattleResults := []*singleBattleResult{
		&singleBattleResult{AttackingArmiesLeft: 5, AttackerWon: true},
		&singleBattleResult{AttackingArmiesLeft: 12, AttackerWon: true}}

	//Assert,Act
	if battleResult := request.calculateBattleResult(singleBattleResults); battleResult.AverageNumberOfAttackersLeft != 8 {
		t.Errorf("Excepted average of 8, instead got %v", battleResult.AverageNumberOfAttackersLeft)
	}
}

func TestMultipleCalculatedBattleResultWithSomeLosers(t *testing.T) {
	//Arrange
	request := new(BattleRequest)
	singleBattleResults := []*singleBattleResult{
		&singleBattleResult{AttackingArmiesLeft: 5, AttackerWon: true},
		&singleBattleResult{AttackingArmiesLeft: 1, AttackerWon: false}}

	//Assert,Act
	if battleResult := request.calculateBattleResult(singleBattleResults); battleResult.PercentageThatWereWins != 50.0 {
		t.Errorf("Excepted average of 50.0, instead got %v", battleResult.PercentageThatWereWins)
	}
}

//test determineSingleBattle(*BattleRequest)
func TestRunSeriesOfBattlesAndVerifyResults(t *testing.T) {
	//Arrange
	request := FetchBattleRequest(4, 2, 1)

	//Act,Assert
	RunSingleBattleTest(t, request, 6, 5, 4) //Attacker should win due to 6's beating 5's
	RunSingleBattleTest(t, request, 5, 5, 1) //Defender should win due to 5's tying 5's
	RunSingleBattleTest(t, request, 5, 6, 1) //Defender should win due to 6's beating 5's
}

//integration tests
func TestFullIntegration(t *testing.T) {
	//Arrange
	request := &BattleRequest{AttackingArmies: 20, DefendingArmies: 12, NumberOfBattles: 1000}

	//Act
	result := request.CalculateBattleResults()

	//Assert
	log.Printf("%#v", result)
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
