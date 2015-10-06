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
Date: Tue, 06 Oct 2015 21:22:19 GMT
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
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3NDAsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.kSqsfCXuOLZHA8y4N3dif5mGfz_rfJJOiTw-Hk5mNsLm5X59P55UvyZ4X0fXQYujwc43d6y3sw7yyQVD7oiAAKj2i5zW1ZGudoW9REyM_E3J0nbVjI5IXL7UBow7nOKGyNZ6UnzyzDiu3w0oUHZucCSGshKMrFhbTb71iFiojFaid1UZVZ_5zRMEArB5dzzFDUglt4sNsLhXam5L34W-B-dar5YNkxQFprVRFkrGxqTS-jjDdGFV7MFCfsKuGdfKbL3l9N4JkU9bBKe1dnuA8ToW6BPD-ZBQDwStysphR6wmuXTaZTa55S6Cr3RBCOwjiM2VuRbfoGjDLqcwU1AhqA
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
Date: Tue, 06 Oct 2015 21:22:20 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f/campaign_types"
    },
    "campaigns": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f/campaigns"
    },
    "self": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f"
    }
  },
  "id": "a5ebf2c2-6b1a-ad6d-8df0-b0265947892f",
  "name": "My Cool App"
}
```

<a name="sender-list"></a>
### List all senders
#### Request **GET** /senders
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3NDAsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.kSqsfCXuOLZHA8y4N3dif5mGfz_rfJJOiTw-Hk5mNsLm5X59P55UvyZ4X0fXQYujwc43d6y3sw7yyQVD7oiAAKj2i5zW1ZGudoW9REyM_E3J0nbVjI5IXL7UBow7nOKGyNZ6UnzyzDiu3w0oUHZucCSGshKMrFhbTb71iFiojFaid1UZVZ_5zRMEArB5dzzFDUglt4sNsLhXam5L34W-B-dar5YNkxQFprVRFkrGxqTS-jjDdGFV7MFCfsKuGdfKbL3l9N4JkU9bBKe1dnuA8ToW6BPD-ZBQDwStysphR6wmuXTaZTa55S6Cr3RBCOwjiM2VuRbfoGjDLqcwU1AhqA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 362
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:20 GMT
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
          "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f/campaign_types"
        },
        "campaigns": {
          "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f/campaigns"
        },
        "self": {
          "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f"
        }
      },
      "id": "a5ebf2c2-6b1a-ad6d-8df0-b0265947892f",
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
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3NDAsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.kSqsfCXuOLZHA8y4N3dif5mGfz_rfJJOiTw-Hk5mNsLm5X59P55UvyZ4X0fXQYujwc43d6y3sw7yyQVD7oiAAKj2i5zW1ZGudoW9REyM_E3J0nbVjI5IXL7UBow7nOKGyNZ6UnzyzDiu3w0oUHZucCSGshKMrFhbTb71iFiojFaid1UZVZ_5zRMEArB5dzzFDUglt4sNsLhXam5L34W-B-dar5YNkxQFprVRFkrGxqTS-jjDdGFV7MFCfsKuGdfKbL3l9N4JkU9bBKe1dnuA8ToW6BPD-ZBQDwStysphR6wmuXTaZTa55S6Cr3RBCOwjiM2VuRbfoGjDLqcwU1AhqA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 310
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:20 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f/campaign_types"
    },
    "campaigns": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f/campaigns"
    },
    "self": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f"
    }
  },
  "id": "a5ebf2c2-6b1a-ad6d-8df0-b0265947892f",
  "name": "My Cool App"
}
```

