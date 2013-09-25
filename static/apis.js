function Controller($scope, $http) {
	$scope.apilists = [
		{
			"tab":"bench",
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
			"tab":"auth",
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
			"tab":"game",
			"path":"game",
			"apis":[
				{
					"name": "list",
					"method": "POST",
					"data": ""
				},{
					"name": "new",
					"method": "POST",
					"data": {"Id":1, "Name":"aa", "Sort":"ASC or DESC"}
				},{
					"name": "del",
					"method": "POST",
					"data": [1]
				},
			] 
		},
		{
			"tab":"match",
			"path":"match",
			"apis":[
				{
					"name": "list",
					"method": "POST",
					"data": ""
				},{
					"name": "new",
					"method": "POST",
					"data": {"Name":"aa", "Gameid":1, "Begin":"2006-01-02 15:04:05", "End":"2014-01-02 15:04:05", "TryMax":1, "TryPrice":10}
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
			"tab":"rank",
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
	sendCodeMirror.setSize("100%", 500)
	recvCodeMirror.setSize("100%", 500)
	historyCodeMirror.setSize("100%", 600)
	sendCodeMirror.addKeyMap({
		"Ctrl-,": function(cm) {
			var hisList = inputHistory[$scope.currUrl]
			if (isdef(hisList)) {
				var idx = Math.max(0, Math.min(hisList.length-1, inputHisIdx-1))
				if (inputHisIdx != idx) {
					inputHisIdx = idx
					sendCodeMirror.doc.setValue(hisList[inputHisIdx][0])
					recvCodeMirror.doc.setValue(hisList[inputHisIdx][1])
				}
			}
		},
		"Ctrl-.": function(cm) {
			var hisList = inputHistory[$scope.currUrl]
			if (isdef(hisList)) {
				var idx = Math.max(0, Math.min(hisList.length-1, inputHisIdx+1))
				if (inputHisIdx != idx) {
					inputHisIdx = idx
					sendCodeMirror.doc.setValue(hisList[inputHisIdx][0])
					recvCodeMirror.doc.setValue(hisList[inputHisIdx][1])
				}
			}
		},
		"Esc":function(cm) {
			var api = $scope.currApi
			if (api && api.data) {
				sendCodeMirror.doc.setValue(JSON.stringify(api.data, null, '\t'))
			} else {
				sendCodeMirror.doc.setValue("")
			}
			recvCodeMirror.doc.setValue("")
		}
	}) 

	CodeMirror.signal(sendCodeMirror, "keydown", 2)

	var inputHistory = {}
	var sendInput = ""
	var sendUrl=""
	var inputHisIdx = 0

	$scope.currApi = null

	$scope.onApiClick = function(api, path) {
		if ($scope.currApi != api) {
			$("#btn-send").removeAttr("disabled")
			$scope.currApi = api
			$scope.currUrl = path+"/"+$scope.currApi.name
			inputHisIdx = 0
			var hisList = inputHistory[$scope.currUrl]
			if (isdef(hisList)) {
				inputHisIdx = hisList.length-1
				sendCodeMirror.doc.setValue(hisList[inputHisIdx][0])
				recvCodeMirror.doc.setValue(hisList[inputHisIdx][1])
			} else {
				if (api.data) {
					sendCodeMirror.doc.setValue(JSON.stringify(api.data, null, '\t'))
				}else{
					sendCodeMirror.doc.setValue("")
				}
				recvCodeMirror.doc.setValue("")
			}
		}
		sendCodeMirror.focus()
	}

	$scope.queryTick = 0
	var lastHisText = ""
	$scope.send = function() {
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

			//history
			var hisDoc = historyCodeMirror.getDoc()
			hisDoc.setCursor({line: 0, ch: 0})

			inputTextTab = "\t"+inputText.replace(/\n/g, "\n\t")
			replyTextTab = "\t"+replyText.replace(/\n/g, "\n\t")

			var hisText = "=> " + $scope.currUrl + "\n" + inputTextTab + "\n<=\n" + replyTextTab + "\n"
			hisText += "------------------------\n"
			if (lastHisText != hisText) {
				lastHisText = hisText
				hisDoc.replaceSelection(hisText, "start")

				//input history
				if (isdef(inputHistory[sendUrl])) {
					var inHisList = inputHistory[sendUrl]
					if (inHisList[inHisList.length-1] != sendInput) {
						inputHistory[sendUrl].push([sendInput, replyText])
					}
				} else {
					inputHistory[sendUrl] = [[sendInput, replyText]]
				}
				inputHisIdx = inputHistory[sendUrl].length-1
			}
			sendCodeMirror.focus()
		}

		var onFail = function(obj) {
			printQueryTick()
			var t = JSON.stringify(obj.responseJSON, null, '\t')
			if (isdef(t))
				t = t.replace(/\\n/g, "\n")
			if (isdef(t))
				t = t.replace(/\\t/g, "  ")
			var text = obj.status + ":" + obj.statusText + "\n\n" + t
			recvCodeMirror.doc.setValue(text)
		}

		function printQueryTick() {
			$scope.$apply(function(){
				$scope.queryTick = Math.round(window.performance.now() - t)
			});
		}

		sendInput = inputText
		sendUrl = $scope.currUrl
		var url = "/"+sendUrl
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



