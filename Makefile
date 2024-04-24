.PHONY: mock
mock:
	#mockgen `-source=./internal/service/authcode.go `-package=service `-destination=./internal/service/authcode.mock.go
	@mockgen -source=./internal/service/authcode.go -package=service -destination=./internal/service//authcode.mock.go

	#mockgen `-source=./internal/service/user.go `-package=service `-destination=./internal/service/user.mock.go
	@mockgen -source=./internal/service/user.go -package=service -destination=./internal/service/mocks/user.mock.go

	#mockgen `-source=./internal/repository/authcode.go `-package=repository `-destination=./internal/repository/authcode.mock.go
	@mockgen -source=./internal/repository/authcode.go -package=repository -destination=./internal/repository/authcode.mock.go

	#mockgen `-source=./internal/repository/user.go `-package=repository `-destination=./internal/repository/user.mock.go
	@mockgen -source=./internal/repository/user.go -package=repository -destination=./internal/repository/user.mock.go

	@go mod tidy