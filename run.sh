#!/bin/bash

go run -ldflags -linkmode=external "${1:-./cmd/lunar-fever}"
