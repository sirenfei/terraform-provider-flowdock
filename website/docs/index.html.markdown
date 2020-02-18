---
layout: "flowdock"
page_title: "Provider: Flowdock"
description: |-
  The Flowdock provider is used to interact with Flowdock organization resources.
---

# Flowdock Provider

The Flowdock provider is used to interact with Flowdock organization resources.

The provider allows you to manage your Flowdock organization's users and flows easily.
It needs to be configured with the proper 'API tokens' before it can be used.
we can get the tokens from https://www.flowdock.com/account/tokens

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the Flowdock Provider
provider "flowdock" {
  version = "1.1.2"
  api_token = "xxxxxx"
}
# Add a user to the organization

resource "flowdock_invitation" "richard_xue-duoduo" {
  # ...
}

```

## Argument Reference

The following arguments are supported in the `provider` block:

* `token` - (Optional) This is the Flowdock personal access token. It can also be
  sourced from the `FLOWDOCK_TOKEN` environment variable.