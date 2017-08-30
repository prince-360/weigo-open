build:
	cd api && go install .
	cd ui && npm run build

build-api:
	cd api && go install .
