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
Date: Wed, 14 Oct 2015 15:53:26 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGU0OTA4YTEtYTViYy0zMGQ3LTcwNDctMTkzMTMzYzljZTM2IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._TEox4tRaWV5qtJR7N1cPRgKM8CvGBM98jURaIjnTb0
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
Date: Wed, 14 Oct 2015 15:53:25 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9/campaign_types"
    },
    "campaigns": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9/campaigns"
    },
    "self": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9"
    }
  },
  "id": "c0a0f312-e311-21de-a722-7dff44b27bc9",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGU0OTA4YTEtYTViYy0zMGQ3LTcwNDctMTkzMTMzYzljZTM2IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._TEox4tRaWV5qtJR7N1cPRgKM8CvGBM98jURaIjnTb0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:25 GMT
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
          "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9/campaign_types"
        },
        "campaigns": {
          "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9/campaigns"
        },
        "self": {
          "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9"
        }
      },
      "id": "c0a0f312-e311-21de-a722-7dff44b27bc9",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGU0OTA4YTEtYTViYy0zMGQ3LTcwNDctMTkzMTMzYzljZTM2IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._TEox4tRaWV5qtJR7N1cPRgKM8CvGBM98jURaIjnTb0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:25 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9/campaign_types"
    },
    "campaigns": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9/campaigns"
    },
    "self": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9"
    }
  },
  "id": "c0a0f312-e311-21de-a722-7dff44b27bc9",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZGU0OTA4YTEtYTViYy0zMGQ3LTcwNDctMTkzMTMzYzljZTM2IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._TEox4tRaWV5qtJR7N1cPRgKM8CvGBM98jURaIjnTb0
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
Date: Wed, 14 Oct 2015 15:53:25 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9/campaign_types"
    },
    "campaigns": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9/campaigns"
    },
    "self": {
      "href": "/senders/c0a0f312-e311-21de-a722-7dff44b27bc9"
    }
  },
  "id": "c0a0f312-e311-21de-a722-7dff44b27bc9",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMjgzMDI5OTctM2I5Yy1iNzFkLWQ0NDMtYjRmYmExMWY4NWU1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.shoFaOsSdCJL78HUoFf-85tCP_2qviP5_RepU1qh9EQ
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Wed, 14 Oct 2015 15:53:24 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjI3NGE0NDItZjNiNC1mNmUxLTg2ZGQtYmIwMDkzMGJjMWIyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._rLnl2LZHZkPghKd6OK9400BcmSCGyFeqP0ukeHKo2s
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
Date: Wed, 14 Oct 2015 15:53:25 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/e26e2e4c-29d9-05dc-1b32-7043498ca3c5"
    }
  },
  "html": "template html",
  "id": "e26e2e4c-29d9-05dc-1b32-7043498ca3c5",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjI3NGE0NDItZjNiNC1mNmUxLTg2ZGQtYmIwMDkzMGJjMWIyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._rLnl2LZHZkPghKd6OK9400BcmSCGyFeqP0ukeHKo2s
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:25 GMT
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
          "href": "/templates/e26e2e4c-29d9-05dc-1b32-7043498ca3c5"
        }
      },
      "html": "html",
      "id": "e26e2e4c-29d9-05dc-1b32-7043498ca3c5",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjI3NGE0NDItZjNiNC1mNmUxLTg2ZGQtYmIwMDkzMGJjMWIyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._rLnl2LZHZkPghKd6OK9400BcmSCGyFeqP0ukeHKo2s
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:25 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/e26e2e4c-29d9-05dc-1b32-7043498ca3c5"
    }
  },
  "html": "template html",
  "id": "e26e2e4c-29d9-05dc-1b32-7043498ca3c5",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjI3NGE0NDItZjNiNC1mNmUxLTg2ZGQtYmIwMDkzMGJjMWIyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._rLnl2LZHZkPghKd6OK9400BcmSCGyFeqP0ukeHKo2s
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
Date: Wed, 14 Oct 2015 15:53:25 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/e26e2e4c-29d9-05dc-1b32-7043498ca3c5"
    }
  },
  "html": "html",
  "id": "e26e2e4c-29d9-05dc-1b32-7043498ca3c5",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZjI3NGE0NDItZjNiNC1mNmUxLTg2ZGQtYmIwMDkzMGJjMWIyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0._rLnl2LZHZkPghKd6OK9400BcmSCGyFeqP0ukeHKo2s
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Wed, 14 Oct 2015 15:53:25 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzIwZjY2ODEtYjIwYi03MTc0LTg4OGQtMzk0YjQ0MDNmYWU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2xLX_7roHnQvNT76FthAXVqO8F8THgQFDuA1e9tA3Zg
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
Date: Wed, 14 Oct 2015 15:53:24 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/a27fdf5a-cb6f-cdeb-0f38-2158f17b32c8"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "a27fdf5a-cb6f-cdeb-0f38-2158f17b32c8",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzIwZjY2ODEtYjIwYi03MTc0LTg4OGQtMzk0YjQ0MDNmYWU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2xLX_7roHnQvNT76FthAXVqO8F8THgQFDuA1e9tA3Zg
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:24 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/4a0bfb75-60c2-982e-c37b-38d6cac0969f/campaign_types"
    },
    "sender": {
      "href": "/senders/4a0bfb75-60c2-982e-c37b-38d6cac0969f"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/a27fdf5a-cb6f-cdeb-0f38-2158f17b32c8"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "a27fdf5a-cb6f-cdeb-0f38-2158f17b32c8",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzIwZjY2ODEtYjIwYi03MTc0LTg4OGQtMzk0YjQ0MDNmYWU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2xLX_7roHnQvNT76FthAXVqO8F8THgQFDuA1e9tA3Zg
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:24 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/a27fdf5a-cb6f-cdeb-0f38-2158f17b32c8"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "a27fdf5a-cb6f-cdeb-0f38-2158f17b32c8",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzIwZjY2ODEtYjIwYi03MTc0LTg4OGQtMzk0YjQ0MDNmYWU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2xLX_7roHnQvNT76FthAXVqO8F8THgQFDuA1e9tA3Zg
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "04a11f3e-3dc6-6ff8-bbae-204661057f11"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:24 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/a27fdf5a-cb6f-cdeb-0f38-2158f17b32c8"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "a27fdf5a-cb6f-cdeb-0f38-2158f17b32c8",
  "name": "updated-campaign-type",
  "template_id": "04a11f3e-3dc6-6ff8-bbae-204661057f11"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMzIwZjY2ODEtYjIwYi03MTc0LTg4OGQtMzk0YjQ0MDNmYWU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2xLX_7roHnQvNT76FthAXVqO8F8THgQFDuA1e9tA3Zg
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Wed, 14 Oct 2015 15:53:24 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZWI1OWE4MGUtZDFlOS0xMzcyLTI5ZTItZGIxNmExZTljNzc0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.ELE2woEYqZlx1ulcgL0lGNvPMGoqNRlEz4BeWYy3kZ0
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "44344860-c779-db79-0325-a9cc46352f5e",
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
      "07bfae6e-cc2b-4811-5e84-d38aa92e476e"
    ]
  },
  "subject": "campaign subject",
  "template_id": "25089281-1338-5431-6527-9c1e27aa70a7",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:25 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/44344860-c779-db79-0325-a9cc46352f5e"
    },
    "self": {
      "href": "/campaigns/01ff18fe-34d2-1343-ed62-cd757f374a7f"
    },
    "status": {
      "href": "/campaigns/01ff18fe-34d2-1343-ed62-cd757f374a7f/status"
    },
    "template": {
      "href": "/templates/25089281-1338-5431-6527-9c1e27aa70a7"
    }
  },
  "campaign_type_id": "44344860-c779-db79-0325-a9cc46352f5e",
  "html": "",
  "id": "01ff18fe-34d2-1343-ed62-cd757f374a7f",
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
      "07bfae6e-cc2b-4811-5e84-d38aa92e476e"
    ]
  },
  "subject": "campaign subject",
  "template_id": "25089281-1338-5431-6527-9c1e27aa70a7",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZWI1OWE4MGUtZDFlOS0xMzcyLTI5ZTItZGIxNmExZTljNzc0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.ELE2woEYqZlx1ulcgL0lGNvPMGoqNRlEz4BeWYy3kZ0
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:25 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/44344860-c779-db79-0325-a9cc46352f5e"
    },
    "self": {
      "href": "/campaigns/01ff18fe-34d2-1343-ed62-cd757f374a7f"
    },
    "status": {
      "href": "/campaigns/01ff18fe-34d2-1343-ed62-cd757f374a7f/status"
    },
    "template": {
      "href": "/templates/25089281-1338-5431-6527-9c1e27aa70a7"
    }
  },
  "campaign_type_id": "44344860-c779-db79-0325-a9cc46352f5e",
  "html": "",
  "id": "01ff18fe-34d2-1343-ed62-cd757f374a7f",
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
      "07bfae6e-cc2b-4811-5e84-d38aa92e476e"
    ]
  },
  "subject": "campaign subject",
  "template_id": "25089281-1338-5431-6527-9c1e27aa70a7",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGU5ZTg2ODQtNDc4NS04NTNlLWVmYTEtOWEzYWYyNGRlMmE1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.0G-fNl8QAhCMQdaiviLzf0HOdpVJvXyZChtFQLaMYoM
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 393
Content-Type: text/plain; charset=utf-8
Date: Wed, 14 Oct 2015 15:53:24 GMT
```
##### Body
```
{
  "_links": {
    "campaign": {
      "href": "/campaigns/376acd79-b909-58e7-22c5-00de68c983d7"
    },
    "self": {
      "href": "/campaigns/376acd79-b909-58e7-22c5-00de68c983d7/status"
    }
  },
  "completed_time": "2015-10-14T15:53:24Z",
  "failed_messages": 0,
  "id": "376acd79-b909-58e7-22c5-00de68c983d7",
  "queued_messages": 0,
  "retry_messages": 0,
  "sent_messages": 1,
  "start_time": "2015-10-14T15:53:24Z",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNmEyODUzMmUtNGNmYi0yNTU0LWFmZjEtYTk5MTg0NWRhMzYyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.gEln6doAkd7waLYvOxXo4WKmj-sApuj6kTRTfO2ygQ4
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Wed, 14 Oct 2015 15:53:22 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiYTU4NDE0NjYtNDdkZC00NDI1LTZjZGEtY2NjMmQyNjJiNTNiIn0.tsm5vtu7qfUD6O1P59y3O1R0HzU1hm9kxUW_cpBH968
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Wed, 14 Oct 2015 15:53:24 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNmEyODUzMmUtNGNmYi0yNTU0LWFmZjEtYTk5MTg0NWRhMzYyIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzI2Mi9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy5hZG1pbiJdfQ.gEln6doAkd7waLYvOxXo4WKmj-sApuj6kTRTfO2ygQ4
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Wed, 14 Oct 2015 15:53:22 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIiLCJpc3MiOiIiLCJzY29wZSI6WyJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMucmVhZCIsIm5vdGlmaWNhdGlvbl9wcmVmZXJlbmNlcy53cml0ZSJdLCJ1c2VyX2lkIjoiYTU4NDE0NjYtNDdkZC00NDI1LTZjZGEtY2NjMmQyNjJiNTNiIn0.tsm5vtu7qfUD6O1P59y3O1R0HzU1hm9kxUW_cpBH968
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Wed, 14 Oct 2015 15:53:24 GMT
```


