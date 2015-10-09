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
Date: Fri, 09 Oct 2015 22:47:44 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDBlNzYxNTEtYjQ3Mi1iZWQzLTcyNzYtZDE4Y2JjN2Q3YWJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.11erqQ2EUonBBko-VeKZ8-BvQ3F79I9vzqA-gaCfYx4
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
Date: Fri, 09 Oct 2015 22:47:43 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91/campaign_types"
    },
    "campaigns": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91/campaigns"
    },
    "self": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91"
    }
  },
  "id": "bdca2ca0-b0b4-2cfb-facc-b673b923ed91",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDBlNzYxNTEtYjQ3Mi1iZWQzLTcyNzYtZDE4Y2JjN2Q3YWJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.11erqQ2EUonBBko-VeKZ8-BvQ3F79I9vzqA-gaCfYx4
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:43 GMT
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
          "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91/campaign_types"
        },
        "campaigns": {
          "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91/campaigns"
        },
        "self": {
          "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91"
        }
      },
      "id": "bdca2ca0-b0b4-2cfb-facc-b673b923ed91",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDBlNzYxNTEtYjQ3Mi1iZWQzLTcyNzYtZDE4Y2JjN2Q3YWJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.11erqQ2EUonBBko-VeKZ8-BvQ3F79I9vzqA-gaCfYx4
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:43 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91/campaign_types"
    },
    "campaigns": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91/campaigns"
    },
    "self": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91"
    }
  },
  "id": "bdca2ca0-b0b4-2cfb-facc-b673b923ed91",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiNDBlNzYxNTEtYjQ3Mi1iZWQzLTcyNzYtZDE4Y2JjN2Q3YWJkIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.11erqQ2EUonBBko-VeKZ8-BvQ3F79I9vzqA-gaCfYx4
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
Date: Fri, 09 Oct 2015 22:47:43 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91/campaign_types"
    },
    "campaigns": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91/campaigns"
    },
    "self": {
      "href": "/senders/bdca2ca0-b0b4-2cfb-facc-b673b923ed91"
    }
  },
  "id": "bdca2ca0-b0b4-2cfb-facc-b673b923ed91",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYjEzMDQ5ZGMtYmU2Ni1mZmUxLTAwY2ItODNmMjQzYWY4ZTM0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.Zqc3x48Bmu41ZYKD-ixDcBgqJbGtcmg0P3z1pVtgXZA
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Fri, 09 Oct 2015 22:47:43 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Q2YjU1ODItNGUyOC1mNjRmLTkzMGMtOGE3OWU5ZWZhYTk5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bWQVh-tOqBSnIf-r8b9WbYWDRAsXRiCAmSwOSEkZhhc
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
Date: Fri, 09 Oct 2015 22:47:43 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/034e9118-3b18-deb4-2ade-8fa171fddd50"
    }
  },
  "html": "template html",
  "id": "034e9118-3b18-deb4-2ade-8fa171fddd50",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Q2YjU1ODItNGUyOC1mNjRmLTkzMGMtOGE3OWU5ZWZhYTk5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bWQVh-tOqBSnIf-r8b9WbYWDRAsXRiCAmSwOSEkZhhc
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:43 GMT
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
          "href": "/templates/034e9118-3b18-deb4-2ade-8fa171fddd50"
        }
      },
      "html": "html",
      "id": "034e9118-3b18-deb4-2ade-8fa171fddd50",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Q2YjU1ODItNGUyOC1mNjRmLTkzMGMtOGE3OWU5ZWZhYTk5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bWQVh-tOqBSnIf-r8b9WbYWDRAsXRiCAmSwOSEkZhhc
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:43 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/034e9118-3b18-deb4-2ade-8fa171fddd50"
    }
  },
  "html": "template html",
  "id": "034e9118-3b18-deb4-2ade-8fa171fddd50",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Q2YjU1ODItNGUyOC1mNjRmLTkzMGMtOGE3OWU5ZWZhYTk5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bWQVh-tOqBSnIf-r8b9WbYWDRAsXRiCAmSwOSEkZhhc
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
Date: Fri, 09 Oct 2015 22:47:43 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/034e9118-3b18-deb4-2ade-8fa171fddd50"
    }
  },
  "html": "html",
  "id": "034e9118-3b18-deb4-2ade-8fa171fddd50",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiN2Q2YjU1ODItNGUyOC1mNjRmLTkzMGMtOGE3OWU5ZWZhYTk5IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.bWQVh-tOqBSnIf-r8b9WbYWDRAsXRiCAmSwOSEkZhhc
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Fri, 09 Oct 2015 22:47:43 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2EzMTVjM2MtN2JmYS1mNDZhLTBiNGQtODI5NTVjYTk1MmU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2gXIzbGOCqX4PXAhmZQBSS7H5JKcy18W7y11jIbWkK8
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
Date: Fri, 09 Oct 2015 22:47:45 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/084832d9-1844-1a2a-b2b1-151b4ac68cdf"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "084832d9-1844-1a2a-b2b1-151b4ac68cdf",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2EzMTVjM2MtN2JmYS1mNDZhLTBiNGQtODI5NTVjYTk1MmU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2gXIzbGOCqX4PXAhmZQBSS7H5JKcy18W7y11jIbWkK8
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:45 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/acaa1979-4d23-46d3-aecb-5e087259b55f/campaign_types"
    },
    "sender": {
      "href": "/senders/acaa1979-4d23-46d3-aecb-5e087259b55f"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/084832d9-1844-1a2a-b2b1-151b4ac68cdf"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "084832d9-1844-1a2a-b2b1-151b4ac68cdf",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2EzMTVjM2MtN2JmYS1mNDZhLTBiNGQtODI5NTVjYTk1MmU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2gXIzbGOCqX4PXAhmZQBSS7H5JKcy18W7y11jIbWkK8
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:45 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/084832d9-1844-1a2a-b2b1-151b4ac68cdf"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "084832d9-1844-1a2a-b2b1-151b4ac68cdf",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2EzMTVjM2MtN2JmYS1mNDZhLTBiNGQtODI5NTVjYTk1MmU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2gXIzbGOCqX4PXAhmZQBSS7H5JKcy18W7y11jIbWkK8
X-Notifications-Version: 2
```
##### Body
```
{
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "1627f793-7a1e-0b88-0978-df3462603044"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 280
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:45 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/084832d9-1844-1a2a-b2b1-151b4ac68cdf"
    }
  },
  "critical": false,
  "description": "still the same great campaign type",
  "id": "084832d9-1844-1a2a-b2b1-151b4ac68cdf",
  "name": "updated-campaign-type",
  "template_id": "1627f793-7a1e-0b88-0978-df3462603044"
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiY2EzMTVjM2MtN2JmYS1mNDZhLTBiNGQtODI5NTVjYTk1MmU0IiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.2gXIzbGOCqX4PXAhmZQBSS7H5JKcy18W7y11jIbWkK8
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Fri, 09 Oct 2015 22:47:45 GMT
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTE2NzNiMDctODhhYy05YzM5LTg4MjQtMDRiZDk3NzlhMTBhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.ObPJUUwM7Nowdxy4OMbylVcB3EoXLapJLaedaIiT5IE
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "ff545cdf-0dab-e047-86db-deeb8ed9b657",
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
      "49392daf-f818-42d8-7869-0143bcef0d27"
    ]
  },
  "subject": "campaign subject",
  "template_id": "b84ec99e-9382-e19a-ef1c-648dd6a1b946",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:43 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/ff545cdf-0dab-e047-86db-deeb8ed9b657"
    },
    "self": {
      "href": "/campaigns/b5282060-6afb-bb07-66d0-61f71c47b393"
    },
    "status": {
      "href": "/campaigns/b5282060-6afb-bb07-66d0-61f71c47b393/status"
    },
    "template": {
      "href": "/templates/b84ec99e-9382-e19a-ef1c-648dd6a1b946"
    }
  },
  "campaign_type_id": "ff545cdf-0dab-e047-86db-deeb8ed9b657",
  "html": "",
  "id": "b5282060-6afb-bb07-66d0-61f71c47b393",
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
      "49392daf-f818-42d8-7869-0143bcef0d27"
    ]
  },
  "subject": "campaign subject",
  "template_id": "b84ec99e-9382-e19a-ef1c-648dd6a1b946",
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
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJzY2ltIHBhc3N3b3JkIiwiY2xpZW50X2lkIjoiYTE2NzNiMDctODhhYy05YzM5LTg4MjQtMDRiZDk3NzlhMTBhIiwiaXNzIjoiaHR0cDovLzEyNy4wLjAuMTo2MjAzOS9vYXV0aC90b2tlbiIsInNjb3BlIjpbIm5vdGlmaWNhdGlvbnMud3JpdGUiXX0.ObPJUUwM7Nowdxy4OMbylVcB3EoXLapJLaedaIiT5IE
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 687
Content-Type: text/plain; charset=utf-8
Date: Fri, 09 Oct 2015 22:47:43 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/ff545cdf-0dab-e047-86db-deeb8ed9b657"
    },
    "self": {
      "href": "/campaigns/b5282060-6afb-bb07-66d0-61f71c47b393"
    },
    "status": {
      "href": "/campaigns/b5282060-6afb-bb07-66d0-61f71c47b393/status"
    },
    "template": {
      "href": "/templates/b84ec99e-9382-e19a-ef1c-648dd6a1b946"
    }
  },
  "campaign_type_id": "ff545cdf-0dab-e047-86db-deeb8ed9b657",
  "html": "",
  "id": "b5282060-6afb-bb07-66d0-61f71c47b393",
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
      "49392daf-f818-42d8-7869-0143bcef0d27"
    ]
  },
  "subject": "campaign subject",
  "template_id": "b84ec99e-9382-e19a-ef1c-648dd6a1b946",
  "text": "campaign body"
}
```


