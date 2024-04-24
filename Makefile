.PHONY: mock
mock:
	@mockgen -source=./internal/service/authcode.go -package=servicemock -destination=./internal/service/mocks/authcode.mock.go