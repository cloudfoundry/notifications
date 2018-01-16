<!---
This document is automatically generated.
DO NOT EDIT THIS BY HAND.
Run the acceptance suite to re-generate the documentation.
-->

# API docs - DEPRECATED
This is version 2 of the Cloud Foundry notifications API and is no longer supported. Please use [version 1](/V1_API.md)
Integrate with this API in order to send messages to developers and billing contacts in your CF deployment.

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
Date: Mon, 19 Oct 2015 16:49:36 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYjc0NWE3OTAtMDc5Ny05NjAyLTMwZTQtMTk1OGUwMzk4ZGI3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CLK8JiZaoZj_bzK0BEpJFFGmI-pOD4hCDDB4FhI5fK0
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
Date: Mon, 19 Oct 2015 16:49:35 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c/campaign_types"
    },
    "campaigns": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c/campaigns"
    },
    "self": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c"
    }
  },
  "id": "fa92611e-d5a9-db99-3cb6-544a87584b3c",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYjc0NWE3OTAtMDc5Ny05NjAyLTMwZTQtMTk1OGUwMzk4ZGI3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CLK8JiZaoZj_bzK0BEpJFFGmI-pOD4hCDDB4FhI5fK0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:35 GMT
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
          "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c/campaign_types"
        },
        "campaigns": {
          "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c/campaigns"
        },
        "self": {
          "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c"
        }
      },
      "id": "fa92611e-d5a9-db99-3cb6-544a87584b3c",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYjc0NWE3OTAtMDc5Ny05NjAyLTMwZTQtMTk1OGUwMzk4ZGI3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CLK8JiZaoZj_bzK0BEpJFFGmI-pOD4hCDDB4FhI5fK0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:35 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c/campaign_types"
    },
    "campaigns": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c/campaigns"
    },
    "self": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c"
    }
  },
  "id": "fa92611e-d5a9-db99-3cb6-544a87584b3c",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYjc0NWE3OTAtMDc5Ny05NjAyLTMwZTQtMTk1OGUwMzk4ZGI3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CLK8JiZaoZj_bzK0BEpJFFGmI-pOD4hCDDB4FhI5fK0
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
Date: Mon, 19 Oct 2015 16:49:35 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c/campaign_types"
    },
    "campaigns": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c/campaigns"
    },
    "self": {
      "href": "/senders/fa92611e-d5a9-db99-3cb6-544a87584b3c"
    }
  },
  "id": "fa92611e-d5a9-db99-3cb6-544a87584b3c",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2RlNzFlMmUtNzg2Yi1mZmU5LTJiMjgtOGM4MmNmNDIxNjRjIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.25tLwzQ4b7DYJ2awXv2NPNv0BDK1z_yKu69QSbKP5fc
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:49:35 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTMyZGU5Y2MtMmQ4OS00NDc1LTM0NzgtMDFkNzEwOTlkYTMxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WF9-uQaUMRRDAPcAABi3ytbr01TwEzHESUhfqA2lsiY
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
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/22331894-12f8-a025-a356-e9e62b14abba"
    }
  },
  "html": "template html",
  "id": "22331894-12f8-a025-a356-e9e62b14abba",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTMyZGU5Y2MtMmQ4OS00NDc1LTM0NzgtMDFkNzEwOTlkYTMxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WF9-uQaUMRRDAPcAABi3ytbr01TwEzHESUhfqA2lsiY
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:36 GMT
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
          "href": "/templates/22331894-12f8-a025-a356-e9e62b14abba"
        }
      },
      "html": "html",
      "id": "22331894-12f8-a025-a356-e9e62b14abba",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTMyZGU5Y2MtMmQ4OS00NDc1LTM0NzgtMDFkNzEwOTlkYTMxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WF9-uQaUMRRDAPcAABi3ytbr01TwEzHESUhfqA2lsiY
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/22331894-12f8-a025-a356-e9e62b14abba"
    }
  },
  "html": "template html",
  "id": "22331894-12f8-a025-a356-e9e62b14abba",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTMyZGU5Y2MtMmQ4OS00NDc1LTM0NzgtMDFkNzEwOTlkYTMxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WF9-uQaUMRRDAPcAABi3ytbr01TwEzHESUhfqA2lsiY
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
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/22331894-12f8-a025-a356-e9e62b14abba"
    }
  },
  "html": "html",
  "id": "22331894-12f8-a025-a356-e9e62b14abba",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTMyZGU5Y2MtMmQ4OS00NDc1LTM0NzgtMDFkNzEwOTlkYTMxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WF9-uQaUMRRDAPcAABi3ytbr01TwEzHESUhfqA2lsiY
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:49:36 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc5OWIzN2QtMDQyZi02MjljLTA0MWItODA5MTZkZTNmYWE0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.IJ3DV895ZBh8LilysGXHmqHSZf5vYwSdg3Tb6Qg3WxA
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
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/173487ae-0c3b-5057-26e1-9cc1ff8dcc37"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "173487ae-0c3b-5057-26e1-9cc1ff8dcc37",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc5OWIzN2QtMDQyZi02MjljLTA0MWItODA5MTZkZTNmYWE0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.IJ3DV895ZBh8LilysGXHmqHSZf5vYwSdg3Tb6Qg3WxA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/21a805bd-acce-d46c-5dfb-2d2567fe2263/campaign_types"
    },
    "sender": {
      "href": "/senders/21a805bd-acce-d46c-5dfb-2d2567fe2263"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/173487ae-0c3b-5057-26e1-9cc1ff8dcc37"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "173487ae-0c3b-5057-26e1-9cc1ff8dcc37",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc5OWIzN2QtMDQyZi02MjljLTA0MWItODA5MTZkZTNmYWE0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.IJ3DV895ZBh8LilysGXHmqHSZf5vYwSdg3Tb6Qg3WxA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/173487ae-0c3b-5057-26e1-9cc1ff8dcc37"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "173487ae-0c3b-5057-26e1-9cc1ff8dcc37",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc5OWIzN2QtMDQyZi02MjljLTA0MWItODA5MTZkZTNmYWE0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.IJ3DV895ZBh8LilysGXHmqHSZf5vYwSdg3Tb6Qg3WxA
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "08296529-1f07-1ac9-78f9-2ee6cd1acfea"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/173487ae-0c3b-5057-26e1-9cc1ff8dcc37"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "173487ae-0c3b-5057-26e1-9cc1ff8dcc37",
  "name": "updated-campaign-type",
  "template_id": "08296529-1f07-1ac9-78f9-2ee6cd1acfea"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc5OWIzN2QtMDQyZi02MjljLTA0MWItODA5MTZkZTNmYWE0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.IJ3DV895ZBh8LilysGXHmqHSZf5vYwSdg3Tb6Qg3WxA
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:49:36 GMT
```


## Campaigns
Campaigns are an email to a set of users using a template provided directly or via a campaign type or via the default template.
<a name="campaign-create"></a>
### Create a new campaign
#### Request **POST** /senders/{id}/campaigns
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiOTQ4NjMzNGQtZThhNC02MjkwLTRmMjItYzg0ZWQxOWU0ZDQ4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.p7XxjWX_n-ARdb_kxiKCoX7HOl3nouhSQy5IKus5CaM
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "de2c5336-f28d-f51a-ca58-ed44cc9536d2",
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
      "aed9b38d-b033-4e3f-5d42-85164a19c32f"
    ]
  },
  "subject": "campaign subject",
  "template_id": "b1d8bf0e-07c1-4983-3763-78c78487d166",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 688
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/de2c5336-f28d-f51a-ca58-ed44cc9536d2"
    },
    "self": {
      "href": "/campaigns/00278db5-2691-aed4-5bd7-50ff82c92e79"
    },
    "status": {
      "href": "/campaigns/00278db5-2691-aed4-5bd7-50ff82c92e79/status"
    },
    "template": {
      "href": "/templates/b1d8bf0e-07c1-4983-3763-78c78487d166"
    }
  },
  "campaign_type_id": "de2c5336-f28d-f51a-ca58-ed44cc9536d2",
  "html": "",
  "id": "00278db5-2691-aed4-5bd7-50ff82c92e79",
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
      "aed9b38d-b033-4e3f-5d42-85164a19c32f"
    ]
  },
  "subject": "campaign subject",
  "template_id": "b1d8bf0e-07c1-4983-3763-78c78487d166",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiOTQ4NjMzNGQtZThhNC02MjkwLTRmMjItYzg0ZWQxOWU0ZDQ4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.p7XxjWX_n-ARdb_kxiKCoX7HOl3nouhSQy5IKus5CaM
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 688
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/de2c5336-f28d-f51a-ca58-ed44cc9536d2"
    },
    "self": {
      "href": "/campaigns/00278db5-2691-aed4-5bd7-50ff82c92e79"
    },
    "status": {
      "href": "/campaigns/00278db5-2691-aed4-5bd7-50ff82c92e79/status"
    },
    "template": {
      "href": "/templates/b1d8bf0e-07c1-4983-3763-78c78487d166"
    }
  },
  "campaign_type_id": "de2c5336-f28d-f51a-ca58-ed44cc9536d2",
  "html": "",
  "id": "00278db5-2691-aed4-5bd7-50ff82c92e79",
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
      "aed9b38d-b033-4e3f-5d42-85164a19c32f"
    ]
  },
  "subject": "campaign subject",
  "template_id": "b1d8bf0e-07c1-4983-3763-78c78487d166",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNThiNTQzZDgtZmY3My0zNjI0LTZmMDAtYmFmM2MwMTRmNWJiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.JXG3R_yGsOeZCBzCe3zw4Sz_D_PrzbOy7vVGdGAP5FA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 420
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:49:36 GMT
```
##### Body
```
{
  "_links": {
    "campaign": {
      "href": "/campaigns/464e92e4-6eb7-a4b3-9475-99878d3e342e"
    },
    "self": {
      "href": "/campaigns/464e92e4-6eb7-a4b3-9475-99878d3e342e/status"
    }
  },
  "completed_time": "2015-10-19T16:49:36Z",
  "failed_messages": 0,
  "id": "464e92e4-6eb7-a4b3-9475-99878d3e342e",
  "queued_messages": 0,
  "retry_messages": 0,
  "sent_messages": 1,
  "start_time": "2015-10-19T16:49:36Z",
  "status": "completed",
  "total_messages": 1,
  "undeliverable_messages": 0
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNmMzYTZkNDQtYzNjYS1mZjhjLWFiY2UtMTYxY2NhYjIyZDg5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.6lk3GJmSRz4_xehpoIhbrV8exN86wFiO-lY5AdaeXAo
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:49:34 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiOGFmZDIyOGQtZWZjZS00ODNhLTc1NGItZGZmMDk0OTQ1ZTc4In0.l4wqCUtkYNyxcnv_uPQNMt-_x__Glw1SI57-8jlik6Q
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:49:32 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNmMzYTZkNDQtYzNjYS1mZjhjLWFiY2UtMTYxY2NhYjIyZDg5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1ODA4Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.6lk3GJmSRz4_xehpoIhbrV8exN86wFiO-lY5AdaeXAo
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:49:34 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiOGFmZDIyOGQtZWZjZS00ODNhLTc1NGItZGZmMDk0OTQ1ZTc4In0.l4wqCUtkYNyxcnv_uPQNMt-_x__Glw1SI57-8jlik6Q
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:49:32 GMT
```


