package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mediocregopher/radix/v4"
)

func resourceKey() *schema.Resource {
	return &schema.Resource{
		Description:   "`redisdb_hash` manages an individual key in redis.",
		CreateContext: resourceKeyWrite,
		ReadContext:   resourceKeyRead,
		UpdateContext: resourceKeyWrite,
		DeleteContext: resourceKeyDelete,
		Schema: map[string]*schema.Schema{
			"key": {
				Description: "The key name.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"value": {
				Description: "The value of the key.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	var diags diag.Diagnostics

	key := d.Get("key").(string)

	var value string
	if err := client.Do(ctx, radix.Cmd(&value, "GET", key)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("value", value); err != nil {
		return diag.FromErr(err)
	}

	d.Set("key", key)

	return diags
}

func resourceKeyWrite(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	key := d.Get("key").(string)
	value := d.Get("value").(string)

	err := client.Do(ctx, radix.FlatCmd(nil, "SET", key, value))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceKeyRead(ctx, d, m)
}

func resourceKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	var diags diag.Diagnostics

	key := d.Get("key")
	if err := client.Do(ctx, radix.FlatCmd(nil, "DEL", key)); err != nil {
		return diag.FromErr(err)
	}

	d.Set("key", "")

	return diags
}
