# CARRIER.MICROSERVICES.GO

This repository holds the source code for a microservice used to store and send communications (such as emails or SMS texts), written in Go using the Serverless framework for cloud management and deployment. The purpose is to allow a service (such as a customer-facing API) to send communications to users reliably while reducing latency. This service can send a communication immediately or queue it up for asynchronous sending, and if there is a failure to send it will re-attempt to send again at increasing intervals.

Local development is run on a local virtual machine managed by Vagrant. To install the virtual machine, make sure you have installed Vagrant (https://www.vagrantup.com/) and a virtual machine provider, such as VirtualBox (https://www.virtualbox.org/).

## Manage Local Development Environment

### Provision Virtual Machine

Sets up the local development environment.

```ssh
> vagrant up
> vagrant ssh
$ cd /vagrant/services/email
$ ./scripts/build.sh
```

#### Configure AWS CLI

In order to use Serverless with AWS, you will need to configure your AWS CLI client from inside the VM:

```ssh
$ aws configure
```

## Service: Email

The service uses Sparkpost (https://www.sparkpost.com/) to deliver emails, but other service providers could be added.

### Configure

The service uses `.env` files to configure custom values in the `serverless.yml` configuration file. It is recommended to create `.env` files for each environment (development, staging, production, etc.). Local development uses `.env`. Use a [template](services/email/.env.template) similar to the following (make sure to change the values to reflect your situation):

```
DOMAIN=domain.com
PREFIX=aws-com-domain
REGION=us-east-1

FUNCTION_TIMEOUT=180
TABLE_READ_CAPACITY_UINTS=1
TABLE_WRITE_CAPACITY_UINTS=1
INDEX_READ_CAPACITY_UINTS=1
INDEX_WRITE_CAPACITY_UINTS=1

LOG_LEVEL=info
LOG_ENCODING=json
API_KEY=
DYNAMODB_ENDPOINT=
SPARKPOST_API_KEY=
JOB_SEND_LIMIT=25
RETRY_LIMIT=5
```

Options for LOG_LEVEL:

* fatal
* panic
* dpanic
* error
* warn
* info
* debug

Options for LOG_ENCODING:

* json
* console

The API_KEY parameter is optional, but if provided will be used during authorization as the "X-API-KEY" header.

The DYNAMODB_ENDPOINT parameter should be set to "http://172.29.5.102:8000" for local development if using the local dynamodb plugin, otherwise it should be left blank.

#### Authentication

If you set an `API_KEY` value in your `.env` file, then you must add an `X-API-KEY` header with each Lambda request set to that value. If you want to use more fine-grained permissions, look into using AWS API Gateway authentication patterns. If you do not want to use API Key authentication, then leave `API_KEY` blank. The examples below assume no authentication for simplicity.

### Install Dependencies

```ssh
$ cd /workspace/services/email
$ ./scripts/build.sh
```

### Start local DynamoDB instance

You can run a local instance of DynamoDB to make local development faster and easier using the [serverless-dynamodb-local](https://github.com/99x/serverless-dynamodb-local) plugin . In a new terminal session in the VM run the following command:

```
$ sls dynamodb start --migrate
```

### Compile

```ssh
$ cd /workspace/services/email
$ make build
```

### Local Invocation

#### GET /emails

Use the following to perform a local smoke test to get a list of emails in the local database:

```ssh
$ cd /workspace/services/email
$ sls invoke local --function email --data '{"httpMethod":"GET", "path":"emails", "queryStringParameters": {}}'
```

#### POST /emails

Use the following to perform a local smoke test to create new emails in the local database:

```ssh
$ cd /workspace/services/email
$ sls invoke local --function email --data '{"httpMethod":"POST", "path":"emails", "body":"{\"emails\":[{\"recipients\":[\"test@test.com\"],\"template\":\"template-1\",\"substitutions\":{\"key1\":\"value 1\",\"key2\":\"value 2\"},\"priority\":1}]}", "queryStringParameters": {}}'
```

#### GET /email/{id}

Use the following to perform a local smoke test to read an existing email from the local database (replace the ID with an actual value from your data set):

```ssh
$ cd /workspace/services/email
$ sls invoke local --function email --data '{"httpMethod":"GET", "path":"email/e0b3ef86-b4a6-4ab5-9036-4c7bdbb9f35d", "queryStringParameters": {}}'
```

#### PUT /email/{id}

Use the following to perform a local smoke test to update an existing email in the local database (replace the ID with an actual value from your data set):

```ssh
$ cd /workspace/services/email
$ sls invoke local --function email --data '{"httpMethod":"PUT", "path":"email/e0b3ef86-b4a6-4ab5-9036-4c7bdbb9f35d", "body":"{\"recipients\":[\"testa@test.com\"],\"template\":\"template-1a\",\"substitutions\":{\"key1\":\"value 1a\",\"key2\":\"value 2a\"},\"send_status\":2,\"queued\":\"0001-01-01T00:00:00+0000\",\"priority\":2}", "queryStringParameters": {}}'
```

#### DELETE /email/{id}

Use the following to perform a local smoke test to delete an existing email from the local database (replace the ID with an actual value from your data set):

```ssh
$ cd /workspace/services/email
$ sls invoke local --function email --data '{"httpMethod":"DELETE", "path":"email/e0b3ef86-b4a6-4ab5-9036-4c7bdbb9f35d", "queryStringParameters": {}}'
```

### Linters

List of linters supplied with project:

* gofmt (https://golang.org/cmd/gofmt/)
* go vet (https://golang.org/cmd/vet/)
* golint (https://github.com/golang/lint)
* gosec (https://github.com/securego/gosec)

```ssh
$ cd /vagrant/service/email
$ ./scripts/lint.sh
```

### Deployment

Deploy to the development environment:

```ssh
$ cd /vagrant/services/email
$ sls deploy --stage development
```

Deploy to the production environment:

```ssh
$ cd /vagrant/services/email
$ sls deploy --stage production
```

## Additional Documentation

* [API Documentation](documentation/api-documentation.md)
* [Working with Vagrant](documentation/vagrant.md)

## Repository Directory Structure

| Directory/File                | Purpose                                                                            |
| ----------------------------- | ---------------------------------------------------------------------------------- |
| `services/`                   | Contains all source code files required for the services                           |
| `└─email/`                    | Contains the source code for the Emails service                                    |
| ` · ├─bin/`                   | Contains compiled service binaries                                                 |
| ` · ├─scripts/`               | Contains scripts to build the service, run linters, and any other useful tools     |
| ` · ├─src/`                   | Contains source code for all of the Emails microservices                           |
| ` · ├─.env.template`          | Template for `.env` files                                                          |
| ` · ├─go.mod`                 | Dependency requirements                                                            |
| ` · ├─Makefile`               | Instructions for `make` to build service binaries                                  |
| ` · └─serverless.yml`         | Serverless framework configuration file                                            |
| `documentation/`              | Documentation files                                                                |
| `provision/`                  | Provision scripts for local virtual machine                                        |
| `LICENSE`                     | The license that governs usage of the this source code                             |
| `README.md`                   | This file                                                                          |
| `Vagranfile`                  | Configuration file for Vagrant when provisioning local development virtual machine |
