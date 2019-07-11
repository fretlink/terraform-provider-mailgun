package mailgun

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mailgun/mailgun-go/v3"
	"os"
	"strconv"
	"testing"
	"time"
	"log"
)

type fullDomain struct {
	domainResponse   mailgun.DomainResponse
	domainConnection mailgun.DomainConnection
	domainTracking   mailgun.DomainTracking
	ipAddress        []string
	credentials      []mailgun.Credential
}

func getFullDomain(mg *mailgun.MailgunImpl, domainName string) (*fullDomain, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	mg = mailgun.NewMailgun(domainName, mg.APIKey())

	var domain fullDomain
	var err error
	domain.domainResponse, err = mg.GetDomain(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("Error Getting mailgun domain Details for %s: Error: %s", domainName, err)
	}

	domain.domainConnection, err = mg.GetDomainConnection(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("Error Getting mailgun domain connection Details for %s: Error: %s", domainName, err)
	}

	domain.domainTracking, err = mg.GetDomainTracking(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("Error Getting mailgun domain tracking Details for %s: Error: %s", domainName, err)
	}

	ipAddress, err := getIps(ctx, mg)

	if err != nil {
		return nil, fmt.Errorf("Error Getting mailgun domain ips2 for %s: Error: %s", domainName, err)
	}
	ips := make([]string, len(ipAddress))
	for i, r := range ipAddress {
		ips[i] = r.IP

	}
	domain.ipAddress = ips
	domain.credentials, err = ListCredentials(domainName, mg.APIKey())
	if err != nil {
		return nil, fmt.Errorf("Error Getting mailgun credentials for %s: Error: %s", domainName, err)
	}
	return &domain, nil
}

func TestAccMailgunDomain_basic(t *testing.T) {
	var domain fullDomain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDomainCheckDestroy(&domain),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplateDomain(testAccDomainConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainCheckExists("mailgun_domain.exemple", &domain),
					testAccDomainCheckAttributes("mailgun_domain.exemple", &domain),
				),
			},
		},
	})
}

func TestAccMailgunDomain_withUpdate(t *testing.T) {
	var domain fullDomain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDomainCheckDestroy(&domain),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplateDomain(testAccDomainConfig_basic),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainCheckExists("mailgun_domain.exemple", &domain),
					testAccDomainCheckAttributes("mailgun_domain.exemple", &domain),
				),
			},

			{
				Config: interpolateTerraformTemplateDomain(testAccDomainConfig_update),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainCheckExists("mailgun_domain.exemple", &domain),
					testAccDomainCheckAttributes("mailgun_domain.exemple", &domain),
				),
			},
		},
	})
}


func TestDomain_importBasic(t *testing.T) {
	var domain fullDomain

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccDomainCheckDestroy(&domain),
		Steps: []resource.TestStep{
			{
				Config: interpolateTerraformTemplateDomain(testAccDomainConfig_import),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainCheckExists("mailgun_domain.exemple",&domain),
				),
			},
			{
				ResourceName:      "mailgun_domain.exemple",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDomainCheckExists(rn string, domain *fullDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("domainID not set")
		}

		mg := testAccProvider.Meta().(*mailgun.MailgunImpl)

		domainId := rs.Primary.ID

		gotDomain, err := getFullDomain(mg, domainId)
		if err != nil {
			return fmt.Errorf("error getting domain: %s", err)
		}

		*domain = *gotDomain

		return nil
	}
}

func testAccDomainCheckAttributes(rn string, domain *fullDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := s.RootModule().Resources[rn].Primary.Attributes

		check := func(key, stateValue, domainValue string) error {
			if domainValue != stateValue {
				return fmt.Errorf("different values for %s in state (%s) and in mailgun (%s)",
					key, stateValue, domainValue)
			}
			return nil
		}

		for key, value := range attrs {
			var err error

			switch key {
			case "name":
				err = check(key, value, domain.domainResponse.Domain.Name)
			case "smtp_password":
				err = check(key, value, domain.domainResponse.Domain.SMTPPassword)
			case "smtp_login":
				err = check(key, value, domain.domainResponse.Domain.SMTPLogin)
			case "wildcard":
				err = check(key, value, strconv.FormatBool(domain.domainResponse.Domain.Wildcard))
			case "state":
				err = check(key, value, domain.domainResponse.Domain.State)
			case "open_tracking_settings_active":
				err = check(key, value, strconv.FormatBool(domain.domainTracking.Open.Active))
			case "click_tracking_settings_active":
				err = check(key, value, strconv.FormatBool(domain.domainTracking.Click.Active))
			case "unsubscribe_tracking_settings_active":
				err = check(key, value, strconv.FormatBool(domain.domainTracking.Unsubscribe.Active))
			case "unsubscribe_tracking_settings_html_footer":
				err = check(key, value, domain.domainTracking.Unsubscribe.HTMLFooter)
			case "unsubscribe_tracking_settings_text_footer":
				err = check(key, value, domain.domainTracking.Unsubscribe.TextFooter)
			case "skip_verification":
				err = check(key, value, strconv.FormatBool(domain.domainConnection.SkipVerification))
			case "require_tls":
				err = check(key, value, strconv.FormatBool(domain.domainConnection.RequireTLS))
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccDomainCheckDestroy(domain *fullDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		mg := testAccProvider.Meta().(*mailgun.MailgunImpl)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		log.Printf("[DEBUG] try to fetch destroyed domain %s",mg.Domain())

		return resource.Retry(1*time.Minute, func() *resource.RetryError {
			_, err := mg.GetDomain(ctx, domain.domainResponse.Domain.Name)
			if err == nil {
				log.Printf("[DEBUG] managed to fetch destroyed domain %s",mg.Domain())
				return resource.RetryableError(err)
			}

			log.Printf("[DEBUG] failed to fetch destroyed domain %s",mg.Domain())

			return nil
		})
	}
}

func interpolateTerraformTemplateDomain(template string) string {
	domainName := ""

	if v := os.Getenv("MAILGUN_DOMAIN"); v != "" {
		domainName = v
	}

	return fmt.Sprintf(template, domainName)
}

const testAccDomainConfig_basic = `
resource "mailgun_domain" "exemple" {
	name="%s"
        wildcard=true
        credentials{
             login="aaaaaaa"
             password="adfshfjqdskjhgfksdgfkqgfk"
        }
}
`

const testAccDomainConfig_update = `
resource "mailgun_domain" "exemple" {
	name="%s"
        credentials{
             login="aaaaaaa"
             password="adfshfjqdskjhgfksdgfkqgfk"
        }

}
`
const testAccDomainConfig_import = `
resource "mailgun_domain" "exemple" {
	name = "%s"
}
`
