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

### Overriding default templates

Any file that exists in `./templates` can be overriden by placing a file of the same name in `./templates/overrides`.

The templates are [go templates](http://golang.org/pkg/text/template/).  The templates have access to the following variables:

| Variable | Description |
| -------- | ----------- |
| KindDescription | Pulled from json posted to endpoint under: kind_description, falls back to kind if not set |
| From | what account is in the from field of the email |
| To | the address the email is going to |
| Subject | Pulled from json posted to endpoint under: subject |
| Text | Pulled from json posted to endpoint under: text |
| HTML | Pulled from json posted to endpoint under: html |
| SourceDescription | Pulled from json posted to endpoint under: source_description, falls back to ClientID if not set |
| ClientID | the access token of the user requesting the email be sent |
| MessageID | unique id for the email being sent |
| Organization | The name of the organization of the space (used for emails to spaces) |
| Space | The name of the space (used for emails to spaces) |

### Example: Overriding space_body.text
To override the plain text template in the email body, write the following in `./templates/overrides/space_body.text`:

```
You are receiving this electronic mail because you are a member of {{.Space}}

All apps in {{.Space}} have had an emergency of type {{.KindDescription}}
```

