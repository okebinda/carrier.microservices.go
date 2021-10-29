# API Documentation
## CARRIER.MICROSERVICES.GO, Version 1.0

The API uses a REST interface using JSON responses. It uses standard HTTP response codes, verbs and authentication. All endpoints should use HTTPS for security and privacy.

<br><br>

## Table of Contents

* [Definitions](#definitions)
* [Authentication](#authentication)
* [Emails](#emails)

<br><br>

## Definitions

### HTTP Status Codes

The following table describes the HTTP status codes supported by this API.

| Code | Short Description     | Long Description                                                                                 |
| ---- | --------------------- | ------------------------------------------------------------------------------------------------ |
| 200  | OK                    | Everything worked correctly as intended.                                                         |
| 201  | Created               | A new resource was successfully created.                                                         |
| 204  | No Content            | There is no content or a resource was successfully removed.                                      |
| 400  | Bad Request           | The request cannot be completed due to client error. Fix the errors before reattempting request. |
| 401  | Unauthorized          | The request cannot be completed because the client is not authenticated.                         |
| 404  | Not Found             | The does not exist or is currently not available.                                                |
| 405  | Method Not Allowed    | The HTTP verb (GET, POST, etc.) is not supported by the requested resource.                      |
| 500  | Internal Server Error | There was an unexpected error on the server.                                                     |

### Endpoints

As a RESTful API, resource endpoints (URLs) are one of the most important parts of the interface. While the application describes "paths" to resources, these are not the complete endpoints. The system also prepends the protocol, domain name, and stage name (e.g.: "development", "production", etc.) to the path to produce the final endpoint.

##### Format

```
{PROTOCOL}://{DOMAIN_NAME}/{STAGE_NAME}{PATH}
```

### Timestamps

Timestamps should all be formatted to the ISO 8601 datetime standard.

Using Moment.js as an example:

```javascript
moment().format("YYYY-MM-DDThh:mm:ssZZ");  // returns similar to: 2020-11-01T12:34:56+0000
```

### Constants

#### Send Status

| Code | Short Description | Long Description                                                                                                    |
| ---- | ----------------- | ------------------------------------------------------------------------------------------------------------------- |
| 1    | Queued            | The message is currently enqueued and will be (re)attempted after its `queued` timestamp according to its priority. |
| 2    | Processing        | The message has been retrieved from the queue and the system is currently attempting to send.                       |
| 3    | Complete          | The message has been successfully sent and is no longer in the send queue.                                          |
| 4    | Failed            | The message could not be sent and will no longer be attmpted, it is no longer in the send queue.                    |

#### Priority

Cients tell the API what priority the message should have. Lower codes have higher priority and will be attempted earlier than higher codes. Prioritization will only become apparent if there are many messages in the send queue. In general, messages that should be sent to a user based on an action they just took should have a higher priority than messages that are addressed to other users who are not currently interacting with the system.

For example, if a user requests a password reset email they are probably waiting for the email, so it should be given a priority of 0 or 1. On the other hand, if a user leaves a reply to another user's post and that user should get an email notification, it is OK if the message is not sent immediately and the priority should be 2 or 3.

| Code | Description                                                                                                                                                                                 |
| ---- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 0    | Attempt to send message immediately, and only add to the send queue if the send fails. This is a synchronous request (the API response will block until the message service has responded). |
| 1    | Add message to the send queue with a high priority.                                                                                                                                         |
| 2    | Add message to the send queue with a medium priority.                                                                                                                                       |
| 3    | Add message to the send queue with a low priority.                                                                                                                                          |

<br><br>

## Authentication

### API Keys

All requests to any endpoint must contain a valid API key as a `X-API-KEY` header. API keys are 32 character strings created and provided by the system owner(s).

##### Request

| HTTP       | Value                                         |
| ---------- | --------------------------------------------- |
| Methods    | *                                             |
| Paths      | *                                             |
| Headers    | - `X-API-KEY`: 32 character string (required) |

##### Errors

| Code | Description        | Notes                                                                                                    |
| ---- | ------------------ | -------------------------------------------------------------------------------------------------------- |
| 401  | Permission denied  | Add or change the `X-API-KEY` header to a correct key. Request a new one from system owner if necessary. |

##### Example

###### Request

```ssh
curl -H "X-API-KEY: Qyk69zBq4ksCCCCr3ZMwgBLqgKgK2UEY" https://1234abcd.execute-api.us-east-1.amazonaws.com/production/emails
```

_For brevity the `X-API-KEY` header will be ignored for the rest of the documentation, but its requirements still apply._

<br><br>

## Emails

### List Emails

Use the following to read a list of emails.

##### Request

| HTTP            | Value                                                                                                                     |
| --------------- | ------------------------------------------------------------------------------------------------------------------------- |
| Method          | GET                                                                                                                       |
| Paths           | /emails                                                                                                                   |
| URL Parameters  | - `page`: Integer; Results page number; Default: 1<br>- `limit`: Integer; Number of results per page to show; Default: 25 |
| Headers         | - `X-API-KEY`                                                                                                             |

##### Response Codes

| Code | Description       | Notes                                                             |
| ---- | ----------------- | ----------------------------------------------------------------- |
| 200  | OK                | Request successful.                                               |
| 400  | Bad Request       | There was a problem with the request, check the query parameters. |
| 401  | Permission denied | Add an API Key header with a valid key, try again.                |
| 500  | Server error      | Generic application error. Check application logs.                |

##### Response Payload

| Key                          | Type      | Value                                                                                                                          |
| ---------------------------- | --------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `emails`                     | object[]  | The top-level email list resource.                                                                                             | 
| `emails`[].`id`              | string    | The email's system ID.                                                                                                         |
| `emails`[].`service_id`      | string    | The ID of the send event supplied by the 3rd party email service.                                                              |
| `emails`[].`recipients`      | string[]  | A list if email addresses to send to.                                                                                          |
| `emails`[].`template`        | string    | The ID of the email template stored in the 3rd party email service.                                                            |
| `emails`[].`substitutions`   | object    | A map of placeholder:values to add dynamic content to the email template.                                                      |
| `emails`[].`send_status`     | integer   | The status of the email: [1, 2, 3, 4].                                                                                         |
| `emails`[].`queued`          | timestamp | The date/time after which a queued email will be sent. Null timestamps (0001-01-01...) indicate the email is not in the queue. |
| `emails`[].`priority`        | integer   | The priority of the email: [0, 1, 2, 3].                                                                                       |
| `emails`[].`attempts`        | integer   | The number of times the system has attempted to send the email.                                                                |
| `emails`[].`accepted`        | integer   | The number of recipients (email addresses) that were accepted for transmission.                                                |
| `emails`[].`rejected`        | integer   | The number of recipients (email addresses) that were rejected for transmission.                                                |
| `emails`[].`last_attempt_at` | timestamp | The date/time of the last attempted transmission.                                                                              |
| `emails`[].`created_at`      | timestamp | The date/time the email record was created.                                                                                    |
| `emails`[].`updated_at`      | timestamp | The date/time the email record was last udpated.                                                                               |
| `limit`                      | integer   | The limit of items to show on a single page.                                                                                   |
| `page`                       | integer   | The current list page number.                                                                                                  |

##### Example

###### Request

```ssh
curl https://1234abcd.execute-api.us-east-1.amazonaws.com/production/emails?page=2&limit=10
```

###### Response

```json
{
    "emails": [
        {
            "id": "d2387b46-17dd-403d-8470-03bb0d648e1d",
            "service_id": "7023322421558193628",
            "recipients": [
                "bsmith@test.com"
            ],
            "template": "welcome-template-1",
            "substitutions": {
                "email_address": "bsmith@test.com",
                "first_name": "Bob",
                "last_name": "Smith",
                "username": "bsmith"
            },
            "send_status": 3,
            "queued": "0001-01-01T00:00:00+0000",
            "priority": 0,
            "attempts": 1,
            "accepted": 1,
            "rejected": 0,
            "last_attempt_at": "2021-10-26T21:11:44+0000",
            "created_at": "2021-10-26T21:11:30+0000",
            "updated_at": "2021-10-26T21:11:44+0000"
        },
        {
            "id": "cca8ebdd-b7ad-4b2b-827c-83353de62262",
            "service_id": "7023353418337201000",
            "recipients": [
                "jdoe@test.com",
                "jdoe2@test.com"
            ],
            "template": "invitation-template-1",
            "substitutions": {
                "name": "Jane Doe",
                "invited_by": "Fred Brown",
                "invitation_code": "ABC123"
            },
            "send_status": 3,
            "queued": "0001-01-01T00:00:00+0000",
            "priority": 2,
            "attempts": 1,
            "accepted": 2,
            "rejected": 0,
            "last_attempt_at": "2021-10-27T01:10:09+0000",
            "created_at": "2021-10-27T01:10:09+0000",
            "updated_at": "2021-10-27T01:10:10+0000"
        }
    ],
    "page": 2,
    "limit": 10
}
```

### Read an Email

Use the following to read the information for a specific email.

##### Request

| HTTP            | Value                                           |
| --------------- | ----------------------------------------------- |
| Method          | GET                                             |
| Path            | /email/{id}                                     |
| Path Parameters | - `id`: String; The system ID for the resource  |
| Headers         | - `X-API-KEY`                                   |

##### Response Codes

| Code | Description       | Notes                                                             |
| ---- | ----------------- | ----------------------------------------------------------------- |
| 200  | OK                | Request successful.                                               |
| 401  | Permission denied | Add an API Key header with a valid key, try again.                |
| 404  | Not Found         | No email matching the supplied ID was found.                      |
| 500  | Server error      | Generic application error. Check application logs.                |

##### Response Payload

| Key                       | Type      | Value                                                                                                                          |
| ------------------------- | --------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `email`                   | object    | The top-level email resource.                                                                                                  | 
| `email`.`id`              | string    | The email's system ID.                                                                                                         |
| `email`.`service_id`      | string    | The ID of the send event supplied by the 3rd party email service.                                                              |
| `email`.`recipients`      | string[]  | A list if email addresses to send to.                                                                                          |
| `email`.`template`        | string    | The ID of the email template stored in the 3rd party email service.                                                            |
| `email`.`substitutions`   | object    | A map of placeholder:values to add dynamic content to the email template.                                                      |
| `email`.`send_status`     | integer   | The status of the email: [1, 2, 3, 4].                                                                                         |
| `email`.`queued`          | timestamp | The date/time after which a queued email will be sent. Null timestamps (0001-01-01...) indicate the email is not in the queue. |
| `email`.`priority`        | integer   | The priority of the email: [0, 1, 2, 3].                                                                                       |
| `email`.`attempts`        | integer   | The number of times the system has attempted to send the email.                                                                |
| `email`.`accepted`        | integer   | The number of recipients (email addresses) that were accepted for transmission.                                                |
| `email`.`rejected`        | integer   | The number of recipients (email addresses) that were rejected for transmission.                                                |
| `email`.`last_attempt_at` | timestamp | The date/time of the last attempted transmission.                                                                              |
| `email`.`created_at`      | timestamp | The date/time the email record was created.                                                                                    |
| `email`.`updated_at`      | timestamp | The date/time the email record was last udpated.                                                                               |

##### Example

###### Request

```ssh
curl https://1234abcd.execute-api.us-east-1.amazonaws.com/production/email/d2387b46-17dd-403d-8470-03bb0d648e1d
```

###### Response

```json
{
    "email": {
        "id": "d2387b46-17dd-403d-8470-03bb0d648e1d",
        "service_id": "7023322421558193628",
        "recipients": [
            "bsmith@test.com"
        ],
        "template": "welcome-template-1",
        "substitutions": {
            "email_address": "bsmith@test.com",
            "first_name": "Bob",
            "last_name": "Smith",
            "username": "bsmith"
        },
        "send_status": 3,
        "queued": "0001-01-01T00:00:00+0000",
        "priority": 0,
        "attempts": 1,
        "accepted": 1,
        "rejected": 0,
        "last_attempt_at": "2021-10-26T21:11:44+0000",
        "created_at": "2021-10-26T21:11:30+0000",
        "updated_at": "2021-10-26T21:11:44+0000"
    }
}
```

### Create an Email

Use the following to create a new email.

##### Request

| HTTP       | Value          |
| ---------- | -------------- |
| Method     | POST           |
| Path       | /emails        |
| Headers    | - `X-API-KEY`  |

##### Request Payload

| Key                   | Type        | Value                                                                     | Validation                                      |
| --------------------- | ----------- | ------------------------------------------------------------------------- | ----------------------------------------------- |
| `recipients`          | string[]    | A list of email addresses to send to.                                     | Required; Minimum 1; Valid email address format |
| `template`            | string      | The ID of the email template to compose content from.                     | Required; Length: 2-255 chars                   |
| `substitutions`       | object      | A map of placeholder:values to add dynamic content to the email template. | -                                               |
| `priority`            | integer     | The priority of the email.                                                | Required; Value: 0-3                            |

##### Response Codes

| Code | Description       | Notes                                                                         |
| ---- | ----------------- | ----------------------------------------------------------------------------- |
| 200  | OK                | Request successful.                                                           |
| 400  | Bad Request       | There was a problem with the request, review errors reported in the response. |
| 401  | Permission denied | Add an API Key header with a valid key, try again.                            |
| 500  | Server error      | Generic application error. Check application logs.                            |

##### Response Payload

| Key                       | Type      | Value                                                                                                                          |
| ------------------------- | --------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `email`                   | object    | The top-level email resource.                                                                                                  | 
| `email`.`id`              | string    | The email's system ID.                                                                                                         |
| `email`.`service_id`      | string    | The ID of the send event supplied by the 3rd party email service.                                                              |
| `email`.`recipients`      | string[]  | A list if email addresses to send to.                                                                                          |
| `email`.`template`        | string    | The ID of the email template stored in the 3rd party email service.                                                            |
| `email`.`substitutions`   | object    | A map of placeholder:values to add dynamic content to the email template.                                                      |
| `email`.`send_status`     | integer   | The status of the email: [1, 2, 3, 4].                                                                                         |
| `email`.`queued`          | timestamp | The date/time after which a queued email will be sent. Null timestamps (0001-01-01...) indicate the email is not in the queue. |
| `email`.`priority`        | integer   | The priority of the email: [0, 1, 2, 3].                                                                                       |
| `email`.`attempts`        | integer   | The number of times the system has attempted to send the email.                                                                |
| `email`.`accepted`        | integer   | The number of recipients (email addresses) that were accepted for transmission.                                                |
| `email`.`rejected`        | integer   | The number of recipients (email addresses) that were rejected for transmission.                                                |
| `email`.`last_attempt_at` | timestamp | The date/time of the last attempted transmission.                                                                              |
| `email`.`created_at`      | timestamp | The date/time the email record was created.                                                                                    |
| `email`.`updated_at`      | timestamp | The date/time the email record was last udpated.                                                                               |

###### Request

```ssh
curl -X POST -H "Content-Type: application/json" \
    -d '{
        "recipients": [
            "jdoe@test.com",
            "jdoe2@test.com"
        ],
        "template": "invitation-template-1",
        "substitutions": {
            "name": "Jane Doe",
            "invited_by": "Fred Brown",
            "invitation_code": "ABC123"
        },
        "priority": 2
    }' \
    https://1234abcd.execute-api.us-east-1.amazonaws.com/production/emails
```

###### Response

```json
{
    "email": {
        "id": "cca8ebdd-b7ad-4b2b-827c-83353de62262",
        "service_id": "",
        "recipients": [
            "jdoe@test.com",
            "jdoe2@test.com"
        ],
        "template": "invitation-template-1",
        "substitutions": {
            "name": "Jane Doe",
            "invited_by": "Fred Brown",
            "invitation_code": "ABC123"
        },
        "send_status": 1,
        "queued": "2021-10-27T01:10:09+0000",
        "priority": 2,
        "attempts": 0,
        "accepted": 0,
        "rejected": 0,
        "last_attempt_at": "0001-01-01T00:00:00+0000",
        "created_at": "2021-10-27T01:10:09+0000",
        "updated_at": "2021-10-27T01:10:10+0000"
    }
}
```

### Update an Email

Use the following to update an existing email.

##### Request

| HTTP            | Value                                           |
| --------------- | ----------------------------------------------- |
| Method          | PUT                                             |
| Path            | /email/{id}                                     |
| Path Parameters | - `id`: String; The system ID for the resource  |
| Headers         | - `X-API-KEY`                                   |

##### Request Payload

| Key                   | Type        | Value                                                                                                                                     | Validation                                      |
| --------------------- | ----------- | ----------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------- |
| `recipients`          | string[]    | A list of email addresses to send to.                                                                                                     | Required; Minimum 1; Valid email address format |
| `template`            | string      | The ID of the email template to compose content from.                                                                                     | Required; Length: 2-255 chars                   |
| `substitutions`       | object      | A map of placeholder:values to add dynamic content to the email template.                                                                 | -                                               |
| `priority`            | integer     | The priority of the email.                                                                                                                | Required; Value: 0-3                            |
| `send_status`         | integer     | The send status of the email.                                                                                                             | Value: 1-4                                      |
| `queued`              | timestamp   | The date/time after which a queued email will be sent. Null timestamps ("0001-01-01T00:00:00+0000") will remove the email from the queue. | Valid timestamp format                          |
| `service_id`          | string      | The ID of the send event supplied by the 3rd party email service.                                                                         | -                                               |

##### Response Codes

| Code | Description       | Notes                                                                         |
| ---- | ----------------- | ----------------------------------------------------------------------------- |
| 200  | OK                | Request successful.                                                           |
| 400  | Bad Request       | There was a problem with the request, review errors reported in the response. |
| 401  | Permission denied | Add an API Key header with a valid key, try again.                            |
| 404  | Not Found         | No email matching the supplied ID was found.                                  |
| 500  | Server error      | Generic application error. Check application logs.                            |

##### Response Payload

| Key                       | Type      | Value                                                                                                                          |
| ------------------------- | --------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `email`                   | object    | The top-level email resource.                                                                                                  | 
| `email`.`id`              | string    | The email's system ID.                                                                                                         |
| `email`.`service_id`      | string    | The ID of the send event supplied by the 3rd party email service.                                                              |
| `email`.`recipients`      | string[]  | A list if email addresses to send to.                                                                                          |
| `email`.`template`        | string    | The ID of the email template stored in the 3rd party email service.                                                            |
| `email`.`substitutions`   | object    | A map of placeholder:values to add dynamic content to the email template.                                                      |
| `email`.`send_status`     | integer   | The status of the email: [1, 2, 3, 4].                                                                                         |
| `email`.`queued`          | timestamp | The date/time after which a queued email will be sent. Null timestamps (0001-01-01...) indicate the email is not in the queue. |
| `email`.`priority`        | integer   | The priority of the email: [0, 1, 2, 3].                                                                                       |
| `email`.`attempts`        | integer   | The number of times the system has attempted to send the email.                                                                |
| `email`.`accepted`        | integer   | The number of recipients (email addresses) that were accepted for transmission.                                                |
| `email`.`rejected`        | integer   | The number of recipients (email addresses) that were rejected for transmission.                                                |
| `email`.`last_attempt_at` | timestamp | The date/time of the last attempted transmission.                                                                              |
| `email`.`created_at`      | timestamp | The date/time the email record was created.                                                                                    |
| `email`.`updated_at`      | timestamp | The date/time the email record was last udpated.                                                                               |

###### Request

```ssh
curl -X PUT -H "Content-Type: application/json" \
    -d '{
        "recipients": [
            "jdoe.a@test.com",
            "jdoe2.b@test.com"
        ],
        "template": "invitation-template-2",
        "substitutions": {
            "name": "Jane A. Doe",
            "invited_by": "Fred B. Brown",
            "invitation_code": "DEF456"
        },
        "priority": 1,
        "send_status": 3,
        "queued": "0001-01-01T00:00:00+0000",
        "service_id": "abcdefg1234567"
    }' \
    https://1234abcd.execute-api.us-east-1.amazonaws.com/production/email/cca8ebdd-b7ad-4b2b-827c-83353de62262
```

###### Response

```json
{
    "email": {
        "id": "cca8ebdd-b7ad-4b2b-827c-83353de62262",
        "service_id": "abcdefg1234567",
        "recipients": [
            "jdoe.a@test.com",
            "jdoe2.b@test.com"
        ],
        "template": "invitation-template-2",
        "substitutions": {
            "name": "Jane A. Doe",
            "invited_by": "Fred B. Brown",
            "invitation_code": "DEF456"
        },
        "send_status": 3,
        "queued": "0001-01-01T00:00:00+0000",
        "priority": 1,
        "attempts": 0,
        "accepted": 0,
        "rejected": 0,
        "last_attempt_at": "0001-01-01T00:00:00+0000",
        "created_at": "2021-10-27T01:10:09+0000",
        "updated_at": "2021-10-27T01:10:12+0000"
    }
}
```

### Delete an Email

Use the following to delete an existing email.

##### Request

| HTTP            | Value                                           |
| --------------- | ----------------------------------------------- |
| Method          | DELETE                                          |
| Path            | /email/{id}                                     |
| Path Parameters | - `id`: String; The system ID for the resource  |
| Headers         | - `X-API-KEY`                                   |

##### Response Codes

| Code | Description       | Notes                                                             |
| ---- | ----------------- | ----------------------------------------------------------------- |
| 204  | No content        | Delete successful.                                                |
| 401  | Permission denied | Add an API Key header with a valid key, try again.                |
| 404  | Not Found         | No email matching the supplied ID was found.                      |
| 500  | Server error      | Generic application error. Check application logs.                |

###### Request

```ssh
curl -X DELETE https://1234abcd.execute-api.us-east-1.amazonaws.com/production/emails
```
