#!/bin/bash
# Kill the worker during the Workflow's 30s durable timer, start a fresh one, and compare logs.
# Needs a dev server: temporal server start-dev
set -e
cd "$(dirname "$0")/.."
go build -o bin/worker ./helloworld/worker
go build -o bin/starter ./helloworld/starter

bin/worker > worker1.log 2>&1 & W1=$!
bin/starter > starter.log 2>&1 & ST=$!
until grep -q 'Sleeping 30s' worker1.log; do sleep 0.5; done
echo "$(date +%T) killing worker 1 during the sleep"; kill -9 $W1
sleep 10
echo "$(date +%T) starting worker 2"
bin/worker > worker2.log 2>&1 & W2=$!
wait $ST; echo "$(date +%T) starter exited"; kill $W2
echo "--- worker 2 log (note: no Activity lines)"; cat worker2.log
