package main

import (
	"github.com/fretlink/terraform-provider-mailgun/mailgun"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: mailgun.Provider})
}
