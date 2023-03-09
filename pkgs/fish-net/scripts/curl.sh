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
5001) # createDeviceSensor
  curl -fsSL "${baseUrl}"/v1/device/"$2"/sensor -X POST -d '{
    "sensorTypeId": 2,
    "remark": "Some Info 2"
  }' | jq
  ;;
5002) # getDeviceSensors
  curl -fsSL "${baseUrl}"/v1/device/"$2"/sensor -X GET | jq
  ;;
5003) # getDeviceSensor
  curl -fsSL "${baseUrl}"/v1/device/"$2"/sensor/"$3" -X GET | jq
  ;;
5004) # deleteSensor
  curl -fsSL "${baseUrl}"/v1/sensor/"$2" -X DELETE | jq
  ;;
5005) # updateSensor
  curl -fsSL "${baseUrl}"/v1/sensor/"$2" -X PUT -d '{
    "remark":"some-info-2"
  }' | jq
  ;;
5006) # querySensor
  curl -fsSL "${baseUrl}"/v1/sensor -X GET | jq
  ;;
5007) # getSensor
  curl -fsSL "${baseUrl}"/v1/sensor/"$2" -X GET | jq
  ;;
# SensorType
6001) # createSensorType
  curl -fsSL "${baseUrl}"/v1/sensorType -X POST -d '{
    "name": "sensorType1",
    "remark": "Some Info"
  }' | jq
  ;;
6002) # deleteSensorType
  curl -fsSL "${baseUrl}"/v1/sensorType/"$2" -X DELETE | jq
  ;;
6003) # updateSensorType
  curl -fsSL "${baseUrl}"/v1/sensorType/"$2" -X PUT -d '{
    "name":"sensorType1-1",
    "remark":"some-info-2"
  }' | jq
  ;;
6004) # querySensorType
  curl -fsSL "${baseUrl}"/v1/sensorType -X GET | jq
  ;;
6005) # getSensorType
  curl -fsSL "${baseUrl}"/v1/sensorType/"$2" -X GET | jq
  ;;
esac
