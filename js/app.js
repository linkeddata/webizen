var Webizen = angular.module('WebizenSearch', ['ui','ui.filters','rt.encodeuri']);
// Main angular controller
function SearchCtrl($scope, $http) {
	$scope.search = [];

	// attempt to find a person using webizen.org
	$scope.lookup = function(query) {
		if (!$scope.search)
			$scope.search;
		$scope.search.selected = false;
		$scope.search.loading = true;
		$scope.search.noresults = false;
		$scope.webidresults = [];

		if (query.length > 0) {
			// get results from server
			$http.get('https://api.webizen.org/v1/search', {
				params: {
					q: query
				}
			}).error(function(res) {
				$scope.search.loading = false;
				$scope.search.noresults = true;
			}).then(function(res){
				angular.forEach(res.data, function(value, key){
					if (value) {
						value.webid = key;
						if (!value.img)
							value.img = ['img/photo.png'];
						value.host = $scope.getHostname(key);
						if ($scope.search.query == query)
							$scope.webidresults.push(value);
					} else {
						$scope.search.noresults = true;
						$scope.search.loading = false;
					}
				});
				$scope.search.loading = false;
			});
		} else {
			$scope.search.loading = false;
		}
	}

	// parse an uri and get the hostname
	$scope.getHostname = function (uri) {
		var l = document.createElement("a");
		l.href = uri;
		return l.hostname;
	};
}

//simple directive to display list of search results
Webizen.directive('searchresults',function(){
  	return {
		replace : true,
		restrict : 'E',
		templateUrl: 'tpl/results.html'
    }; 
})
