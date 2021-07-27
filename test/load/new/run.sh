#!/bin/sh

# Enable more open files
ulimit -n 20000

# Run the test under k6
k6 run test.js
