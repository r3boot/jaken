#!/usr/bin/env bash

COMMAND_NAME='hello'

mosquitto_sub -h localhost -v -t "from/irc/+/+/${COMMAND_NAME}" | while read MSG;do
  TOPIC="$(echo "${MSG}" | cut -d" " -f1)"
  CHANNEL="$(echo "${TOPIC}" | cut -d/ -f3)"
  NICKNAME="$(echo "${TOPIC}" | cut -d/ -f4)"
  NAME="$(echo "${MSG}" | cut -d" " -f2)"

  if [[ "${NAME}" == '(null)' ]]; then
    mosquitto_pub -h localhost -t "to/irc/${CHANNEL}/notice" -m "Hi ${NICKNAME} :)"
  else
    mosquitto_pub -h localhost -t "to/irc/${CHANNEL}/notice" -m "Pleased to meet you ${NAME} :)"
  fi
done