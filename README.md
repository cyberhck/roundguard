### Roundguard
This is a simple Proof of Concept of a roundrobin load balancing mechanism. This is done as a part of tech assessment. It wouldn't be a wise idea to use this on actual production, although the reverse proxy should work (but then why wouldn't you use different proxy like nginx or let aws handle it)

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

### Thought process
- I think it'll be easy if I built a package on this project that can apply round-robin scheduling on any given item.
- The round-robin package should be able to "replace" the existing items with new items (this will help us remove the node that has gone down automatically)
- We'll also need a way to take the items out of the load-balancer, for that I can try to write a balancing package.
