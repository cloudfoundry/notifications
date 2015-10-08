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
  * [Retrieve a senders](#sender-get)
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
Date: Thu, 08 Oct 2015 17:44:33 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTk4MzA0YzQtZjAyZS00Yjc0LTdjMmMtN2RkZWMxMDc5ZWRhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.zBODY4xwxge2teVsBIBqWVyaeAhYZSfMPPrlNuGaeeA
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
Date: Thu, 08 Oct 2015 17:44:33 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239/campaign_types"
    },
    "campaigns": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239/campaigns"
    },
    "self": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239"
    }
  },
  "id": "6c0ba4bf-d244-8d8e-0a48-1af5412e0239",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTk4MzA0YzQtZjAyZS00Yjc0LTdjMmMtN2RkZWMxMDc5ZWRhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.zBODY4xwxge2teVsBIBqWVyaeAhYZSfMPPrlNuGaeeA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:33 GMT
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
          "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239/campaign_types"
        },
        "campaigns": {
          "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239/campaigns"
        },
        "self": {
          "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239"
        }
      },
      "id": "6c0ba4bf-d244-8d8e-0a48-1af5412e0239",
      "name": "My Cool App"
    }
  ]
}
```

<a name="sender-get"></a>
### Retrieve a senders
#### Request **GET** /senders/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTk4MzA0YzQtZjAyZS00Yjc0LTdjMmMtN2RkZWMxMDc5ZWRhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.zBODY4xwxge2teVsBIBqWVyaeAhYZSfMPPrlNuGaeeA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:33 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239/campaign_types"
    },
    "campaigns": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239/campaigns"
    },
    "self": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239"
    }
  },
  "id": "6c0ba4bf-d244-8d8e-0a48-1af5412e0239",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMTk4MzA0YzQtZjAyZS00Yjc0LTdjMmMtN2RkZWMxMDc5ZWRhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.zBODY4xwxge2teVsBIBqWVyaeAhYZSfMPPrlNuGaeeA
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
Date: Thu, 08 Oct 2015 17:44:33 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239/campaign_types"
    },
    "campaigns": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239/campaigns"
    },
    "self": {
      "href": "/senders/6c0ba4bf-d244-8d8e-0a48-1af5412e0239"
    }
  },
  "id": "6c0ba4bf-d244-8d8e-0a48-1af5412e0239",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTc0NDMzOTMtMWUzNS00YzBlLTZhYWEtMjgwMWM4MGZiYmUyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.zUdlNIVBHLW9I6JJhDVUmoWPevuv6npXimypob8t2OU
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Thu, 08 Oct 2015 17:44:35 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjc4Y2M0ZjMtZDA2ZS00ZWEwLTY2NzItM2U0ODYxMjM1MzNhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.qp_T-qUfuxwLl1XPbcy4fd2ItJwaFkTKqpvtz-W7ABU
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
Date: Thu, 08 Oct 2015 17:44:34 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/9a92141d-8c8f-f70b-3f76-386c8d4a15e8"
    }
  },
  "html": "template html",
  "id": "9a92141d-8c8f-f70b-3f76-386c8d4a15e8",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjc4Y2M0ZjMtZDA2ZS00ZWEwLTY2NzItM2U0ODYxMjM1MzNhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.qp_T-qUfuxwLl1XPbcy4fd2ItJwaFkTKqpvtz-W7ABU
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:34 GMT
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
          "href": "/templates/9a92141d-8c8f-f70b-3f76-386c8d4a15e8"
        }
      },
      "html": "html",
      "id": "9a92141d-8c8f-f70b-3f76-386c8d4a15e8",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjc4Y2M0ZjMtZDA2ZS00ZWEwLTY2NzItM2U0ODYxMjM1MzNhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.qp_T-qUfuxwLl1XPbcy4fd2ItJwaFkTKqpvtz-W7ABU
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:34 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/9a92141d-8c8f-f70b-3f76-386c8d4a15e8"
    }
  },
  "html": "template html",
  "id": "9a92141d-8c8f-f70b-3f76-386c8d4a15e8",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjc4Y2M0ZjMtZDA2ZS00ZWEwLTY2NzItM2U0ODYxMjM1MzNhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.qp_T-qUfuxwLl1XPbcy4fd2ItJwaFkTKqpvtz-W7ABU
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
Date: Thu, 08 Oct 2015 17:44:34 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/9a92141d-8c8f-f70b-3f76-386c8d4a15e8"
    }
  },
  "html": "html",
  "id": "9a92141d-8c8f-f70b-3f76-386c8d4a15e8",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjc4Y2M0ZjMtZDA2ZS00ZWEwLTY2NzItM2U0ODYxMjM1MzNhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.qp_T-qUfuxwLl1XPbcy4fd2ItJwaFkTKqpvtz-W7ABU
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Thu, 08 Oct 2015 17:44:34 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNTAxMDk0MzItZGEwMy00NzEwLTQ0YzAtOWVhNjg3YmUxYmZhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.VRKgRbaSK-xMxHo5r3lVDWWCJNxyhhuPZG3J10MIUEk
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
Date: Thu, 08 Oct 2015 17:44:35 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/600b9769-bdbd-1ef2-6078-4c3738ab0b4f"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "600b9769-bdbd-1ef2-6078-4c3738ab0b4f",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNTAxMDk0MzItZGEwMy00NzEwLTQ0YzAtOWVhNjg3YmUxYmZhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.VRKgRbaSK-xMxHo5r3lVDWWCJNxyhhuPZG3J10MIUEk
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:35 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/0e59fa37-0700-8421-905a-24bf9f8e1a3d/campaign_types"
    },
    "sender": {
      "href": "/senders/0e59fa37-0700-8421-905a-24bf9f8e1a3d"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/600b9769-bdbd-1ef2-6078-4c3738ab0b4f"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "600b9769-bdbd-1ef2-6078-4c3738ab0b4f",
      "name": "some-campaign-type",
      "template_id": ""
    }
  ]
}
```

