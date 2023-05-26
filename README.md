# Notice

If you are trying to deploy notifications, **DO NOT** use this repo to deploy it directly. Please use the [bosh-release and accompanying directions](https://github.com/cloudfoundry/notifications-release).

# Notifications [![CI Status](https://wings.pivotal.io/api/v1/teams/cf-notifications/pipelines/cf-notifications/badge)](https://wings.pivotal.io/teams/cf-notifications/pipelines/cf-notifications)

A notifications component that parts of CF can use to send email to end users.
There is a [walkthrough](https://github.com/cloudfoundry/notifications/blob/master/walkthrough.md) of using the Notifications Service which will show how to set it up and use it.

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
  authorities: notifications.write,critical_notifications.write,emails.write
  autoapprove:
```

#### View and Edit User Preferences
To view and edit a user's preferences for receiving non-critical notifications, a client will need to be configured with notification_preferences.read scope and notification_preferences.write scope.

A client with notification_preferences.admin scope has the ability to retrieve an arbitrary user's preferences.

```yaml
notifications-client-name:
  scope: notification_preferences.read,notification_preferences.write,openid
  resource_ids: none
  authorized_grant_types: authorization_code,client_credentials,refresh_token
  redirect_uri: http://example.com/sessions/create
  authorities: notification_preferences.admin
  autoapprove:
```

If you are unfamiliar with UAA consult the [UAA token overview](https://github.com/cloudfoundry/uaa/blob/master/docs/UAA-Tokens.md).

## Configuring Environment Variables

| Variable                     | Description                                 | Default  |
|------------------------------|---------------------------------------------|----------|
| CC_HOST\*                    | Cloud Controller Host                       | \<none\> |
| CORS_ORIGIN                  | Value to use for CORS Origin Header         | *        |
| DB_LOGGING_ENABLED           | Logs DB interactions when set to true       | false    |
| DB_MAX_OPEN_CONNS            | Maximum number of open DB connections       | 0 (unlimited) |
| DATABASE_URL\*               | URL to your Database                        | \<none\> |
| DEFAULT_UAA_SCOPES\*         | Comma separated list of scopes              | \<none\> |
| ENCRYPTION_KEY\*             | Key used to encrypt the unsubscribe ID      | \<none\> |
| GOBBLE_MIGRATIONS_DIR\*      | Location of the gobble migrations directory | \<none\> |
| PORT                         | Port that application will bind to          | 3000     |
| ROOT_PATH\*                  | Root path of your application               | \<none\> |
| SMTP_AUTH_MECHANISM\*        | SMTP Authentication (none, plain, cram-md5). Most users will want to use `plain`. | \<none\> |
| SMTP_CRAMMD5_SECRET          | Secret value used for CRAMMD5 SMTP auth     | \<none\> |
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

 - Users via the `/users/:id` endpoint
 - Spaces via the `/spaces/:id` endpoint
 - Organizations via the `/organizations/:id` endpoint
 - All users in the system via the `/everyone` endpoint
 - UAA Scopes via the `/uaa_scopes/:scope` endpoint
 - Emails via the `/emails` endpoint

The Users, Spaces, Organizations, Everyone, and UAA Scopes endpoints expect a json body to be posted with following keys:

| Key                  | Description                                    |
|----------------------|------------------------------------------------|
| kind_id\*            | a key to identify the type of email to be sent |
| text\*\*             | the text version of the email                  |
| html\*\*             | the html version of the email                  |
| subject\*            | the text of the subject                        |
| reply_to             | the Reply-To address for the email             |

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


<a name="api-docs"></a>
### API Documentation
- [Version 1 Documentation](/V1_API.md)

## Configuring Email Templates
You can do a whole lot to configure templates for your notifications, see [API Docs](#api-docs) for specific endpoints available!

<a name="unsubscribe-id"></a>
#### UnsubscribeID

AES Encryption is used to encrypt a token value for unsubscribing a user from a
notification. The format of the token is the `user_guid|client_id|kind_id`. The
key used to instantiate a cipher is a 16 byte MD5 sum of the text given to the
`ENCRYPTION_KEY` environment variable.

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

#### Running tests

Docker is needed to run tests. 

To get the required image, change into the `docker` directory and run `docker-compose up -d`. 

If this is successful `docker ps` should show a mariadb image running on port 3306 and mysql should have a database called `notifications_test`.

Move up a directory to the root of the project and run `./bin/test` to run tests.
