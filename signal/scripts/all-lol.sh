#!/usr/bin/env bash

on_signal() {
    echo "signal received - ignoring lol"
}

trap on_sigint SIGINT
trap on_signal SIGTERM
trap on_signal SIGTSTP

while true; do
    echo tick
    sleep 1
done