<a name="campaign-type-get"></a>
### Retrieve a campaign type
#### Request **GET** /senders/{id}/campaign_types/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNTAxMDk0MzItZGEwMy00NzEwLTQ0YzAtOWVhNjg3YmUxYmZhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.VRKgRbaSK-xMxHo5r3lVDWWCJNxyhhuPZG3J10MIUEk
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:35 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/600b9769-bdbd-1ef2-6078-4c3738ab0b4f"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "600b9769-bdbd-1ef2-6078-4c3738ab0b4f",
  "name": "some-campaign-type",
  "template_id": ""
}
```

<a name="campaign-type-update"></a>
### Update a campaign type
#### Request **PUT** /senders/{id}/campaign_types/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNTAxMDk0MzItZGEwMy00NzEwLTQ0YzAtOWVhNjg3YmUxYmZhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.VRKgRbaSK-xMxHo5r3lVDWWCJNxyhhuPZG3J10MIUEk
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "4def041e-b710-1802-a697-8710096e67b0"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:35 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/600b9769-bdbd-1ef2-6078-4c3738ab0b4f"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "600b9769-bdbd-1ef2-6078-4c3738ab0b4f",
  "name": "updated-campaign-type",
  "template_id": "4def041e-b710-1802-a697-8710096e67b0"
}
```

<a name="campaign-type-delete"></a>
### Delete a campaign type
#### Request **DELETE** /senders/{id}/campaign_types/{id}
##### Required Scopes
```
notifications.write
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNTAxMDk0MzItZGEwMy00NzEwLTQ0YzAtOWVhNjg3YmUxYmZhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.VRKgRbaSK-xMxHo5r3lVDWWCJNxyhhuPZG3J10MIUEk
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Thu, 08 Oct 2015 17:44:35 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiOGE5MDBmYTQtY2E3Yy00OTc4LTQwMGUtNDc2MmY4Nzc5NGJiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2_SaE8K_fqKfR--jxWqqz65afMY6gKmViNJFicY_V20
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "3ad101f9-15d8-44c5-ca76-4c1fa87ba6a4",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "email": "test@example.com"
  },
  "subject": "campaign subject",
  "template_id": "61e6bfa1-0ff2-79b0-852a-7729c398ab18",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 593
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:36 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/3ad101f9-15d8-44c5-ca76-4c1fa87ba6a4"
    },
    "self": {
      "href": "/campaigns/0c4f07bd-99b5-0015-5f10-201ef070efb4"
    },
    "status": {
      "href": "/campaigns/0c4f07bd-99b5-0015-5f10-201ef070efb4/status"
    },
    "template": {
      "href": "/templates/61e6bfa1-0ff2-79b0-852a-7729c398ab18"
    }
  },
  "campaign_type_id": "3ad101f9-15d8-44c5-ca76-4c1fa87ba6a4",
  "html": "",
  "id": "0c4f07bd-99b5-0015-5f10-201ef070efb4",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "email": "test@example.com"
  },
  "subject": "campaign subject",
  "template_id": "61e6bfa1-0ff2-79b0-852a-7729c398ab18",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiOGE5MDBmYTQtY2E3Yy00OTc4LTQwMGUtNDc2MmY4Nzc5NGJiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1NTY1MC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2_SaE8K_fqKfR--jxWqqz65afMY6gKmViNJFicY_V20
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 593
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 17:44:36 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/3ad101f9-15d8-44c5-ca76-4c1fa87ba6a4"
    },
    "self": {
      "href": "/campaigns/0c4f07bd-99b5-0015-5f10-201ef070efb4"
    },
    "status": {
      "href": "/campaigns/0c4f07bd-99b5-0015-5f10-201ef070efb4/status"
    },
    "template": {
      "href": "/templates/61e6bfa1-0ff2-79b0-852a-7729c398ab18"
    }
  },
  "campaign_type_id": "3ad101f9-15d8-44c5-ca76-4c1fa87ba6a4",
  "html": "",
  "id": "0c4f07bd-99b5-0015-5f10-201ef070efb4",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "email": "test@example.com"
  },
  "subject": "campaign subject",
  "template_id": "61e6bfa1-0ff2-79b0-852a-7729c398ab18",
  "text": "campaign body"
}
```


