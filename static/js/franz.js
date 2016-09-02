/*
 Copyright 2011 The Go Authors.  All rights reserved.
 Use of this source code is governed by a BSD-style
 license that can be found in the LICENSE file.
 */

(function (angular) {
    'use strict';
    angular.module('franz', ['ngMessages']);
})(window.angular);

function FranzCtrl($scope, $http) {
    $scope.tasks = [];
    $scope.working = false;
    $scope.errorSending = false;
    $scope.maximumLength = 400;

    var handleError = function (data, status) {
        console.log('code ' + status + ': ' + data);
        $scope.errorSending = true;
        $scope.working = false;
        $scope.date = new Date();
    };

    var refresh = function () {
        return $http.get('/task/').success(function (data) {
            $scope.tasks = data.Tasks;
        }).error(handleError);
    };

    $scope.sendMessage = function () {
        $scope.working = true;
        if (confirm('Do you really want to send this message?')) {
            $http.post('/task/', {Title: $scope.todoText})
                .error(handleError)
                .success(function () {
                    refresh().then(function () {
                        $scope.errorSending = false;
                        $scope.working = false;
                        $scope.todoText = '';
                    })
                });
        } else {
            $scope.working = false;
        }
    };

    $scope.toggleDone = function (task) {
        data = {ID: task.ID, Title: task.Title, Done: !task.Done}
        $http.put('/task/' + task.ID, data)
            .error(handleError)
            .success(function () {
                task.Done = !task.Done
            });
    };

    refresh().then(function () {
        $scope.working = false;
    });
}