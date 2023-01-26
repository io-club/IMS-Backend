#!/usr/bin/env sh

default_module="easy_note"

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
    -h, --help                Print this usage information.
    -v, --version             Print this version information.
commands:
    start                     Start project
    stop                      Stop project
    new <rpc | api> <name>  Create a new rpc or api service
                                rpc : create rpc service
                                api : create api service
                                name: service name
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
"new")
  [ $# -gt 2 ] || exit_with_usage
  case $2 in
  "rpc")
    name=$3
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
  "api")
    #
    ;;
  esac
  ;;
*)
  echo '???'
  ;;
esac
