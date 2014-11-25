# Notifications

A notifications component that parts of CF can use to send email to end users

## UAA Client
Notifications itself needs the following UAA client in order to retrieve arbitrary users' email addresses and check the members of  arbitrary spaces:

```yaml
properties:
  uaa:
    clients:
      notifications:
        secret: my-secret
        authorities: cloud_controller.admin,scim.read
        authorized-grant-types: client_credentials
```

### Client Configurations
#### Send Notifications
The following client configurations are needed for sending messages to individual users, users in a specific space and arbitrary email addresses. 

To send non-critical notifications, notifications.write scope is required. Sending critical notifications requires critical_notifications.write scope. To send notifications to an arbitrary email address requires emails.write scope.

```yaml
notifications-client-name:
  scope: uaa.none
  resource_ids: none
  authorized_grant_types: client_credentials
  authorities: notifications.write critical_notifications.write emails.write
  autoapprove:
```

#### View and Edit User Preferences
To view and edit a user's preferences for receiving non-critical notifications, a client will need to be configured with notification_preferences.read scope and notification_preferences.write scope.

A client with notification_preferences.admin scope has the ability to retrieve an arbitrary user's preferences.

```yaml
notifications-client-name:
  scope: notification_preferences.read notification_preferences.write openid
  resource_ids: none
  authorized_grant_types: authorization_code client_credentials refresh_token
  redirect_uri: http://example.com/sessions/create
  authorities: notification_preferences.admin
  autoapprove:
```

If you are unfamiliar with UAA consult the [UAA token overview](https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-Tokens.md).

##Configuring Environment Variables

| Variable                     | Description                                 | Default  |
|------------------------------|---------------------------------------------|----------|
| CC_HOST\*                    | Cloud Controller Host                       | \<none\> |
| CORS_ORIGIN                  | Value to use for CORS Origin Header         | *        |
| DB_LOGGING_ENABLED           | Logs DB interactions when set to true       | false    |
| DATABASE_URL\*               | URL to your Database                        | \<none\> |
| ENCRYPTION_KEY\*             | Key used to encrypt the unsubscribe ID      | \<none\> |
| GOBBLE_MIGRATIONS_DIR\*      | Location of the gobble migrations directory | \<none\> |
| MODEL_MIGRATIONS_DIRECTORY\* | Location of the model migrations directory  | \<none\> |
| PORT                         | Port that application will bind to          | 3000     |
| ROOT_PATH\*                  | Root path of your application               | \<none\> |
| SMTP_LOGGING_ENABLED         | Logs SMTP interactions when set to true     | \<none\> |
| SMTP_HOST\*                  | SMTP Host                                   | \<none\> |
| SMTP_PASS                    | SMTP Password                               | \<none\> |
| SMTP_PORT\*                  | SMTP Port                                   | \<none\> |
| SMTP_TLS                     | Use TLS when talking to SMTP server         | true     |
| SMTP_USER                    | SMTP Username                               | \<none\> |
| SENDER\*                     | Emails are sent from this address           | \<none\> |
| TEST_MODE                    | Run in test mode                            | false    |
| UAA_CLIENT_ID\*              | The UAA client ID                           | \<none\> |
| UAA_CLIENT_SECRET\*          | The UAA client secret                       | \<none\> |
| UAA_HOST\*                   | The UAA Host                                | \<none\> |
| VERIFY_SSL                   | Verifies SSL                                | true     |


\* required

## Posting to a notifications endpoint

Notifications currently supports several different types of messages.  Messages can be sent to:

 - Users via the `/users/id` endpoint
 - Spaces via the `/spaces/id` endpoint
 - Organizations via the `/organizations/id` endpoint
 - All users in the system via the `/everyone` endpoint
 - UAA Scopes via the `/uaa_scopes/scope` endpoint
 - Emails via the `/emails` endpoint

The Users, Spaces, Organizations, Everyone, and UAA Scopes endpoints expect a json body to be posted with following keys:

| Key                  | Description                                    |
|----------------------|------------------------------------------------|
| kind_id\*            | a key to identify the type of email to be sent |
| text\*\*             | the text version of the email                  |
| html\*\*             | the html version of the email                  |
| kind_description     | a description of the kind_id                   |
| subject\*            | the text of the subject                        |
| reply_to             | the Reply-To address for the email             |
| source_description   | a description of the sender                    |

\* required

\*\* either text or html have to be set, not both

The Emails endpoint expects a json body to be posted with the following keys:

| Key                | Description                                    |
|--------------------|------------------------------------------------|
| to\*               | the recipient of the email                     |
| subject\*          | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |
| text\**            | the text version of the email                  |
| html\**            | the html version of the email                  |

\* required

\*\* either text or html have to be set, not both


[Further API Documentation](/API.md)

## Configuring Email Templates
The default templates are located in **./templates**. The templates directory should look similar to this:

	overrides/
	user_body.text
	user_body.html
	space_body.text
	space_body.html
	organization_body.text
	organization_body.html
	everyone_body.text
	everyone_body.html
	uaa_scope_body.html
	uaa_scope_body.text
	email_body.text
	email_body.html
	subject.provided


