package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mediocregopher/radix/v4"
)

func dataSourceHash() *schema.Resource {
	return &schema.Resource{
		Description: "`redisdb_hash` data source can be used to get a map of all the fields and values in a Hash.",
		ReadContext: dataSourceHashRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The ID (key) of the Hash.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"key": {
				Description: "The Hash key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"hash": {
				Description: "The fields and values of the Hash.",
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceHashRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	var diags diag.Diagnostics

	key := d.Get("key").(string)

	var hashMap map[string]string
	if err := client.Do(ctx, radix.Cmd(&hashMap, "HGETALL", key)); err != nil {
		return diag.FromErr(err)
	}

	if len(hashMap) == 0 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Hash not found with key %q", key),
		})
	} else {
		d.SetId(key)
		d.Set("key", key)
		if err := d.Set("hash", hashMap); err != nil {
			return diag.FromErr(err)
		}
	}

	return diags
}
