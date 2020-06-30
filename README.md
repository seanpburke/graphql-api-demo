# GraphQL API Demo

## Overview

This project demonstrates a GraphQL API backed by a Fargate container and DynamoDB.
The GraphQL API is implemented in Go using:

https://github.com/graph-gophers/graphql-go

> The goal of this project is to provide full support of the GraphQL draft specification with a set of idiomatic, easy to use Go packages.

This demonstration uses the AWS services:

- ECS Fargate
- DynamoDB

## DynamoDB Schema

This demo is based on a single-table DynamoDB schema. For background on the single-table approach,
you can refer to these excellent talks:

[Deep Dive: Advanced design patterns: Rick Houlihan - AWS re:Invent 2018](https://www.youtube.com/watch?v=HaEPXoXVf2k)

[AWS re:Invent 2019: Data modeling with Amazon DynamoDB (CMY304)](https://www.youtube.com/watch?v=DIQVJqiSUkE)

This demonstration uses the sample Movie dataset from Amazon's DynamoDB Developer Guide:

https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/samples/moviedata.zip

To the _Movie_ data set, we add _Stores_ and _Customers_, who rent _Movies_ from the _Store_.
To represent these transactions, we populate the DynamoDB table **Rentals** with six item types:


Item Type   | PK            | SK             | GSI2PK        | Attributes
------------|---------------|----------------|---------------|--------------------------------------------
Store       | STO#Phone     | LOCATION       |               | Phone, Name, Location
Inventory   | STO#Phone     | MOV#Year#Title |               | Phone, Year, Title, Count
Customer    | CUS#Phone     | CONTACT        | STO#Phone     | Phone, Contact, StorePhone
Rental      | CUS#Phone     | REN#Phone#Date |               | Phone, Date
Movie       | MOV#Year#Title| INFO           |               | Year, Title, Info
MovieRental | MOV#Year#Title| REN#Phone#Date |               | Year, Title, Phone, Date, DueDate, ReturnDate

Every item has attributes _PK_ (Partition Key) and _SK_ (Sort Key) attribute. The Global Secondary Index _GSI1_
is an inverted index with SK as the partition key and PK as the sort key.  _Customer_ items also have the attribute _GSI2PK_,
the primary key for _GSI2_, which enables us to find all _Customers_ of a given _Store_.
This table is created by `cmd/table_create/main.go`, which calls `pkg/table/table.go:CreateTable()`.

More details on these item types and the DynamoDB queries we make on them, can be found in the source code:

    pkg/schema/customer.go
    pkg/schema/movie.go
    pkg/schema/store.go

The Makefile also demonstrates some of the queries that you can make via the GraphQL API.

### Using a local DynamoDB

The DynamoDB portions of this app are set up to run with a local DynamoDB service.
To learn more about Amazon's "local" DynamoDB, refer to these links:

https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.DownloadingAndRunning.html
https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/DynamoDBLocal.UsageNotes.html


## Running the Demo

### Install Tools

You will need some prerequisites to run this demo.
Before creating the prerequesites, ensure that you install the following tools:

- Install the AWS CLI
- Set up AWS credentials via `aws configure`
- Install golang
- Install jq
- On MacOS, install XCode comand line tools
- Install Docker

This is a quick guide to setting up Docker for use with ECS: [Docker basics for Amazon ECS](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/docker-basics.html)

### AWS Prerequisites

In your AWS account, You will need to set up the following:

- AWS IAM User privileges
- AWS IAM role `ECS-Task-DynamoDB`
- AWS DynamoDB Table
- AWS ECR Repository `graphql-api-demo`
- AWS ECS Default Cluster

Most of these things will fall into the free tier, so there should only be minor charges.

#### AWS IAM User

Your IAM user will need access to these IAM policies:

- AWSCloudFormationFullAccess
- AmazonDynamoDBFullAccess
- AmazonEC2ContainerRegistryPowerUser
- AmazonECS_FullAccess
- AmazonS3FullAccess

#### AWS IAM Role

This demo  also requires you to define the IAM role `ECS-Task-DynamoDB` to be assumed
by the ECS Task. Define the role for the AWS service (trusted entity) **ECS Task**:
```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
```
You must create and attach this policy to the role:
```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "dynamodb:*"
            ],
            "Resource": "arn:aws:dynamodb:*:*:table/*",
            "Effect": "Allow"
        }
    ]
}
```
For more information on defining IAM roles for ECS tasks, see: [IAM Roles for Tasks](https://docs.aws.amazon.com/AmazonECS/latest/developerguide/task-iam-roles.html)

#### AWS DynamoDB Table

To create and populate the DynamoDB table, do:

    make create-table

#### AWS ECR Repository

In the Elastic Container Registry, create a new repository named `graphql-api-demo`,
copy the repository URI, and edit the Makefile to set the correnct REGION and REGISTRY:

    REGION=   us-west-1
    REGISTRY= 12345678900.dkr.ecr.us-west-1.amazonaws.com

Now, build the docker image and push it to your ECR Repository:

    make docker-push

#### AWS ECS Default Cluster

Go to the AWS ECS console, and ensure that you have a Cluster named 'default'.
There should be a button to create it automatically, if it is not already present.

#### AWS CloudWatch Log Group

Go to the AWS CloudWatch console, and create a Log Group named `/ecs/graphql-api-demo`.
(Yes, I should add this to the CloudFormation template.)

### Deploying the Demo

When all of the prerequisites in place, you are ready to create the CloudFormation stack.

    make ecs-deploy

This will create a CloudFormation stack named 'graphql-api-demo-stack'. The stack will
create an ECS Fargate Service 'graphql-api-demo-service' in your default cluster,
and the ECS Task Definition 'graphql-api-demo-task'. The service initially has one
task with its DesiredCount set to zero. To start a task to serve the GraphQL API:

    make ecs-start

Once the task has started, you can exercise the GraphQL API with some sample curl queries:

    make test-api-fargate

### Removing the Demo

When you wish to delete the ECS Task:

    make ecs-stop

To delete the entire CloudWatch stack:

    make ecs-delete

To delete the DynamoDB table:

    make delete-table

### Conclusions

When you use github.com/graph-gophers/graphql-go, the code that you must write to
back the GraphQL API, is little more than you would need for an equivalent REST API.
GraphQL is much more flexible than REST, and also allows you to fetch various related
data in a single API call, which improves client responsiveness.
