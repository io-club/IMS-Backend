#! /usr/bin/env bash
# This script is used to test the API
# Usage: ./curl.sh [API] [ID]

set -eux

host=http://localhost
port=8081
baseUrl="${host}:$port"

# code: service_code : api_code

case $1 in
# Auth
1001) # Register
  ;;
# User
2001) # createUser
  curl -fsSL -X POST "${baseUrl}"/v1/user -d '{
    "username": "test2"
  }' | jq
  ;;
2002) # deleteUser
  curl -fsSL "${baseUrl}"/v1/user/"$2" -X DELETE | jq
  ;;
2003) # updateUser
  curl -fsSL "${baseUrl}"/v1/user/"$2" -X PUT -d '{
    "nickname":"test21"
  }' | jq
  ;;
2004) # queryUser
  curl -fsSL "${baseUrl}"/v1/user -X GET -d '{
    "limit": 1
  }' | jq
  ;;
2005) # getUser
  curl -fsSL "${baseUrl}"/v1/user/"$2" -X GET | jq
  ;;
# Wordcase
3001) # createWordcase
  curl -fsSL "${baseUrl}"/v1/wordcase -X POST -d '{
    "group": "test1",
    "key": "test_key2",
    "value": "test_value2",
    "order": 1,
    "disable": false,
    "remark": "Some Info"
  }' | jq
  ;;
3002) # deleteWordcase
  curl -fsSL "${baseUrl}"/v1/wordcase/"$2" -X DELETE | jq
  ;;
3003) # updateWordcase
  curl -fsSL "${baseUrl}"/v1/wordcase/"$2" -X PUT -d '{
    "value":"test2_value2"
  }' | jq
  ;;
3004) # queryWordcase
  curl -fsSL "${baseUrl}"/v1/wordcase -X GET | jq
  ;;
3005) # getWordcase
  curl -fsSL "${baseUrl}"/v1/wordcase/"$2" -X GET | jq
  ;;
# Device
4001) # createDevice
  curl -fsSL "${baseUrl}"/v1/device -X POST -d '{
    "name": "device1",
    "remark": "Some Info"
  }' | jq
  ;;
4002) # deleteDevice
  curl -fsSL "${baseUrl}"/v1/device/"$2" -X DELETE | jq
  ;;
4003) # updateDevice
  curl -fsSL "${baseUrl}"/v1/device/"$2" -X PUT -d '{
    "name":"device1-1",
    "remark":"some-info-2"
  }' | jq
  ;;
4004) # queryDevice
  curl -fsSL "${baseUrl}"/v1/device -X GET | jq
  ;;
4005) # getDevice
  curl -fsSL "${baseUrl}"/v1/device/"$2" -X GET | jq
  ;;
# Sensor
5001) # createSensor
  curl -fsSL "${baseUrl}"/v1/sensor -X POST -d '{
    "name": "test1",
    "type": "test_key2",
    "remark": "Some Info"
  }' | jq
  ;;
5002) # deleteSensor
  curl -fsSL "${baseUrl}"/v1/sensor/"$2" -X DELETE | jq
  ;;
5003) # updateSensor
  curl -fsSL "${baseUrl}"/v1/sensor/"$2" -X PUT -d '{
    "name":"test"
  }' | jq
  ;;
5004) # querySensor
  curl -fsSL "${baseUrl}"/v1/sensor -X GET | jq
  ;;
5005) # getSensor
  curl -fsSL "${baseUrl}"/v1/sensor/"$2" -X GET | jq
  ;;
# Data
6001) # createData
  curl -fsSL "${baseUrl}"/v1/data -X POST -d '{
    "device_id": 1,
    "sensor_id": 1,
    "value": 1,
    "remark": "Some Info"
  }' | jq
  ;;
6002) # deleteData
  curl -fsSL "${baseUrl}"/v1/data/"$2" -X DELETE | jq
  ;;
6003) # updateData
  curl -fsSL "${baseUrl}"/v1/data/"$2" -X PUT -d '{
    "value": 2
  }' | jq
  ;;
6004) # queryData
  curl -fsSL "${baseUrl}"/v1/data -X GET | jq
  ;;
6005) # getData
  curl -fsSL "${baseUrl}"/v1/data/"$2" -X GET | jq
  ;;
esac
