local key = KEYS[1]
local cntKey = key .. ":cnt"
local expectedCode = ARGV[1]

local cnt = tonumber(redis.call("get", cntKey))
local code = redis.call("get", key)
if cnt <= 0 then
    return 300
end

-- 即使匹配上了，也不能立刻删除验证码，攻击者匹配成功后可能会再次请求发送验证码
if code == expectedCode then
    -- 标记为-1，表示验证码不可用
    redis.call("set", cntKey, -1)
    return 0
else
    redis.call("decr", cntKey, -1)
    return 400
end