<a name="sender-update"></a>
### Update a sender
#### Request **PUT** /senders/{id}
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3NDAsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.kSqsfCXuOLZHA8y4N3dif5mGfz_rfJJOiTw-Hk5mNsLm5X59P55UvyZ4X0fXQYujwc43d6y3sw7yyQVD7oiAAKj2i5zW1ZGudoW9REyM_E3J0nbVjI5IXL7UBow7nOKGyNZ6UnzyzDiu3w0oUHZucCSGshKMrFhbTb71iFiojFaid1UZVZ_5zRMEArB5dzzFDUglt4sNsLhXam5L34W-B-dar5YNkxQFprVRFkrGxqTS-jjDdGFV7MFCfsKuGdfKbL3l9N4JkU9bBKe1dnuA8ToW6BPD-ZBQDwStysphR6wmuXTaZTa55S6Cr3RBCOwjiM2VuRbfoGjDLqcwU1AhqA
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
Date: Tue, 06 Oct 2015 21:22:20 GMT
```
##### Body
```
{
  "_links": {
    "campaign_types": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f/campaign_types"
    },
    "campaigns": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f/campaigns"
    },
    "self": {
      "href": "/senders/a5ebf2c2-6b1a-ad6d-8df0-b0265947892f"
    }
  },
  "id": "a5ebf2c2-6b1a-ad6d-8df0-b0265947892f",
  "name": "My Not Cool App"
}
```

<a name="sender-delete"></a>
### Delete a sender
#### Request **DELETE** /senders/{id}
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3NDAsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.kSqsfCXuOLZHA8y4N3dif5mGfz_rfJJOiTw-Hk5mNsLm5X59P55UvyZ4X0fXQYujwc43d6y3sw7yyQVD7oiAAKj2i5zW1ZGudoW9REyM_E3J0nbVjI5IXL7UBow7nOKGyNZ6UnzyzDiu3w0oUHZucCSGshKMrFhbTb71iFiojFaid1UZVZ_5zRMEArB5dzzFDUglt4sNsLhXam5L34W-B-dar5YNkxQFprVRFkrGxqTS-jjDdGFV7MFCfsKuGdfKbL3l9N4JkU9bBKe1dnuA8ToW6BPD-ZBQDwStysphR6wmuXTaZTa55S6Cr3RBCOwjiM2VuRbfoGjDLqcwU1AhqA
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 06 Oct 2015 21:22:20 GMT
```


