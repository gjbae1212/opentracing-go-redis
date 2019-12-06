#!/bin/bash

go test -v $(go list ./... | grep -v vendor) --count 1 -covermode=atomic
