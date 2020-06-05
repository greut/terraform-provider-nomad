package nomad

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: regionsDataSourceRead,

		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func regionsDataSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(ProviderConfig).client

	log.Printf("[DEBUG] Reading regions from Nomad")
	regions, err := client.Regions().List()
	if err != nil {
		return diag.Errorf("error reading regions from Nomad: %s", err.Error())
	}
	log.Printf("[DEBUG] Read %d regions from Nomad", len(regions))
	d.SetId(client.Address() + "/regions")

	if err := d.Set("regions", regions); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