When emailing a single user, `user_body.html` and `user_body.text` are used as templates in the email body for the html and plaintext, respectively. When emailing a space, `space_body.html` and `space_body.text` are used as templates in the email body for the html and plaintext, respectively. When emailing an organization, `organization_body.html` and `organization_body.text` are used as templates in the email body for the html and plaintext, respectively. When emailing all users in a system, `everyone_body.html` and `everyone_body.text` are used as templates in the email body for the html and plaintext, respectively. When emailing a UAA scope, `uaa_scope_body.html` and `uaa_scope_body.text` are used as templates in the email body for the html and plaintext, respectively. 

The default subject template is found in `subject.provided`.

When using the `/emails` endpoint, `email_body.html` and `email_body.text` are used as templates in the email body for the html and plaintext, respectively.

The files located in the templates directory are defaults and will be used if there is no override of the same name.


### Overriding default templates
There are three types of overrides.

| Type                     | Effect                                                                         |
|--------------------------|--------------------------------------------------------------------------------|
| global                   | applies to any email sent, has least precedence                                |
| client specific          | applies to any email sent by a client, higher precedence than global           |
| client and kind specific | applies to any email sent by a client with the kind_id, has highest precedence |

All types of overrides have access to the same variables described in global overrides.  The only difference is which email types they are applied to.

#### global overrides

Any file that exists in `./templates` can be overriden by placing a file of the same name in `./templates/overrides`.

The templates are [go templates](http://golang.org/pkg/text/template/).  The templates have access to the following variables:

| Variable          | Description                                                                                      |
|-------------------|--------------------------------------------------------------------------------------------------|
| KindDescription   | Pulled from json posted to endpoint under: kind_description, falls back to kind if not set       |
| From              | Which account is in the from field of the email                                                  |
| ReplyTo           | The address is in the reply to field of the email                                                |
| To                | The address the email is going to                                                                |
| Subject           | Pulled from json posted to endpoint under: subject                                               |
| Text              | Pulled from json posted to endpoint under: text                                                  |
| HTML              | Pulled from json posted to endpoint under: html                                                  |
| SourceDescription | Pulled from json posted to endpoint under: source_description, falls back to ClientID if not set |
| ClientID          | The access token of the user requesting the email be sent                                        |
| MessageID         | Unique id for the email being sent                                                               |
| UserGUID          | Unique id for the user the email is sent to                                                      |
| Organization      | The name of the organization of the space (used for emails to spaces)                            |
| OrganizationGUID  | The guid of the organization of the space (used for emails to spaces)                            |
| Space             | The name of the space (used for emails to spaces)                                                |
| SpaceGUID         | The guid of the space (used for emails to spaces)                                                |
| UnsubscribeID     | The id for unsubscribing a user from a notification. See [here](#unsubscribe-id) for details     |

##### Example: Overriding space_body.text
To override the plain text template in the email body, write the following in `./templates/overrides/space_body.text`:

```
You are receiving this electronic mail because you are a member of {{.Space}}

All apps in {{.Space}} have had an emergency of type {{.KindDescription}}
```

#### overriding by clientID

To override a template for a client with id `banana`. You would place a file in the `./templates/overrides` directory with the following name scheme:

	clientID.templateName.templateExtension
	
So to override the subject.missing the file name would be:

	banana.subject.missing
	
This override would only be applied to emails sent by the client banana.
	

#### overriding by clientID and kind_id

To override a template for client `banana` with the kind_id `damage` you would place a file in the `./templates/overrides` directory with the following name scheme:

	clientID.kind_id.templateName.templateExtension
	
So to override the user_body.text template the file name would be:

	banana.damage.user_body.text
	
This override only applies to requests that match both the clientId and the kind.  This has the most precedence and overrides all other overrides.

[Further Template API Documentation](/API.md#get_template)

<a name="unsubscribe-id"></a>
#### UnsubscribeID

AES Encryption is used to encrypt a token value for unsubscribing a user from a 
notification. The format of the token is the `user_guid|client_id|kind_id`. The key used to instantiate a cipher is a 16 byte MD5 sum of the text given to the `ENCRYPTION_KEY` environment variable.

Encrypting:

1. Concatenate user GUID, client ID, and kind ID into a single string, delimited by a `|` character.
1. Base64 encode the concatenated string.
1. Encrypt the encoded text using AES cipher in CFB mode.
1. Base64 encode the cipher text.

Decrypting:

1. Base64 decode the unsubscribe token.
1. Decrypt the decoded text using AES cipher in CFB mode.
1. Base64 decode the decrypted text.
1. Split the text at the `|` characters.



### Development

#### Running locally

The application can be run locally by executing the `./bin/run` script. This script will look for a file called `./bin/env/development` to load environment variables. Setting the `TEST_MODE` env var to true will disable the requirement for a running SMTP server.
