package nomad

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNamespaces() *schema.Resource {
	return &schema.Resource{
		ReadContext: namespacesDataSourceRead,

		Schema: map[string]*schema.Schema{
			"namespaces": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func namespacesDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(ProviderConfig).client

	log.Printf("[DEBUG] Reading namespaces from Nomad")
	resp, _, err := client.Namespaces().List(nil)
	if err != nil {
		return diag.Errorf("error reading namespaces from Nomad: %s", err)
	}
	namespaces := make([]string, 0, len(resp))
	for _, v := range resp {
		namespaces = append(namespaces, v.Name)
	}
	log.Printf("[DEBUG] Read %d namespaces from Nomad", len(namespaces))
	d.SetId(client.Address() + "/namespaces")

	if err := d.Set("namespaces", namespaces); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
