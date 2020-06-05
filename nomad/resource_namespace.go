package nomad

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNamespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNamespaceWrite,
		UpdateContext: resourceNamespaceWrite,
		DeleteContext: resourceNamespaceDelete,
		ReadContext:   resourceNamespaceRead,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Unique name for this namespace.",
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
			},

			"description": {
				Description: "Description for this namespace.",
				Optional:    true,
				Type:        schema.TypeString,
			},

			"quota": {
				Description: "Quota to set for this namespace.",
				Optional:    true,
				Type:        schema.TypeString,
			},
		},
	}
}

func resourceNamespaceWrite(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(ProviderConfig).client

	namespace := api.Namespace{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Quota:       d.Get("quota").(string),
	}

	log.Printf("[DEBUG] Upserting namespace %q", namespace.Name)
	_, err := client.Namespaces().Register(&namespace, nil)
	if err != nil {
		return diag.Errorf("error inserting namespace %q: %s", namespace.Name, err)
	}
	log.Printf("[DEBUG] Created namespace %q", namespace.Name)
	d.SetId(namespace.Name)

	return resourceNamespaceRead(ctx, d, meta)
}

func resourceNamespaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(ProviderConfig).client
	name := d.Id()

	log.Printf("[DEBUG] Deleting namespace %q", name)
	retries := 0
	for {
		var err diag.Diagnostics
		if name == api.DefaultNamespace {
			log.Printf("[DEBUG] Can't delete default namespace, clearing attributes instead")
			d.Set("description", "Default shared namespace")
			d.Set("quota", "")
			err = resourceNamespaceWrite(ctx, d, meta)
		} else {
			if _, er := client.Namespaces().Delete(name, nil); er != nil {
				err = diag.FromErr(er)
			}
		}

		if len(err) == 0 {
			break
		} else if retries < 10 {
			if strings.Contains(err[0].Summary, "has non-terminal jobs") {
				log.Printf("[WARN] could not delete namespace %q because of non-terminal jobs, will pause and retry", name)
				time.Sleep(5 * time.Second)
				retries++
				continue
			}
			return append(err, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("error deleting namespace %q.", name),
			})
		} else {
			return append(err, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("too many failures attempting to delete namespace %q", name),
			})
		}
	}

	if name == api.DefaultNamespace {
		log.Printf("[DEBUG] %s namespace reset", name)
	} else {
		log.Printf("[DEBUG] Deleted namespace %q", name)
	}

	return nil
}

func resourceNamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(ProviderConfig).client
	name := d.Id()

	log.Printf("[DEBUG] Reading namespace %q", name)
	namespace, _, err := client.Namespaces().Info(name, nil)
	if err != nil {
		// we have Exists, so no need to handle 404
		return diag.Errorf("error reading namespace %q: %s", name, err)
	}
	log.Printf("[DEBUG] Read namespace %q", name)

	d.Set("name", namespace.Name)
	d.Set("description", namespace.Description)
	d.Set("quota", namespace.Quota)

	return nil
}
