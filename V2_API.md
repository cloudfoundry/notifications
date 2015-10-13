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
  * [Unsubscribe a user](#unsubscriber-put)
  * [Resubscribe a user](#unsubscriber-delete)

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
Date: Tue, 13 Oct 2015 15:57:39 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzBhZWU1NGItYjRiYy1lNTEzLTM4NDQtNzBmNmViM2E1MGRlIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.y6hXh2NTAWT1b1TrPlGh02uqmg4Hg_ThSVKx27zOm_I
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
Date: Tue, 13 Oct 2015 15:57:38 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834/campaign_types"
    },
    "campaigns": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834/campaigns"
    },
    "self": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834"
    }
  },
  "id": "3e4bbcf7-f25c-146d-deb3-6c4ad3e06834",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzBhZWU1NGItYjRiYy1lNTEzLTM4NDQtNzBmNmViM2E1MGRlIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.y6hXh2NTAWT1b1TrPlGh02uqmg4Hg_ThSVKx27zOm_I
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:38 GMT
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
          "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834/campaign_types"
        },
        "campaigns": {
          "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834/campaigns"
        },
        "self": {
          "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834"
        }
      },
      "id": "3e4bbcf7-f25c-146d-deb3-6c4ad3e06834",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzBhZWU1NGItYjRiYy1lNTEzLTM4NDQtNzBmNmViM2E1MGRlIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.y6hXh2NTAWT1b1TrPlGh02uqmg4Hg_ThSVKx27zOm_I
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:38 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834/campaign_types"
    },
    "campaigns": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834/campaigns"
    },
    "self": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834"
    }
  },
  "id": "3e4bbcf7-f25c-146d-deb3-6c4ad3e06834",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzBhZWU1NGItYjRiYy1lNTEzLTM4NDQtNzBmNmViM2E1MGRlIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.y6hXh2NTAWT1b1TrPlGh02uqmg4Hg_ThSVKx27zOm_I
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
Date: Tue, 13 Oct 2015 15:57:38 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834/campaign_types"
    },
    "campaigns": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834/campaigns"
    },
    "self": {
      "href": "/senders/3e4bbcf7-f25c-146d-deb3-6c4ad3e06834"
    }
  },
  "id": "3e4bbcf7-f25c-146d-deb3-6c4ad3e06834",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiOTljMTZhNTctNjkxZi1kYzRhLWQwM2MtMzZkMTI1ZDk4MjFjIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.mCDqUAUYjBlvdkTwoaP6mCoVHz7tJ2YNDFJF5_79Sio
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 15:57:39 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGZjMzI3NDEtZGQzMy03NWJiLTFiYzgtMWM3MTliOTQ2NzA4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cfP1aRaISbprUhkUh6UlxSQH4mk77JzpzaAElN_Us2s
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
Date: Tue, 13 Oct 2015 15:57:39 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/d0fc938a-c124-dc33-3234-ebc76577681b"
    }
  },
  "html": "template html",
  "id": "d0fc938a-c124-dc33-3234-ebc76577681b",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGZjMzI3NDEtZGQzMy03NWJiLTFiYzgtMWM3MTliOTQ2NzA4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cfP1aRaISbprUhkUh6UlxSQH4mk77JzpzaAElN_Us2s
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:39 GMT
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
          "href": "/templates/d0fc938a-c124-dc33-3234-ebc76577681b"
        }
      },
      "html": "html",
      "id": "d0fc938a-c124-dc33-3234-ebc76577681b",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGZjMzI3NDEtZGQzMy03NWJiLTFiYzgtMWM3MTliOTQ2NzA4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cfP1aRaISbprUhkUh6UlxSQH4mk77JzpzaAElN_Us2s
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:39 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/d0fc938a-c124-dc33-3234-ebc76577681b"
    }
  },
  "html": "template html",
  "id": "d0fc938a-c124-dc33-3234-ebc76577681b",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGZjMzI3NDEtZGQzMy03NWJiLTFiYzgtMWM3MTliOTQ2NzA4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cfP1aRaISbprUhkUh6UlxSQH4mk77JzpzaAElN_Us2s
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
Date: Tue, 13 Oct 2015 15:57:39 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/d0fc938a-c124-dc33-3234-ebc76577681b"
    }
  },
  "html": "html",
  "id": "d0fc938a-c124-dc33-3234-ebc76577681b",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGZjMzI3NDEtZGQzMy03NWJiLTFiYzgtMWM3MTliOTQ2NzA4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.cfP1aRaISbprUhkUh6UlxSQH4mk77JzpzaAElN_Us2s
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 15:57:39 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZmNjZWQ1M2EtNWMzZS1mN2Q5LTlhZGMtYmJlZjQxZDZmMDAxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CwUDlz1-UaVqWp-XC00yQM4-50-ti85KJl5hvlASoMY
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
Date: Tue, 13 Oct 2015 15:57:37 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/5a542f9b-c0b1-7728-f94e-97f356def40f"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "5a542f9b-c0b1-7728-f94e-97f356def40f",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZmNjZWQ1M2EtNWMzZS1mN2Q5LTlhZGMtYmJlZjQxZDZmMDAxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CwUDlz1-UaVqWp-XC00yQM4-50-ti85KJl5hvlASoMY
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:37 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/2bc12812-36d6-e409-b5d7-2fa3b4751df8/campaign_types"
    },
    "sender": {
      "href": "/senders/2bc12812-36d6-e409-b5d7-2fa3b4751df8"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/5a542f9b-c0b1-7728-f94e-97f356def40f"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "5a542f9b-c0b1-7728-f94e-97f356def40f",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZmNjZWQ1M2EtNWMzZS1mN2Q5LTlhZGMtYmJlZjQxZDZmMDAxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CwUDlz1-UaVqWp-XC00yQM4-50-ti85KJl5hvlASoMY
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:37 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/5a542f9b-c0b1-7728-f94e-97f356def40f"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "5a542f9b-c0b1-7728-f94e-97f356def40f",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZmNjZWQ1M2EtNWMzZS1mN2Q5LTlhZGMtYmJlZjQxZDZmMDAxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CwUDlz1-UaVqWp-XC00yQM4-50-ti85KJl5hvlASoMY
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "11559fc2-db5f-b3d4-c6c4-931d52cdc835"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:37 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/5a542f9b-c0b1-7728-f94e-97f356def40f"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "5a542f9b-c0b1-7728-f94e-97f356def40f",
  "name": "updated-campaign-type",
  "template_id": "11559fc2-db5f-b3d4-c6c4-931d52cdc835"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZmNjZWQ1M2EtNWMzZS1mN2Q5LTlhZGMtYmJlZjQxZDZmMDAxIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.CwUDlz1-UaVqWp-XC00yQM4-50-ti85KJl5hvlASoMY
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 15:57:37 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODI4NTRhMzgtZGM4YS03NThhLWE2Y2YtZjdjMjY0MjBkNThkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.57BvXEasfG8DOGC_5PQtIkbpgY8SvHlTJILs6Z1Iw5k
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "c3e443e5-64d5-5d77-349f-e2c2da51eed6",
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
      "8d72af9b-bfa5-4d87-7953-64739c584169"
    ]
  },
  "subject": "campaign subject",
  "template_id": "e673fb66-2d7a-7401-52fe-30cac8ee2dd3",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:34 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/c3e443e5-64d5-5d77-349f-e2c2da51eed6"
    },
    "self": {
      "href": "/campaigns/c652c6fa-07d7-db2d-628b-c5f1625cc9e3"
    },
    "status": {
      "href": "/campaigns/c652c6fa-07d7-db2d-628b-c5f1625cc9e3/status"
    },
    "template": {
      "href": "/templates/e673fb66-2d7a-7401-52fe-30cac8ee2dd3"
    }
  },
  "campaign_type_id": "c3e443e5-64d5-5d77-349f-e2c2da51eed6",
  "html": "",
  "id": "c652c6fa-07d7-db2d-628b-c5f1625cc9e3",
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
      "8d72af9b-bfa5-4d87-7953-64739c584169"
    ]
  },
  "subject": "campaign subject",
  "template_id": "e673fb66-2d7a-7401-52fe-30cac8ee2dd3",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiODI4NTRhMzgtZGM4YS03NThhLWE2Y2YtZjdjMjY0MjBkNThkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.57BvXEasfG8DOGC_5PQtIkbpgY8SvHlTJILs6Z1Iw5k
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:34 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/c3e443e5-64d5-5d77-349f-e2c2da51eed6"
    },
    "self": {
      "href": "/campaigns/c652c6fa-07d7-db2d-628b-c5f1625cc9e3"
    },
    "status": {
      "href": "/campaigns/c652c6fa-07d7-db2d-628b-c5f1625cc9e3/status"
    },
    "template": {
      "href": "/templates/e673fb66-2d7a-7401-52fe-30cac8ee2dd3"
    }
  },
  "campaign_type_id": "c3e443e5-64d5-5d77-349f-e2c2da51eed6",
  "html": "",
  "id": "c652c6fa-07d7-db2d-628b-c5f1625cc9e3",
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
      "8d72af9b-bfa5-4d87-7953-64739c584169"
    ]
  },
  "subject": "campaign subject",
  "template_id": "e673fb66-2d7a-7401-52fe-30cac8ee2dd3",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGQ1OGQwOTItYmM5ZS1iNGVjLWZiOTItMjBmZWM5MDM4ZTZlIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.tGeDkTz9ftD8XZc1g_GtppiwakWYsB4efp9yamMJUro
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 393
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:57:34 GMT
```
##### Body
```
{
  "_links": {
    "campaign": {
      "href": "/campaigns/5a5dbaa5-3ed9-62ff-b365-ee1448d9ed8b"
    },
    "self": {
      "href": "/campaigns/5a5dbaa5-3ed9-62ff-b365-ee1448d9ed8b/status"
    }
  },
  "completed_time": "2015-10-13T15:57:34Z",
  "failed_messages": 0,
  "id": "5a5dbaa5-3ed9-62ff-b365-ee1448d9ed8b",
  "queued_messages": 0,
  "retry_messages": 0,
  "sent_messages": 1,
  "start_time": "2015-10-13T15:57:34Z",
  "status": "completed",
  "total_messages": 1
}
```


## Unsubscribing
EVAN WILL EDIT THIS
<a name="unsubscriber-put"></a>
### Unsubscribe a user
#### Request **PUT** /senders/{id}/campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.admin
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiOTE1ODVjYzMtZmUwNy02MmYwLTBmNDItMmY1ODE2YTVjZGExIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.i01gnFatXQgosapruUaRyhww748m9MgRAj4wwK9tPoU
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 15:57:38 GMT
```

<a name="unsubscriber-delete"></a>
### Resubscribe a user
#### Request **DELETE** /senders/{id}/campaign_types/{id}/unsubscribers/{id}
##### Required Scopes
```
notification_preferences.admin
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiOTE1ODVjYzMtZmUwNy02MmYwLTBmNDItMmY1ODE2YTVjZGExIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MTMwMi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.i01gnFatXQgosapruUaRyhww748m9MgRAj4wwK9tPoU
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 15:57:38 GMT
```


