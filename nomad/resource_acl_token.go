package nomad

import (
	"context"
	"log"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceACLToken() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceACLTokenCreate,
		UpdateContext: resourceACLTokenUpdate,
		DeleteContext: resourceACLTokenDelete,
		ReadContext:   resourceACLTokenRead,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"accessor_id": {
				Description: "Nomad-generated ID for this token.",
				Computed:    true,
				Type:        schema.TypeString,
			},

			"secret_id": {
				Description: "The value that grants access to Nomad.",
				Computed:    true,
				Sensitive:   true,
				Type:        schema.TypeString,
			},

			"name": {
				Description: "Human-readable name for this token.",
				Optional:    true,
				Type:        schema.TypeString,
			},

			"type": {
				Description: "The type of token to create, 'client' or 'management'.",
				Required:    true,
				Type:        schema.TypeString,
			},

			"policies": {
				Description: "The ACL policies to associate with the token, if it's a 'client' type.",
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},

			"global": {
				Description: "Whether the token should be replicated to all regions or not.",
				Optional:    true,
				Type:        schema.TypeBool,
				ForceNew:    true,
				Default:     false,
			},

			"create_time": {
				Description: "The timestamp the token was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceACLTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	policies := make([]string, 0, len(d.Get("policies").(*schema.Set).List()))
	for _, pol := range d.Get("policies").(*schema.Set).List() {
		policies = append(policies, pol.(string))
	}

	token := api.ACLToken{
		Name:     d.Get("name").(string),
		Type:     d.Get("type").(string),
		Policies: policies,
		Global:   d.Get("global").(bool),
	}

	// create our token
	log.Println("[DEBUG] Creating ACL token")
	resp, _, err := client.ACLTokens().Create(&token, nil)
	if err != nil {
		return diag.Errorf("error creating ACL token: %s", err)
	}
	log.Printf("[DEBUG] Created ACL token %q", resp.AccessorID)
	d.SetId(resp.AccessorID)

	d.Set("accessor_id", resp.AccessorID)
	d.Set("secret_id", resp.SecretID)
	d.Set("name", resp.Name)
	d.Set("type", resp.Type)
	d.Set("policies", resp.Policies)
	d.Set("global", resp.Global)
	d.Set("create_time", resp.CreateTime.UTC().String())

	return resourceACLTokenRead(ctx, d, meta)
}

func resourceACLTokenUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	policies := make([]string, 0, len(d.Get("policies").(*schema.Set).List()))
	for _, pol := range d.Get("policies").(*schema.Set).List() {
		policies = append(policies, pol.(string))
	}

	token := api.ACLToken{
		AccessorID: d.Id(),
		Name:       d.Get("name").(string),
		Type:       d.Get("type").(string),
		Policies:   policies,
	}

	// update the token
	log.Printf("[DEBUG] Updating ACL token %q", d.Id())
	_, _, err := client.ACLTokens().Update(&token, nil)
	if err != nil {
		return diag.Errorf("error updating ACL token %q. %s", d.Id(), err)
	}
	log.Printf("[DEBUG] Updated ACL token %q", d.Id())

	return resourceACLTokenRead(ctx, d, meta)
}

func resourceACLTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client
	accessor := d.Id()

	// delete the token
	log.Printf("[DEBUG] Deleting ACL token %q", accessor)
	_, err := client.ACLTokens().Delete(accessor, nil)
	if err != nil {
		return diag.Errorf("error deleting ACL token %q: %s", accessor, err)
	}
	log.Printf("[DEBUG] Deleted ACL token %q", accessor)

	return nil
}

func resourceACLTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client
	accessor := d.Id()

	// retrieve the token
	log.Printf("[DEBUG] Reading ACL token %q", accessor)
	token, _, err := client.ACLTokens().Info(accessor, nil)
	if err != nil {
		// we have Exists, so no need to handle 404
		return diag.Errorf("error reading ACL token %q: %s", accessor, err)
	}
	log.Printf("[DEBUG] Read ACL token %q", accessor)

	d.Set("name", token.Name)
	d.Set("type", token.Type)
	d.Set("policies", token.Policies)
	d.Set("accessor_id", token.AccessorID)
	d.Set("secret_id", token.SecretID)
	d.Set("global", token.Global)
	d.Set("create_time", token.CreateTime.UTC().String())

	return nil
}
