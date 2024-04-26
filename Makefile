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

	# mockgen `-source=./internal/repository/cache/user.go `-package=cache `-destination=./internal/repository/cache/user.mock.go
	@mockgen -source=./internal/repository/cache/user.go -package=cache -destination=./internal/repository/cache/user.mock.go

	# mockgen `-source=./internal/repository/cache/authcode.go `-package=cache `-destination=./internal/repository/cache/authcode.mock.go
	@mockgen -source=./internal/repository/cache/authcode.go -package=cache -destination=./internal/repository/cache/authcode.go

	# mockgen `-source=./internal/repository/dao/user.go `-package=dao `-destination=./internal/repository/dao/user.mock.go
	@mockgen -source=./internal/repository/dao/user.go -package=dao -destination=./internal/repository/dao/user.mock.go

	@go mod tidy