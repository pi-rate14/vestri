1. each request is made to loadbalancer running on localhost:8080/api/v1
2. request path (api/v1) is looked up to find the service it is mapped to
3. all the server replicas of this service are looked up
4. request is forwarded to one of these services using a load balancing technique (round robin)
