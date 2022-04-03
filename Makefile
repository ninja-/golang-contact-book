build:
	(cd server/main && go build) && (cd client/main && go build)
unit-test:
	(cd server && go test .)
integration-test:
	(cd client && go test .)