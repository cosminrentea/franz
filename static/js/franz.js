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
    $scope.messages = [];
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
        return $http.get('/message/')
            .error(handleError)
            .success(function (data) {
                $scope.messages = data.Messages;
            });
    };

    $scope.sendMessage = function () {
        $scope.working = true;
        if (confirm('Do you really want to send this message?')) {
            $http.post('/message/', {Title: $scope.messageText})
                .error(handleError)
                .success(function () {
                    refresh().then(function () {
                        $scope.errorSending = false;
                        $scope.working = false;
                        $scope.messageText = '';
                    })
                });
        } else {
            $scope.working = false;
        }
    };

    $scope.toggleDone = function (message) {
        data = {ID: message.ID, Title: message.Title, Done: !message.Done}
        $http.put('/message/' + message.ID, data)
            .error(handleError)
            .success(function () {
                message.Done = !message.Done
            });
    };

    refresh().then(function () {
        $scope.working = false;
    });
}