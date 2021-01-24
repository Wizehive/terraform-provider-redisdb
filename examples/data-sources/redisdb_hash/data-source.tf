# Read Hash fields/values by key
data "redisdb_hash" "name" {
  key = "hash_key"
}
