.PHONY: mock
mock:
	@mockgen -source=./internal/service/authcode.go -package=servicemock -destination=./internal/service/mocks/authcode.mock.go
	@mockgen -source=./internal/service/user.go -package=servicemock -destination=./internal/service/mocks/user.mock.go