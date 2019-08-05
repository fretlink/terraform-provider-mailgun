package mailgun

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/mailgun/mailgun-go/v3"
	"log"
	"time"
)

func resourceMailgunDomain() *schema.Resource {
	return &schema.Resource{
		Create: CreateDomain,
		Update: UpdateDomain,
		Delete: DeleteDomain,
		Read:   ReadDomain,
		Importer: &schema.ResourceImporter{
			State: ImportStatePassthroughDomain,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"spam_action": &schema.Schema{
				Type:     schema.TypeString,
				Default:  "disabled",
				ForceNew: true,
				Optional: true,
			},

			"smtp_password": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},

			"smtp_login": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"wildcard": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				ForceNew: true,
				Optional: true,
			},

			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},

			"state": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"force_dkim_authority": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"dkim_key_size": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  1024,
				ForceNew: true,
				Optional: true,
			},

			"ips": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				ForceNew: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"credentials": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"login": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"password": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"open_tracking_settings_active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"click_tracking_settings_active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"unsubscribe_tracking_settings_active": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"unsubscribe_tracking_settings_html_footer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "\n<br>\n<p><a href=\"%unsubscribe_url%\">unsubscribe</a></p>\n",
			},
			"unsubscribe_tracking_settings_text_footer": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "\n\nTo unsubscribe click: <%unsubscribe_url%>\n\n",
			},

			"require_tls": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"skip_verification": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},

			"receiving_records": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"valid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"sending_records": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"valid": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func interfaceToStringTab(i interface{}) []string {
	aInterface := i.([]interface{})
	aString := make([]string, len(aInterface))
	for i, v := range aInterface {
		aString[i] = v.(string)
	}
	return aString
}
func CreateDomain(d *schema.ResourceData, meta interface{}) error {
	mg := meta.(*mailgun.MailgunImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log.Printf("[DEBUG] creating  mailgun domain: %s", d.Id())

	creationResponse, err := mg.CreateDomain(ctx, d.Get("name").(string), &mailgun.CreateDomainOptions{
		Password:           d.Get("smtp_password").(string),
		SpamAction:         mailgun.SpamAction(d.Get("spam_action").(string)),
		Wildcard:           d.Get("wildcard").(bool),
		ForceDKIMAuthority: d.Get("force_dkim_authority").(bool),
		DKIMKeySize:        d.Get("dkim_key_size").(int),
		IPS:                interfaceToStringTab(d.Get("ips")),
	})

	if err != nil {
		return fmt.Errorf("Error creating mailgun domain: %s", err.Error())
	}

	mg = mailgun.NewMailgun(creationResponse.Domain.Name, mg.APIKey())

	for _, i := range d.Get("credentials").([]interface{}) {
		credential := i.(map[string]interface{})
		err = mg.CreateCredential(ctx, credential["login"].(string), credential["password"].(string))
		if err != nil {
			return fmt.Errorf("Error creating mailgun credential: %s", err.Error())
		}
	}

	err = mg.UpdateUnsubscribeTracking(ctx, creationResponse.Domain.Name, boolToString(d.Get("unsubscribe_tracking_settings_active").(bool)), d.Get("unsubscribe_tracking_settings_html_footer").(string), d.Get("unsubscribe_tracking_settings_text_footer").(string))
	if err != nil {
		return fmt.Errorf("Error updating mailgun unsubscribe tracking settings: %s", err.Error())
	}

	err = mg.UpdateOpenTracking(ctx, creationResponse.Domain.Name, boolToString(d.Get("open_tracking_settings_active").(bool)))
	if err != nil {
		return fmt.Errorf("Error updating mailgun open tracking settings: %s", err.Error())
	}

	err = mg.UpdateClickTracking(ctx, creationResponse.Domain.Name, boolToString(d.Get("click_tracking_settings_active").(bool)))
	if err != nil {
		return fmt.Errorf("Error updating mailgun click tracking settings: %s", err.Error())
	}

	err = mg.UpdateDomainConnection(ctx, creationResponse.Domain.Name, mailgun.DomainConnection{RequireTLS: d.Get("require_tls").(bool), SkipVerification: d.Get("skip_verification").(bool)})
	if err != nil {
		return fmt.Errorf("Error updating mailgun connexion settings: %s", err.Error())
	}

	d.SetId(creationResponse.Domain.Name)

	return ReadDomain(d, meta)
}

func UpdateDomain(d *schema.ResourceData, meta interface{}) error {
	mg := meta.(*mailgun.MailgunImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	domainName := d.Get("name").(string)
	mg = mailgun.NewMailgun(domainName, mg.APIKey())

	log.Printf("[DEBUG] updating  mailgun domain: %s", d.Id())

	if d.HasChange("unsubscribe_tracking_settings_active") || d.HasChange("unsubscribe_tracking_settings_html_footer") || d.HasChange("unsubscribe_tracking_settings_text_footer") {
		err := mg.UpdateUnsubscribeTracking(ctx, domainName, boolToString(d.Get("unsubscribe_tracking_settings_active").(bool)), d.Get("unsubscribe_tracking_settings_html_footer").(string), d.Get("unsubscribe_tracking_settings_text_footer").(string))
		if err != nil {
			return fmt.Errorf("Error updating mailgun unsubscribe tracking settings: %s", err.Error())
		}
	}
	if d.HasChange("open_tracking_settings_active") {
		err := mg.UpdateOpenTracking(ctx, domainName, boolToString(d.Get("open_tracking_settings_active").(bool)))
		if err != nil {
			return fmt.Errorf("Error updating mailgun open tracking settings: %s", err.Error())
		}
	}

	if d.HasChange("click_tracking_settings_active") {
		err := mg.UpdateClickTracking(ctx, domainName, boolToString(d.Get("click_tracking_settings_active").(bool)))
		if err != nil {
			return fmt.Errorf("Error updating mailgun click tracking settings: %s", err.Error())
		}
	}

	if d.HasChange("require_tls") || d.HasChange("skip_verification") {
		err := mg.UpdateDomainConnection(ctx, domainName, mailgun.DomainConnection{RequireTLS: d.Get("require_tls").(bool), SkipVerification: d.Get("skip_verification").(bool)})
		if err != nil {
			return fmt.Errorf("Error updating mailgun connexion settings: %s", err.Error())
		}
	}

	if d.HasChange("credentials") {
		old, new := d.GetChange("credentials")
		for _, i := range old.([]interface{}) {
			oldCredential := i.(map[string]interface{})
			found := false
			for _, j := range new.([]interface{}) {
				newCredential := j.(map[string]interface{})
				if oldCredential["login"] == newCredential["login"] {
					found = true
					if oldCredential["password"] != newCredential["password"] && newCredential["password"] != "" {
						err := mg.ChangeCredentialPassword(ctx, oldCredential["login"].(string), newCredential["password"].(string))
						if err != nil {
							return fmt.Errorf("Error updating mailgun credential password: %s", err.Error())
						}
					}
					break
				}
			}
			if !found {
				err := mg.DeleteCredential(ctx, oldCredential["login"].(string))
				if err != nil {
					return fmt.Errorf("Error deleting mailgun credential : %s", err.Error())
				}
			}
		}

		for _, i := range new.([]interface{}) {
			newCredential := i.(map[string]interface{})
			found := false
			for _, j := range old.([]interface{}) {
				oldCredential := j.(map[string]interface{})
				if oldCredential["login"] == newCredential["login"] {
					found = true
					break
				}
			}
			if !found {
				err := mg.CreateCredential(ctx, newCredential["login"].(string), newCredential["password"].(string))
				if err != nil {
					return fmt.Errorf("Error creating  mailgun credential : %s", err.Error())
				}
			}
		}
	}

	return ReadDomain(d, meta)
}

func DeleteDomain(d *schema.ResourceData, meta interface{}) error {
	mg := meta.(*mailgun.MailgunImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	log.Printf("[DEBUG] Deleting mailgun domain: %s", d.Id())

	err := mg.DeleteDomain(ctx, d.Get("name").(string))

	return err
}

func ReadDomain(d *schema.ResourceData, meta interface{}) error {
	mg := meta.(*mailgun.MailgunImpl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()
	domainName := d.Id()
	mg = mailgun.NewMailgun(domainName, mg.APIKey())

	domainResponse, err := mg.GetDomain(ctx, domainName)
	if err != nil {
		return fmt.Errorf("Error Getting mailgun domain Details for %s: Error: %s", d.Id(), err)
	}
	d.Set("created_at", domainResponse.Domain.CreatedAt.String())
	d.Set("smtp_login", domainResponse.Domain.SMTPLogin)
	d.Set("name", domainResponse.Domain.Name)
	d.Set("smtp_password", domainResponse.Domain.SMTPPassword)
	d.Set("wildcard", domainResponse.Domain.Wildcard)
	d.Set("spam_action", domainResponse.Domain.SpamAction)
	d.Set("state", domainResponse.Domain.State)

	simpleReceivingRecords := make([]map[string]interface{}, len(domainResponse.ReceivingDNSRecords))
	for i, r := range domainResponse.ReceivingDNSRecords {
		simpleReceivingRecords[i] = make(map[string]interface{})
		simpleReceivingRecords[i]["priority"] = r.Priority
		simpleReceivingRecords[i]["name"] = r.Name
		simpleReceivingRecords[i]["valid"] = r.Valid
		simpleReceivingRecords[i]["value"] = r.Value
		simpleReceivingRecords[i]["record_type"] = r.RecordType
	}
	d.Set("receiving_records", simpleReceivingRecords)

	simpleSendingRecords := make([]map[string]interface{}, len(domainResponse.SendingDNSRecords))
	for i, r := range domainResponse.SendingDNSRecords {
		simpleSendingRecords[i] = make(map[string]interface{})
		simpleSendingRecords[i]["name"] = r.Name
		simpleSendingRecords[i]["priority"] = r.Priority
		simpleSendingRecords[i]["valid"] = r.Valid
		simpleSendingRecords[i]["value"] = r.Value
		simpleSendingRecords[i]["record_type"] = r.RecordType
	}
	d.Set("sending_records", simpleSendingRecords)

	domainConnection, err := mg.GetDomainConnection(ctx, domainName)
	if err != nil {
		return fmt.Errorf("Error Getting mailgun domain connection  Details for %s: Error: %s", d.Id(), err)
	}
	d.Set("require_tls", domainConnection.RequireTLS)
	d.Set("skip_verification", domainConnection.SkipVerification)

	domainTracking, err := mg.GetDomainTracking(ctx, domainName)
	if err != nil {
		return fmt.Errorf("Error Getting mailgun domain tracking Details for %s: Error: %s", d.Id(), err)
	}

	d.Set("open_tracking_settings_active", domainTracking.Open.Active)

	d.Set("click_tracking_settings_active", domainTracking.Click.Active)
	d.Set("unsubscribe_tracking_settings_active", domainTracking.Unsubscribe.Active)
	d.Set("unsubscribe_tracking_settings_html_footer", domainTracking.Unsubscribe.HTMLFooter)
	d.Set("unsubscribe_tracking_settings_text_footer", domainTracking.Unsubscribe.TextFooter)

	ipAddress, err := getIps(ctx, mg)

	if err != nil {
		return fmt.Errorf("Error Getting mailgun domain ips1 for %s: Error: %s", d.Id(), err)
	}
	ips := make([]string, len(ipAddress))
	for i, r := range ipAddress {
		ips[i] = r.IP

	}
	d.Set("ips", ips)

	credentialsResponse, err := ListCredentials(domainName, mg.APIKey())
	if err != nil {
		return fmt.Errorf("Error Getting mailgun credentials for %s: Error: %s", d.Id(), err)
	}

	credentials := make([]map[string]interface{}, len(credentialsResponse))
	credentialsConf := d.Get("credentials").([]interface{})
	for i, r := range credentialsResponse {
		credentials[i] = make(map[string]interface{})
		credentials[i]["created_at"] = r.CreatedAt.String()
		credentials[i]["login"] = r.Login
		for _, c:= range credentialsConf {
			conf:=c.(map[string]interface{})
			if conf["login"] == credentials[i]["login"] {
				credentials[i]["password"] = conf["password"]
			}
		}
	}

	d.Set("credentials", credentials)

	d.SetId(domainName)

	return nil
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func ListCredentials(domain, apiKey string) ([]mailgun.Credential, error) {
	mg := mailgun.NewMailgun(domain, apiKey)
	it := mg.ListCredentials(nil)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var page, result []mailgun.Credential
	for it.Next(ctx, &page) {
		result = append(result, page...)
	}

	if it.Err() != nil {
		return nil, it.Err()
	}
	return result, nil
}

func ImportStatePassthroughDomain(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if _, ok := d.GetOk("dkim_key_size"); !ok {
		d.Set("dkim_key_size", 1024)
	}

	if _, ok := d.GetOk("force_dkim_authority"); !ok {
		d.Set("force_dkim_authority", false)
	}
	return []*schema.ResourceData{d}, nil
}

func getIps(ctx context.Context,mg *mailgun.MailgunImpl) ([]mailgun.IPAddress, error){
	var ipAddress []mailgun.IPAddress
	log.Printf("[DEBUG] begin to fetch ips for %s",mg.Domain())
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		ipAddress, err = mg.ListDomainIPS(ctx)
		if err != nil {
			log.Printf("[DEBUG] failed to fetch ips for %s",mg.Domain())
			return resource.RetryableError(err)
		}
		log.Printf("[DEBUG] managed to fetch ips for %s",mg.Domain())

		return nil
	})
	return ipAddress, err
}