## Templates
EVAN WILL EDIT THIS
<a name="template-create"></a>
### Create a new template
#### Request **POST** /templates
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzgsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.mM4oj3NYciRFkmYUfT9HrC0Pl19zGFQXTPpS82PZ8oZhmSkCtlg_BQM0xJl37nQslX9jGaphr_eYGhk08ENvOL6HOeglkBDSr4CC3c4KtuMCYOWaK9NQesAaozmjpug1_QxkngRRoka62UeTRQkKlBPsTghS6GN7F6uELdzKuhIZYgUSF0DVX2t8sPMr0-pomLxFDFCT94Dwc5HOmjngzlsljNoPoJFIFc1Nyb8SX4zw7GDe6LcVSQMklThhx-0VUhQKPHLkXP61hFKU7SOl_ZprNvN1NfJmpT8UPEa4kQMo-P9wCTNGskSGweQzwTlmyNkTphfHg4ldeuHRB4s_7Q
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
Date: Tue, 06 Oct 2015 21:22:18 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/4b88c433-28fd-deaf-da92-5a1cf6abc09a"
    }
  },
  "html": "template html",
  "id": "4b88c433-28fd-deaf-da92-5a1cf6abc09a",
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
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzgsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.mM4oj3NYciRFkmYUfT9HrC0Pl19zGFQXTPpS82PZ8oZhmSkCtlg_BQM0xJl37nQslX9jGaphr_eYGhk08ENvOL6HOeglkBDSr4CC3c4KtuMCYOWaK9NQesAaozmjpug1_QxkngRRoka62UeTRQkKlBPsTghS6GN7F6uELdzKuhIZYgUSF0DVX2t8sPMr0-pomLxFDFCT94Dwc5HOmjngzlsljNoPoJFIFc1Nyb8SX4zw7GDe6LcVSQMklThhx-0VUhQKPHLkXP61hFKU7SOl_ZprNvN1NfJmpT8UPEa4kQMo-P9wCTNGskSGweQzwTlmyNkTphfHg4ldeuHRB4s_7Q
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 298
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:18 GMT
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
          "href": "/templates/4b88c433-28fd-deaf-da92-5a1cf6abc09a"
        }
      },
      "html": "html",
      "id": "4b88c433-28fd-deaf-da92-5a1cf6abc09a",
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
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzgsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.mM4oj3NYciRFkmYUfT9HrC0Pl19zGFQXTPpS82PZ8oZhmSkCtlg_BQM0xJl37nQslX9jGaphr_eYGhk08ENvOL6HOeglkBDSr4CC3c4KtuMCYOWaK9NQesAaozmjpug1_QxkngRRoka62UeTRQkKlBPsTghS6GN7F6uELdzKuhIZYgUSF0DVX2t8sPMr0-pomLxFDFCT94Dwc5HOmjngzlsljNoPoJFIFc1Nyb8SX4zw7GDe6LcVSQMklThhx-0VUhQKPHLkXP61hFKU7SOl_ZprNvN1NfJmpT8UPEa4kQMo-P9wCTNGskSGweQzwTlmyNkTphfHg4ldeuHRB4s_7Q
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 266
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:18 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/4b88c433-28fd-deaf-da92-5a1cf6abc09a"
    }
  },
  "html": "template html",
  "id": "4b88c433-28fd-deaf-da92-5a1cf6abc09a",
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
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzgsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.mM4oj3NYciRFkmYUfT9HrC0Pl19zGFQXTPpS82PZ8oZhmSkCtlg_BQM0xJl37nQslX9jGaphr_eYGhk08ENvOL6HOeglkBDSr4CC3c4KtuMCYOWaK9NQesAaozmjpug1_QxkngRRoka62UeTRQkKlBPsTghS6GN7F6uELdzKuhIZYgUSF0DVX2t8sPMr0-pomLxFDFCT94Dwc5HOmjngzlsljNoPoJFIFc1Nyb8SX4zw7GDe6LcVSQMklThhx-0VUhQKPHLkXP61hFKU7SOl_ZprNvN1NfJmpT8UPEa4kQMo-P9wCTNGskSGweQzwTlmyNkTphfHg4ldeuHRB4s_7Q
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
Date: Tue, 06 Oct 2015 21:22:18 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/templates/4b88c433-28fd-deaf-da92-5a1cf6abc09a"
    }
  },
  "html": "html",
  "id": "4b88c433-28fd-deaf-da92-5a1cf6abc09a",
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
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzgsImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.mM4oj3NYciRFkmYUfT9HrC0Pl19zGFQXTPpS82PZ8oZhmSkCtlg_BQM0xJl37nQslX9jGaphr_eYGhk08ENvOL6HOeglkBDSr4CC3c4KtuMCYOWaK9NQesAaozmjpug1_QxkngRRoka62UeTRQkKlBPsTghS6GN7F6uELdzKuhIZYgUSF0DVX2t8sPMr0-pomLxFDFCT94Dwc5HOmjngzlsljNoPoJFIFc1Nyb8SX4zw7GDe6LcVSQMklThhx-0VUhQKPHLkXP61hFKU7SOl_ZprNvN1NfJmpT8UPEa4kQMo-P9wCTNGskSGweQzwTlmyNkTphfHg4ldeuHRB4s_7Q
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 06 Oct 2015 21:22:18 GMT
```


## Campaign Types
EVAN WILL EDIT THIS
<a name="campaign-type-create"></a>
### Create a new campaign type
#### Request **POST** /senders/{id}/campaign_types
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzksImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.eo13IleXhSxNJItX4qa8G-zq5JkZI5bu-aRg9hAwGE-5OpJS2h4rTFksWmCUsqjDk3xVHpJZoIjC2wqO8mO1VokzAMosDQu0P9ux11DBmze2HfeQd72cpcjKqD2zVFImNFJarQWvToykiKGyMx5gVgpwWT-pDgaAAwPSZXbHQyN-E3TbCFRtrTvTgC994r6xah36zj9bbCyPvxU6dy2tWkBbkhYfwOTW5-qntiIg-_R91_dLQpGoKZV-T7W48yJ6kgFElq1FDjOLzvgSI_YdnygWm6asjrucxM-rbv07_ObiTqC9GK1eieel07yHvUrpIwvfm0MnTQK-ICijtiOSPA
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
Date: Tue, 06 Oct 2015 21:22:19 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/906104fe-4208-5c46-6894-1fa0e65c3d73"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "906104fe-4208-5c46-6894-1fa0e65c3d73",
  "name": "some-campaign-type",
  "template_id": ""
}
```

