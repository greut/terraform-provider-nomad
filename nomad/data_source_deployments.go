package nomad

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDeployments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDeploymentsRead,
		Schema: map[string]*schema.Schema{

			"deployments": {
				Description: "Deployments",
				Computed:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeMap},
			},
		},
	}
}

func dataSourceDeploymentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConfig := meta.(ProviderConfig)
	client := providerConfig.client

	log.Printf("[DEBUG] Getting deployments...")
	deployment_list, _, err := client.Deployments().List(nil)
	if err != nil {
		// As of Nomad 0.4.1, the API client returns an error for 404
		// rather than a nil result, so we must check this way.
		if strings.Contains(err.Error(), "404") {
			return diag.FromErr(err)
		}

		return diag.Errorf("error checking for deployments: %#v", err)
	}

	var deployments []map[string]interface{}

	for _, deployment := range deployment_list {
		entry := make(map[string]interface{})
		entry["ID"] = deployment.ID
		entry["JobID"] = deployment.JobID
		entry["JobVersion"] = strconv.Itoa(int(deployment.JobVersion))
		entry["Status"] = deployment.Status
		entry["StatusDescription"] = deployment.StatusDescription
		deployments = append(deployments, entry)
	}

	d.SetId(client.Address() + "/deployments")

	if err := d.Set("deployments", deployments); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
