<!---
This document is automatically generated.
DO NOT EDIT THIS BY HAND.
Run the acceptance suite to re-generate the documentation.
-->

# API docs
This is version 2 of the Cloud Foundry notifications API. Integrate with this API in order to send messages to developers and billing contacts in your CF deployment.

## Authorization Tokens
See [here](https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-Tokens.md) for more information about UAA tokens.

## Table of Contents
* Info
  * [Retrieve information about the API](#info-get)
* Senders
  * [Create a new sender](#sender-create)
  * [List all senders](#sender-list)
  * [Retrieve a sender](#sender-get)
  * [Update a sender](#sender-update)
  * [Delete a sender](#sender-delete)
* Templates
  * [Create a new template](#template-create)
  * [Retrieve a list of templates](#template-list)
  * [Retrieve a template](#template-get)
  * [Update a template](#template-update)
  * [Delete a template](#template-delete)
* Campaign Types
  * [Create a new campaign type](#campaign-type-create)
  * [Retrieve a list of campaign types](#campaign-type-list)
  * [Retrieve a campaign type](#campaign-type-get)
  * [Update a campaign type](#campaign-type-update)
  * [Delete a campaign type](#campaign-type-delete)
* Campaigns
  * [Create a new campaign](#campaign-create)
  * [Retrieve a campaign](#campaign-get)
  * [Retrieve the status of a campaign](#campaign-status)
* Unsubscribing
  * [Unsubscribe a user (with a client token)](#unsubscriber-put-client)
  * [Unsubscribe a user (with a user token)](#unsubscriber-put-user)
  * [Resubscribe a user (with a client token)](#unsubscriber-delete-client)
  * [Resubscribe a user (with a user token)](#unsubscriber-delete-user)

## Info
EVAN WILL EDIT THIS
<a name="info-get"></a>
### Retrieve information about the API
#### Request **GET** /info
##### Headers
```
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 13
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "version": 2
}
```


## Senders
EVAN WILL EDIT THIS
<a name="sender-create"></a>
### Create a new sender
#### Request **POST** /senders
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGMwYjgzYmItYTc4Ny1mNzAyLWNkOTktZjgwZjg2M2I2NDJiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cpCWKRRrqGGfnjXin5_DQNNxqFbuT0fZekuYODEcivk
X-Notifications-Version: 2
```
##### Body
```
{
  "name": "My Cool App"
}
```
#### Response 201 Created
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d/campaign_types"
    },
    "campaigns": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d/campaigns"
    },
    "self": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d"
    }
  },
  "id": "70ed79b1-71a4-c95e-760f-d1a2d823bf5d",
  "name": "My Cool App"
}
```

<a name="sender-list"></a>
### List all senders
#### Request **GET** /senders
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGMwYjgzYmItYTc4Ny1mNzAyLWNkOTktZjgwZjg2M2I2NDJiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cpCWKRRrqGGfnjXin5_DQNNxqFbuT0fZekuYODEcivk
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders"
    }
  },
  "senders": [
    {
      "_links": {
        "campaign_types": {
          "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d/campaign_types"
        },
        "campaigns": {
          "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d/campaigns"
        },
        "self": {
          "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d"
        }
      },
      "id": "70ed79b1-71a4-c95e-760f-d1a2d823bf5d",
      "name": "My Cool App"
    }
  ]
}
```

<a name="sender-get"></a>
### Retrieve a sender
#### Request **GET** /senders/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGMwYjgzYmItYTc4Ny1mNzAyLWNkOTktZjgwZjg2M2I2NDJiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cpCWKRRrqGGfnjXin5_DQNNxqFbuT0fZekuYODEcivk
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d/campaign_types"
    },
    "campaigns": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d/campaigns"
    },
    "self": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d"
    }
  },
  "id": "70ed79b1-71a4-c95e-760f-d1a2d823bf5d",
  "name": "My Cool App"
}
```

<a name="sender-update"></a>
### Update a sender
#### Request **PUT** /senders/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGMwYjgzYmItYTc4Ny1mNzAyLWNkOTktZjgwZjg2M2I2NDJiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cpCWKRRrqGGfnjXin5_DQNNxqFbuT0fZekuYODEcivk
X-Notifications-Version: 2
```
##### Body
```
{
  "name": "My Not Cool App"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 314
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d/campaign_types"
    },
    "campaigns": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d/campaigns"
    },
    "self": {
      "href": "/senders/70ed79b1-71a4-c95e-760f-d1a2d823bf5d"
    }
  },
  "id": "70ed79b1-71a4-c95e-760f-d1a2d823bf5d",
  "name": "My Not Cool App"
}
```

<a name="sender-delete"></a>
### Delete a sender
#### Request **DELETE** /senders/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiOWZkMGRlOGMtMTBmMi1kY2UyLTM2ZWUtNDdiMGZhNjhkOGExIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.49fvhmIqG7AAo7kq5cEA4XlfeUx6aEEDxJ7yFMMH8io
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 21:48:59 GMT
```


## Templates
EVAN WILL EDIT THIS
<a name="template-create"></a>
### Create a new template
#### Request **POST** /templates
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTc2OTE5OGEtMmQ3NC1kY2JjLTc3MTAtYzUzMGRlMTc2NTY4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.1YtDB-c3XDOxh8o2QeMsQSw0eJa9m1xd7V5r3jZHZXo
X-Notifications-Version: 2
```
##### Body
```
{
  "html": "template html",
  "metadata": {
    "template": "metadata"
  },
  "name": "An interesting template",
  "subject": "template subject",
  "text": "template text"
}
```
#### Response 201 Created
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/e26f22c3-bb55-1274-f19f-4808d73df6ce"
    }
  },
  "html": "template html",
  "id": "e26f22c3-bb55-1274-f19f-4808d73df6ce",
  "metadata": {
    "template": "metadata"
  },
  "name": "An interesting template",
  "subject": "template subject",
  "text": "template text"
}
```

<a name="template-list"></a>
### Retrieve a list of templates
#### Request **GET** /templates
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTc2OTE5OGEtMmQ3NC1kY2JjLTc3MTAtYzUzMGRlMTc2NTY4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.1YtDB-c3XDOxh8o2QeMsQSw0eJa9m1xd7V5r3jZHZXo
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates"
    }
  },
  "templates": [
    {
      "_links": {
        "self": {
          "href": "/templates/e26f22c3-bb55-1274-f19f-4808d73df6ce"
        }
      },
      "html": "html",
      "id": "e26f22c3-bb55-1274-f19f-4808d73df6ce",
      "metadata": {
        "banana": "something"
      },
      "name": "A more interesting template",
      "subject": "subject",
      "text": "text"
    }
  ]
}
```

<a name="template-get"></a>
### Retrieve a template
#### Request **GET** /templates/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTc2OTE5OGEtMmQ3NC1kY2JjLTc3MTAtYzUzMGRlMTc2NTY4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.1YtDB-c3XDOxh8o2QeMsQSw0eJa9m1xd7V5r3jZHZXo
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/e26f22c3-bb55-1274-f19f-4808d73df6ce"
    }
  },
  "html": "template html",
  "id": "e26f22c3-bb55-1274-f19f-4808d73df6ce",
  "metadata": {
    "template": "metadata"
  },
  "name": "An interesting template",
  "subject": "template subject",
  "text": "template text"
}
```

<a name="template-update"></a>
### Update a template
#### Request **PUT** /templates/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTc2OTE5OGEtMmQ3NC1kY2JjLTc3MTAtYzUzMGRlMTc2NTY4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.1YtDB-c3XDOxh8o2QeMsQSw0eJa9m1xd7V5r3jZHZXo
X-Notifications-Version: 2
```
##### Body
```
{
  "html": "html",
  "metadata": {
    "banana": "something"
  },
  "name": "A more interesting template",
  "subject": "subject",
  "text": "text"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 242
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:03 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/e26f22c3-bb55-1274-f19f-4808d73df6ce"
    }
  },
  "html": "html",
  "id": "e26f22c3-bb55-1274-f19f-4808d73df6ce",
  "metadata": {
    "banana": "something"
  },
  "name": "A more interesting template",
  "subject": "subject",
  "text": "text"
}
```

<a name="template-delete"></a>
### Delete a template
#### Request **DELETE** /templates/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTc2OTE5OGEtMmQ3NC1kY2JjLTc3MTAtYzUzMGRlMTc2NTY4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.1YtDB-c3XDOxh8o2QeMsQSw0eJa9m1xd7V5r3jZHZXo
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 21:49:03 GMT
```


## Campaign Types
EVAN WILL EDIT THIS
<a name="campaign-type-create"></a>
### Create a new campaign type
#### Request **POST** /senders/{id}/campaign_types
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODVjNTQ5YjAtYTlkMS04OGNhLTY5NjYtNDZmOGZjOGRiMGY3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dIXdFTTHTEtVME7dT8XIz38c9eZw9qgP1A5cHqGP5Y0
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "a great campaign type",
  "name": "some-campaign-type"
}
```
#### Response 201 Created
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:04 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/c93d803b-62ac-bb8e-802f-9fadc8050288"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "c93d803b-62ac-bb8e-802f-9fadc8050288",
  "name": "some-campaign-type",
  "template_id": ""
}
```

<a name="campaign-type-list"></a>
### Retrieve a list of campaign types
#### Request **GET** /senders/{id}/campaign_types
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODVjNTQ5YjAtYTlkMS04OGNhLTY5NjYtNDZmOGZjOGRiMGY3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dIXdFTTHTEtVME7dT8XIz38c9eZw9qgP1A5cHqGP5Y0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:04 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/4da74e92-03ce-6dab-48ea-ca6f7fc89e16/campaign_types"
    },
    "sender": {
      "href": "/senders/4da74e92-03ce-6dab-48ea-ca6f7fc89e16"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/c93d803b-62ac-bb8e-802f-9fadc8050288"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "c93d803b-62ac-bb8e-802f-9fadc8050288",
      "name": "some-campaign-type",
      "template_id": ""
    }
  ]
}
```

<a name="campaign-type-get"></a>
### Retrieve a campaign type
#### Request **GET** /campaign_types/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODVjNTQ5YjAtYTlkMS04OGNhLTY5NjYtNDZmOGZjOGRiMGY3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dIXdFTTHTEtVME7dT8XIz38c9eZw9qgP1A5cHqGP5Y0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:04 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/c93d803b-62ac-bb8e-802f-9fadc8050288"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "c93d803b-62ac-bb8e-802f-9fadc8050288",
  "name": "some-campaign-type",
  "template_id": ""
}
```

<a name="campaign-type-update"></a>
### Update a campaign type
#### Request **PUT** /campaign_types/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODVjNTQ5YjAtYTlkMS04OGNhLTY5NjYtNDZmOGZjOGRiMGY3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dIXdFTTHTEtVME7dT8XIz38c9eZw9qgP1A5cHqGP5Y0
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "c689d530-bdd9-6727-bfec-1df0bf2c1cf2"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:04 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/c93d803b-62ac-bb8e-802f-9fadc8050288"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "c93d803b-62ac-bb8e-802f-9fadc8050288",
  "name": "updated-campaign-type",
  "template_id": "c689d530-bdd9-6727-bfec-1df0bf2c1cf2"
}
```

<a name="campaign-type-delete"></a>
### Delete a campaign type
#### Request **DELETE** /campaign_types/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODVjNTQ5YjAtYTlkMS04OGNhLTY5NjYtNDZmOGZjOGRiMGY3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dIXdFTTHTEtVME7dT8XIz38c9eZw9qgP1A5cHqGP5Y0
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 21:49:04 GMT
```


## Campaigns
EVAN WILL EDIT THIS
<a name="campaign-create"></a>
### Create a new campaign
#### Request **POST** /senders/{id}/campaigns
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMmI4ZTRiYzEtNmM2Yi1kZThkLTU2ZjMtMmU5NjYzNjkwOGY5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.r6MixpYR1sd1SHUipNevK0zcyM3mplfnq3obHt_KzW4
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "3be6534d-d593-4acb-99fe-3427fab82fa5",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "emails": [
      "test@example.com"
    ],
    "orgs": [
      "org-123"
    ],
    "spaces": [
      "space-123"
    ],
    "users": [
      "1638b139-9b7a-4323-7c02-fd9575a4b8b2"
    ]
  },
  "subject": "campaign subject",
  "template_id": "630bd81b-80e7-40b7-9a0b-fd41a493d302",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:04 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/3be6534d-d593-4acb-99fe-3427fab82fa5"
    },
    "self": {
      "href": "/campaigns/68215535-8186-d91d-cc96-33200def12e8"
    },
    "status": {
      "href": "/campaigns/68215535-8186-d91d-cc96-33200def12e8/status"
    },
    "template": {
      "href": "/templates/630bd81b-80e7-40b7-9a0b-fd41a493d302"
    }
  },
  "campaign_type_id": "3be6534d-d593-4acb-99fe-3427fab82fa5",
  "html": "",
  "id": "68215535-8186-d91d-cc96-33200def12e8",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "emails": [
      "test@example.com"
    ],
    "orgs": [
      "org-123"
    ],
    "spaces": [
      "space-123"
    ],
    "users": [
      "1638b139-9b7a-4323-7c02-fd9575a4b8b2"
    ]
  },
  "subject": "campaign subject",
  "template_id": "630bd81b-80e7-40b7-9a0b-fd41a493d302",
  "text": "campaign body"
}
```

<a name="campaign-get"></a>
### Retrieve a campaign
#### Request **GET** /campaigns/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMmI4ZTRiYzEtNmM2Yi1kZThkLTU2ZjMtMmU5NjYzNjkwOGY5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.r6MixpYR1sd1SHUipNevK0zcyM3mplfnq3obHt_KzW4
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:04 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/3be6534d-d593-4acb-99fe-3427fab82fa5"
    },
    "self": {
      "href": "/campaigns/68215535-8186-d91d-cc96-33200def12e8"
    },
    "status": {
      "href": "/campaigns/68215535-8186-d91d-cc96-33200def12e8/status"
    },
    "template": {
      "href": "/templates/630bd81b-80e7-40b7-9a0b-fd41a493d302"
    }
  },
  "campaign_type_id": "3be6534d-d593-4acb-99fe-3427fab82fa5",
  "html": "",
  "id": "68215535-8186-d91d-cc96-33200def12e8",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "emails": [
      "test@example.com"
    ],
    "orgs": [
      "org-123"
    ],
    "spaces": [
      "space-123"
    ],
    "users": [
      "1638b139-9b7a-4323-7c02-fd9575a4b8b2"
    ]
  },
  "subject": "campaign subject",
  "template_id": "630bd81b-80e7-40b7-9a0b-fd41a493d302",
  "text": "campaign body"
}
```

<a name="campaign-status"></a>
### Retrieve the status of a campaign
#### Request **GET** /campaigns/{id}/status
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTYyMjkzNDUtZmI5ZC0xN2ZjLTcyYjMtZjA4ZGY1MThmNTI5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.FD08XcjMlyr8BKAZdjCyBMmOOuApdXDfTxDKy0y1Qeo
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 393
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 21:49:04 GMT
```
##### Body
```
{
  "_links": {
    "campaign": {
      "href": "/campaigns/be2c920f-b99b-bd18-1dc4-de3f3ef5c6e0"
    },
    "self": {
      "href": "/campaigns/be2c920f-b99b-bd18-1dc4-de3f3ef5c6e0/status"
    }
  },
  "completed_time": "2015-10-13T21:49:04Z",
  "failed_messages": 0,
  "id": "be2c920f-b99b-bd18-1dc4-de3f3ef5c6e0",
  "queued_messages": 0,
  "retry_messages": 0,
  "sent_messages": 1,
  "start_time": "2015-10-13T21:49:04Z",
  "status": "completed",
  "total_messages": 1
}
```


## Unsubscribing
EVAN WILL EDIT THIS
<a name="unsubscriber-put-client"></a>
### Unsubscribe a user (with a client token)
#### Request **PUT** /senders/{id}/campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.admin
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTU1YzgzMmMtNmQ0ZS01ODExLTkyMjItMmIwYmRmYjM4MTJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ._vNpqVfrzMU8klza8dITncFJKwViRz9FPK-86UPm-B8
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 21:49:02 GMT
```

<a name="unsubscriber-put-user"></a>
### Unsubscribe a user (with a user token)
#### Request **PUT** /senders/{id}/campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.read notification_preferences.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiZjliZTY1NTItZjA0Ny00YmIwLTQ3NjctMzU1MjUzZWZkMTkyIn0.LxtpWJG8fHgOq9sNS79Dz24KbUySqDkLTs1CMeYxKis
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 21:49:02 GMT
```

<a name="unsubscriber-delete-client"></a>
### Resubscribe a user (with a client token)
#### Request **DELETE** /senders/{id}/campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.admin
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTU1YzgzMmMtNmQ0ZS01ODExLTkyMjItMmIwYmRmYjM4MTJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkzOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ._vNpqVfrzMU8klza8dITncFJKwViRz9FPK-86UPm-B8
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 21:49:02 GMT
```

<a name="unsubscriber-delete-user"></a>
### Resubscribe a user (with a user token)
#### Request **DELETE** /senders/{id}/campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.read notification_preferences.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiZjliZTY1NTItZjA0Ny00YmIwLTQ3NjctMzU1MjUzZWZkMTkyIn0.LxtpWJG8fHgOq9sNS79Dz24KbUySqDkLTs1CMeYxKis
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 21:49:02 GMT
```


