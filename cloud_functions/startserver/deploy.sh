# Deploy to GCP
gcloud functions deploy start-server\
	--gen2 \
	--runtime=go121 \
	--region=us-central1 \
	--source=. \
	--entry-point=StartServer \
	--trigger-http 

