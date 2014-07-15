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
