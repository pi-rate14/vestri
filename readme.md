## Walkthrough

1. each request is made to loadbalancer running on localhost:8080/api/v1
2. request path (api/v1) is looked up to find the service it is mapped to
3. all the server replicas of this service are looked up
4. Load Balancing Strategies for this Service is looked up. Defaultes to Round Robin
5. request is forwarded to one of these replicas using the LB strategy

### TODOS:

1. Error Handling
2. Optimising WRR
3. Converting metadata to struct
4. Refactor unnecessaty packages
5. Configure period based on config file

### Potential Features

1. Add more LB Strategies
2. Dynamically update weights if a server goes down
3. If no server available, put requests in a queue for a timeout
