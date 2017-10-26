# So you want to send a notification

This is a walkthrough of using the Notifications Service. We'll show you how to set it up, through the following example: Darla is a cloud operator who needs to take the system down for maintenance and wants to notify everybody about her intentions.

## Prerequisites
We're assuming that Darla has

- CloudFoundry set up somewhere.  Darla uses the [Bosh Lite](https://github.com/cloudfoundry/bosh-lite) install as her development environment.
- An `admin` account on her CloudFoundry setup
- The Notifications service installed and running - She can use the [BOSH release](https://github.com/cloudfoundry-incubator/notifications-release) to simplify deployment to her environment
- The [`cf`](https://github.com/cloudfoundry/cli) and [`uaac`](https://rubygems.org/gems/cf-uaac) command line tools available.

## Create your Client and get a token
To interact with the Notifications service, Darla needs certain [UAA](http://docs.cloudfoundry.org/concepts/architecture/uaa.html) scopes (authorities).  Rather than use her `admin` user account directly, she creates a `notifications-admin` client with the required scopes:

```
uaac client add notifications-admin --authorized_grant_types client_credentials --authorities \
    notifications.manage,notifications.write,notification_templates.write,notification_templates.read,critical_notifications.write
```

It's worth noting that she doesn't need all of these scopes just to send a notification. `notifications.manage` is used to update notifications and assign templates for that notification. `notification_templates.write` allows Darla to custom make her own template for a notification, and `notification_templates.read` allows her to check which templates are saved in the database. Finally, `notifications.write` is the scope necessary to send a notification to a user, space, everyone in the system, and more!

Now, Darla logs in via her newly created client (stay logged in with this client for the rest of the examples below):

```
uaac token client get notifications-admin
```

## Registering Notifications
Darla can't send a notification unless she has registered it first. Again, registering notifications requires the `notifications.manage` scope on her client.

```
uaac curl https://notifications.darla.example.com/notifications -X PUT --data '{  "source_name": "Cloud Ops Team",
  "notifications": {
     "system-going-down": {"critical": true, "description": "Cloud going down" },
     "system-up": { "critical": false, "description": "Cloud back up" }
     }
 }'
```
Darla has registered two different notifications, `system-going-down` and `system-up.` In addition, she gives the `notifications-admin` client the human-friendly description "Cloud Ops Team." We want to note that she has made the "system-going-down" notification `critical`. This means that no users can unsubscribe from that notification. Setting notifications as critical requires the `critical_notifications.write` scope.

## Create a custom template
The system provides a default template for all notifications, but Darla has decided to forgo this luxury.  Darla wants to include her own branding and has opted to create her own custom template using the curl below (note this action requires the `notification_templates.write` scope):

```
uaac curl https://notifications.darla.example.com/templates -X POST --data \
'{"name":"site-maintenance","subject":"Maintenance: {{.Subject}}","text":"The site has gone down for maintenance.  More information to follow {{.Text}}","html":"<p>The site has gone down for maintenance.  More information to follow {{.HTML}}"}'
```
A template is made up of a human readable name, a subject, a text representation of the the template you are sending (for mail clients that don't support HTML), and an HTML version of the template.

Special attention and care should be paid to the variables that take this form `{{.}}`.  These variables will interpolate data provided in the send step below into the template before a notification is sent.  Data that can be inserted into a template during the send step includes, `{{.Text}}`, `{{.HTML}}`, and `{{.Subject}}`.

The result of this curl returns a unique template ID that can be used in subsequent calls to refer to your custom template - it will look similar to this:

`{"template-id": "E3710280-954B-4147-B7E2-AF5BF62772B5"}`

P.S. Darla can check all of the saved templates by curling

```
uaac curl https://notifications.darla.example.com/templates -X GET
```
To view a list of all templates you must have the `notifications_templates.read` scope.

## Associate custom template to your notification
Darla now wants to associate her custom template with the `system-going-down` notification.  Any notification that does not have a custom template applied, like her `system-up` notification, defaults to a system-provided template.

```
uaac curl https://notifications.darla.example.com/clients/notifications-admin/notifications/system-going-down/template \
-X PUT --data '{"template": "E3710280-954B-4147-B7E2-AF5BF62772B5"}'
```
Here, Darla has associated the `system-going-down` notification belonging to the `notifications-admin` client with the template ID `E3710280-954B-4147-B7E2-AF5BF62772B5`. This is the template id of the template we created in the previously step.

This action requires the `notifications.manage` scope.

## Send your notification to all users
Darla is ready to send her `system-going-down` notification to all users of the system.  She performs this curl and includes some other pertinent information that gets directly inserted into the template:

```
uaac curl https://notifications.darla.example.com/everyone -X POST --data \
'{"kind_id":"system-going-down","text":"The system is going down while we upgrade our storage","html":"<h1>THE SYSTEM IS DOWN</h1><p>The system is going down while we upgrade our storage</p>","subject":"Upgrade to Storage","reply_to":"no-reply@example.com"}'
```

The data included in the post body above gets interpolated into the variables we previously inserted into our created template (remember they had the special syntax similar to `{{.Text}}`).

Sending a critical notification requires the scope `critical_notifications.write` whereas sending a non-critical notification requires the scope of `notifications_write`.

Darla could have also chosen to send the above notification to one specific user, an email address or possibly just a particular space.  For more information on targeting your notification at particular audiences see our [api docs](https://github.com/cloudfoundry-incubator/notifications/blob/master/API.md).
