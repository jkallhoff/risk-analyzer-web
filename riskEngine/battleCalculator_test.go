package riskEngine

import (
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
