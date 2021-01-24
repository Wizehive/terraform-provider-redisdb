package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mediocregopher/radix/v4"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"redisdb": func() (*schema.Provider, error) {
		return New("acctest")(), nil
	},
}

var testClient radix.Client

func TestMain(m *testing.M) {
	if os.Getenv("TF_ACC") == "" {
		// short circuit non acceptance test runs
		os.Exit(m.Run())
	}

	hostname := os.Getenv("REDISDB_HOSTNAME")
	port := os.Getenv("REDISDB_PORT")
	database := os.Getenv("REDISDB_DATABASE")

	cfg := radix.PoolConfig{
		Dialer: radix.Dialer{
			SelectDB: database,
		},
	}

	client, err := cfg.New(context.TODO(), "tcp", fmt.Sprintf("%s:%s", hostname, port))
	if err != nil {
		panic(err)
	}

	testClient = client

	resource.TestMain(m)
}

func TestProvider(t *testing.T) {
	if err := New("acctest")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func preCheck(t *testing.T) {
	variables := []string{
		"REDISDB_HOSTNAME",
		"REDISDB_PORT",
		"REDISDB_DATABASE",
	}

	for _, variable := range variables {
		value := os.Getenv(variable)
		if value == "" {
			t.Fatalf("`%s` must be set for acceptance tests!", variable)
		}
	}
}

func seedRedis(command string, args ...string) {
	testClient.Do(context.TODO(), radix.Cmd(nil, command, args...))
}
