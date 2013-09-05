function Controller($scope, $http) {

	$scope.apilists = [
		{
			"path":"bench",
			"apis":[
				{
					"name": "login",
					"method": "GET",
					"body": ""
				},
				{
					"name": "hello",
					"method": "GET",
					"body": ""
				}
			] 
		},
		{
			"path":"auth",
			"apis":[
				{
					"name": "login",
					"method": "POST",
					"body": '{"username":"?", password:"?"}'
				},
				{
					"name": "register",
					"method": "POST",
					"body": '{"username":"?", password:"?"}'
				}
			] 
		}
	]
}

