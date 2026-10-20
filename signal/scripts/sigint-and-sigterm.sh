#!/usr/bin/env bash

on_sigint() {
    echo "sigint received - ignore lol"
}

on_sigterm() {
    echo "sigterm received - ignore lol"
}

trap on_sigint SIGINT
trap on_sigterm SIGTERM

while true; do
    echo tick
    sleep 1
done
