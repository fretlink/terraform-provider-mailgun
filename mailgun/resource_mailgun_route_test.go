package mailgun

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mailgun/mailgun-go"
	"strconv"
	"testing"
	"time"
)

func TestAccMailgunRoute_basic(t *testing.T) {
	var route mailgun.Route

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccRouteCheckDestroy(&route),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccRouteCheckExists("mailgun_route.exemple", &route),
					testAccRouteCheckAttributes("mailgun_route.exemple", &route),
				),
			},
		},
	})
}

func TestAccMailgunRoute_withUpdate(t *testing.T) {
	var route mailgun.Route

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccRouteCheckDestroy(&route),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccRouteCheckExists("mailgun_route.exemple", &route),
					testAccRouteCheckAttributes("mailgun_route.exemple", &route),
				),
			},

			{
				Config: testAccRouteConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccRouteCheckExists("mailgun_route.exemple", &route),
					testAccRouteCheckAttributes("mailgun_route.exemple", &route),
				),
			},
		},
	})
}

func testAccRouteCheckExists(rn string, route *mailgun.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("routeID not set")
		}

		mg := testAccProvider.Meta().(*mailgun.MailgunImpl)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		gotRoute, err := mg.GetRoute(ctx, rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error getting route: %s", err)
		}

		*route = gotRoute

		return nil
	}
}

func testAccRouteCheckAttributes(rn string, route *mailgun.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		attrs := s.RootModule().Resources[rn].Primary.Attributes

		check := func(key, stateValue, routeValue string) error {
			if routeValue != stateValue {
				return fmt.Errorf("different values for %s in state (%s) and in mailgun (%s)",
					key, stateValue, routeValue)
			}
			return nil
		}

		for key, value := range attrs {
			var err error

			switch key {
			case "priority":
				err = check(key, value, strconv.Itoa(route.Priority))
			case "description":
				err = check(key, value, route.Description)
			case "expression":
				err = check(key, value, route.Expression)
			case "created_at":
				err = check(key, value, route.CreatedAt.String())
			case "route_id":
				err = check(key, value, route.Id)
			case "actions":
				for _, k := range route.Actions {
					err = check(key, value, k)
					if err != nil {
						return err
					}
				}
			}
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccRouteCheckDestroy(route *mailgun.Route) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		mg := testAccProvider.Meta().(*mailgun.MailgunImpl)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()

		_, err := mg.GetRoute(ctx, route.Id)
		if err == nil {
			return fmt.Errorf("route still exists")
		}

		return nil
	}
}

const testAccRouteConfig_basic = `
resource "mailgun_route" "exemple" {
	priority=5
        description="ho ho hoh"
        expression="match_recipient(\".*@samples.mailgun.org\")"
        actions=[
          "forward(\"http://myhost.com/messages/\")",
          "stop()"
        ]
}
`

const testAccRouteConfig_update = `
resource "mailgun_route" "exemple" {
        priority=4
        description="ho ho hohf"
        expression="match_recipient(\".*@samples.mailgun.org\")"
        actions=[
          "forward(\"http://myhost.com/messages/\")",
          "stop()"
        ]
}
`
