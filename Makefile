
# These settings are written to config.json.
# You will need to change REGION and REGISTRY.
#
APPNAME	= $(shell basename `pwd`)
REGION	= us-west-2
REGISTRY= 519021347300.dkr.ecr.us-west-2.amazonaws.com
TABLE	= Rentals

default: run

lint:
	go fmt ./...
	go fix ./...
	go vet ./...
	go vet -vettool=$$(which shadow) ./...

generate: config.json pkg/schema/schema_graphql.go

config.json: Makefile
	echo '{"appname":"$(APPNAME)","region":"$(REGION)","registry":"$(REGISTRY)","table":"$(TABLE)"}' \
	| jq . > $@

pkg/schema/schema_graphql.go: api/schema.graphql
	go generate ./pkg/schema

# Run the unit tests against the local DynamoDB
test:	generate ddb-start
	env CONFIG=`pwd`/config.json \
	    AWS_DDB_ENDPOINT=http://localhost:8000 \
	go test -test.v ./...

build:	test
	go build

run:	build
	env CONFIG=`pwd`/config.json \
	    AWS_SDK_LOAD_CONFIG=1 \
	    AWS_DDB_ENDPOINT=http://localhost:8000 \
	./$(APPNAME)

#===================================================================================================
# Docker targets

TIER	= latest
IMAGE	= $(APPNAME):$(TIER)

# The CGO_ENABLED=0 is needed to ensure that we build a static executable,
# which is required to build a Docker image from 'scratch'.
build-linux: generate
	CGO_ENABLED=0 GOOS=linux go build

docker-build: build-linux
	docker build --tag $(IMAGE) .

# Re https://docs.aws.amazon.com/AmazonECS/latest/developerguide/docker-basics.html
docker-login:
	aws ecr get-login-password --region $(REGION) \
	| docker login --username AWS --password-stdin $(REGISTRY)

docker-push: docker-build docker-login
	docker tag  $(IMAGE) $(REGISTRY)/$(IMAGE)
	docker push $(REGISTRY)/$(IMAGE)

docker-pull: docker-login
	docker pull $(REGISTRY)/$(IMAGE)

docker-run: docker-build
	docker run --rm \
	--env     AWS_SDK_LOAD_CONFIG=1 \
	--volume  ~/.aws:/root/.aws:ro  \
	--publish 8080:8080 \
	$(IMAGE)

docker-stop:
	docker stop `docker ps | grep $(IMAGE) | cut -d ' ' -f 1`


#=================================================================================
# Test the GraphQL API with curl

API_IP	= 127.0.0.1
API_URL	= http://$(API_IP):8080/query

test-api: get-customer get-movie get-store get-store-customers get-store-movies get-store-movies-year get-store-movies-title

get-customer:
	curl -s -XPOST -d '{"query": "{ customer(phone: \"815-717-3861\") { phone storephone contact { firstname lastname } } }"}' $(API_URL) \
	| jq -M .

get-movie:
	curl -s -XPOST -d '{"query": "{ movie(year: 2013, title: \"Rush\") { year title info { directors rating genres plot rank actors } } }"}' $(API_URL) \
	| jq -M .

get-store:
	curl -s -XPOST -d '{"query": "{ store(phone: \"828-555-1249\") { phone name location  { address city state zip } } }"}' $(API_URL) \
	| jq -M .

get-store-customers:
	curl -s -XPOST -d '{"query": "{ store(phone: \"828-555-1249\") { phone name customers { phone contact { firstname lastname } } } }"}' $(API_URL) \
	| jq -M .

get-store-movies:
	curl -s -XPOST -d '{"query": "{ store(phone: \"828-555-1249\") { phone name movies(year:0, title:\"\") { year title count } } }"}' $(API_URL) \
	| jq -M .

get-store-movies-year:
	curl -s -XPOST -d '{"query": "{ store(phone: \"828-555-1249\") { phone name movies(year:2014, title:\"\") { year title count } } }"}' $(API_URL) \
	| jq -M .

get-store-movies-title:
	curl -s -XPOST -d '{"query": "{ store(phone: \"828-555-1249\") { phone name movies(year:2014, title:\"X\") { year title count } } }"}' $(API_URL) \
	| jq -M .

.PHONY: get-customer get-movie get-store get-store-customers get-store-movies get-store-movies-year get-store-movies-title

#=================================================================================
# Local DynamoDB
# https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.DownloadingAndRunning.html
#
# 'make ddb-start' downloads the movie data and the local DynamoDB to ./data/,
# runs the local DynamoDB in the background, creates the table "Rentals",
# and populates it with data.
#
PID=	/tmp/dynamo_db.pid

