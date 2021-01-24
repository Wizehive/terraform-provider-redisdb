resource "redisdb_hash" "name" {
  key = "prefix:name" # Key can use redis namespaces with `:` delimeter
  hash = {
    key1    = "value1"
    key2    = "value2"
    "key:3" = "value3"
  }
}
