# CARRIER.MICROSERVICES.GO
## Design Document

### Introduction

This project acts as middleware between a customer-facing service (API) and a 3rd party service that sends messages (emails, SMS texts, etc.). I accepts requests to send messages and saves them in a queue to increase reliability and response latency.

## Services

### Email

* REST API
* Cron Job

## Tech Stack

* Go
* AWS Lambda
* AWS DynamoDB
* Serverless Framework
