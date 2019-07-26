# buildstatus

Simple Cloud Build status notification based build status changes published by Cloud Build to Cloud PubSub topic

* [Cloud Build](https://cloud.google.com/cloud-build/)
* [Cloud Rub](https://cloud.google.com/run/)
* [Cloud Build Trigger Demo](github.com/mchmarny/knative-gitops-using-cloud-build)

## Configuration

Start by capturing few configuration values required to deploy and configuring the notification service:

```shell
PRJ=$(gcloud config get-value project)
PRJ_NUM=$(gcloud projects list --filter="${PRJ}" --format="value(PROJECT_NUMBER)")
```

In addition to the above, we are also going to need a couple of Slack API information.

> Note, `SLACK_BUILD_STATUS_CHANNEL` is the ID of the channel, not its name

```shell
SLACK_API_TOKEN=
SLACK_BUILD_STATUS_CHANNEL=
```

## Service Deploy

Once you define the above variables, you can now deploy that image to Cloud Run using this command:

> Note, for the Slack token, consider using something like [berglas](https://github.com/GoogleCloudPlatform/berglas) to avoid storing secrets in environment variables

```shell
gcloud beta run deploy cloud-build-status \
	--image=gcr.io/cloudylabs-public/cloud-build-status:0.1.1 \
  --region us-central1 \
  --set-env-vars=SLACK_API_TOKEN=$SLACK_API_TOKEN,SLACK_BUILD_STATUS_CHANNEL=$SLACK_BUILD_STATUS_CHANNEL
```

> When prompted to "allow unauthenticated" select "n" for No. That will set that service private so that we can use service account later to enable PubSub to "push" evens to this service.

Once deployment is completed, make sure to capture the URL that's generated for this service.

## Service Account & Authorization

Because this Cloud Run service will only process events pushed from Cloud PubSub, we will need to create a service account and ensure that only that service account is able to invoke the Cloud Run service.

First, let's though enable your project to create Cloud Pub/Sub authentication tokens:

```shell
gcloud projects add-iam-policy-binding $PRJ \
    --member="serviceAccount:service-${PRJ_NUM}@gcp-sa-pubsub.iam.gserviceaccount.com" \
    --role=roles/iam.serviceAccountTokenCreator
```

Then let's create a service account (`buildstatusinvoker`) that will be used by PubSub to invoke our Cloud Run service:

```shell
gcloud iam service-accounts create buildstatusinvoker \
    --display-name "Cloud Run Notification Service Invoker"
```

Now we can create a policy binding for that service account to access our Cloud Run service

```shell
gcloud beta run services add-iam-policy-binding buildstatus \
	--member=serviceAccount:buildstatusinvoker@${PRJ}.iam.gserviceaccount.com \
	--role=roles/run.invoker
```

## PubSub Subscription

To enable PubSub to send topic data to Cloud Run service we will need to create a PubSub topic subscription called `cloud-builds-sub` for `cloud-builds` topic.

Also, since Cloud Run generates service URL including random portion of the service name so let's capture the URL:

```shell
SURL=$(gcloud beta run services describe cloud-build-status --region us-central1 --format 'value(status.domain)')
```

Now we are create the PubSub push subscription:

```shell
gcloud beta pubsub subscriptions create cloud-builds-sub \
	--topic cloud-builds \
	--push-endpoint="${SURL}/" \
	--push-auth-service-account="buildstatusinvoker@${PRJ}.iam.gserviceaccount.com"
```

## Log

When running the `cloud-build-status` service you can see in the Cloud Run service log tab the raw data that was pushed by PubSub subscription to the service and the processed data that was pushed onto the target topic

## Disclaimer

This is my personal project and it does not represent my employer. I take no responsibility for issues caused by this code. I do my best to ensure that everything works, but if something goes wrong, my apologies is all you will get.
