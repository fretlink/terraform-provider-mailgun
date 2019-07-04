---
layout: "mailgun"
page_title: "Provider: Mailgun"
sidebar_current: "docs-mailgun-index"
description: |-
  The Mailgun provider configures domains and routes in Mailgun.
---

# Mailgun Provider

The Mailgun provider allows Terraform to create and configure domains and routes in [Mailgun](https://www.mailgun.com/).

The provider configuration block accepts the following arguments:

* ``domain`` - (Required) The domain name for the ressources created with the provider. May alternatively be set via the
  ``MAILGUN_DOMAIN`` environment variable.

* ``apikey`` - (Required) The API auth token to use when making requests. May alternatively
  be set via the ``MAILGUN_APIKEY`` environment variable.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
provider "mailgun" {
  domain = "domain.com"
  apikey   = "15ee99178cc7q6325df7ff8a15211228-2f778ta3-e04c2946"
}

resource "mailgun_domain" "example" {
      name="domain.com"
      spam_action="block"
      smtp_password="password"
      wildcard=true
      force_dkim_authority=true
      dkim_key_size=1024
      ips=["192.161.0.1", "192.168.0.2"]
      credentials{
           login="login"
           password="password"
      }
      open_tracking_settings_active=true
      click_tracking_settings_active=true
      unsubscribe_tracking_settings_active=true
      unsubscribe_tracking_settings_html_footer="<p>footer</p>"
      unsubscribe_tracking_settings_text_footer="footer"
      require_tls=true
      skip_verification=true
}

resource "mailgun_route" "example" {
        depends_on = [mailgun_domain.example]
        priority=5
        description="description"
        expression="match_recipient(\".*@samples.mailgun.org\")"
        actions=[
          "forward(\"http://myhost.com/messages/\")",
          "stop()"
        ]
}

```
