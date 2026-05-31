package app

import "github.com/redis/go-redis/v9"

var routeFailureScript = redis.NewScript(`
local count = redis.call("INCR", KEYS[1])
if count == 1 then
  redis.call("EXPIRE", KEYS[1], ARGV[2])
end
if count >= tonumber(ARGV[1]) then
  redis.call("SET", KEYS[2], "1", "EX", ARGV[3])
  redis.call("DEL", KEYS[1])
end
return count
`)
