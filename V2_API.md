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
Date: Tue, 13 Oct 2015 22:17:04 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDEwYmYzZTctMzU3MS01YzMxLThkODAtYTg0NDllMzQzOTQ4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.RHDBOHvpOETXHyGMxE76ZzrHxdAdf7XAPfAqGAcp1dQ
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
Date: Tue, 13 Oct 2015 22:16:59 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a/campaign_types"
    },
    "campaigns": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a/campaigns"
    },
    "self": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a"
    }
  },
  "id": "02331d50-a1e3-4b86-4e9d-6c757a0d2e1a",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDEwYmYzZTctMzU3MS01YzMxLThkODAtYTg0NDllMzQzOTQ4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.RHDBOHvpOETXHyGMxE76ZzrHxdAdf7XAPfAqGAcp1dQ
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:16:59 GMT
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
          "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a/campaign_types"
        },
        "campaigns": {
          "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a/campaigns"
        },
        "self": {
          "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a"
        }
      },
      "id": "02331d50-a1e3-4b86-4e9d-6c757a0d2e1a",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDEwYmYzZTctMzU3MS01YzMxLThkODAtYTg0NDllMzQzOTQ4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.RHDBOHvpOETXHyGMxE76ZzrHxdAdf7XAPfAqGAcp1dQ
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:16:59 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a/campaign_types"
    },
    "campaigns": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a/campaigns"
    },
    "self": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a"
    }
  },
  "id": "02331d50-a1e3-4b86-4e9d-6c757a0d2e1a",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDEwYmYzZTctMzU3MS01YzMxLThkODAtYTg0NDllMzQzOTQ4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.RHDBOHvpOETXHyGMxE76ZzrHxdAdf7XAPfAqGAcp1dQ
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
Date: Tue, 13 Oct 2015 22:16:59 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a/campaign_types"
    },
    "campaigns": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a/campaigns"
    },
    "self": {
      "href": "/senders/02331d50-a1e3-4b86-4e9d-6c757a0d2e1a"
    }
  },
  "id": "02331d50-a1e3-4b86-4e9d-6c757a0d2e1a",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2ExZmIwZDAtMDZkMi1hYjI1LWJlYWMtNTI4MDUzMDgwNThiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.G3ageTNU4aac0fhdEiB_Gc16ehq5TEL0LZfGHiF0Dyc
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 22:17:02 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZDA1OGNjZjYtYjVhMy04MjkyLTZkZGQtNDI5ZWQyMWNmZWNjIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.5IWCxWnFdTVgEhTBFrgCsRC1Z9UR4rSuzsttfLpTL9s
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
Date: Tue, 13 Oct 2015 22:17:00 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/ea770b0a-05eb-1c38-31f1-8e80b6de7b40"
    }
  },
  "html": "template html",
  "id": "ea770b0a-05eb-1c38-31f1-8e80b6de7b40",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZDA1OGNjZjYtYjVhMy04MjkyLTZkZGQtNDI5ZWQyMWNmZWNjIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.5IWCxWnFdTVgEhTBFrgCsRC1Z9UR4rSuzsttfLpTL9s
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:17:00 GMT
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
          "href": "/templates/ea770b0a-05eb-1c38-31f1-8e80b6de7b40"
        }
      },
      "html": "html",
      "id": "ea770b0a-05eb-1c38-31f1-8e80b6de7b40",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZDA1OGNjZjYtYjVhMy04MjkyLTZkZGQtNDI5ZWQyMWNmZWNjIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.5IWCxWnFdTVgEhTBFrgCsRC1Z9UR4rSuzsttfLpTL9s
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:17:00 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/ea770b0a-05eb-1c38-31f1-8e80b6de7b40"
    }
  },
  "html": "template html",
  "id": "ea770b0a-05eb-1c38-31f1-8e80b6de7b40",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZDA1OGNjZjYtYjVhMy04MjkyLTZkZGQtNDI5ZWQyMWNmZWNjIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.5IWCxWnFdTVgEhTBFrgCsRC1Z9UR4rSuzsttfLpTL9s
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
Date: Tue, 13 Oct 2015 22:17:00 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/ea770b0a-05eb-1c38-31f1-8e80b6de7b40"
    }
  },
  "html": "html",
  "id": "ea770b0a-05eb-1c38-31f1-8e80b6de7b40",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZDA1OGNjZjYtYjVhMy04MjkyLTZkZGQtNDI5ZWQyMWNmZWNjIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.5IWCxWnFdTVgEhTBFrgCsRC1Z9UR4rSuzsttfLpTL9s
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 22:17:00 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMmNhMTkzYmItYTFiOS0yMzU4LWQ5MzItYzM2YWVlNDU3MGY1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MNPSe3m46_ogwN95HXfkdIEH3pUcobmgEeYyPI7hgKo
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
Date: Tue, 13 Oct 2015 22:17:01 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/09f2f769-fd11-0815-2b52-851115a41887"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "09f2f769-fd11-0815-2b52-851115a41887",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMmNhMTkzYmItYTFiOS0yMzU4LWQ5MzItYzM2YWVlNDU3MGY1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MNPSe3m46_ogwN95HXfkdIEH3pUcobmgEeYyPI7hgKo
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:17:01 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/99d12ae4-912a-13d6-36d9-61f91d74b91c/campaign_types"
    },
    "sender": {
      "href": "/senders/99d12ae4-912a-13d6-36d9-61f91d74b91c"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/09f2f769-fd11-0815-2b52-851115a41887"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "09f2f769-fd11-0815-2b52-851115a41887",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMmNhMTkzYmItYTFiOS0yMzU4LWQ5MzItYzM2YWVlNDU3MGY1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MNPSe3m46_ogwN95HXfkdIEH3pUcobmgEeYyPI7hgKo
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:17:01 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/09f2f769-fd11-0815-2b52-851115a41887"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "09f2f769-fd11-0815-2b52-851115a41887",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMmNhMTkzYmItYTFiOS0yMzU4LWQ5MzItYzM2YWVlNDU3MGY1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MNPSe3m46_ogwN95HXfkdIEH3pUcobmgEeYyPI7hgKo
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "dc83a608-9db3-7e30-ea58-73f7b134e293"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:17:01 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/09f2f769-fd11-0815-2b52-851115a41887"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "09f2f769-fd11-0815-2b52-851115a41887",
  "name": "updated-campaign-type",
  "template_id": "dc83a608-9db3-7e30-ea58-73f7b134e293"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMmNhMTkzYmItYTFiOS0yMzU4LWQ5MzItYzM2YWVlNDU3MGY1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MNPSe3m46_ogwN95HXfkdIEH3pUcobmgEeYyPI7hgKo
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 22:17:01 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTdhNmU3NWMtZjI3MC1iNTA5LTUxZTEtYmJiZDZhMDhmODAyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.Ffh5zMWBXR54ZzFnR3RTrunyeMrZZCi4E6lvR5Oq-Vc
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "4f57ac2c-d7e4-a32e-3dc0-12c0690f08a5",
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
      "2f7d2229-565a-47c1-53e8-e9d4c2cade58"
    ]
  },
  "subject": "campaign subject",
  "template_id": "86af93fe-b9df-22b2-3d2a-3a08e8fdf211",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:16:59 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/4f57ac2c-d7e4-a32e-3dc0-12c0690f08a5"
    },
    "self": {
      "href": "/campaigns/c14bdda5-2bfc-7158-e6f3-b87c15f7685c"
    },
    "status": {
      "href": "/campaigns/c14bdda5-2bfc-7158-e6f3-b87c15f7685c/status"
    },
    "template": {
      "href": "/templates/86af93fe-b9df-22b2-3d2a-3a08e8fdf211"
    }
  },
  "campaign_type_id": "4f57ac2c-d7e4-a32e-3dc0-12c0690f08a5",
  "html": "",
  "id": "c14bdda5-2bfc-7158-e6f3-b87c15f7685c",
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
      "2f7d2229-565a-47c1-53e8-e9d4c2cade58"
    ]
  },
  "subject": "campaign subject",
  "template_id": "86af93fe-b9df-22b2-3d2a-3a08e8fdf211",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTdhNmU3NWMtZjI3MC1iNTA5LTUxZTEtYmJiZDZhMDhmODAyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.Ffh5zMWBXR54ZzFnR3RTrunyeMrZZCi4E6lvR5Oq-Vc
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:16:59 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/4f57ac2c-d7e4-a32e-3dc0-12c0690f08a5"
    },
    "self": {
      "href": "/campaigns/c14bdda5-2bfc-7158-e6f3-b87c15f7685c"
    },
    "status": {
      "href": "/campaigns/c14bdda5-2bfc-7158-e6f3-b87c15f7685c/status"
    },
    "template": {
      "href": "/templates/86af93fe-b9df-22b2-3d2a-3a08e8fdf211"
    }
  },
  "campaign_type_id": "4f57ac2c-d7e4-a32e-3dc0-12c0690f08a5",
  "html": "",
  "id": "c14bdda5-2bfc-7158-e6f3-b87c15f7685c",
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
      "2f7d2229-565a-47c1-53e8-e9d4c2cade58"
    ]
  },
  "subject": "campaign subject",
  "template_id": "86af93fe-b9df-22b2-3d2a-3a08e8fdf211",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYjU5OThkMDAtNTM4Yy1hOGEyLWNiNzMtZWY4MDExNDBhYWQwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.T0BvOdqvz3dRmDrUxmyFvWsxosBRzk4iuoZX2briKpw
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 393
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 22:17:00 GMT
```
##### Body
```
{
  "_links": {
    "campaign": {
      "href": "/campaigns/5d2497c0-005c-663a-29b0-51943aa5b019"
    },
    "self": {
      "href": "/campaigns/5d2497c0-005c-663a-29b0-51943aa5b019/status"
    }
  },
  "completed_time": "2015-10-13T22:17:00Z",
  "failed_messages": 0,
  "id": "5d2497c0-005c-663a-29b0-51943aa5b019",
  "queued_messages": 0,
  "retry_messages": 0,
  "sent_messages": 1,
  "start_time": "2015-10-13T22:17:00Z",
  "status": "completed",
  "total_messages": 1
}
```


## Unsubscribing
EVAN WILL EDIT THIS
<a name="unsubscriber-put-client"></a>
### Unsubscribe a user (with a client token)
#### Request **PUT** /campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.admin
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMDY4YWY0MGUtMGIyYi04MDQ4LWJkN2ItMTg1MTUyNmZkNWMzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.DD91qsrGWIbGOrc6dA50KjiheNIruawzAojziJQ-BkU
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 22:16:58 GMT
```

<a name="unsubscriber-put-user"></a>
### Unsubscribe a user (with a user token)
#### Request **PUT** /campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.read notification_preferences.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiMjE1OWY2YjgtODIyOS00ZWIyLTdiODYtMDFkNzQzOTJlYmQ2In0.SRHoUu4QAU_54RtzrivEq3ymWABxR8fOruWRsIUWKek
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 22:17:00 GMT
```

<a name="unsubscriber-delete-client"></a>
### Resubscribe a user (with a client token)
#### Request **DELETE** /campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.admin
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMDY4YWY0MGUtMGIyYi04MDQ4LWJkN2ItMTg1MTUyNmZkNWMzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTEyOC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.DD91qsrGWIbGOrc6dA50KjiheNIruawzAojziJQ-BkU
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 22:16:58 GMT
```

<a name="unsubscriber-delete-user"></a>
### Resubscribe a user (with a user token)
#### Request **DELETE** /campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.read notification_preferences.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiMjE1OWY2YjgtODIyOS00ZWIyLTdiODYtMDFkNzQzOTJlYmQ2In0.SRHoUu4QAU_54RtzrivEq3ymWABxR8fOruWRsIUWKek
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 22:17:00 GMT
```


