# opentracing-go-redis

<p align="left">
<a href="https://hits.seeyoufarm.com"/><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Fopentracing-go-redis"/></a>
<a href="https://goreportcard.com/report/github.com/gjbae1212/opentracing-go-redis"><img src="https://goreportcard.com/badge/github.com/gjbae1212/opentracing-go-redis" alt="Go Report Card" /></a> 
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-GREEN.svg" alt="license" /></a>
</p>

## OVERVIEW
[OpenTracing](http://opentracing.io/) before and after hook for [go-redis v7](https://github.com/go-redis/redis).
It's to support **redis.Client**, **redis.ClusterClient**, **redis.Ring**.

## HOW TO USE
```
// select *redis.ClusterClient or *redis.Ring or *redis.Client  
var client redis.UniversalClient

// wrap redis client
ctx := context.Background()
client := otredis.WrapClient(ctx, client)

// call redis command
client.Get("test")
```

## LICENSE
This project is following The MIT.
