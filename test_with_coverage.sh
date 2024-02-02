#!/bin/bash

coverageFile="c.out"
# To create a coverage report you have to run the test with this flags
go test ./13-writing-tests/... -coverprofile=$coverageFile
# Then you need to run this to open it in the browser
go tool cover -html=$coverageFile
