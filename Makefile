prepare:
	go install github.com/swaggo/swag/cmd/swag@latest
swag:
	~/go/bin/swag init -g main.go --output docs/f_swag
run:
	go run main.go
open-docs:
	open http://localhost:3000/swagger/index.html
