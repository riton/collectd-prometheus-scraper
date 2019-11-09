#!/bin/bash

alias tail_values="docker-compose exec collectd_server tail -f /tmp/value_list.txt"
alias rebuild="docker-compose build"
alias restart="docker-compose up -d --scale whoami=3 --scale dummy_client=3"