<a name="campaign-type-list"></a>
### Retrieve a list of campaign types
#### Request **GET** /senders/{id}/campaign_types
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzksImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.eo13IleXhSxNJItX4qa8G-zq5JkZI5bu-aRg9hAwGE-5OpJS2h4rTFksWmCUsqjDk3xVHpJZoIjC2wqO8mO1VokzAMosDQu0P9ux11DBmze2HfeQd72cpcjKqD2zVFImNFJarQWvToykiKGyMx5gVgpwWT-pDgaAAwPSZXbHQyN-E3TbCFRtrTvTgC994r6xah36zj9bbCyPvxU6dy2tWkBbkhYfwOTW5-qntiIg-_R91_dLQpGoKZV-T7W48yJ6kgFElq1FDjOLzvgSI_YdnygWm6asjrucxM-rbv07_ObiTqC9GK1eieel07yHvUrpIwvfm0MnTQK-ICijtiOSPA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 405
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:19 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/senders/679de949-023a-abf9-bad0-32dd1f531bb9/campaign_types"
    },
    "sender": {
      "href": "/senders/679de949-023a-abf9-bad0-32dd1f531bb9"
    }
  },
  "campaign_types": [
    {
      "_links": {
        "self": {
          "href": "/campaign_types/906104fe-4208-5c46-6894-1fa0e65c3d73"
        }
      },
      "critical": false,
      "description": "a great campaign type",
      "id": "906104fe-4208-5c46-6894-1fa0e65c3d73",
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
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzksImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.eo13IleXhSxNJItX4qa8G-zq5JkZI5bu-aRg9hAwGE-5OpJS2h4rTFksWmCUsqjDk3xVHpJZoIjC2wqO8mO1VokzAMosDQu0P9ux11DBmze2HfeQd72cpcjKqD2zVFImNFJarQWvToykiKGyMx5gVgpwWT-pDgaAAwPSZXbHQyN-E3TbCFRtrTvTgC994r6xah36zj9bbCyPvxU6dy2tWkBbkhYfwOTW5-qntiIg-_R91_dLQpGoKZV-T7W48yJ6kgFElq1FDjOLzvgSI_YdnygWm6asjrucxM-rbv07_ObiTqC9GK1eieel07yHvUrpIwvfm0MnTQK-ICijtiOSPA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 228
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:19 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/906104fe-4208-5c46-6894-1fa0e65c3d73"
    }
  },
  "critical": false,
  "description": "a great campaign type",
  "id": "906104fe-4208-5c46-6894-1fa0e65c3d73",
  "name": "some-campaign-type",
  "template_id": ""
}
```

<a name="campaign-type-update"></a>
### Update a campaign type
#### Request **PUT** /senders/{id}/campaign_types/{id}
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzksImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.eo13IleXhSxNJItX4qa8G-zq5JkZI5bu-aRg9hAwGE-5OpJS2h4rTFksWmCUsqjDk3xVHpJZoIjC2wqO8mO1VokzAMosDQu0P9ux11DBmze2HfeQd72cpcjKqD2zVFImNFJarQWvToykiKGyMx5gVgpwWT-pDgaAAwPSZXbHQyN-E3TbCFRtrTvTgC994r6xah36zj9bbCyPvxU6dy2tWkBbkhYfwOTW5-qntiIg-_R91_dLQpGoKZV-T7W48yJ6kgFElq1FDjOLzvgSI_YdnygWm6asjrucxM-rbv07_ObiTqC9GK1eieel07yHvUrpIwvfm0MnTQK-ICijtiOSPA
X-Notifications-Version: 2
```
##### Body
```
{
  "critical": true,
  "description": "still the same great campaign type",
  "name": "updated-campaign-type",
  "template_id": "0a12e553-d83d-07d0-5e72-e2a3f43399f6"
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 279
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:19 GMT
```
##### Body
```
{
  "_links": {
    "self": {
      "href": "/campaign_types/906104fe-4208-5c46-6894-1fa0e65c3d73"
    }
  },
  "critical": true,
  "description": "still the same great campaign type",
  "id": "906104fe-4208-5c46-6894-1fa0e65c3d73",
  "name": "updated-campaign-type",
  "template_id": "0a12e553-d83d-07d0-5e72-e2a3f43399f6"
}
```

