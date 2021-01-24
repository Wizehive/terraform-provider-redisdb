package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/mediocregopher/radix/v4"
)

func TestAccResourceHash_basic(t *testing.T) {
	hashKey := fmt.Sprintf("testacc:resource:hash:basic-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { preCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccResourceHashDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceHash_basic(hashKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("redisdb_hash.basic", "key", hashKey),
					resource.TestCheckResourceAttr("redisdb_hash.basic", "hash.%", "2"),
					resource.TestCheckResourceAttr("redisdb_hash.basic", "hash.key", "value"),
					resource.TestCheckResourceAttr("redisdb_hash.basic", "hash.foo", "bar"),
				),
			},
		},
	})
}

func testAccResourceHashDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "redisdb_hash" {
			continue
		}

		var exists bool
		err := testClient.Do(context.TODO(), radix.Cmd(&exists, "EXISTS", rs.Primary.ID))
		if err != nil {
			if exists {
				return fmt.Errorf("Hash still exists")
			}
		}

	}
	return nil
}

func testAccResourceHash_basic(key string) string {
	return fmt.Sprintf(`
		resource "redisdb_hash" "basic" {
			key = %q
			hash = {
				key = "value"
				foo = "bar"
			}
		}
	`, key)
}
