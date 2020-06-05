package nomad

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSentinelPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSentinelPolicyWrite,
		UpdateContext: resourceSentinelPolicyWrite,
		DeleteContext: resourceSentinelPolicyDelete,
		ReadContext:   resourceSentinelPolicyRead,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Unique name for this policy.",
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
			},

			"description": {
				Description: "Description for this policy.",
				Optional:    true,
				Type:        schema.TypeString,
			},

			"scope": {
				Description:  "Specifies the scope for this policy. Only 'submit-job' is currently supported.",
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{"submit-job"}, false),
			},

			"enforcement_level": {
				Description: "Specifies the enforcement level of the policy.",
				Required:    true,
				Type:        schema.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"advisory",
					"hard-mandatory",
					"soft-mandatory",
				}, false),
			},

			"policy": {
				Description: "The Sentinel policy.",
				Required:    true,
				Type:        schema.TypeString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// TODO: this should probably parse the AST to avoid false positives
					return strings.TrimSpace(old) == strings.TrimSpace(new)
				},
			},
		},
	}
}

func resourceSentinelPolicyWrite(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(ProviderConfig).client

	policy := api.SentinelPolicy{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Scope:            d.Get("scope").(string),
		EnforcementLevel: d.Get("enforcement_level").(string),
		Policy:           d.Get("policy").(string),
	}

	log.Printf("[DEBUG] Creating Sentinel policy %q", policy.Name)
	_, err := client.SentinelPolicies().Upsert(&policy, nil)
	if err != nil {
		return diag.Errorf("error upserting Sentinel policy %q: %s", policy.Name, err)
	}
	log.Printf("[DEBUG] Upserted Sentinel policy %q", policy.Name)
	d.SetId(policy.Name)

	return resourceSentinelPolicyRead(ctx, d, meta)
}

func resourceSentinelPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(ProviderConfig).client
	name := d.Id()

	log.Printf("[DEBUG] Deleting Sentinel policy %q", name)
	_, err := client.SentinelPolicies().Delete(name, nil)
	if err != nil {
		return diag.Errorf("error deleting Sentinel policy %q: %s", name, err)
	}
	log.Printf("[DEBUG] Deleted Sentinel policy %q", name)

	return nil
}

func resourceSentinelPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(ProviderConfig).client
	name := d.Id()

	log.Printf("[DEBUG] Reading Sentinel policy %q", name)
	policy, _, err := client.SentinelPolicies().Info(name, nil)
	if err != nil {
		// we have Exists, so no need to handle 404
		return diag.Errorf("error reading Sentinel policy %q: %s", name, err)
	}
	log.Printf("[DEBUG] Read Sentinel policy %q", name)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("scope", policy.Scope)
	d.Set("enforcement_level", policy.EnforcementLevel)
	d.Set("policy", policy.Policy)

	return nil
}
