$( document ).ready(function() {
	$(".btn").click(function(e){
		e.preventDefault();
		battleRequester.showRollingAnimation();
		battleRequester.submitRequest(battleRequester.showResults);
		return false;
	})
});

var battleRequester = {
	attackingArmies: $("#attackingArmies"),
	defendingArmies: $("#defendingArmies"),

	showRollingAnimation: function(){
		$("#resultsContainer").show();
		$("#result").hide();
		$("#loadingAnimation").show();
	},

	showResults: function(){
		$("#resultsContainer").show();
		$("#loadingAnimation").hide();
		$("#result").show();
	},

	submitRequest: function(handler){
	var values = {
		attackingArmies: this.attackingArmies.val(),
		defendingArmies: this.defendingArmies.val()
	};
	$.ajax({
	      url: "/BattleRequest",
	      type: "POST",
	      dataType: "json",
	      data: values,
	      success: function(response){
	      	$("#percentage").html(response.PercentageThatWereWins);
	      	$("#avgLeft").html(response.AverageNumberOfAttackersLeft);
	      	handler();
	      }
	    });
	}
}