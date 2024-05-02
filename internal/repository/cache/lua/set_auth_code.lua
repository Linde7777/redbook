local key = KEYS[1]
local val = ARGV[1]
local ttl = tonumber(redis.call("ttl", key))


-- key存在，但无过期时间
if ttl == -1 then
    return 'err not expire'

    -- 当key不存在或者距离上次发送验证码已经超过一分钟，可以重新发送
elseif ttl == -2 or ttl < 540 then
    local cntKey = key..":cnt"
    redis.call("set", key, val)
    redis.call("expire", key, 600)

    -- 设置验证码验证次数为3，在verify_code.lua里面会用这个字段，每验证一次，这个字段的值减1，直到为0
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    return 'ok'

    -- 已经发送了一个验证码，但是还不到一分钟
else
    return 'err exceed send limit'
end