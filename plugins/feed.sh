#!/usr/bin/env bash

mosquitto_sub -h localhost -v -t "from/irc/+/+/message" | while read MSG;
  do echo "${MSG}";
done