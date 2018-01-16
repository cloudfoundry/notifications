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
* Bananas
  * [Create Bananas](#bananas-create)
  * [Update a Banana](#bananas-update)
* Bicycles
  * [Create Bicycles](#bicycles-create)
  * [Delete a Bicycle](#bicycles-delete)

## Bananas
These are really tasty!
<a name="bananas-create"></a>
### Create Bananas
#### Request **POST** /bananas
##### Required Scopes
```
bananas.grow
```
##### Headers
```
Authorization: Bearer some-token
X-Notifications-Version: 2
```
##### Body
```
{
  "organic": true
}
```
#### Response 201 Created
##### Headers
```
Content-Length: 123
Content-Type: application/json
```
##### Body
```
{
  "id": "banana-1",
  "organic": true
}
```

<a name="bananas-update"></a>
### Update a Banana
#### Request **PUT** /bananas/{id}
##### Headers
```
Authorization: Bearer some-token
X-Notifications-Version: 2
```
##### Body
```
{
  "organic": false
}
```
#### Response 200 OK
##### Headers
```
Content-Length: 123
Content-Type: application/json
```
##### Body
```
{
  "id": "banana-1",
  "organic": false
}
```


## Bicycles
Fun when they go fast.
<a name="bicycles-create"></a>
### Create Bicycles
#### Request **POST** /bicycles
##### Required Scopes
```
bicycles.manufacture
```
##### Headers
```
Authorization: Bearer some-token
X-Notifications-Version: 2
```
##### Body
```
{
  "color": "blue",
  "size": 46
}
```
#### Response 201 Created
##### Headers
```
Content-Length: 123
Content-Type: application/json
```
##### Body
```
{
  "color": "blue",
  "id": "bicycle-15",
  "size": 46
}
```

<a name="bicycles-delete"></a>
### Delete a Bicycle
#### Request **DELETE** /bicycles/{id}
##### Required Scopes
```
bicycles.crash
```
##### Headers
```
Authorization: Bearer some-token
X-Notifications-Version: 2
```
#### Response 204 No Content
##### Headers
```
Date: today
```
