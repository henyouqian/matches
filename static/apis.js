function Controller($scope, $http) {
	$scope.apilists = [
		{
			"path":"bench",
			"apis":[
				{
					"name": "login",
					"method": "GET",
					"data": ""
				},{
					"name": "hello",
					"method": "GET",
					"data": ""
				},{
					"name": "dbsingleselect",
					"method": "GET",
					"data": ""
				},{
					"name": "dbinsert",
					"method": "GET",
					"data": ""
				},{
					"name": "dbinserttx",
					"method": "GET",
					"data": ""
				},{
					"name": "redisget",
					"method": "GET",
					"data": ""
				},{
					"name": "json",
					"method": "GET",
					"data": ""
				},
				{
					"name": "json2",
					"method": "GET",
					"data": ""
				},

			]
		},
		{
			"path":"auth",
			"apis":[
				{
					"name": "login",
					"method": "POST",
					"data": {"Username":"admin", "Password":"admin", "Appsecret":"app_secret_here"}
				},{
					"name": "logout",
					"method": "POST",
					"data": ""
				},{
					"name": "info",
					"method": "POST",
					"data": ""
				},{
					"name": "register",
					"method": "POST",
					"data": {"Username":"?", "Password":"?"}
				},{
					"name": "newapp",
					"method": "POST",
					"data": {"Name":"?"}
				},{
					"name": "listapp",
					"method": "POST",
					"data": ""
				},
			] 
		},
		{
			"path":"game",
			"apis":[
				{
					"name": "list",
					"method": "POST",
					"data": ""
				},{
					"name": "new",
					"method": "POST",
					"data": {"Id":1, "Name":"aa", "Sort":"DESC"}
				},{
					"name": "del",
					"method": "POST",
					"data": [1]
				},
			] 
		},
		{
			"path":"match",
			"apis":[
				{
					"name": "list",
					"method": "POST",
					"data": ""
				},{
					"name": "new",
					"method": "POST",
					"data": {"Name":"aa", "Gameid":1, "Begin":"2006-01-02 15:04:05", "End":"2014-01-02 15:04:05"}
				},{
					"name": "del",
					"method": "POST",
					"data": [1]
				},{
					"name": "start",
					"method": "POST",
					"data": {"MatchId":1}
				},{
					"name": "addscore",
					"method": "POST",
					"data": {"TrySecret":"xxxyyyzzz", "Score":12345}
				},
			] 
		},
		{
			"path":"rank",
			"apis":[
				{
					"name": "mine",
					"method": "POST",
					"data": {"MatchId":1}
				},
			] 
		},
	]

	var sendCodeMirror = CodeMirror.fromTextArea(sendTextArea, 
		{
			theme: "elegant",
		}
	);
	var recvCodeMirror = CodeMirror.fromTextArea(recvTextArea, 
		{
			theme: "elegant",
		}
	);
	var historyCodeMirror = CodeMirror.fromTextArea(historyTextArea, 
		{
			theme: "elegant",
			readOnly: true
		}
	);
	historyCodeMirror.setSize("100%", 600)

	$scope.selectedApiPath = ""
	$scope.currApi = null

	$scope.onApiClick = function(api, path) {
		if ($scope.currApi != api) {
			$("#btn-send").removeAttr("disabled")
			$scope.currApi = api
			$scope.currApi.path = path
			$scope.currUrl = $scope.currApi.path+"/"+$scope.currApi.name
			if (api.data) {
				sendCodeMirror.doc.setValue(JSON.stringify(api.data, null, '\t'))
			} else {
				sendCodeMirror.doc.setValue("")
			}
			
			// recvCodeMirror.doc.setValue("")
		}
	}

	$scope.queryTick = 0
	var lastHisText = ""
	$scope.send = function() {
		var url = "../"+$scope.currUrl
		var input = sendCodeMirror.doc.getValue()
		var inputText = input
		if (input) {
			try {
				input = JSON.parse(input)
			} catch(err) {
				alert("parse json error")
				return
			}	
		}

		var onReceive = function(json) {
			printQueryTick()
			var replyText = JSON.stringify(json, null, '\t')
			recvCodeMirror.doc.setValue(replyText)

			var hisDoc = historyCodeMirror.getDoc()
			hisDoc.setCursor({line: 0, ch: 0})

			inputText = "\t"+inputText.replace(/\n/g, "\n\t");
			replyText = "\t"+replyText.replace(/\n/g, "\n\t");

			var hisText = "=> " + $scope.currUrl + "\n" + inputText + "\n<=\n" + replyText + "\n"
			hisText += "------------------------\n"
			console.log(lastHisText, hisText)
			if (lastHisText != hisText) {
				lastHisText = hisText
				hisDoc.replaceSelection(hisText, "start")
			}
		}

		var onFail = function(obj) {
			printQueryTick()
			var text = obj.status + ":" + obj.statusText + "\n\n" + JSON.stringify(obj.responseJSON, null, '\t')
			recvCodeMirror.doc.setValue(text)
		}

		function printQueryTick() {
			$scope.$apply(function(){
				$scope.queryTick = Math.round(window.performance.now() - t)
			});
		}
		var t = window.performance.now()
		if ($scope.currApi.method == "GET") {
			$.getJSON(url, input, onReceive)
			.fail(onFail)
		}else if ($scope.currApi.method == "POST") {
			$.post(url, sendCodeMirror.doc.getValue(), onReceive, "json")
			.fail(onFail)
		}
	}

	$('#collapseOne').on('shown.bs.collapse', function () {
		historyCodeMirror.refresh()
	})
}



