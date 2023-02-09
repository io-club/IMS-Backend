#!/usr/bin/env sh

default_module="fish_net"

set -eu

stop_and_clean() {
  pkill -SIGINT api
  pkill -SIGINT demouser
  pkill -SIGINT demonote
  exit 0
}

version() {
  echo "fish_net_control v0.1.0"
}

usage() {
  version
  cat <<EndOfMessage
USAGE:
  ./control <command> [OPTIONS] [INPUT]...

options:
    -h, --help                 Print this usage information.
    -v, --version              Print this version information.
commands:
    start                      Start project
    stop                       Stop project
    api <new | update>         Create or Update hertz api service
    rpc <new | update> <name>  Create or update a kitex rpc service

EndOfMessage
}

exit_with_usage() {
  usage
  exit 1
}

trap 'stop_and_clean' INT

[ $# -gt 0 ] || exit_with_usage
case $1 in
"-h" | "--help")
  usage
  ;;
"-v" | "--version")
  version
  ;;
"run")
  cd cmd/user &&
    sh build.sh &&
    sh output/bootstrap.sh &

  cd cmd/note &&
    sh build.sh &&
    sh output/bootstrap.sh &

  cd cmd/api &&
    go run . &

  while true; do
    read -r INPUT
    [ "$INPUT" = "stop" ] && break
    bash -c "$INPUT"
  done

  stop_and_clean

  ;;
"stop")
  stop_and_clean
  ;;
"rpc")
  [ $# -gt 1 ] || exit_with_usage
  case $2 in
  "new")
    [ $# -gt 2 ] || exit_with_usage
    name=$3
    # execute in the project root directory
    kitex \
      --thrift-plugin validator \
      -module "$default_module" \
      idl/"$name".thrift

    # execute in cmd/user
    mkdir -p "cmd/$name" && cd "cmd/$name" && kitex \
      --thrift-plugin validator \
      -module "$default_module" \
      -service "${default_module}$name" \
      -use "$default_module"/kitex_gen \
      ../../idl/"$name".thrift
    ;;
  "update")
    #
    ;;
  esac
  ;;
"api")
  [ $# -gt 1 ] || exit_with_usage
  case $2 in
  "new")
    hz new -mod "$default_module" -idl ../../idl/api.thrift
    ;;
  "update")
	  hz update -mod "$default_module" -idl ../../idl/api.thrift
    ;;
  esac
  ;;
*)
  echo '???'
  ;;
esac
