# Notifications

- System Status
	- [Check service status](#get-info)
- Sending Notifications
	- [Send a notification to a user](#post-users-guid)
	- [Send a notification to a space](#post-spaces-guid)
	- [Send a notification to an organization](#post-organizations-guid)
	- [Send a notification to a UAA-scope](#post-uaa-scopes)
	- [Send a notification to an email address](#post-emails)
- Registering Notifications
	- [Registering client notifications](#put-registration)
- Managing User Preferences
	- [Retrieve options for /user_preferences endpoints](#options-user-preferences)
	- [Retrieve user preferences with a user token](#get-user-preferences)
	- [Update user preferences with a user token](#patch-user-preferences)
	- [Retrieve options for /user_preferences/{user-guid} endpoints](#options-user-preferences-guid)
	- [Retrieve user preferences with a client token](#get-user-preferences-guid)
	- [Update user preferences with a client token](#patch-user-preferences-guid)
- Managing Templates
    - [Retrieve Templates](#get-template)
    - [Set Templates](#set-template)
    - [Delete Templates](#delete-template)


## System Status

<a name="get-info"></a>
#### Check service status

##### Request

###### Route
```
GET /info
```

###### CURL example
```
$ curl -i -X GET http://notifications.example.com/info

HTTP/1.1 200 OK
Connection: close
Content-Length: 2
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 21:29:36 GMT
X-Cf-Requestid: 2cf01258-ccff-41e9-6d82-41a4441af4af

{}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields   | Description |
| -------- | ----------- |
| \<None\> |             |


## Sending Notifications

<a name="post-users-guid"></a>
#### Send a notification to a user

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /users/{user-guid}
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| kind_description   | a description of the kind_id                   |
| subject            | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |
| source_description | a description of the sender                    |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
curl -i -X POST \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "html":"this is a test"}' \
  http://notifications.example.com/users/user-guid

HTTP/1.1 200 OK
Connection: close
Content-Length: 129
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 21:50:13 GMT
X-Cf-Requestid: 5c9bca88-280e-41d1-6e80-26a2a97adf4a

[{
	"notification_id":"451dd96a-ab8f-4a0b-5c3cb3bfe8ac1732",
	"recipient":"user-guid",
	"status":"queued"
}]
```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----
<a name="post-spaces-guid"></a>
#### Send a notification to a space

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /spaces/{space-guid}
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| kind_description   | a description of the kind_id                   |
| subject            | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |
| source_description | a description of the sender                    |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "html":"this is a test"}' \
  http://notifications.example.com/spaces/space-guid

HTTP/1.1 200 OK
Connection: close
Content-Length: 641
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 22:01:34 GMT
X-Cf-Requestid: 4dcfc91c-9cf6-4a51-497a-8ae506ce37f5

[{
	"notification_id":"f44da2ff-e402-435d-54e8-8703970d5917",
	"recipient":"user-guid-1",
	"status":"queued"
 },
 {
 	"notification_id":"253305c8-eb72-4430-690e-76cbd8eae8ee",
 	"recipient":"user-guid-2",
 	"status":"queued"
}]
```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----
<a name="post-organizations-guid"></a>
#### Send a notification to an organization

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /organizations/{organization-guid}
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| kind_description   | a description of the kind_id                   |
| subject            | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |
| source_description | a description of the sender                    |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "html":"this is a test"}' \
  http://notifications.example.com/organizations/organization-guid

Connection: close
Content-Length: 897
Content-Type: text/plain; charset=utf-8
Date: Thu, 06 Nov 2014 20:06:27 GMT
X-Cf-Requestid: 3a564cd9-74c8-46f6-5d31-8a8b600fc43f

[{
	"notification_id":"344f4b28-07d5-4490-468f-0a2f6fb4a65c",
	"recipient":"55498729-5749-4a4c-9e13-6893b795561b",
	"status":"queued"
	},{
	"notification_id":"96e633ef-8749-4dec-411a-f38a87f3fe79",
	"recipient":"d55067b8-cf2d-44ab-b70c-03dfd577a465",
	"status":"queued"
}]
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----
<a name="post-uaa-scopes"></a>
#### Send a notification to a UAA Scope

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Sending __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
POST /uaa_scopes/{scope}
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| kind_description   | a description of the kind_id                   |
| subject            | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |
| source_description | a description of the sender                    |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"kind_id":"example-kind-id", "html":"this is a test"}' \
  http://notifications.example.com/uaa_scopes/uaa.scope

Connection: close
Content-Length: 897
Content-Type: text/plain; charset=utf-8
Date: Thu, 06 Nov 2014 20:06:27 GMT
X-Cf-Requestid: 3a564cd9-74c8-46f6-5d31-8a8b600fc43f

[{
	"notification_id":"344f4b28-07d5-4490-468f-0a2f6fb4a65c",
	"recipient":"55498729-5749-4a4c-9e13-6893b795561b",
	"status":"queued"
	},{
	"notification_id":"96e633ef-8749-4dec-411a-f38a87f3fe79",
	"recipient":"d55067b8-cf2d-44ab-b70c-03dfd577a465",
	"status":"queued"
}]
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| recipient       | User GUID of notification recipient       |
| status          | Current delivery status of notification   |

----
<a name="post-emails"></a>
#### Send a notification to an email address

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `emails.write` scope

###### Route
```
POST /emails
```
###### Params

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| to\*          | The email address (and possibly full name) of the intended recipient in SMTP compatible format. |
| subject | The desired subject line of the notification.  The final subject may be prefixed, suffixed, or truncated by the notifier, all dependent on the templates.|
| reply_to | The email address to be included as the Reply-To address of the outgoing message. |
| text ** | The message body, in plain text  (required if html is absent) |
| html ** | The message body, in HTML  (required if text is absent) |

\* required

\*\* either text or html have to be set, not both

###### CURL example
```
$ curl -i -X POST \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"to":"user@example.com", "html":"this is a test"}' \
  http://notifications.example.com/emails

HTTP/1.1 200 OK
Connection: close
Content-Length: 108
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 22:27:48 GMT
X-Cf-Requestid: eb7ee46c-2142-4a74-5b73-e4971eea511a

[{
	"email":"user@example.com",
	"notification_id":"86ad7892-8217-4359-54b1-fe3ca60d8ac9",
	"status":"queued"
}]
```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields          | Description                               |
| --------------- | ----------------------------------------- |
| notification_id | Random GUID assigned to notification sent |
| email           | Email address of notification recipient   |
| status          | Current delivery status of notification   |

## Registering Notifications

<a name="put-registration"></a>
#### Registering client notifications

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notifications.write` scope. Registering __critical__ notifications requires the `critical_notifications.write` scope.

###### Route
```
PUT /registration
```
###### Params

| Key                 | Description                                    |
| ------------------- | ---------------------------------------------- |
| source_description* | A description of the sender, to be displayed in messages to users instead of the raw "client_id" field (which is derived from UAA) |
| kinds               | A complete list of all notification kinds that this client plans on using.  If passed, the notifier will add and remove kinds from its internal datastore to match the provided list. See table below for kinds fields |

\* required

###### Kinds

| Key                       | Description |
| ------------------------- | ----------- |
| id*                       | A simple machine readable string that identifies this type of notification.  It should be in the format /[0-9a-z_-.]+/i The notifier can use the ID to determine whether and how to notify a user. It’s recommended to use a GUID that doesn’t change for this field. |
| description*              | A description of the kind, to be displayed in messages to users instead of the raw “id” field |
| critical (default: false) | A boolean describing whether this kind of notification is to be considered “critical”, usually meaning that it cannot be unsubscribed from.  Because critical notifications can be annoying to end-users, registering a critical notification kind requires the client to have an access token with the critical_notifications.write scope. |

\* required

###### CURL example
```
$ curl -i -X PUT \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"source_description": "Galactic Empire", "kinds": [{"id": "example-kind-id", "description":"Example Kind Description", "critical": true}]}' \
  http://notifications.example.com/registration


HTTP/1.1 200 OK
Connection: close
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 22:47:50 GMT
X-Cf-Requestid: f39e22a4-6693-4a6d-6b27-006aecc924d4
```
##### Response

###### Status
```
200 OK
```

###### Body
| Fields   | Description |
| -------- | ----------- |
| \<none\> |             |

## Managing User Preferences

<a name="options-user-preferences"></a>
#### Retrieve Options for /user_preferences endpoints

##### Request

###### Route
```
OPTIONS /user_preferences
```

###### CURL example
```
$ curl -i -X OPTIONS \
  http://notifications.example.com/user_preferences

HTTP/1.1 204 No Content
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 22:54:40 GMT
X-Cf-Requestid: 686f601e-b6c7-4849-5699-6eed1a72004b
```
##### Response

###### Status
```
204 No Content
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences` path support the specified headers from any origin.

###### Body
| Fields   | Description |
| -------- | ----------- |
| \<none\> |             |

----
<a name="get-user-preferences"></a>
#### Retrieve user preferences with a user token

##### Request

###### Headers
```
Authorization: bearer <USER-TOKEN>
```
\* The user token requires `notification_preferences.write` scope.

###### Route
```
GET /user_preferences
```

###### CURL example
```
$ curl -i -X GET \
  -H "Authorization: Bearer <USER-TOKEN>" \
  http://notifications.example.com/user_preferences

HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 450
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 23:19:11 GMT
X-Cf-Requestid: 92cffe86-16fe-41a8-4b80-b10987b11060

{
	"clients" : {
		"login-service": {
			"effa96de-2349-423a-b5e4-b1e84712a714": {
				"count": 8,
				"email": true,
				"kind_description": "Forgot Password",
				"source_description": "Login Service"
			}
		},
		"MySQL Service": {
			"6236f606-627d-4079-b0bd-f0b7e8d3d2a9": {
				"count": 1,
				"email": false,
				"kind_description": "Downtime Notification",
				"source_description": "Galactic Empire Datastore"
			},
			"fb89e98a-a1f5-47e5-9e2d-d95940b32d3d": {
				"count": 18,
				"email": true,
				"kind_description": "Provision Notification",
				"source_description": "Galactic Empire Datastore"
			}
		}
	}
}
```
##### Response

###### Status
```
200 OK
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences` path support the specified headers from any origin.

###### Body
| Fields             | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| client_id          | Top-level keys are client identifiers                           |
| kind_id            | 2nd-level keys are kind identifiers                             |
| count              | Number of notifications sent for this kind                      |
| email              | Indicates if the user is subscribed to receive the notification |
| kind_description   | A human-friendly description of the notification                |
| source_description | A human-friendly description of the sending client              |

----
<a name="patch-user-preferences"></a>
#### Update user preferences with a user token

##### Request

###### Headers
```
Authorization: bearer <USER-TOKEN>
```
\* The user token requires `notification_preferences.write` scope.

###### Route
```
PATCH /user_preferences
```

###### Params
| Fields             | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| client_id          | Top-level keys are client identifiers                           |
| kind_id            | 2nd-level keys are kind identifiers                             |
| email              | Indicates if the user is subscribed to receive the notification |

###### CURL example
```
$ curl -i -X PATCH \
  -H "Authorization: Bearer <USER-TOKEN>" \
  -d '{"clients": {"login-service":{"effa96de-2349-423a-b5e4-b1e84712a714":{"email":true}}}}'
  http://notifications.example.com/user_preferences

HTTP/1.1 204 No Content
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 23:19:11 GMT
X-Cf-Requestid: 92cffe86-16fe-41a8-4b80-b10987b11060
```
##### Response

###### Status
```
204 No Content
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences` path support the specified headers from any origin.

----
<a name="options-user-preferences-guid"></a>
#### Retrieve Options for /user_preferences/{user-guid} endpoints

##### Request

###### Route
```
OPTIONS /user_preferences/{user-guid}
```

###### CURL example
```
$ curl -i -X OPTIONS \
  http://notifications.example.com/user_preferences/user-guid

HTTP/1.1 204 No Content
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 23:07:22 GMT
X-Cf-Requestid: bfb28efe-757e-4b65-4d48-1d2c6d7a9ce6
```
##### Response

###### Status
```
204 No Content
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences/{user-guid}` path support the specified headers from any origin.

###### Body
| Fields   | Description |
| -------- | ----------- |
| \<none\> |             |

----
<a name="get-user-preferences-guid"></a>
#### Retrieve user preferences with a client token

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_preferences.admin` scope.

###### Route
```
GET /user_preferences/{user-guid}
```

###### CURL example
```
$ curl -i -X GET \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/user_preferences/user-guid

HTTP/1.1 200 OK
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 450
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 23:19:11 GMT
X-Cf-Requestid: 92cffe86-16fe-41a8-4b80-b10987b11060

{
	"clients": {
		"login-service": {
			"effa96de-2349-423a-b5e4-b1e84712a714": {
				"count": 8,
				"email": true,
				"kind_description": "Forgot Password",
				"source_description": "Login Service"
			}
		},
		"mysql-service": {
			"6236f606-627d-4079-b0bd-f0b7e8d3d2a9": {
				"count": 1,
				"email": false,
				"kind_description": "Downtime Notification",
				"source_description": "Galactic Empire Datastore"
			},
			"fb89e98a-a1f5-47e5-9e2d-d95940b32d3d": {
				"count": 18,
				"email": true,
				"kind_description": "Provision Notification",
				"source_description": "Galactic Empire Datastore"
			}
		}
	}
}
```
##### Response

###### Status
```
200 OK
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences/{user-guid}` path support the specified headers from any origin.

###### Body
| Fields             | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| client_id          | Top-level keys are client identifiers                           |
| kind_id            | 2nd-level keys are kind identifiers                             |
| count              | Number of notifications sent for this kind                      |
| email              | Indicates if the user is subscribed to receive the notification |
| kind_description   | A human-friendly description of the notification                |
| source_description | A human-friendly description of the sending client              |

----
<a name="patch-user-preferences-guid"></a>
#### Update user preferences with a client token

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token requires `notification_preferences.admin` scope.

###### Route
```
PATCH /user_preferences/user-guid
```

###### Params
| Fields             | Description                                                     |
| ------------------ | --------------------------------------------------------------- |
| client_id          | Top-level keys are client identifiers                           |
| kind_id            | 2nd-level keys are kind identifiers                             |
| email              | Indicates if the user is subscribed to receive the notification |

###### CURL example
```
$ curl -i -X PATCH \
  -H "Authorization: Bearer <USER-TOKEN>" \
  -d '{"clients": {"login-service":{"effa96de-2349-423a-b5e4-b1e84712a714":{"email":true}}}}'
  http://notifications.example.com/user_preferences/user-guid

HTTP/1.1 204 No Content
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
Connection: close
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 30 Sep 2014 23:19:11 GMT
X-Cf-Requestid: 92cffe86-16fe-41a8-4b80-b10987b11060
```
##### Response

###### Status
```
204 No Content
```

###### Headers
```
Access-Control-Allow-Headers: Accept, Authorization, Content-Type
Access-Control-Allow-Methods: GET, PATCH
Access-Control-Allow-Origin: *
```
The above headers constitute a CORS contract. They indicate that the GET and PATCH endpoints for the `/user_preferences/user-guid` path support the specified headers from any origin.

## Managing Templates

<a name="get-template"></a>

### Retrieve Templates

This endpoint is used to retrieve the templates that are stored using the template filename as a parameter for the request.

##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client token does not require any scopes at this time.

###### Route
```
GET /templates/{template-filename}
```

###### CURL example
```
$ curl -i -X GET \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  http://notifications.example.com/templates/template-filename

HTTP/1.1 200 OK
Connection: close
Content-Length: 1126
Content-Type: text/plain; charset=utf-8
Date: Mon, 27 Oct 2014 23:52:31 GMT
X-Cf-Requestid: 5d19a080-2c88-4fe6-5eb5-42f9bda2d073

{
	"text": "Message to: {{.To}}, sent from the {{.ClientID}} UAA Client",
	"html": "<p>Message to: {{.To}}, sent from the {{.ClientID}} UAA Client</p>",
	"overridden": false
}
```

##### Response

###### Status
```
200 OK
```

###### Body
| Fields      | Description                                                     |
| ----------- | --------------------------------------------------------------- |
| text        | The template used for the text portion of the notification      |
| html        | The template used for the HTML portion of the notification      |
| overridden  | True if this template is overridden                             |

<a name="set-template"></a>
### Set Templates

This endpoint is used to change the templates attached to particular kinds of notifications.


##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client does not require any scopes at this time

###### Route
```
PUT /templates/{template-filename}
```
###### Params

| Key      | Description                                                     |
| -------- | --------------------------------------------------------------- |
| text*    | The template used for the text portion of the notification      |
| html*    | The template used for the HTML portion of the notification      |

\* required

###### CURL example
```
$ curl -i -X PUT \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"text":"Message to: {{.To}}, sent from the {{.ClientID}} UAA Client", "html": "<p>Message to: {{.To}}, sent from the {{.ClientID}} UAA Client</p>",}' \
  http://notifications.example.com/templates/template-filename

204 No Content
Connection: close
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603
```

##### Response

###### Status
```
204 No Content
```

###### Body
| Fields   | Description |
| -------- | ----------- |
| \<none\> |             |

<a name="delete-template"></a>
### Delete Templates

This endpoint is used to delete an override template attached to a particular kind of notifications.


##### Request

###### Headers
```
Authorization: bearer <CLIENT-TOKEN>
```
\* The client does not require any scopes at this time

###### Route
```
DELETE /templates/{template-filename}
```
###### Params

| Key                   | Description                                                     |
| --------------------- | --------------------------------------------------------------- |
| template-filename*    | The template used for the text portion of the notification      |

\* required

###### CURL example
```
$ curl -X DELETE \
  http://notifications.example.com/templates/template-filename

204 No Content
Connection: close
Content-Length: 0
Content-Type: text/plain; charset=utf-8
Date: Tue, 28 Oct 2014 00:18:48 GMT
X-Cf-Requestid: 8938a949-66b1-43f5-4fad-a91fc050b603
```

##### Response

###### Status
```
204 No Content
```

###### Body
| Fields   | Description |
| -------- | ----------- |
| \<none\> |             |


