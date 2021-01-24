package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceHash_basic(t *testing.T) {
	hashKey := "testacc:datasource:hash:basic"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			preCheck(t)
			seedRedis("HSET", hashKey, "foo", "bar", "lorem", "ipsum")
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceHash_basic(hashKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.redisdb_hash.basic", "key", hashKey),
					resource.TestCheckResourceAttr("data.redisdb_hash.basic", "hash.%", "2"),
					resource.TestCheckResourceAttr("data.redisdb_hash.basic", "hash.foo", "bar"),
					resource.TestCheckResourceAttr("data.redisdb_hash.basic", "hash.lorem", "ipsum"),
				),
			},
		},
	})
}

func testAccDataSourceHash_basic(key string) string {
	return fmt.Sprintf(`
	data "redisdb_hash" "basic" {
		key = %q
	}
	`, key)
}
