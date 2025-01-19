### Roundguard
This is a simple Proof of Concept of a roundrobin load balancing mechanism. I'm doing this as a part of a challenge. It wouldn't be a wise idea to use this on actual production, although the reverse proxy should work (but then why wouldn't you use different proxy like nginx or let aws handle it)

### Scope
For this problem, I decided to look at problem statement and write down scope of what will be included and what won't be. I also decided to add some things that seemed fun to write but may not necessarily be requirements.

For this project, I'm going to scope to the following:

#### In Scope
- Prometheus metrics will be exposed, but alerts won't be setup, there also won't be a prometheus scraper to scrape this.
- Structured logging that is capable of both JSON and rich text will be built.
- For the load balancing API, we'll try to build a reverse proxy.
- We will add a "liveness" check from the load balancing server to remove that endpoint from round-robin.

#### Out of scope
- No error budgets if a liveness check fails, it'll be immediately removed, but added as soon as it starts responding with 200 status code.
- Load balancer won't call `/ready` endpoint to try and detect if the application can start receiving traffic, we'll simply start the load balancer to keep this simple.
- Load balancer will NOT retry any failed request (network issues or otherwise)
- We will also keep the load tests out of scope as of now. While it's easy to write a quick k6 load testing script, I want to take that time to polish the things I do have in scope.

### Assumptions
- for echo server, we'll require the request body to be in JSON format (according to the problem statement), but it'd be fairly trivial to create something like a "reflect" endpoint that can respond with exactly what was sent (even if it was not JSON).
- since we'll have echo server and load balancing server, we'll accept only POST in the echo server, but any kind of request in the proxy (load balancing server)
- This should be capable of handling large amount of traffic (so the locks needs to be carefully used)


### How to set up
- First clone this repository locally
- Ensure you have latest golang compiler installed
- Once cloned, simply run `go mod tidy && go mod vendor` (This should first download all dependencies locally, and then move them to a `vendor` directory in the project)
- For this project I'll be commiting the generated mocks, but if you need to generate, simply run `go generate ./...`
- Start first echo server by running the command in new terminal: `go run cmd/main.go echo start --port 8001`. This will start 1 echo server that is listening on http://localhost:8001
- Start second echo server by running the command: `go run cmd/main.go echo start --port 8002`
- Start third echo server by running the command: `go run cmd/main.go echo start --port 8003`
- You can continue starting multiple servers like this, feel free to do as you wish.
- Now start the load balancer by running this command: `go run cmd/main.go lb start --hosts "http://localhost:8001" --hosts "http://localhost:8002" --hosts "http://localhost:8003" --port 3000 `

By now you should be seeing some logs that indicates that the server has started.
You should be able to fire some requests to `http://localhost:3000/reflect` (since proxy server is in port `3000`)
The response header by the name `X-Backend-Server` will indicate which instance is responding to your request.
You should also be able to go to the load balancer's /metrics endpoint which is designed to be scraped by prometheus. This should give us much more information on the health of the system.
There are some alerts that are recommended to be created:
- Any increase in `proxy_count` with the label `availability` that's set to `unavailable`, this indicates one of the service is down. If we get the alert, it's expected to check the logs of the load balancer, find out which instance has gone down, investigate why it's down and restore it.
- Any latency `request_duration_milliseconds` should also be monitored.

There are 3 configuration that you can provide using the environment variables, you can modify to your liking:
```bash
LOG_LEVEL=info LOG_FORMAT=json PROXY_CHECK_DURATION=5s go run cmd/main.go lb start --hosts "http://localhost:8001" --hosts "http://localhost:8002" --hosts "http://localhost:8003" --port 3000 
```

Now try killing one of the echo server that you started in the beginning, after waiting for 1 cycle, the load balancer shouldn't route any traffic to that instance anymore.

Once that's confirmed, try bringing the same instance back up, after it goes through the next check phase, it will begin routing to that instance again.

### Thought process
- I think it'll be easy if I built a package on this project that can apply round-robin scheduling on any given item.
- The round-robin package should be able to "replace" the existing items with new items (this will help us remove the node that has gone down automatically)
- We'll also need a way to take the items out of the load-balancer, for that I can try to write a balancing package.
- Let's try to build a rebalancing system because that might be tricky, echo server can be done next since it's very simple.
- Let's add some logging, for this the best performing library would be the zap logger from uber, but let's go with logrus since I want to check that out anyway.
- While we could do a simple loop, I want to use the lo package to keep the code more readable.
- Now that I think rebalancing and roundrobin is implemented to my liking, let's try to quickly whip up an echo server.
- For echo server, the problem statement says it's going to have to be JSON, so I'm going to add a JSON decoder and encoder.
- Now with the test passing, I'm fairly confident that this will work, next up is going to be writing a liveness handler.
- Now that we have a liveness handler, we can then start by creating a server.
- I think we now have a good enough server that's capable of responding to `POST /reflect` endpoint as well as `GET /live` endpoint.
- We can now start working on command line API that will actually start this server for us to make some API calls.
- The echo server works correctly, it's now time to go for the load balancing server.
- So the resulting command looks like this: `go run main.go echo start --port 8888` (for a compiled binary, it'd look like: `./roundguard echo start --port 8888`)
- I'm now able to start as many instance of echo server as I'd like. I'd also be able to specify any port I want.
- Adding a simple reverse proxy now.
- Now that the reverse proxy works, let's add some more logs first
