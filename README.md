# Notifications

A notifications component that parts of CF can use to send email to end users

## UAA Client

```yaml
properties:
  uaa:
    clients:
      notifications:
        secret: my-secret
        authorities: cloud_controller.admin,scim.read
        authorized-grant-types: client_credentials
```

##Configuring Environment Variables

| Variable              | Description                                 | Default  |
| --------------------- | ------------------------------------------- | -------- |
| CC_HOST\*             | Cloud Controller Host                       | \<none\> |
| CORS_ORIGIN           | Value to use for CORS Origin Header         | *        |
| DATABASE_URL\*        | URL to your Database                        | \<none\> |
| PORT                  | Port that application will bind to          | 3000     |
| ROOT_PATH\*           | Root path of your application               | \<none\> |
| SMTP_HOST\*           | SMTP Host                                   | \<none\> |
| SMTP_PASS\*           | SMTP Password                               | \<none\> |
| SMTP_PORT\*           | SMTP Port                                   | \<none\> |
| SMTP_TLS              | Use TLS when talking to SMTP server         | true     |
| SMTP_USER\*           | SMTP Username                               | \<none\> |
| SENDER\*              | Emails are sent from this address           | \<none\> |
| UAA_CLIENT_ID\*       | The UAA client ID                           | \<none\> |
| UAA_CLIENT_SECRET\*   | The UAA client secret                       | \<none\> |
| UAA_HOST\*            | The UAA Host                                | \<none\> |
| VERIFY_SSL            | Verifies SSL                                | true     |
| TEST_MODE             | Run in test mode                            | false    |
| GOBBLE_MIGRATIONS_DIR | Location of the gobble migrations directory | \<none\> |

\* required

## Posting to a notifications endpoint

Notifications currently supports two different types of messages.  Messages to individual users and messages to spaces.

Users are messaged via the `/users/id` endpoint and spaces are messaged via `/spaces/id` endpoint.

Both endpoints expect a json body to be posted with following keys:

| Key                | Description                                    |
| ------------------ | ---------------------------------------------- |
| kind_id\*          | a key to identify the type of email to be sent |
| text\*\*           | the text version of the email                  |
| html\*\*           | the html version of the email                  |
| kind_description   | a description of the kind_id                   |
| subject            | the text of the subject                        |
| reply_to           | the Reply-To address for the email             |
| source_description | a description of the sender                    |

\* required

\*\* either text or html have to be set, not both

[Further API Documentation](/API.md)

## Configuring Email Templates
The default templates are located in **./templates**. The templates directory should look similar to this:

	overrides/
	email_body.html
	email_body.text
	space_body.html
	space_body.text
	subject.missing
	subject.provided
	user_body.html
	user_body.text

When emailing a space, `space_body.html` and `space_body.text` are used as templates in the email body for the html and plaintext, respectively. When emailing a single user, `user_body.html` and `user_body.text` are used as templates in the email body for the html and plaintext, respectively. When the subject is provided, the default subject template is found in `subject.provided`, while `subject.missing` is used when the email subject is not provided. When using the `/emails` endpoint, `email_body.html` and `email_body.text` are used as templates in the email body for the html and plaintext, respectively.

The files located in the templates directory are defaults and will be used if there is no override of the same name.


### Overriding default templates
There are three types of overrides.

| Type                     | Effect                                                                         |
| ------------------------ | ------------------------------------------------------------------------------ |
| global                   | applies to any email sent, has least precedence                                |
| client specific          | applies to any email sent by a client, higher precedence than global           |
| client and kind specific | applies to any email sent by a client with the kind_id, has highest precedence |

All types of overrides have access to the same variables described in global overrides.  The only difference is which email types they are applied to.

#### global overrides

Any file that exists in `./templates` can be overriden by placing a file of the same name in `./templates/overrides`.

The templates are [go templates](http://golang.org/pkg/text/template/).  The templates have access to the following variables:

| Variable          | Description                                                                                      |
| ----------------- | ------------------------------------------------------------------------------------------------ |
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
| Space             | The name of the space (used for emails to spaces)                                                |

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

### Development

#### Running locally

The application can be run locally by executing the `./bin/run` script. This script will look for a file called `./bin/env/development` to load environment variables. Setting the `TEST_MODE` env var to true will disable the requirement for a running SMTP server.