<a name="campaign-type-delete"></a>
### Delete a campaign type
#### Request **DELETE** /senders/{id}/campaign_types/{id}
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzksImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.eo13IleXhSxNJItX4qa8G-zq5JkZI5bu-aRg9hAwGE-5OpJS2h4rTFksWmCUsqjDk3xVHpJZoIjC2wqO8mO1VokzAMosDQu0P9ux11DBmze2HfeQd72cpcjKqD2zVFImNFJarQWvToykiKGyMx5gVgpwWT-pDgaAAwPSZXbHQyN-E3TbCFRtrTvTgC994r6xah36zj9bbCyPvxU6dy2tWkBbkhYfwOTW5-qntiIg-_R91_dLQpGoKZV-T7W48yJ6kgFElq1FDjOLzvgSI_YdnygWm6asjrucxM-rbv07_ObiTqC9GK1eieel07yHvUrpIwvfm0MnTQK-ICijtiOSPA
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: Tue, 06 Oct 2015 21:22:19 GMT
```


## Campaigns
EVAN WILL EDIT THIS
<a name="campaign-create"></a>
### Create a new campaign
#### Request **POST** /senders/{id}/campaigns
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzksImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.eo13IleXhSxNJItX4qa8G-zq5JkZI5bu-aRg9hAwGE-5OpJS2h4rTFksWmCUsqjDk3xVHpJZoIjC2wqO8mO1VokzAMosDQu0P9ux11DBmze2HfeQd72cpcjKqD2zVFImNFJarQWvToykiKGyMx5gVgpwWT-pDgaAAwPSZXbHQyN-E3TbCFRtrTvTgC994r6xah36zj9bbCyPvxU6dy2tWkBbkhYfwOTW5-qntiIg-_R91_dLQpGoKZV-T7W48yJ6kgFElq1FDjOLzvgSI_YdnygWm6asjrucxM-rbv07_ObiTqC9GK1eieel07yHvUrpIwvfm0MnTQK-ICijtiOSPA
X-Notifications-Version: 2
```
##### Body
```
{
  "campaign_type_id": "a18abeb5-3550-a383-ef35-9d8f14ae5476",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "email": "test@example.com"
  },
  "subject": "campaign subject",
  "template_id": "505ba60f-38e6-2efd-42b3-24dae6362acd",
  "text": "campaign body"
}
```
#### Response 202 Accepted
##### Headers
```
Content-Length: 593
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:19 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/a18abeb5-3550-a383-ef35-9d8f14ae5476"
    },
    "self": {
      "href": "/campaigns/b3a56525-a8c3-c357-a7f2-662afda384f4"
    },
    "status": {
      "href": "/campaigns/b3a56525-a8c3-c357-a7f2-662afda384f4/status"
    },
    "template": {
      "href": "/templates/505ba60f-38e6-2efd-42b3-24dae6362acd"
    }
  },
  "campaign_type_id": "a18abeb5-3550-a383-ef35-9d8f14ae5476",
  "html": "",
  "id": "b3a56525-a8c3-c357-a7f2-662afda384f4",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "email": "test@example.com"
  },
  "subject": "campaign subject",
  "template_id": "505ba60f-38e6-2efd-42b3-24dae6362acd",
  "text": "campaign body"
}
```

