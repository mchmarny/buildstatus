
.PHONY: event mod


mod:
	go mod tidy
	go mod vendor


image: mod
	gcloud builds submit \
		--project cloudylabs-public \
		--tag gcr.io/cloudylabs-public/cloud-build-status:0.1.1


