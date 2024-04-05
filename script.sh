#!/bin/bash

# This script can be used to test rate limiting of the gateway APIs
num_requests=20

# Define the authorization token
# Get it from /login API
authorization_token="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTIyMTQ5MTIsInVzZXJuYW1lIjoiZGVlcGFrIn0.0Ah4J9HvcHZii7KcJFj-BjRY9y7GxKp0ziP_7sopa58"

# Send multiple requests
for ((i=1; i<=$num_requests; i++)); do
    curl -X GET http://localhost:<ENTER PORT HERE> -H "Authorization: Bearer $authorization_token" &
done

# Wait for all requests to finish
wait