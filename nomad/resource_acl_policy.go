package nomad

import (
	"context"
	"log"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceACLPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceACLPolicyCreate,
		UpdateContext: resourceACLPolicyUpdate,
		DeleteContext: resourceACLPolicyDelete,
		ReadContext:   resourceACLPolicyRead,

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

			"rules_hcl": {
				Description: "HCL or JSON representation of the rules to enforce on this policy. Use file() to specify a file as input.",
				Required:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceACLPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	policy := api.ACLPolicy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       d.Get("rules_hcl").(string),
	}

	// upsert our policy
	log.Printf("[DEBUG] Creating ACL policy %q", policy.Name)
	_, err := client.ACLPolicies().Upsert(&policy, nil)
	if err != nil {
		return diag.Errorf("error inserting ACLPolicy %q: %s", policy.Name, err)
	}
	log.Printf("[DEBUG] Created ACL policy %q", policy.Name)
	d.SetId(policy.Name)

	return resourceACLPolicyRead(ctx, d, meta)
}

func resourceACLPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	policy := api.ACLPolicy{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Rules:       d.Get("rules_hcl").(string),
	}

	// upsert our policy
	log.Printf("[DEBUG] Updating ACL policy %q", policy.Name)
	_, err := client.ACLPolicies().Upsert(&policy, nil)
	if err != nil {
		return diag.Errorf("error updating ACLPolicy %q: %s", policy.Name, err)
	}
	log.Printf("[DEBUG] Updated ACL policy %q", policy.Name)

	return resourceACLPolicyRead(ctx, d, meta)
}

func resourceACLPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client
	name := d.Id()

	// delete the policy
	log.Printf("[DEBUG] Deleting ACL policy %q", name)
	_, err := client.ACLPolicies().Delete(name, nil)
	if err != nil {
		return diag.Errorf("error deleting ACLPolicy %q: %s", name, err)
	}
	log.Printf("[DEBUG] Deleted ACL policy %q", name)

	return nil
}

func resourceACLPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client
	name := d.Id()

	// retrieve the policy
	log.Printf("[DEBUG] Reading ACL policy %q", name)
	policy, _, err := client.ACLPolicies().Info(name, nil)
	if err != nil {
		// we have Exists, so no need to handle 404
		return diag.Errorf("error reading ACLPolicy %q: %s", name, err)
	}
	log.Printf("[DEBUG] Read ACL policy %q", name)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("rules_hcl", policy.Rules)

	return nil
}
