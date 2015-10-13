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
Date: Tue, 13 Oct 2015 00:01:50 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYmE5NzQ2ODUtMzU5Yy04YmRlLTRiN2QtOTQwNmE4MmQ5OWNiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WzShIRO3PnESXKVhoGBd_gaKzVZPncI2Y6KCl9e5I4s
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
Date: Tue, 13 Oct 2015 00:01:49 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1/campaign_types"
    },
    "campaigns": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1/campaigns"
    },
    "self": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1"
    }
  },
  "id": "c975cf89-024c-ee1c-ace0-5ce3ee14a2f1",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYmE5NzQ2ODUtMzU5Yy04YmRlLTRiN2QtOTQwNmE4MmQ5OWNiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WzShIRO3PnESXKVhoGBd_gaKzVZPncI2Y6KCl9e5I4s
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:49 GMT
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
          "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1/campaign_types"
        },
        "campaigns": {
          "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1/campaigns"
        },
        "self": {
          "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1"
        }
      },
      "id": "c975cf89-024c-ee1c-ace0-5ce3ee14a2f1",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYmE5NzQ2ODUtMzU5Yy04YmRlLTRiN2QtOTQwNmE4MmQ5OWNiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WzShIRO3PnESXKVhoGBd_gaKzVZPncI2Y6KCl9e5I4s
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:49 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1/campaign_types"
    },
    "campaigns": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1/campaigns"
    },
    "self": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1"
    }
  },
  "id": "c975cf89-024c-ee1c-ace0-5ce3ee14a2f1",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYmE5NzQ2ODUtMzU5Yy04YmRlLTRiN2QtOTQwNmE4MmQ5OWNiIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.WzShIRO3PnESXKVhoGBd_gaKzVZPncI2Y6KCl9e5I4s
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
Date: Tue, 13 Oct 2015 00:01:49 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1/campaign_types"
    },
    "campaigns": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1/campaigns"
    },
    "self": {
      "href": "/senders/c975cf89-024c-ee1c-ace0-5ce3ee14a2f1"
    }
  },
  "id": "c975cf89-024c-ee1c-ace0-5ce3ee14a2f1",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTNlMjBlMzQtYmI0Yy03YjM4LWRkOGMtZDdiYTEwM2YyYmRkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.GT9QRcFM4Svg_ndHGoRXQjBpVdU616Rrol_vk4WenQU
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 00:01:51 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTVhM2FmZjItNzVjZS1lMjgxLTJiOGYtYTk3ZDM3NDZkOTNkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.4pRwz7x_y7Y4ZcHaK4Y9HssvD7fdN5V4TD7NgcbDig8
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
Date: Tue, 13 Oct 2015 00:01:46 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/ba150bde-0f4a-305e-479c-4776f9fc59ad"
    }
  },
  "html": "template html",
  "id": "ba150bde-0f4a-305e-479c-4776f9fc59ad",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTVhM2FmZjItNzVjZS1lMjgxLTJiOGYtYTk3ZDM3NDZkOTNkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.4pRwz7x_y7Y4ZcHaK4Y9HssvD7fdN5V4TD7NgcbDig8
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:46 GMT
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
          "href": "/templates/ba150bde-0f4a-305e-479c-4776f9fc59ad"
        }
      },
      "html": "html",
      "id": "ba150bde-0f4a-305e-479c-4776f9fc59ad",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTVhM2FmZjItNzVjZS1lMjgxLTJiOGYtYTk3ZDM3NDZkOTNkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.4pRwz7x_y7Y4ZcHaK4Y9HssvD7fdN5V4TD7NgcbDig8
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:46 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/ba150bde-0f4a-305e-479c-4776f9fc59ad"
    }
  },
  "html": "template html",
  "id": "ba150bde-0f4a-305e-479c-4776f9fc59ad",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTVhM2FmZjItNzVjZS1lMjgxLTJiOGYtYTk3ZDM3NDZkOTNkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.4pRwz7x_y7Y4ZcHaK4Y9HssvD7fdN5V4TD7NgcbDig8
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
Date: Tue, 13 Oct 2015 00:01:46 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/ba150bde-0f4a-305e-479c-4776f9fc59ad"
    }
  },
  "html": "html",
  "id": "ba150bde-0f4a-305e-479c-4776f9fc59ad",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTVhM2FmZjItNzVjZS1lMjgxLTJiOGYtYTk3ZDM3NDZkOTNkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.4pRwz7x_y7Y4ZcHaK4Y9HssvD7fdN5V4TD7NgcbDig8
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 00:01:46 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMDhmMzU3MWItMDkwNC1mMjA1LWU0MTEtNTViYTNkNjJhY2IwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.47SRWIZNFl4ECSJmZq47ypPBGuQR4Az8MaXRV52qNZg
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
Date: Tue, 13 Oct 2015 00:01:48 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/3e46f93b-1dce-da30-64c4-9674ca43e6c1"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "3e46f93b-1dce-da30-64c4-9674ca43e6c1",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMDhmMzU3MWItMDkwNC1mMjA1LWU0MTEtNTViYTNkNjJhY2IwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.47SRWIZNFl4ECSJmZq47ypPBGuQR4Az8MaXRV52qNZg
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:48 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/0ba1da08-fff5-db10-626d-84499d52ce5d/campaign_types"
    },
    "sender": {
      "href": "/senders/0ba1da08-fff5-db10-626d-84499d52ce5d"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/3e46f93b-1dce-da30-64c4-9674ca43e6c1"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "3e46f93b-1dce-da30-64c4-9674ca43e6c1",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMDhmMzU3MWItMDkwNC1mMjA1LWU0MTEtNTViYTNkNjJhY2IwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.47SRWIZNFl4ECSJmZq47ypPBGuQR4Az8MaXRV52qNZg
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:48 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/3e46f93b-1dce-da30-64c4-9674ca43e6c1"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "3e46f93b-1dce-da30-64c4-9674ca43e6c1",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMDhmMzU3MWItMDkwNC1mMjA1LWU0MTEtNTViYTNkNjJhY2IwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.47SRWIZNFl4ECSJmZq47ypPBGuQR4Az8MaXRV52qNZg
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "3c077b8a-b15b-a9ee-1a20-7ee7defc7bff"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:48 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/3e46f93b-1dce-da30-64c4-9674ca43e6c1"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "3e46f93b-1dce-da30-64c4-9674ca43e6c1",
  "name": "updated-campaign-type",
  "template_id": "3c077b8a-b15b-a9ee-1a20-7ee7defc7bff"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiMDhmMzU3MWItMDkwNC1mMjA1LWU0MTEtNTViYTNkNjJhY2IwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.47SRWIZNFl4ECSJmZq47ypPBGuQR4Az8MaXRV52qNZg
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 00:01:48 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTViMDZlZjUtZTA0MS0wYzljLTVkMmMtZmZkYzFhMGY1YmQ3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.FDd_5VIY-JIugD6J7TDuydmBdoZxbFW2b7AJ_dfu7dg
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "cf8590a3-f2b8-6c6f-677e-f928ce6ed15d",
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
      "05ee5262-01d1-4d18-5466-75eadfcac3d0"
    ]
  },
  "subject": "campaign subject",
  "template_id": "728f7294-8538-10d2-888b-af956ec8595b",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:49 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/cf8590a3-f2b8-6c6f-677e-f928ce6ed15d"
    },
    "self": {
      "href": "/campaigns/f19a1be5-7330-18e2-2e7d-5f2137ea5db9"
    },
    "status": {
      "href": "/campaigns/f19a1be5-7330-18e2-2e7d-5f2137ea5db9/status"
    },
    "template": {
      "href": "/templates/728f7294-8538-10d2-888b-af956ec8595b"
    }
  },
  "campaign_type_id": "cf8590a3-f2b8-6c6f-677e-f928ce6ed15d",
  "html": "",
  "id": "f19a1be5-7330-18e2-2e7d-5f2137ea5db9",
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
      "05ee5262-01d1-4d18-5466-75eadfcac3d0"
    ]
  },
  "subject": "campaign subject",
  "template_id": "728f7294-8538-10d2-888b-af956ec8595b",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiZTViMDZlZjUtZTA0MS0wYzljLTVkMmMtZmZkYzFhMGY1YmQ3IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MzY5Ny9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.FDd_5VIY-JIugD6J7TDuydmBdoZxbFW2b7AJ_dfu7dg
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 00:01:49 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/cf8590a3-f2b8-6c6f-677e-f928ce6ed15d"
    },
    "self": {
      "href": "/campaigns/f19a1be5-7330-18e2-2e7d-5f2137ea5db9"
    },
    "status": {
      "href": "/campaigns/f19a1be5-7330-18e2-2e7d-5f2137ea5db9/status"
    },
    "template": {
      "href": "/templates/728f7294-8538-10d2-888b-af956ec8595b"
    }
  },
  "campaign_type_id": "cf8590a3-f2b8-6c6f-677e-f928ce6ed15d",
  "html": "",
  "id": "f19a1be5-7330-18e2-2e7d-5f2137ea5db9",
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
      "05ee5262-01d1-4d18-5466-75eadfcac3d0"
    ]
  },
  "subject": "campaign subject",
  "template_id": "728f7294-8538-10d2-888b-af956ec8595b",
  "text": "campaign body"
}
```


