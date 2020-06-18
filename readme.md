# Client side rate limiting

This repository is to support [this](https://www.prakharsrivastav.com/posts/client-side-ratelimiting/
) blog post.

There are several integration scenarios where a middleware just acts as a proxy between two applications. This behaves good under low traffic conditions but might start spamming the target application under high traffic conditions. 

This setup is one of the several ways in which a middleware can throttle the requests to target applications. 

## Tools required. 

We use vegeta to send bulk requests to the middleware. Please install it using the instructions [here](https://github.com/tsenart/vegeta) .

On linux,  you can install using the following command:

```bash
wget https://github.com/tsenart/vegeta/releases/download/v12.8.3/vegeta-12.8.3-linux-amd64.tar.gz
tar xfz vegeta-12.8.3-linux-amd64.tar.gz

# this might require sudo login
mv vegeta /usr/bin/vegeta
```

## Running vegeta 

In the root of your application create a file called `target.list` with following contents
```bash
GET http://localhost:10000/root
Content-Type: application/json
```

where port 10000 belongs to your middleware. We will use port 10001 for the target server which we dont want to spam.

to run the vegeta attack run the below command 
```bash
vegeta attack -duration=1s -rate=20 -targets=target.list --output=resp.bin && vegeta report resp.bin
```

where rate is 20 requests per second, the attack runs for the duration 1 second.

You can try different combinations of attack parameters to validate the rate limit.


## Running application 

### Run the target server 
```bash
prakhar@tardis (master)✗ % go run target/main.go 
2020/06/18 19:58:43 starting target server
```

### Run the standard middleware
```bash
prakhar@tardis (master)✗ [1] % go run middleware/main.go 
2020/06/18 20:00:54 starting standard middleware
```

### Run the vegeta attack 
```bash
prakhar@tardis (master)✗ % vegeta attack -duration=2s -rate=10 -targets=target.list --output=resp.bin && vegeta report resp.bin 
Requests      [total, rate, throughput]         20, 10.53, 10.53
Duration      [total, attack, wait]             1.9s, 1.9s, 349.681µs
Latencies     [min, mean, 50, 90, 95, 99, max]  317.815µs, 404.47µs, 360.902µs, 518.276µs, 775.443µs, 943.369µs, 943.369µs
Bytes In      [total, mean]                     811, 40.55
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:20  
Error Set:
```

This will spam your server with a lot of requests. The more you increase the duration and rate, the more spammy middleware will be.

### Run the rate limited middleware

first setup the ratelimit for your application (check constant `rateLimit` in `middleware_rate_limited/main.go`), by default it will block every second.

```bash
prakhar@tardis (master)✗ [1] % go run middleware_rate_limited/main.go
2020/06/18 20:28:24 starting rate limited middleware
```

### Run the attack

you will observe different response from the target server based on duration, rate and rateLimit
- for duration = 2s, rate = 10 and rateLimit = time.Second ~ 20 seconds.
- for duration = 2s, rate = 10 and rateLimit = time.Second/2 ~ 10 seconds.
- for duration = 2s, rate = 10 and rateLimit = time.Second*2 ~ 40 seconds.

check the logs for `middleware_rate_limited/main.go` for timing details.

```bash
prakhar@tardis (master)✗ % vegeta attack -duration=2s -rate=10 -targets=target.list --output=resp.bin && vegeta report resp.bin                                          ~/Workspace/examples/blog.examples/workers/c_collect
Requests      [total, rate, throughput]         20, 10.53, 10.53
Duration      [total, attack, wait]             1.9s, 1.9s, 216.628µs
Latencies     [min, mean, 50, 90, 95, 99, max]  149.895µs, 236.448µs, 213.307µs, 282.508µs, 462.509µs, 633.837µs, 633.837µs
Bytes In      [total, mean]                     0, 0.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:20  
```