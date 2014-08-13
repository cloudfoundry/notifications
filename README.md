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

## Configuring Email Templates
The default templates are located in **./templates**. The templates directory should look similar to this:

	overrides/
	space_body.html
	space_body.text
	subject.missing
	subject.provided
	user_body.html
	user_body.text

When emailing a space, `space_body.html` and `space_body.text` are used as templates in the email body for the html and plaintext, respectively. When emailing a single user, `user_body.html` and `user_body.text` are used as templates in the email body for the html and plaintext, respectively. When the subject is provided, the default subject template is found in `subject.provided`, while `subject.missing` is used when the email subject is not provided.

The files located in the templates directory are defaults and will be used if there is no override of the same name.

## Posting to a notifications endpoint

Notifications currently supports two different types of messages.  Messages to individual users and messages to spaces.

Users are messaged via the `/users/id` endpoint and spaces are messaged via `/spaces/id` endpoint.

Both endpoints expect a json body to be posted with following keys:

| Key | Description |
| ---- | ------ |
| kind_id\* | a key to identify the type of email to be sent |
| text\*\* | the text version of the email |
| html\*\* | the html version of the email |
| kind_description | a description of the kind_id |
| subject | the text of the subject |
| reply_to | the Reply-To address for the email |
| source_description | a description of the sender |

\* required 
 
\*\* either text or html have to be set, not both


### Overriding default templates
There are three types of overrides. 

| Type | Effect |
| ---- | ------ |
| global | applies to any email sent, has least precedence |
| client specific | applies to any email sent by a client, higher precedence than global |
| client and kind specific | applies to any email sent by a client with the kind_id, has highest precedence |

All types of overrides have access to the same variables described in global overrides.  The only difference is which email types they are applied to.

#### global overrides

Any file that exists in `./templates` can be overriden by placing a file of the same name in `./templates/overrides`.

The templates are [go templates](http://golang.org/pkg/text/template/).  The templates have access to the following variables:

| Variable | Description |
| -------- | ----------- |
| KindDescription | Pulled from json posted to endpoint under: kind_description, falls back to kind if not set |
| From | what account is in the from field of the email |
| ReplyTo | the address is in the reply to field of the email |
| To | the address the email is going to |
| Subject | Pulled from json posted to endpoint under: subject |
| Text | Pulled from json posted to endpoint under: text |
| HTML | Pulled from json posted to endpoint under: html |
| SourceDescription | Pulled from json posted to endpoint under: source_description, falls back to ClientID if not set |
| ClientID | the access token of the user requesting the email be sent |
| MessageID | unique id for the email being sent |
| Organization | The name of the organization of the space (used for emails to spaces) |
| Space | The name of the space (used for emails to spaces) |

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

