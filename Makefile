RELEASE=0.0.1
APP_NAME=app
IMAGE_PROJECT=cloudylabs-public

test:
	go test -v ./...

run:
	go run -v main.go

mod:
	go mod tidy
	go mod vendor

image: mod
	gcloud builds submit \
		--project $(IMAGE_PROJECT) \
		--tag gcr.io/$(IMAGE_PROJECT)/$(APP_NAME):$(RELEASE)


