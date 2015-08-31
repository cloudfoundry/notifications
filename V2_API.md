# Notifications V2 Documentation

- System Status
  - [Check service status](#get-info)
- Senders
  - [Creating a sender](#create-sender)
  - [Listing the senders](#list-senders)
  - [Retrieving a sender](#retrieve-sender)
  - [Updating a sender](#update-sender)
- Campaign types
  - [Creating a campaign type](#create-campaign-type)
  - [Retrieving a campaign type](#retrieve-campaign-type)
  - [Listing the campaign types](#list-campaign-types)
  - [Updating a campaign type](#update-campaign-type)
  - [Deleting a campaign type](#delete-campaign-type)
- Templates
  - [Creating a template](#create-template)
  - [Retrieving a template](#retrieve-template)
  - [Deleting a template](#delete-template)
- Campaigns
  - [Sending a campaign](#send-campaign)

## System Status

<a name="get-info"></a>
#### Check service status

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
```

###### Route
```
GET /info
```

###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  http://notifications.example.com/info

HTTP/1.1 200 OK
Connection: close
Content-Length: 13
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 21:29:36 GMT
X-Cf-Requestid: 2cf01258-ccff-41e9-6d82-41a4441af4af

{"version": 2}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields  | Description        |
| ------- | ------------------ |
| version | API version number |

## Senders

<a name="create-sender"></a>
#### Creating a sender

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.

###### Route
```
POST /senders
```
###### Params

| Key    | Description                               |
| ------ | ----------------------------------------- |
| name\* | the human-readable name given to a sender |

\* required

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"name":"my-sender"}'
  http://notifications.example.com/senders

HTTP/1.1 201 Created
Content-Length: 64
Content-Type: text/plain; charset=utf-8
Date: Fri, 17 Jul 2015 19:30:32 GMT
X-Cf-Requestid: ce9f6b5a-317d-4d0f-7197-df63540c7f22

{"id":"4bbd0431-9f5b-49bb-701d-8c2caa755ed0","name":"my-sender"}
```

##### Response

###### Status
```
201 Created
```

###### Body
| Fields | Description                  |
| ------ | ---------------------------- |
| id     | System-generated sender ID   |
| name   | Sender name                  |

<a name="list-senders"></a>
#### Listing Senders

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.

###### Route
```
GET /senders
```

###### CURL Example
```
$ curl -i -X GET \
  -H "Authorization: bearer <CLIENT-TOKEN>" \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  http://notifications.example.com/senders

HTTP/1.1 200 OK
Date: Mon, 17 Aug 2015 21:39:31 GMT
Content-Length: 145
Content-Type: text/plain; charset=utf-8

{"senders":[{"id":"abb9b009-d3c3-4de2-43f0-341e671e2f3d","name":"sender one"},{"id":"379c5b3d-d3ec-4148-6608-68638ca977c5","name":"sender two"}]}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields | Description |
| ------ | ----------- |
| id     | Sender ID   |
| name   | Sender name |

<a name="retrieve-sender"></a>
#### Retrieving a sender

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.

###### Route
```
GET /senders/{senderID}
```

###### Params
| Key          | Description                              |
| -------------| ---------------------------------------- |
| senderID\*   | The "id" returned when creating a sender |

\* required

###### CURL Example
```
$ curl -i -X GET \
  -H "Authorization: bearer <CLIENT-TOKEN>" \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  http://notifications.example.com/senders/4bbd0431-9f5b-49bb-701d-8c2caa755ed0

HTTP/1.1 200 OK
Content-Length: 64
Content-Type: text/plain; charset=utf-8
Date: Fri, 17 Jul 2015 21:00:06 GMT
X-Cf-Requestid: 4fab7338-11ba-44d2-75fd-c34046518dae

{"id":"4bbd0431-9f5b-49bb-701d-8c2caa755ed0","name":"my-sender"}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields | Description |
| ------ | ----------- |
| id     | Sender ID   |
| name   | Sender name |


<a name="update-sender"></a>
#### Updating a sender

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.

###### Route
```
PUT /senders/{senderID}
```

###### Params
| Key          | Description                              |
| -------------| ---------------------------------------- |
| senderID\*   | The "id" returned when creating a sender |

###### Body
| Fields | Description |
| ------ | ----------- |
| name   | Sender name |

\* required

###### CURL Example
```
$ curl -i -X GET \
  -H "Authorization: bearer <CLIENT-TOKEN>" \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -d '{"name": "my-updated-sender"}'
  http://notifications.example.com/senders/64816c21-9eb3-49f0-6489-699cb2defe6b

HTTP/1.1 200 OK
Content-Length: 72
Content-Type: text/plain; charset=utf-8
Date: Fri, 17 Jul 2015 21:00:06 GMT
X-Cf-Requestid: 4fab7338-11ba-44d2-75fd-c34046518dae

{"id":"64816c21-9eb3-49f0-6489-699cb2defe6b","name":"my-updated-sender"}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields | Description |
| ------ | ----------- |
| id     | Sender ID   |
| name   | Sender name |

## Campaign types

<a name="create-campaign-type"></a>
#### Creating a campaign type

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.
\*\* Creation of a critical campaign type requires `critical_notifications.write` scope.

###### Route
```
POST /senders/<sender-id>/campaign-types
```
###### Params

| Key                       | Description                                                         |
| ------------------------- | ------------------------------------------------------------------- |
| name\*                    | the human-readable name given to a campaign type                |
| description\*             | the human-readable description given to a campaign type         |
| critical (default: false) | a flag to indicate whether the campaign type is critical or not |
| template_id               | the ID of a template to use for this campaign type              |

\* required

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"name":"my-campaign-type","description":"campaign type description","critical":false,"template_id":""}'
  http://notifications.mrorange.cfla.cf-app.com/senders/4bbd0431-9f5b-49bb-701d-8c2caa755ed0/campaign_types

HTTP/1.1 201 Created
Content-Length: 155
Content-Type: text/plain; charset=utf-8
Date: Wed, 22 Jul 2015 16:00:37 GMT
X-Cf-Requestid: 6106873b-14ea-4fd9-6418-946c1651e4ac

{"critical":false,"description":"campaign type description","id":"3d9aa963-97bb-4b48-4c3c-ecccad6314f8","name":"my-campaign-type","template_id":""}
```

##### Response

###### Status
```
201 Created
```

###### Body
| Fields        | Description                           |
| ------------- | ------------------------------------- |
| id            | System-generated campaign type ID |
| name          | Campaign type name                |
| description   | Campaign type description         |
| critical      | Critical campaign type flag       |
| template_id   | Template ID                           |

<a name="retrieve-campaign-type"></a>
#### Retrieving A Campaign type

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.

###### Route
```
GET /senders/<sender-id>/campaign-types/<campaign-type-id>
```
###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/senders/4bbd0431-9f5b-49bb-701d-8c2caa755ed0/campaign_types/3369a6ae-22c5-4da9-7081-b35350c79c4c

200 OK
RESPONSE HEADERS:
  Date: Tue, 28 Jul 2015 00:54:54 GMT
  Content-Length: 155
  Content-Type: text/plain; charset=utf-8
  Connection: close
RESPONSE BODY:
{"critical":false,"description":"campaign type description","id":"3369a6ae-22c5-4da9-7081-b35350c79c4c","name":"my-campaign-type","template_id":""}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields             | Description                           |
| ------------------ | ------------------------------------- |
| id                 | System-generated campaign type ID |
| name               | Campaign type name                |
| description        | Campaign type description         |
| critical           | Critical campaign type flag       |
| template_id        | Template ID                           |

<a name="list-campaign-types"></a>
#### Listing Campaign types

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.
\*\* Creation of a critical campaign type requires `critical_notifications.write` scope.

###### Route
```
GET /senders/<sender-id>/campaign-types
```
###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/senders/4bbd0431-9f5b-49bb-701d-8c2caa755ed0/campaign_types

HTTP/1.1 200 OK
Date: Thu, 23 Jul 2015 19:22:46 GMT
Content-Length: 180
Content-Type: text/plain; charset=utf-8

{"campaign_types":[{"critical":false,"description":"campaign type description","id":"702ce4c7-93a0-42b5-4fd5-4d0ed68e2cd7","name":"my-campaign-type","template_id":""}]}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields             | Description                           |
| ------------------ | ------------------------------------- |
| campaign_types | The array of campaign types       |
| id                 | System-generated campaign type ID |
| name               | Campaign type name                |
| description        | Campaign type description         |
| critical           | Critical campaign type flag       |
| template_id        | Template ID                          |

<a name="update-campaign-type"></a>
##### Update a Campaign Type

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```

\* The token requires `notifications.write` scope.
\*\* Updating a critical campaign type requires `critical_notifications.write` scope.

###### Route
```
PUT /senders/<sender-id>/campaign_types/<campaign-type-id>
```

###### CURL example
```
$ curl -i -X PUT \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  --data '{"name": "new campaign type", "description": "new campaign description", "critical": true}' \
  http://notifications.example.com/senders/a6c38f92-8fa9-488b-4f4c-7f4d4e0c0fd2/campaign_types/5cbc4458-3dba-481b-74c3-4548114b830b

HTTP/1.1 200 OK
Content-Length: 146
Content-Type: text/plain; charset=utf-8
Date: Tue, 04 Aug 2015 20:47:35 GMT

{"critical":true,"description":"new campaign description","id":"5cbc4458-3dba-481b-74c3-4548114b830b","name":"new campaign type","template_id":""}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields             | Description                           |
| ------------------ | ------------------------------------- |
| id                 | System-generated campaign type ID |
| name               | Campaign type name                |
| description        | Campaign type description         |
| critical           | Critical campaign type flag       |
| template_id        | Template ID                           |

<a name="delete-campaign-type"></a>
##### Deleting a Campaign Type

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```

\* The token requires `notifications.write` scope.

###### Route
```
DELETE /senders/<sender-id>/campaign_types/<campaign-type-id>
```

###### CURL example
```
$ curl -i -X DELETE \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/senders/a6c38f92-8fa9-488b-4f4c-7f4d4e0c0fd2/campaign_types/5cbc4458-3dba-481b-74c3-4548114b830b

204 No Content
RESPONSE HEADERS:
  Date: Wed, 05 Aug 2015 22:24:15 GMT
  Connection: close
RESPONSE BODY:
```

##### Response

###### Status
```
204 No Content
```

## Templates

<a name="create-template"></a>
#### Creating a template

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.

###### Route
```
POST /template
```
###### Params

| Key      | Description                                 |
| -------- | ------------------------------------------- |
| name\*   | the human-readable name given to a template |
| html\*\* | Template html body                          |
| text\*\* | Template text body                          |
| subject  | Template subject                            |
| metadata | Template metadata in JSON format            |

\* required
\*\* either html or text is required

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"name":"my-template", "text": "template text"}'
  http://notifications.example.com/templates

HTTP/1.1 201 Created
Content-Length: 138
Content-Type: text/plain; charset=utf-8
Date: Fri, 17 Jul 2015 19:30:32 GMT
X-Cf-Requestid: ce9f6b5a-317d-4d0f-7197-df63540c7f22

{"html":"","id":"54b39e0a-3c90-11e5-b0b3-10ddb1cec8d5","metadata":{},"name":"my-template","subject":"{{.Subject}}","text":"template text"}
```

##### Response

###### Status
```
201 Created
```

###### Body
| Fields   | Description                      |
| -------- | -------------------------------- |
| id       | System-generated template ID     |
| name     | Template name                    |
| html     | Template html body               |
| text     | Template text body               |
| subject  | Template subject                 |
| metadata | Template metadata in JSON format |

<a name="retrieve-template"></a>
#### Retrieving a template

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.

###### Route
```
GET /template/{templateID}
```
###### Params

| Key          | Description                                |
| ------------ | ------------------------------------------ |
| templateID\* | the "id" returned when creating a template |

\* required

###### CURL example
```
$ curl -i -X GET \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/templates/54b39e0a-3c90-11e5-b0b3-10ddb1cec8d5

HTTP/1.1 200 OK
Content-Length: 138
Content-Type: text/plain; charset=utf-8
Date: Fri, 17 Jul 2015 19:30:32 GMT
X-Cf-Requestid: ce9f6b5a-317d-4d0f-7197-df63540c7f22

{"html":"","id":"54b39e0a-3c90-11e5-b0b3-10ddb1cec8d5","metadata":{},"name":"my-template","subject":"{{.Subject}}","text":"template text"}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields   | Description                      |
| -------- | -------------------------------- |
| id       | System-generated template ID     |
| name     | Template name                    |
| html     | Template html body               |
| text     | Template text body               |
| subject  | Template subject                 |
| metadata | Template metadata in JSON format |

<a name="delete-template"></a>
#### Deleting a template

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.

###### Route
```
DELETE /template/{templateID}
```
###### Params

| Key          | Description                                |
| ------------ | ------------------------------------------ |
| templateID\* | the "id" returned when creating a template |

\* required

###### CURL example
```
$ curl -i -X DELETE \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/templates/54b39e0a-3c90-11e5-b0b3-10ddb1cec8d5

HTTP/1.1 204 No Content
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Fri, 17 Jul 2015 19:30:32 GMT
X-Cf-Requestid: ce9f6b5a-317d-4d0f-7197-df63540c7f22
```

##### Response

###### Status
```
200 OK
```

## Campaigns

<a name="send-campaign"></a>
#### Sending a campaign

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The token requires `notifications.write` scope.
\*\* Sending a critical campaign type requires `critical_notifications.write` scope.

###### Route
```
POST /sender/{sender-id}/campaigns
```
###### Params

| Key                | Description                               |
| ------------------ | ----------------------------------------- |
| send_to\*          | audience to deliver to |
| campaign_type_id\* | id of the previously created campaign |
| text\*\*		     | the text of your email |
| html\*\*		     | the html of your email |
| subject\*		     | subject of the email |
| template_id	     | the id of the template you would like to use |
| reply_to		     | email address used for replies |

\* required
\*\* either text or html have to be set, not both

###### Supported Audience Types
- user (provide a user GUID)
- space (provide a space GUID)
- org (provide an org GUID)
- email (provide an email address)

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{
  	"send_to": {
		"user": "c033fc5a-5878-45ca-8f7b-66f1857cfabc"
	},
	"campaign_type_id":"49b3bad1-a897-44eb-ab38-05c5725dfcb8",
  	"text":"this is an email",
	"subject":"this is a subject",
	"template_id":"8a947854-68d0-4914-9740-12e60743b0b9"
  }'
  http://notifications.example.com/senders/555a8e36-89da-48a2-8091-01881acd5051/campaigns

HTTP/1.1 202 Accepted
Content-Length: 57
Content-Type: text/plain; charset=utf-8
Date: Fri, 17 Jul 2015 19:30:32 GMT
X-Cf-Requestid: ce9f6b5a-317d-4d0f-7197-df63540c7f22

{ "campaign_id": "7e45da15-acac-441d-912f-e18d306eae83"}
```

##### Response

###### Status
```
202 Accepted
```

###### Body
| Fields      | Description                  |
| ----------- | ---------------------------- |
| campaign_id | System-generated campaign ID |