ddb-start: data $(PID)

data:	data/moviedata.json.gz data/DynamoDBLocal.jar

data/moviedata.json.gz:
	cd data \
	; curl -s -O https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/samples/moviedata.zip \
	; unzip moviedata.zip ; gzip moviedata.json ; rm moviedata.zip

data/DynamoDBLocal.jar:
	cd data \
	; curl -s -O https://s3.us-west-2.amazonaws.com/dynamodb-local/dynamodb_local_latest.tar.gz \
	; tar xf dynamodb_local_latest.tar.gz \
	; rm dynamodb_local_latest.tar.gz

$(PID):	# Run the local DynamoDB under nohup. Use -sharedDb to persist the database to a file.
	nohup java -Djava.library.path=data/DynamoDBLocal_lib -jar data/DynamoDBLocal.jar -inMemory \
	& echo $$! > $@
	sleep 1
	kill -0 `cat $@`
	env AWS_DDB_ENDPOINT=http://localhost:8000 make create-table

ddb-stop:
	if [ -e $(PID) ] ; then kill `cat $(PID)` ; rm $(PID) ; fi
	cat /dev/null > nohup.out

#=================================================================================
# DynamoDB
# To operate on the local DynamoDB, set AWS_DDB_ENDPOINT=http://localhost:8000 or:
#
#     env AWS_DDB_ENDPOINT=http://localhost:8000 make ...
#
# This environment variable also affects the Go code via pkg/config/config.go.
#
END=	$(shell if [ -n "$$AWS_DDB_ENDPOINT" ] ; then echo "--endpoint $$AWS_DDB_ENDPOINT" ; else echo "" ; fi )

create-table: config.json data/moviedata.json.gz
	go run cmd/table_create/main.go
	aws dynamodb wait table-exists --table-name $(TABLE) $(END) | cat
	go run cmd/movies_load/main.go < data/moviedata.json.gz
	go run cmd/stores_load/main.go
	go run cmd/customers_load/main.go
	go run cmd/movie_rent/main.go

delete-table:
	aws dynamodb delete-table --table-name $(TABLE) $(END) \
	| cat
	aws dynamodb wait table-not-exists --table-name $(TABLE) $(END) \
	| cat

describe-table:
	aws dynamodb describe-table $(END) \
	--table-name $(TABLE) \
	| cat

list-tables:
	aws dynamodb list-tables $(END) \
	| cat

get-item:
	aws dynamodb get-item $(END) \
	--table-name $(TABLE) \
	--key '{"PK": {"S": "CUS#828-234-1717"}, "SK": {"S": "CONTACT"}}' \
	| cat


query:	# Query for GSI2PK = STO#<phone> to get the store's customers.
	aws dynamodb query $(END) \
	--table-name $(TABLE) \
        --index-name GSI2 \
	--key-condition-expression "GSI2PK = :k" \
	--expression-attribute-values  '{":k":{"S":"STO#310-555-8800"}}' \
	| cat

scan:	# Scan for SK = "INFO" to get all movies
	aws dynamodb scan $(END) \
	--table-name $(TABLE) \
        --filter-expression "SK = :sk" \
	--expression-attribute-values '{":sk":{"S": "INFO"} }' \
	| cat

#=================================================================================
# ECS Fargate
#
# The ecs-deploy target creates or updates the CloudFormation stack,
# but the ECS Service initially has a task with DesiredCount equal to zero.
# The ecs-start target changes the DesiredCount to one. 
# 
ecs-deploy:
	scripts/cf_stack_create | cat
	@ echo Waiting for stack creation to complete...
	scripts/cf_stack_wait stack-create-complete | cat

ecs-start:
	scripts/ecs_service_update 1 | cat

# Use this target to test the Fargate service.
# We extract the load balancer ARN from the stack outputs,
# and use that to obtain the ELB's DNS name.
test-api-fargate:
	make API_IP=$$(aws cloudformation describe-stacks --stack-name "graphql-api-demo-stack" \
	| jq -M -r '.Stacks[0].Outputs[] | select(.OutputKey == "EcsElbName") | .OutputValue' \
	| xargs aws elbv2 describe-load-balancers --load-balancer-arns \
	| jq -M -r '.LoadBalancers[0].DNSName' \
	) test-api

ecs-stop:
	scripts/ecs_service_update 0 | cat

ecs-delete:
	scripts/cf_stack_delete | cat
	@ echo Waiting for stack deletion to complete...
	scripts/cf_stack_wait stack-delete-complete | cat
