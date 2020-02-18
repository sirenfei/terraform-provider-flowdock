---
layout: "flowdock"
page_title: "Flowdock: flowdock_invitation"
description: |-
  Provides a Flowdock invitation resource.
---

# flowdock_user

Provides a Flowdock invitation resource.

This resource allows you to invite/remove users from your organization. When applied,
a new invitation will be lunched or the existing user's id will be added plus "u" as a prefix. When destroyed, that user will be removed from the organisation.

## Example Usage

```hcl
# Add a user to the organization
provider "flowdock" {
  version = "1.1.7"
  api_token = "xxxxxxxx"
}
resource "flowdock_invitation" "richard_mouse-duoduo" {
   org = "smart-mouse"
   flow = "ops-projects"
   email = "richard.mouse@gmail.com"
   message = "hello mickie mouse"
}
```

## Argument Reference

The following arguments are supported:

* `org` - (Required) The name of the organisation.
* `flow` - (Required) The name of the flow.
* `email` - (Required) The email of the user's.
* `message` - (Optional) A description of the invitation.
## Attributes Reference

The following attributes are exported:

* `id` - The ID of the created invitation.


## Import

Admin of the organisation can be import the users id using the command below e.g.

```
$ terraform import flowdock_invitation.resoueceInstaneName userId_indexfOFlow_flowName_orgName
$ terraform import flowdock_invitation.resoueceInstaneName userId_1_flowName_orgName
$ terraform import flowdock_invitation.resoueceInstaneName userId_2_flowName_orgName
```
