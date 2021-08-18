package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mediocregopher/radix/v4"
)

func dataSourceKey() *schema.Resource {
	return &schema.Resource{
		Description: "`redisdb_key` data source can be used to get a value of a specific key.",
		ReadContext: dataSourceKeyRead,
		Schema: map[string]*schema.Schema{
			"key": {
				Description: "The key name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"values": {
				Description: "The value of the Key.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	var diags diag.Diagnostics

	key := d.Get("key").(string)

	var value string
	if err := client.Do(ctx, radix.Cmd(&value, "GET", key)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("key", key)
	if err := d.Set("value", value); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
