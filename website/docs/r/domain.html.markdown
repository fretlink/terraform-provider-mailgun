---
layout: "mailgun"
page_title: "Mailgun: mailgun_domain"
sidebar_current: "docs-mailgun-domain"
description: |-
  The domain_resource allows mailgun domain to be managed by Terraform.
---

# mailgun\_domain

The domain resource allows Mailgun domain to be managed by Terraform.

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the domain
* `spam_action` - (Optional) "disabled", "block", or "tag".If "disabled", no spam filtering will occur for inbound messages.If "block", inbound spam messages will not be delivered.If "tag", inbound messages will be tagged with a spam header. See Spam Filter.Defaults to disabled.
* `smtp_password` - (Optional) Password for SMTP authentication
* `wildcard` - (Optional) Determines whether the domain will accept email for sub-domains when sending messages.Defaults to false.
* `force_dkim_authority` - (Optional) If set to true, the domain will be the DKIM authority for itself even if the root domain is registered on the same mailgun account.If set to false, the domain will have the same DKIM authority as the root domain registered on the same mailgun account. Defaults to false
* `dkim_key_size` - (Optional) 1024 or 2048. Set the length of your domain’s generated DKIM key. Defaults to 1024.
* `ips` - (Optional) An optional, comma-separated list of IP addresses to be assigned to this domain. If not specified, all dedicated IP addresses on the account will be assigned. If the request cannot be fulfilled (e.g. a requested IP is not assigned to the account, etc), a 400 will be returned.
* `credentials` - (Optional) SMTP credentials for the domain
* `open_tracking_settings_active` - (Optional) true to enable open tracking. Defauls to false
* `click_tracking_settings_active` - (Optional) true to enable click tracking. Defauls to false
* `unsubscribe_tracking_settings_active` - (Optional) true to enable unsubscribe tracking. Defauls to false
* `unsubscribe_tracking_settings_html_footer` - (Optional)Custom HTML version of unsubscribe footer.Defaults to "\n<br>\n<p><a hre=\"%unsubscribe_url%\">unsubscribe</a></p>\n"
* `unsubscribe_tracking_settings_text_footer` - (Optional) Custom text version of unsubscribe footer. Defaults to "\n\nTo unsubscribe click: <%unsubscribe_url%>\n\n"
* `require_tls` - (Optional) If set to true, this requires the message only be sent over a TLS connection. If a TLS connection can not be established, Mailgun will not deliver the message.If set to false, Mailgun will still try and upgrade the connection, but if Mailgun cannot, the message will be delivered over a plaintext SMTP connection. Defaults to false.
* `skip_verification` - (Optional)If set to true, the certificate and hostname will not be verified when trying to establish a TLS connection and Mailgun will accept any certificate during delivery. If set to false, Mailgun will verify the certificate and hostname. If either one can not be verified, a TLS connection will not be established. Defaults to false.
The `credentials`  object supports the following:
* `login` - (Required) The user name
* `password` - (Required) A password for the SMTP credentials. (Length Min 5, Max 32)

## Attributes Reference

The following attribute is exported:

* `smtp_login` - An username for the SMTP credentials.
* `created_at` - The date of creation of the domain.
* `state` - The state of the domain.
* `receiving_records` - DNS records for receiving.
* `sending_records` - DNS records for sending.
The `receiving_records` `sending_records` and object exports the following:
* `name` - The name of the record.
* `priority` - The priority of the record lower value means a more important priority.
* `record_type` - The type of record.
* `valid` - Wether the record is valid or not.
* `value` - The value of the record.

## Import

Mailgun domain can be imported using the domain name, e.g.

```
tf import mailgun_domain.example domain.com

```
