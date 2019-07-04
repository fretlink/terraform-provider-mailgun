package mailgun

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mailgun/mailgun-go/v3"
	"log"
	"time"
)

func resourceMailgunRoute() *schema.Resource {
	return &schema.Resource{
		Create: CreateRoute,
		Update: UpdateRoute,
		Delete: DeleteRoute,
		Read:   ReadRoute,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"route_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"actions": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func CreateRoute(d *schema.ResourceData, meta interface{}) error {
	mg := meta.(*mailgun.MailgunImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log.Printf("[DEBUG] creating  mailgun route: %s", d.Id())

	creationResponse, err := mg.CreateRoute(ctx, mailgun.Route{
		Priority:    d.Get("priority").(int),
		Description: d.Get("description").(string),
		Expression:  d.Get("expression").(string),
		Actions:     interfaceToStringTab(d.Get("actions")),
	})

	if err != nil {
		return fmt.Errorf("Error creating mailgun route: %s", err.Error())
	}

	d.SetId(creationResponse.Id)
	return ReadRoute(d, meta)
}

func UpdateRoute(d *schema.ResourceData, meta interface{}) error {
	mg := meta.(*mailgun.MailgunImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log.Printf("[DEBUG] updating  mailgun route: %s", d.Id())

	_, err := mg.UpdateRoute(ctx, d.Id(), mailgun.Route{
		Priority:    d.Get("priority").(int),
		Description: d.Get("description").(string),
		Expression:  d.Get("expression").(string),
		Actions:     interfaceToStringTab(d.Get("actions")),
	})

	if err != nil {
		return fmt.Errorf("Error updating mailgun route: %s", err.Error())
	}

	return ReadRoute(d, meta)
}

func DeleteRoute(d *schema.ResourceData, meta interface{}) error {
	mg := meta.(*mailgun.MailgunImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log.Printf("[DEBUG] Deleting mailgun route: %s", d.Id())

	err := mg.DeleteRoute(ctx, d.Id())

	return err
}

func ReadRoute(d *schema.ResourceData, meta interface{}) error {
	mg := meta.(*mailgun.MailgunImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	route, err := mg.GetRoute(ctx, d.Id())

	if err != nil {
		return fmt.Errorf("Error Getting mailgun route Details for %s: Error: %s", d.Id(), err)
	}

	d.Set("priority", route.Priority)
	d.Set("description", route.Description)
	d.Set("expression", route.Expression)
	d.Set("actions", route.Actions)
	d.Set("created_at", route.CreatedAt)
	d.Set("route_id", route.Id)

	d.SetId(route.Id)

	return nil
}
