#!/bin/sh

# start redis in background
redis-server &

# start tracker
./tracker