<a name="campaign-get"></a>
### Retrieve a campaign
#### Request **GET** /campaigns/{id}
##### Required Scopes
```
notifications.manage notifications.write emails.write notification_preferences.admin critical_notifications.write notification_templates.admin notification_templates.write notification_templates.read
```
##### Headers
```
Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJjbGllbnRfaWQiOiJteS1jbGllbnQiLCJleHAiOjE0NDQ0MjU3MzksImlzcyI6Imh0dHA6Ly8xMjcuMC4wLjE6NTY2NTMvb2F1dGgvdG9rZW4iLCJzY29wZSI6WyJub3RpZmljYXRpb25zLm1hbmFnZSIsIm5vdGlmaWNhdGlvbnMud3JpdGUiLCJlbWFpbHMud3JpdGUiLCJub3RpZmljYXRpb25fcHJlZmVyZW5jZXMuYWRtaW4iLCJjcml0aWNhbF9ub3RpZmljYXRpb25zLndyaXRlIiwibm90aWZpY2F0aW9uX3RlbXBsYXRlcy5hZG1pbiIsIm5vdGlmaWNhdGlvbl90ZW1wbGF0ZXMud3JpdGUiLCJub3RpZmljYXRpb25fdGVtcGxhdGVzLnJlYWQiXX0.eo13IleXhSxNJItX4qa8G-zq5JkZI5bu-aRg9hAwGE-5OpJS2h4rTFksWmCUsqjDk3xVHpJZoIjC2wqO8mO1VokzAMosDQu0P9ux11DBmze2HfeQd72cpcjKqD2zVFImNFJarQWvToykiKGyMx5gVgpwWT-pDgaAAwPSZXbHQyN-E3TbCFRtrTvTgC994r6xah36zj9bbCyPvxU6dy2tWkBbkhYfwOTW5-qntiIg-_R91_dLQpGoKZV-T7W48yJ6kgFElq1FDjOLzvgSI_YdnygWm6asjrucxM-rbv07_ObiTqC9GK1eieel07yHvUrpIwvfm0MnTQK-ICijtiOSPA
X-Notifications-Version: 2
```
#### Response 200 OK
##### Headers
```
Content-Length: 593
Content-Type: text/plain; charset=utf-8
Date: Tue, 06 Oct 2015 21:22:19 GMT
```
##### Body
```
{
  "_links": {
    "campaign_type": {
      "href": "/campaign_types/a18abeb5-3550-a383-ef35-9d8f14ae5476"
    },
    "self": {
      "href": "/campaigns/b3a56525-a8c3-c357-a7f2-662afda384f4"
    },
    "status": {
      "href": "/campaigns/b3a56525-a8c3-c357-a7f2-662afda384f4/status"
    },
    "template": {
      "href": "/templates/505ba60f-38e6-2efd-42b3-24dae6362acd"
    }
  },
  "campaign_type_id": "a18abeb5-3550-a383-ef35-9d8f14ae5476",
  "html": "",
  "id": "b3a56525-a8c3-c357-a7f2-662afda384f4",
  "reply_to": "no-reply@example.com",
  "send_to": {
    "email": "test@example.com"
  },
  "subject": "campaign subject",
  "template_id": "505ba60f-38e6-2efd-42b3-24dae6362acd",
  "text": "campaign body"
}
```


