# Notifications V2 Documentation

- System Status
	- [Check service status](#get-info)
- Senders
	- [Creating a sender](#create-sender)
	- [Retrieving a sender](#retrieve-sender)

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
\* The user token requires `notifications.write` scope.

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

<a name="retrieve-sender"></a>
#### Retrieving a sender

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The user token requires `notifications.write` scope.

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


## Notification Types

<a name="create-notification-type"></a>
#### Creating a notification type

##### Request

###### Headers
```
X-NOTIFICATIONS-VERSION: 2
Authorization: Bearer <CLIENT-TOKEN>
```
\* The user token requires `notifications.write` scope.
\*\* Creation of a critical notification type requires `critical_notifications.write` scope.

###### Route
```
POST /senders/<sender-id>/notification-types
```
###### Params

| Key                       | Description                                                         |
| ------------------------- | ------------------------------------------------------------------- |
| name\*                    | the human-readable name given to a notification type                |
| description\*             | the human-readable description given to a notification type         |
| critical (default: false) | a flag to indicate whether the notification type is critical or not |
| template_id               | the ID of a template to use for this notification type              |

\* required

###### CURL example
```
$ curl -i -X POST \
  -H "X-NOTIFICATIONS-VERSION: 2" \
  -H "Authorization: Bearer <CLIENT-TOKEN>" \
  -d '{"name":"my-notification-type","description":"notification type description","critical":false,"template_id":""}'
  http://notifications.mrorange.cfla.cf-app.com/senders/4bbd0431-9f5b-49bb-701d-8c2caa755ed0/notification_types

HTTP/1.1 201 Created
Content-Length: 155
Content-Type: text/plain; charset=utf-8
Date: Wed, 22 Jul 2015 16:00:37 GMT
X-Cf-Requestid: 6106873b-14ea-4fd9-6418-946c1651e4ac

{"critical":false,"description":"notification type description","id":"3d9aa963-97bb-4b48-4c3c-ecccad6314f8","name":"my-notification-type","template_id":""}
```

##### Response

###### Status
```
201 Created
```

###### Body
| Fields        | Description                           |
| ------------- | ------------------------------------- |
| id            | System-generated notification type ID |
| name          | Notification type name                |
| description   | Notification type description         |
| critical      | Critical notification type flag       |
| template_id   | Template ID                           |
