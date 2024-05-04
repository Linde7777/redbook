local key = KEYS[1]
local currentTimeInUnixMilli = tonumber(ARGV[1])
local windowSizeInUnixMilli = tonumber(ARGV[2])
local threshold = tonumber(ARGV[3])

local windowStartInUnixMilli = currentTimeInUnixMilli - windowSizeInUnixMilli
redis.Call("ZREMRANGEBYSCORE", key, 0, windowStartInUnixMilli)

local reqCount = redis.call("ZCOUNT", key, windowStartInUnixMilli, currentTimeInUnixMilli)
if reqCount < threshold then
    redis.call("ZADD", key, currentTimeInUnixMilli, '')
    redis.call("EXPIRE", key, windowSizeInUnixMilli)
    return 'false'
else
    return 'true'
end
