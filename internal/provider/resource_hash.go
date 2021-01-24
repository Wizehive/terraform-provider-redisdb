package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/mediocregopher/radix/v4"
	"github.com/r3labs/diff/v2"
)

func resourceHash() *schema.Resource {
	return &schema.Resource{
		Description:   "`redisdb_hash` manages an individual Hash in redis.",
		CreateContext: resourceHashCreate,
		ReadContext:   resourceHashRead,
		UpdateContext: resourceHashUpdate,
		DeleteContext: resourceHashDelete,
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
				ForceNew:    true,
			},
			"hash": {
				Description: "The fields and values of the Hash.",
				Type:        schema.TypeMap,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceHashCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	key := d.Get("key").(string)
	hash := d.Get("hash").(map[string]interface{})

	err := client.Do(ctx, radix.FlatCmd(nil, "HSET", key, hash))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(key)

	return resourceHashRead(ctx, d, m)
}

func resourceHashRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	var diags diag.Diagnostics

	id := d.Id()

	var hashMap map[string]string
	if err := client.Do(ctx, radix.Cmd(&hashMap, "HGETALL", id)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("hash", hashMap); err != nil {
		return diag.FromErr(err)
	}

	d.Set("key", id)

	return diags
}

func resourceHashUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	id := d.Id()
	if d.HasChanges("hash") {
		err := client.Do(ctx, radix.WithConn(id, func(ctx context.Context, c radix.Conn) error {
			if err := c.Do(ctx, radix.Cmd(nil, "MULTI")); err != nil {
				return err
			}

			var err error
			defer func() {
				if err != nil {
					c.Do(ctx, radix.Cmd(nil, "DISCARD"))
				}
			}()

			changes, _ := diff.Diff(d.GetChange("hash"))
			for _, change := range changes {
				switch change.Type {
				case "create":
					err = c.Do(ctx, radix.Cmd(nil, "HSET", id, change.Path[0], change.To.(string)))
				case "update":
					err = c.Do(ctx, radix.Cmd(nil, "HSET", id, change.Path[0], change.To.(string)))
				case "delete":
					err = c.Do(ctx, radix.Cmd(nil, "HDEL", id, change.Path[0]))
				}
			}

			return c.Do(ctx, radix.Cmd(nil, "EXEC"))
		}))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceHashRead(ctx, d, m)
}

func resourceHashDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(radix.Client)

	var diags diag.Diagnostics

	id := d.Id()
	if err := client.Do(ctx, radix.Cmd(nil, "DEL", id)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return diags
}
