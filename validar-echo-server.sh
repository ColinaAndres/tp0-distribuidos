#!/bin/bash

TEST_MESSAGE="Testing Message from Script"

SERVER_RESPONSE=$(docker run --rm --network tp0_testing_net busybox /bin/sh -c "echo '$TEST_MESSAGE' | nc server 12345")

if [ "$SERVER_RESPONSE" == "$TEST_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi



