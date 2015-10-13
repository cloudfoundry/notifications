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
Date: Tue, 13 Oct 2015 15:44:55 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlZmM2NGItOTMyYi0yYTRmLWNhODYtMjI1Yzg3NTlmNDlmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.vPJO81GyRfaD15h3I6PgQFxF9wyjVXPR1zfgBs5wnvM
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
Date: Tue, 13 Oct 2015 15:44:58 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535/campaign_types"
    },
    "campaigns": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535/campaigns"
    },
    "self": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535"
    }
  },
  "id": "931e035d-dc1f-62f1-b209-e89983ea9535",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlZmM2NGItOTMyYi0yYTRmLWNhODYtMjI1Yzg3NTlmNDlmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.vPJO81GyRfaD15h3I6PgQFxF9wyjVXPR1zfgBs5wnvM
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:44:58 GMT
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
          "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535/campaign_types"
        },
        "campaigns": {
          "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535/campaigns"
        },
        "self": {
          "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535"
        }
      },
      "id": "931e035d-dc1f-62f1-b209-e89983ea9535",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlZmM2NGItOTMyYi0yYTRmLWNhODYtMjI1Yzg3NTlmNDlmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.vPJO81GyRfaD15h3I6PgQFxF9wyjVXPR1zfgBs5wnvM
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:44:58 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535/campaign_types"
    },
    "campaigns": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535/campaigns"
    },
    "self": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535"
    }
  },
  "id": "931e035d-dc1f-62f1-b209-e89983ea9535",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGVlZmM2NGItOTMyYi0yYTRmLWNhODYtMjI1Yzg3NTlmNDlmIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.vPJO81GyRfaD15h3I6PgQFxF9wyjVXPR1zfgBs5wnvM
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
Date: Tue, 13 Oct 2015 15:44:58 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535/campaign_types"
    },
    "campaigns": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535/campaigns"
    },
    "self": {
      "href": "/senders/931e035d-dc1f-62f1-b209-e89983ea9535"
    }
  },
  "id": "931e035d-dc1f-62f1-b209-e89983ea9535",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDNmZWZkYTAtMGQzMi0zMWUxLTcwM2EtZjcwNGM1MjE5ZWU1IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.xB7p4d3Qxi1Hht5Et7e6nYoh1UaLfFbOYP7vqAUwBqQ
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 15:44:56 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTAzNWE1NWEtNTAxZi03Mzk5LTRjNTItMDllM2RlZjgwYzJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dSQn0Qtkt0LXn9dBeDK3H0P_X7GOZOMPQSXDlvvPNHU
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
Date: Tue, 13 Oct 2015 15:44:59 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/9e665fbf-0f01-d518-9b8f-06f9a376b835"
    }
  },
  "html": "template html",
  "id": "9e665fbf-0f01-d518-9b8f-06f9a376b835",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTAzNWE1NWEtNTAxZi03Mzk5LTRjNTItMDllM2RlZjgwYzJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dSQn0Qtkt0LXn9dBeDK3H0P_X7GOZOMPQSXDlvvPNHU
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:44:59 GMT
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
          "href": "/templates/9e665fbf-0f01-d518-9b8f-06f9a376b835"
        }
      },
      "html": "html",
      "id": "9e665fbf-0f01-d518-9b8f-06f9a376b835",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTAzNWE1NWEtNTAxZi03Mzk5LTRjNTItMDllM2RlZjgwYzJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dSQn0Qtkt0LXn9dBeDK3H0P_X7GOZOMPQSXDlvvPNHU
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:44:59 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/9e665fbf-0f01-d518-9b8f-06f9a376b835"
    }
  },
  "html": "template html",
  "id": "9e665fbf-0f01-d518-9b8f-06f9a376b835",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTAzNWE1NWEtNTAxZi03Mzk5LTRjNTItMDllM2RlZjgwYzJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dSQn0Qtkt0LXn9dBeDK3H0P_X7GOZOMPQSXDlvvPNHU
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
Date: Tue, 13 Oct 2015 15:44:59 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/9e665fbf-0f01-d518-9b8f-06f9a376b835"
    }
  },
  "html": "html",
  "id": "9e665fbf-0f01-d518-9b8f-06f9a376b835",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTAzNWE1NWEtNTAxZi03Mzk5LTRjNTItMDllM2RlZjgwYzJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.dSQn0Qtkt0LXn9dBeDK3H0P_X7GOZOMPQSXDlvvPNHU
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 15:44:59 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTVlMzAyZjgtZDRlMi02ZDg5LTExODItMWYzZGVhODFmM2FkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.c53P8dRPeF-KKRxvzhBvwpU2KiZMA3U10LQNjnLRS30
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
Date: Tue, 13 Oct 2015 15:44:56 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/edd19fdd-8afa-528f-e8ad-9c2c83080944"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "edd19fdd-8afa-528f-e8ad-9c2c83080944",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTVlMzAyZjgtZDRlMi02ZDg5LTExODItMWYzZGVhODFmM2FkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.c53P8dRPeF-KKRxvzhBvwpU2KiZMA3U10LQNjnLRS30
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:44:56 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/f60feedb-401c-0787-966b-7747105dbf93/campaign_types"
    },
    "sender": {
      "href": "/senders/f60feedb-401c-0787-966b-7747105dbf93"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/edd19fdd-8afa-528f-e8ad-9c2c83080944"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "edd19fdd-8afa-528f-e8ad-9c2c83080944",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTVlMzAyZjgtZDRlMi02ZDg5LTExODItMWYzZGVhODFmM2FkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.c53P8dRPeF-KKRxvzhBvwpU2KiZMA3U10LQNjnLRS30
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:44:56 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/edd19fdd-8afa-528f-e8ad-9c2c83080944"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "edd19fdd-8afa-528f-e8ad-9c2c83080944",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTVlMzAyZjgtZDRlMi02ZDg5LTExODItMWYzZGVhODFmM2FkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.c53P8dRPeF-KKRxvzhBvwpU2KiZMA3U10LQNjnLRS30
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "2eac326e-9d65-382a-a1f7-afa7a585e228"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:44:56 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/edd19fdd-8afa-528f-e8ad-9c2c83080944"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "edd19fdd-8afa-528f-e8ad-9c2c83080944",
  "name": "updated-campaign-type",
  "template_id": "2eac326e-9d65-382a-a1f7-afa7a585e228"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTVlMzAyZjgtZDRlMi02ZDg5LTExODItMWYzZGVhODFmM2FkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.c53P8dRPeF-KKRxvzhBvwpU2KiZMA3U10LQNjnLRS30
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 13 Oct 2015 15:44:56 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGM0ZDU2NjgtOTRkYi1kMWFmLTQ1YmUtMWZhOTVjOTkxOGYwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.3haLI0XSqXdAtgmK7aCW_a5J8c59qKNjHwY-Km0o10Y
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "9ff6f2f0-a926-9c97-ba77-1426820a2955",
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
      "3a99fee3-a53e-4fdb-7be1-e1cc6cab18c3"
    ]
  },
  "subject": "campaign subject",
  "template_id": "c5a3f30d-0f6e-60a8-4e03-7bda6a17a4b3",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:45:00 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/9ff6f2f0-a926-9c97-ba77-1426820a2955"
    },
    "self": {
      "href": "/campaigns/04a388ec-dfc6-4f8d-61a0-5c7d70605387"
    },
    "status": {
      "href": "/campaigns/04a388ec-dfc6-4f8d-61a0-5c7d70605387/status"
    },
    "template": {
      "href": "/templates/c5a3f30d-0f6e-60a8-4e03-7bda6a17a4b3"
    }
  },
  "campaign_type_id": "9ff6f2f0-a926-9c97-ba77-1426820a2955",
  "html": "",
  "id": "04a388ec-dfc6-4f8d-61a0-5c7d70605387",
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
      "3a99fee3-a53e-4fdb-7be1-e1cc6cab18c3"
    ]
  },
  "subject": "campaign subject",
  "template_id": "c5a3f30d-0f6e-60a8-4e03-7bda6a17a4b3",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNGM0ZDU2NjgtOTRkYi1kMWFmLTQ1YmUtMWZhOTVjOTkxOGYwIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.3haLI0XSqXdAtgmK7aCW_a5J8c59qKNjHwY-Km0o10Y
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:45:00 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/9ff6f2f0-a926-9c97-ba77-1426820a2955"
    },
    "self": {
      "href": "/campaigns/04a388ec-dfc6-4f8d-61a0-5c7d70605387"
    },
    "status": {
      "href": "/campaigns/04a388ec-dfc6-4f8d-61a0-5c7d70605387/status"
    },
    "template": {
      "href": "/templates/c5a3f30d-0f6e-60a8-4e03-7bda6a17a4b3"
    }
  },
  "campaign_type_id": "9ff6f2f0-a926-9c97-ba77-1426820a2955",
  "html": "",
  "id": "04a388ec-dfc6-4f8d-61a0-5c7d70605387",
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
      "3a99fee3-a53e-4fdb-7be1-e1cc6cab18c3"
    ]
  },
  "subject": "campaign subject",
  "template_id": "c5a3f30d-0f6e-60a8-4e03-7bda6a17a4b3",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNTM3ZjBmY2UtOWYyNi0wZmFjLWFiMjEtNTg2NTVmNDEzY2E4IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo1MDA0MS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.4BC7ZasO3pvNyi5KVIW0GG6dw3kHiyfjfj5po6nc8eA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 393
Content-Type: text/plain; charset=utf-8
Date: Tue, 13 Oct 2015 15:45:00 GMT
```
##### Body
```
{
  "_links": {
    "campaign": {
      "href": "/campaigns/b08ce097-d331-156c-eea9-0f70b3bfc75e"
    },
    "self": {
      "href": "/campaigns/b08ce097-d331-156c-eea9-0f70b3bfc75e/status"
    }
  },
  "completed_time": "2015-10-13T15:45:00Z",
  "failed_messages": 0,
  "id": "b08ce097-d331-156c-eea9-0f70b3bfc75e",
  "queued_messages": 0,
  "retry_messages": 0,
  "sent_messages": 1,
  "start_time": "2015-10-13T15:45:00Z",
  "status": "completed",
  "total_messages": 1
}
```


