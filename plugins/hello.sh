#!/usr/bin/env bash

COMMAND_NAME='hello'

mosquitto_sub -h localhost -v -t "from/irc/+/+/${COMMAND_NAME}" | while read ARGS;
  do echo "${ARGS}";
done