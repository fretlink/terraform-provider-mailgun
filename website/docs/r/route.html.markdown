---
layout: "mailgun"
page_title: "Mailgun: mailgun_route"
sidebar_current: "docs-mailgun-route"
description: |-
  The route_resource allows mailgun route to be managed by Terraform.
---

# mailgun\_route

The route resource allows Mailgun route to be managed by Terraform.

## Example Usage

```hcl
resource "mailgun_route" "example" {
        priority=5
        description="description"
        expression="match_recipient(\".*@samples.mailgun.org\")"
        actions=[
          "forward(\"http://myhost.com/messages/\")",
          "stop()"
        ]
}
```

## Argument Reference

The following arguments are supported:

* `priority` - (Required)Integer: smaller number indicates higher priority. Higher priority routes are handled first.
* `expression` - (Required) An arbitrary string.
* `description` - (Required) A filter expression like match_recipient('.*@gmail.com')
* `actions` - (Required) Route action. This action is executed when the expression evaluates to True. Example: forward("alice@example.com") You can pass multiple action parameters.


## Attributes Reference

The following attribute is exported:

* `route_id` - ID of the route.
* `created_at` - The date of creation of the route.

## Import

Mailgun  can be imported using the route ID, e.g.

```
tf import mailgun_route.example 4f3bad2335335426750048c6

```
