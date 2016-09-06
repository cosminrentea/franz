/*
 Copyright 2011 The Go Authors.  All rights reserved.
 Use of this source code is governed by a BSD-style
 license that can be found in the LICENSE file.
 */

var angular = require('angular');

angular.module('franz', [require('angular-messages')])
    .controller('FranzCtrl', function ($scope, $http) {
    $scope.messages = [];
    $scope.working = false;
    $scope.errorSending = false;
    $scope.maximumLength = 400;
    $scope.minimumLength = 5;

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
        var data = {
            ID: message.ID,
            Content: message.Content,
            Done: !message.Done
        };
        $http.put('/message/' + message.ID, data)
            .error(handleError)
            .success(function () {
                message.Done = !message.Done
            });
    };

    refresh().then(function () {
        $scope.working = false;
    });
});