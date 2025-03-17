#!/bin/bash

MESSAGE="testMessage"

docker run -d --quiet --name server_validator --network tp0_testing_net alpine:latest tail -f /dev/null

RESPONSE=$(echo "$MESSAGE" | docker exec -i server_validator nc server:12345)

if [ "$RESPONSE" = "$MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi

docker rm -f server_validator