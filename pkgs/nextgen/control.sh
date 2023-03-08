#!/usr/bin/env bash

default_module="fish_net"

rpc_service_list=("user" "note")
http_service_list=("api")

set -eu

version() {
  echo "fish_net_control v0.1.0"
}

usage() {
  version
  cat <<EOF
USAGE:
  ./control <command> [OPTIONS] [INPUT]...

options:default_module
    -h, --help                            Print this usage information.
    -v, --version                         Print this version information.
commands:
    run <all | rpc | http> [NAME]         Start project
    update <all | rpc | http> [NAME]      Update project
    stop <all | user>                     Stop project
    new <rpc | http> <NAME>               Create a new rpc or http service
                                            rpc: create rpc service
                                            http: create http service
                                            name: service name
EOF
}

control_rpc() {
  [ $# -gt 1 ] || exit_with_usage
  case $1 in
  "new")
    name=$2
    # execute in the project root directory
    kitex \
      --thrift-plugin validator \
      -module "$default_module" \
      idl/"$name".thrift

    # execute in cmd/user
    mkdir "cmd/$name" && cd "cmd/$name" && kitex \
      --thrift-plugin validator \
      -module "$default_module" \
      -service "${default_module}$name" \
      -use "$default_module"/kitex_gen \
      ../../idl/"$name".thrift
    ;;
  "update")
    name=$2
    cd "cmd/$name" && kitex \
      --thrift-plugin validator \
      -module "$default_module" \
      -service "${default_module}$name" \
      -use "$default_module"/kitex_gen \
      ../../idl/"$name".thrift
    ;;
  "run")
    cd cmd/"${2}" && sh build.sh && sh output/bootstrap.sh &
    ;;
  "stop")
    pkill -SIGINT "$2"
    ;;
  "reload")
    control_rpc stop "$2"
    control_rpc start "$2"
    ;;
  *)
    exit_with_usage
    ;;
  esac
}

control_http() {
  [ $# -gt 1 ] || exit_with_usage
  case $1 in
  "new")
    name=$2
    mkdir "cmd/$name" &&
      cd "cmd/$name" &&
      hz new \
        --mod "$default_module/cmd/$name" \
        --idl ../../idl/"$name".thrift &&
      rm "./go.mod"
    ;;
  "update")
    name=$2
    cd "cmd/$name" &&
      hz update \
        -mod "$default_module/cmd/$name" \
        --idl ../../idl/"$name".thrift
    ;;
  "run")
    cd cmd/"${2}" && go run . &
    ;;
  "stop")
    pkill -SIGINT "$2"
    ;;
  "reload")
    control_http stop "$2"
    control_http start "$2"
    ;;
  *)
    exit_with_usage
    ;;
  esac
}

control_service() {
  [ $# -gt 1 ] || exit_with_usage

  # $1: [run, stop, reload, new ... ]

  case "${2}" in
  "all")
    [ $# -gt 2 ] || {
      control_service "$1" rpc
      control_service "$1" http
    }

    case $3 in
    "rpc")
      for service in "${rpc_service_list[@]}"; do
        control_rpc "$1" "$service"
      done
      ;;
    "http")
      for service in "${http_service_list[@]}"; do
        control_http "$1" "$service"
      done
      ;;
    esac
    ;;
  "rpc")
    control_rpc "$1" "$3"
    ;;
  "http")
    control_http "$1" "$3"
    ;;
  *)
    exit_with_usage
    ;;
  esac
}

exit_with_usage() {
  usage
  exit 0
}

exit_success() {
  control_service stop all
  exit 0
}

trap 'exit_success' INT

[ $# -gt 0 ] || exit_with_usage

case $1 in
"-h" | "--help")
  usage
  ;;
"-v" | "--version")
  version
  ;;
*)
  control_service "${@:1}"
  ;;
esac
