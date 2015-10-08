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
Date: Thu, 08 Oct 2015 19:10:44 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc4NGVjNzAtM2JlMC00ZjJiLTU2NDYtNWQzOTFmZGM0Njc1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MBdqS73iR6jaMWLVYSYQvpo8DIfTOgcTQjtUaEzIM58
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
Date: Thu, 08 Oct 2015 19:10:45 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671/campaign_types"
    },
    "campaigns": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671/campaigns"
    },
    "self": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671"
    }
  },
  "id": "61cdfb00-2db2-813b-de4d-e52649c2f671",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc4NGVjNzAtM2JlMC00ZjJiLTU2NDYtNWQzOTFmZGM0Njc1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MBdqS73iR6jaMWLVYSYQvpo8DIfTOgcTQjtUaEzIM58
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:45 GMT
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
          "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671/campaign_types"
        },
        "campaigns": {
          "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671/campaigns"
        },
        "self": {
          "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671"
        }
      },
      "id": "61cdfb00-2db2-813b-de4d-e52649c2f671",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc4NGVjNzAtM2JlMC00ZjJiLTU2NDYtNWQzOTFmZGM0Njc1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MBdqS73iR6jaMWLVYSYQvpo8DIfTOgcTQjtUaEzIM58
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:45 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671/campaign_types"
    },
    "campaigns": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671/campaigns"
    },
    "self": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671"
    }
  },
  "id": "61cdfb00-2db2-813b-de4d-e52649c2f671",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODc4NGVjNzAtM2JlMC00ZjJiLTU2NDYtNWQzOTFmZGM0Njc1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.MBdqS73iR6jaMWLVYSYQvpo8DIfTOgcTQjtUaEzIM58
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
Date: Thu, 08 Oct 2015 19:10:45 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671/campaign_types"
    },
    "campaigns": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671/campaigns"
    },
    "self": {
      "href": "/senders/61cdfb00-2db2-813b-de4d-e52649c2f671"
    }
  },
  "id": "61cdfb00-2db2-813b-de4d-e52649c2f671",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNzllOTQ1ZWYtZGFkYy00NDAwLTcwZDAtNzUyNmRlODkxMDIwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.qmgYwjnQlLJd0TJY1i3nqtBF_GVu9CB4uxobIFHe3A0
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Thu, 08 Oct 2015 19:10:44 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2I4ODE5NWQtODkzMi00MDY5LTUxMGUtN2M2MDQ2ZDdiYmU4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dz5AZQu3An5JikmJ5CbAgW5ydKwfqNe1AIJla6h9mMY
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
Date: Thu, 08 Oct 2015 19:10:45 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/450f3dbd-eb96-4278-f3ed-a5af7dc67990"
    }
  },
  "html": "template html",
  "id": "450f3dbd-eb96-4278-f3ed-a5af7dc67990",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2I4ODE5NWQtODkzMi00MDY5LTUxMGUtN2M2MDQ2ZDdiYmU4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dz5AZQu3An5JikmJ5CbAgW5ydKwfqNe1AIJla6h9mMY
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:45 GMT
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
          "href": "/templates/450f3dbd-eb96-4278-f3ed-a5af7dc67990"
        }
      },
      "html": "html",
      "id": "450f3dbd-eb96-4278-f3ed-a5af7dc67990",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2I4ODE5NWQtODkzMi00MDY5LTUxMGUtN2M2MDQ2ZDdiYmU4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dz5AZQu3An5JikmJ5CbAgW5ydKwfqNe1AIJla6h9mMY
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:45 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/450f3dbd-eb96-4278-f3ed-a5af7dc67990"
    }
  },
  "html": "template html",
  "id": "450f3dbd-eb96-4278-f3ed-a5af7dc67990",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2I4ODE5NWQtODkzMi00MDY5LTUxMGUtN2M2MDQ2ZDdiYmU4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dz5AZQu3An5JikmJ5CbAgW5ydKwfqNe1AIJla6h9mMY
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
Date: Thu, 08 Oct 2015 19:10:45 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/450f3dbd-eb96-4278-f3ed-a5af7dc67990"
    }
  },
  "html": "html",
  "id": "450f3dbd-eb96-4278-f3ed-a5af7dc67990",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2I4ODE5NWQtODkzMi00MDY5LTUxMGUtN2M2MDQ2ZDdiYmU4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dz5AZQu3An5JikmJ5CbAgW5ydKwfqNe1AIJla6h9mMY
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Thu, 08 Oct 2015 19:10:45 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Y1MWQwNzMtM2NiMS00NDM5LTU0NGMtNGMwMjU3MmJiMWRmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.XPkMousqYx9wuv4GsFwNI0rOBDCsPUMpD8ikk50ak8o
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
Date: Thu, 08 Oct 2015 19:10:43 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/0cbc8499-c49d-c841-caf5-67ae2314b661"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "0cbc8499-c49d-c841-caf5-67ae2314b661",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Y1MWQwNzMtM2NiMS00NDM5LTU0NGMtNGMwMjU3MmJiMWRmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.XPkMousqYx9wuv4GsFwNI0rOBDCsPUMpD8ikk50ak8o
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:43 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/f4a5549b-9e50-87e9-bbb5-9e083db2d6f6/campaign_types"
    },
    "sender": {
      "href": "/senders/f4a5549b-9e50-87e9-bbb5-9e083db2d6f6"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/0cbc8499-c49d-c841-caf5-67ae2314b661"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "0cbc8499-c49d-c841-caf5-67ae2314b661",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Y1MWQwNzMtM2NiMS00NDM5LTU0NGMtNGMwMjU3MmJiMWRmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.XPkMousqYx9wuv4GsFwNI0rOBDCsPUMpD8ikk50ak8o
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:43 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/0cbc8499-c49d-c841-caf5-67ae2314b661"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "0cbc8499-c49d-c841-caf5-67ae2314b661",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Y1MWQwNzMtM2NiMS00NDM5LTU0NGMtNGMwMjU3MmJiMWRmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.XPkMousqYx9wuv4GsFwNI0rOBDCsPUMpD8ikk50ak8o
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "9b2f1281-4ac0-28da-ccc7-4460db1358b8"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:43 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/0cbc8499-c49d-c841-caf5-67ae2314b661"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "0cbc8499-c49d-c841-caf5-67ae2314b661",
  "name": "updated-campaign-type",
  "template_id": "9b2f1281-4ac0-28da-ccc7-4460db1358b8"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Y1MWQwNzMtM2NiMS00NDM5LTU0NGMtNGMwMjU3MmJiMWRmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.XPkMousqYx9wuv4GsFwNI0rOBDCsPUMpD8ikk50ak8o
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Thu, 08 Oct 2015 19:10:43 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZWY1MmFlMjMtODdlYS00MTlkLTcyZmMtMWNmNzNmOTU1N2E5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.GxIKYQRb1TGlkKE-Am28gacYVPA5mXVq9k-Bk0hT15U
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "8e69a00a-da26-36ed-b4fc-28a2150c6593",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "emails": [
      "test@example.com"
    ]
  },
  "subject": "campaign subject",
  "template_id": "82248ef0-4b71-9438-d08a-2742879c5698",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 596
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:44 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/8e69a00a-da26-36ed-b4fc-28a2150c6593"
    },
    "self": {
      "href": "/campaigns/beb846d2-3a45-3577-5939-210660df6ffc"
    },
    "status": {
      "href": "/campaigns/beb846d2-3a45-3577-5939-210660df6ffc/status"
    },
    "template": {
      "href": "/templates/82248ef0-4b71-9438-d08a-2742879c5698"
    }
  },
  "campaign_type_id": "8e69a00a-da26-36ed-b4fc-28a2150c6593",
  "html": "",
  "id": "beb846d2-3a45-3577-5939-210660df6ffc",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "emails": [
      "test@example.com"
    ]
  },
  "subject": "campaign subject",
  "template_id": "82248ef0-4b71-9438-d08a-2742879c5698",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZWY1MmFlMjMtODdlYS00MTlkLTcyZmMtMWNmNzNmOTU1N2E5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTI4OC9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.GxIKYQRb1TGlkKE-Am28gacYVPA5mXVq9k-Bk0hT15U
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 596
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2015 19:10:44 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/8e69a00a-da26-36ed-b4fc-28a2150c6593"
    },
    "self": {
      "href": "/campaigns/beb846d2-3a45-3577-5939-210660df6ffc"
    },
    "status": {
      "href": "/campaigns/beb846d2-3a45-3577-5939-210660df6ffc/status"
    },
    "template": {
      "href": "/templates/82248ef0-4b71-9438-d08a-2742879c5698"
    }
  },
  "campaign_type_id": "8e69a00a-da26-36ed-b4fc-28a2150c6593",
  "html": "",
  "id": "beb846d2-3a45-3577-5939-210660df6ffc",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "emails": [
      "test@example.com"
    ]
  },
  "subject": "campaign subject",
  "template_id": "82248ef0-4b71-9438-d08a-2742879c5698",
  "text": "campaign body"
}
```


