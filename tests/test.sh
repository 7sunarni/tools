# !/bin/bash

cd gopls

go run main.go references ../tests/main.go:20:34

go run main.go references ../tests/main.go:13:38
