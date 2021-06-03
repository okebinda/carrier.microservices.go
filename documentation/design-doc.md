# Carrier Design Document

## Services

### Email

### Endpoints

* POST /emails
  * Source: SQS
  * Fields:
    * To (list)
    * CC (list)
    * Subject (text)
    * From (text)
    * Reply to (text)
    * Body (text)
* GET /emails

### Tech Stack

* AWS Lambda
* Go
* PostgreSQL

