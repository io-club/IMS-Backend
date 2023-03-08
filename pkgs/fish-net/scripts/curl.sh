#! /usr/bin/env bash
# This script is used to test the API
# Usage: ./curl.sh [API] [ID]

set -eux

host=http://localhost
port=8081
baseUrl="${host}:$port"

case $1 in
# Auth
101) # Register
  ;;
# User
201) # createUser
  curl -fsSL -X POST "${baseUrl}"/v1/user -d '{
    "username": "test"
  }' | jq
  ;;
202) # deleteUser
  curl -fsSL "${baseUrl}"/v1/user/"$2" -X DELETE | jq
  ;;
203) # updateUser
  curl -fsSL "${baseUrl}"/v1/user/"$2" -X PUT -d '{
    "nickname":"test",
    "icon":"test",
  }' | jq
  ;;
204) # queryUser
  curl -fsSL "${baseUrl}"/v1/user -X GET -d '{
    "limit": 1
  }' | jq
  ;;
205) # getUser
  curl -fsSL "${baseUrl}"/v1/user/"$2" -X GET | jq
  ;;
# Wordcase
301)
  curl -fsSL "${baseUrl}"/v1/wordcase -X POST | jq
  ;;
302)
  curl -fsSL "${baseUrl}"/v1/wordcase/"$2" -X DELETE | jq
  ;;
303)
  curl -fsSL "${baseUrl}"/v1/wordcase/"$2" -X PUT -d '{
    "word":"test",
    "meaning":"test",
    "example":"test",
  }' | jq
  ;;
304)
  curl -fsSL "${baseUrl}"/v1/wordcase -X GET | jq
  ;;
305)
  curl -fsSL "${baseUrl}"/v1/wordcase/"$2" -X GET | jq
  ;;
# ...
esac
