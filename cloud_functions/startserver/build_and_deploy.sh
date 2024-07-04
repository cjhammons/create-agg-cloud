# Build the go binary
go build -o main

# Deploy to GCP
gcloud functions deploy go-http-function \
	--gen2 \
	--runtime=go121 \
	--region=REGION \
	--source=. \
	--entry-point=HelloGet \
	--trigger-http 

