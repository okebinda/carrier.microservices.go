# CARRIER.MICROSERVICES.GO

This repository holds the source code for a microservice used to store and send emails, written in Go using the Serverless framework for cloud management and deployment.

Local development is run on a local virtual machine managed by Vagrant. To install the virtual machine, make sure you have installed Vagrant (https://www.vagrantup.com/) and a virtual machine provider, such as VirtualBox (https://www.virtualbox.org/).

## Manage Local Development Environment

### Provision Virtual Machine

Sets up the local development environment.

```ssh
> vagrant up
> vagrant ssh
```

#### Configure AWS CLI

In order to use Serverless with AWS, you will need to configure your AWS CLI client from inside the VM:

```ssh
$ aws configure
```

### Start Virtual Machine

Starts the local development environment and logs in to the virtual machine. This is a prerequisite for any following steps if the machine is not already booted.

```ssh
> vagrant up
> vagrant ssh
```

### Stop Virtual Machine

Stops the local development environment. Run this command from the host (i.e. log out of any virtual machine SSH sessions).

```ssh
> vagrant halt
```

### Delete Virtual Machine

Deletes the local development environment. Run this command from the host (i.e. log out of any virtual machine SSH sessions).

```ssh
> vagrant destroy
```

Sometimes it is useful to completely remove all residual Vagrant files after destroying the box, in this case run the additional command:

```ssh
> rm -rf ./vagrant
```

## Service: Email

### Configure

The service uses `.env` files to configure custom values in the `serverless.yml` configuration file. It is recommended to create `.env` files for each environment (dev, stage, prod, etc.) using a template similar to the following (make sure to change the values to reflect your situation):

```
DOMAIN=domain.com
PREFIX=aws-com-domain
REGION=us-east-1
API_KEY=
```

### Install Dependencies

```ssh
$ cd /workspace/services/email
$ ./scripts/build.sh
```

### Compile

```ssh
$ cd /workspace/services/email
$ make build
```

### Local Invocation

#### GET Emails

Use the following to perform a local smoke test for the Upload URL function:

```ssh
$ cd /workspace/services/email
$ sls invoke local --function email --data '{"httpMethod":"GET", "path":"emails", "queryStringParameters": {}}'
```

```ssh
$ cd /workspace/services/email
$ sls invoke local --function email --data '{"httpMethod":"GET", "path":"email/ef7f232d-e100-4552-9b80-a6fd587ade36", "queryStringParameters": {}}'
```

### Use

#### 0) Authentication

If you set an `API_KEY` value in your `.env` file, then you must add an `X-API-KEY` header with each Lambda request set to that value. If you want to use more fine-grained permissions, look into using AWS API Gateway authentication patterns. If you do not want to use API Key authentication, then leave `API_KEY` blank. The examples below assume no authentication for simplicity.

#### 1) Get a List of Emails

TBD


### Deployment

Deploy to the development environment:

```ssh
$ cd /vagrant/services/email
$ sls deploy --stage dev
```

Deploy to the production environment:

```ssh
$ cd /vagrant/services/email
$ sls deploy --stage prod
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

## Repository Directory Structure

| Directory/File                | Purpose                                                                            |
| ----------------------------- | ---------------------------------------------------------------------------------- |
| `services/`                   | Contains all source code files required for the services                           |
| `└─email/`                    | Contains the source code for the Emails service                                    |
| ` · ├─bin/`                   | Contains compiled service binaries                                                 |
| ` · ├─scripts/`               | Contains scripts to build the service, run linters, and any other useful tools     |
| ` · ├─src/`                   | Contains source code for all of the Emails microservices                           |
| ` · ├─go.mod`                 | Dependency requirements                                                            |
| ` · ├─Makefile`               | Instructions for `make` to build service binaries                                  |
| ` · └─serverless.yml`         | Serverless framework configuration file                                            |
| `data/`                       | Contains additional resources, such as sample images                               |
| `documentation/`              | Documentation files                                                                |
| `provision/`                  | Provision scripts for local virtual machine                                        |
| `scripts/`                    | Contains various scripts                                                           |
| `LICENSE`                     | The license that governs usage of the this source code                             |
| `README.md`                   | This file                                                                          |
| `Vagranfile`                  | Configuration file for Vagrant when provisioning local development virtual machine |
