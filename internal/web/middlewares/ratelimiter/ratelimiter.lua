local key = KEYS[1]
local now = tonumber(ARGV[1])
local windowSize = tonumber(ARGV[2])
local threshold = tonumber(ARGV[3])

local windowStart = now - windowSize
redis.Call("ZREMRANGEBYSCORE", key, 0, windowStart)

local reqCount = redis.call("ZCOUNT", key, windowStart, now)
if reqCount < threshold then
    redis.call("ZADD", key, now, '')
    redis.call("EXPIRE", key, windowSize)
    return 'false'
else
    return 'true'
end
