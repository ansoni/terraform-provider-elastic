# Examples

	docker-compose up -d

You should now have an elasticsearch and kibana running on 127.0.0.1:9200 and 127.0.0.1:5601 respectively.

# Run Basic Example

You should have installed the terraform-provider-elastic first.

	cd basic
	terrform init && terraform apply

# Run X-Pack Example

X-Pack Trial or Purchase required (You can sign-up for trial in Kibana GUI - Management::License)

	cd xpack
	terraform init && terraform apply
	
