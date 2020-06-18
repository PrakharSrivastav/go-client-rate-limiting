# Client Side rate limiting

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