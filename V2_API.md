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
Date: Mon, 19 Oct 2015 16:13:09 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDUyZTUzMDctYWRmYi03NmU5LTFmNzQtZmUyYzhhZTAyZWEwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bDqJAMI1vyKqcB71xhcVOKsCzw5FHkD5YLE0_bulkss
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
Date: Mon, 19 Oct 2015 16:13:12 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f/campaign_types"
    },
    "campaigns": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f/campaigns"
    },
    "self": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f"
    }
  },
  "id": "23667e74-e644-cd2d-c466-3a174bd6cd5f",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDUyZTUzMDctYWRmYi03NmU5LTFmNzQtZmUyYzhhZTAyZWEwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bDqJAMI1vyKqcB71xhcVOKsCzw5FHkD5YLE0_bulkss
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:12 GMT
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
          "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f/campaign_types"
        },
        "campaigns": {
          "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f/campaigns"
        },
        "self": {
          "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f"
        }
      },
      "id": "23667e74-e644-cd2d-c466-3a174bd6cd5f",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDUyZTUzMDctYWRmYi03NmU5LTFmNzQtZmUyYzhhZTAyZWEwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bDqJAMI1vyKqcB71xhcVOKsCzw5FHkD5YLE0_bulkss
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:12 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f/campaign_types"
    },
    "campaigns": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f/campaigns"
    },
    "self": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f"
    }
  },
  "id": "23667e74-e644-cd2d-c466-3a174bd6cd5f",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDUyZTUzMDctYWRmYi03NmU5LTFmNzQtZmUyYzhhZTAyZWEwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bDqJAMI1vyKqcB71xhcVOKsCzw5FHkD5YLE0_bulkss
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
Date: Mon, 19 Oct 2015 16:13:12 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f/campaign_types"
    },
    "campaigns": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f/campaigns"
    },
    "self": {
      "href": "/senders/23667e74-e644-cd2d-c466-3a174bd6cd5f"
    }
  },
  "id": "23667e74-e644-cd2d-c466-3a174bd6cd5f",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDMyM2U0YTItZjljYy1kOWQyLWQzYzUtOTdiMDdkYTJiMTY3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.c0idlgqtxIytH9Z8_Cvzxae5y44JGr-zx-Wh4yW1GIE
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:13:09 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMjc1NzliMTgtYTMwYi1jNDc4LTEwYmQtNTZkZWE0OTJjOTAzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.421FmDlouxhlGOX_bVJGcPKO_ep8Ip_JUVM2uoDjCs8
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
Date: Mon, 19 Oct 2015 16:13:11 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/51085550-c848-e675-465d-6bf18b1a2701"
    }
  },
  "html": "template html",
  "id": "51085550-c848-e675-465d-6bf18b1a2701",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMjc1NzliMTgtYTMwYi1jNDc4LTEwYmQtNTZkZWE0OTJjOTAzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.421FmDlouxhlGOX_bVJGcPKO_ep8Ip_JUVM2uoDjCs8
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:11 GMT
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
          "href": "/templates/51085550-c848-e675-465d-6bf18b1a2701"
        }
      },
      "html": "html",
      "id": "51085550-c848-e675-465d-6bf18b1a2701",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMjc1NzliMTgtYTMwYi1jNDc4LTEwYmQtNTZkZWE0OTJjOTAzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.421FmDlouxhlGOX_bVJGcPKO_ep8Ip_JUVM2uoDjCs8
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:11 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/51085550-c848-e675-465d-6bf18b1a2701"
    }
  },
  "html": "template html",
  "id": "51085550-c848-e675-465d-6bf18b1a2701",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMjc1NzliMTgtYTMwYi1jNDc4LTEwYmQtNTZkZWE0OTJjOTAzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.421FmDlouxhlGOX_bVJGcPKO_ep8Ip_JUVM2uoDjCs8
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
Date: Mon, 19 Oct 2015 16:13:11 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/51085550-c848-e675-465d-6bf18b1a2701"
    }
  },
  "html": "html",
  "id": "51085550-c848-e675-465d-6bf18b1a2701",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMjc1NzliMTgtYTMwYi1jNDc4LTEwYmQtNTZkZWE0OTJjOTAzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.421FmDlouxhlGOX_bVJGcPKO_ep8Ip_JUVM2uoDjCs8
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:13:11 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlNjE3MDAtZTZhOC0yYzdiLWQ3MTItYjc0ZjM4OTUwZjQzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.mJY3eVxjRn3f7Ln1c_8IktEgT7bbKRhC-A0yp3bBsM0
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
Date: Mon, 19 Oct 2015 16:13:09 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/b1573ef5-d086-8138-1026-3fb0b2f01d08"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "b1573ef5-d086-8138-1026-3fb0b2f01d08",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlNjE3MDAtZTZhOC0yYzdiLWQ3MTItYjc0ZjM4OTUwZjQzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.mJY3eVxjRn3f7Ln1c_8IktEgT7bbKRhC-A0yp3bBsM0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:09 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/484ed6e8-8ec6-3701-5ca9-af337afb589d/campaign_types"
    },
    "sender": {
      "href": "/senders/484ed6e8-8ec6-3701-5ca9-af337afb589d"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/b1573ef5-d086-8138-1026-3fb0b2f01d08"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "b1573ef5-d086-8138-1026-3fb0b2f01d08",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlNjE3MDAtZTZhOC0yYzdiLWQ3MTItYjc0ZjM4OTUwZjQzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.mJY3eVxjRn3f7Ln1c_8IktEgT7bbKRhC-A0yp3bBsM0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:09 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/b1573ef5-d086-8138-1026-3fb0b2f01d08"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "b1573ef5-d086-8138-1026-3fb0b2f01d08",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlNjE3MDAtZTZhOC0yYzdiLWQ3MTItYjc0ZjM4OTUwZjQzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.mJY3eVxjRn3f7Ln1c_8IktEgT7bbKRhC-A0yp3bBsM0
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "75207c6d-409b-68f6-e19f-f99b2ff677c7"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:09 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/b1573ef5-d086-8138-1026-3fb0b2f01d08"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "b1573ef5-d086-8138-1026-3fb0b2f01d08",
  "name": "updated-campaign-type",
  "template_id": "75207c6d-409b-68f6-e19f-f99b2ff677c7"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlNjE3MDAtZTZhOC0yYzdiLWQ3MTItYjc0ZjM4OTUwZjQzIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.mJY3eVxjRn3f7Ln1c_8IktEgT7bbKRhC-A0yp3bBsM0
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:13:09 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZWU1NTcwYTktMjI2Yy04MzE2LTFhN2EtZTNmNGVkYTgxMjM4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.hLJWZXejckku6xX72a0-fhwcXlwPT7anfNQ7U7z9idQ
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "79db33f4-d1b1-146f-d26b-0e4341bd7ac0",
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
      "d04373bb-6c13-42f3-48b5-7f959e7bc7d5"
    ]
  },
  "subject": "campaign subject",
  "template_id": "17244696-bce9-4ed7-66d2-01e3c594070b",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:10 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/79db33f4-d1b1-146f-d26b-0e4341bd7ac0"
    },
    "self": {
      "href": "/campaigns/c6c28f6c-828f-67ce-bced-9eb2b2701d11"
    },
    "status": {
      "href": "/campaigns/c6c28f6c-828f-67ce-bced-9eb2b2701d11/status"
    },
    "template": {
      "href": "/templates/17244696-bce9-4ed7-66d2-01e3c594070b"
    }
  },
  "campaign_type_id": "79db33f4-d1b1-146f-d26b-0e4341bd7ac0",
  "html": "",
  "id": "c6c28f6c-828f-67ce-bced-9eb2b2701d11",
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
      "d04373bb-6c13-42f3-48b5-7f959e7bc7d5"
    ]
  },
  "subject": "campaign subject",
  "template_id": "17244696-bce9-4ed7-66d2-01e3c594070b",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZWU1NTcwYTktMjI2Yy04MzE2LTFhN2EtZTNmNGVkYTgxMjM4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.hLJWZXejckku6xX72a0-fhwcXlwPT7anfNQ7U7z9idQ
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:10 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/79db33f4-d1b1-146f-d26b-0e4341bd7ac0"
    },
    "self": {
      "href": "/campaigns/c6c28f6c-828f-67ce-bced-9eb2b2701d11"
    },
    "status": {
      "href": "/campaigns/c6c28f6c-828f-67ce-bced-9eb2b2701d11/status"
    },
    "template": {
      "href": "/templates/17244696-bce9-4ed7-66d2-01e3c594070b"
    }
  },
  "campaign_type_id": "79db33f4-d1b1-146f-d26b-0e4341bd7ac0",
  "html": "",
  "id": "c6c28f6c-828f-67ce-bced-9eb2b2701d11",
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
      "d04373bb-6c13-42f3-48b5-7f959e7bc7d5"
    ]
  },
  "subject": "campaign subject",
  "template_id": "17244696-bce9-4ed7-66d2-01e3c594070b",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTI0NzczOTYtMjA4Ni1lMzVhLTEzMmEtZWRjOGRjOWJjZDY1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.ggtHMDIS_XXKCbD2Z5HES5onj_hh7FRBNqx6hyQWP5w
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 420
Content-Type: text/plain; charset=utf-8
Date: Mon, 19 Oct 2015 16:13:10 GMT
```
##### Body
```
{
  "_links": {
    "campaign": {
      "href": "/campaigns/79e96bda-7d58-a415-ad02-2eb60c5387b3"
    },
    "self": {
      "href": "/campaigns/79e96bda-7d58-a415-ad02-2eb60c5387b3/status"
    }
  },
  "completed_time": "2015-10-19T16:13:10Z",
  "failed_messages": 0,
  "id": "79e96bda-7d58-a415-ad02-2eb60c5387b3",
  "queued_messages": 0,
  "retry_messages": 0,
  "sent_messages": 1,
  "start_time": "2015-10-19T16:13:10Z",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYzAxMzRlN2UtYWNkYS05NGI0LWZiOTAtYzYwNWI3NjY5OGNmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.bZMJxHu8pBmRuGopg74FrXq7_laDgdLKY3zv6fMeysY
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:13:09 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiMzNhMzVhNTgtNTUyYy00NjFlLTVhYzUtZjkwOTdmZDFkYjJhIn0.ExFsYwzUMAAwbtW3tk-ZT3vrOulxykR3VQ5TVIdx89o
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:13:10 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYzAxMzRlN2UtYWNkYS05NGI0LWZiOTAtYzYwNWI3NjY5OGNmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjkwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.bZMJxHu8pBmRuGopg74FrXq7_laDgdLKY3zv6fMeysY
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:13:09 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiMzNhMzVhNTgtNTUyYy00NjFlLTVhYzUtZjkwOTdmZDFkYjJhIn0.ExFsYwzUMAAwbtW3tk-ZT3vrOulxykR3VQ5TVIdx89o
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Mon, 19 Oct 2015 16:13:10 GMT
```


