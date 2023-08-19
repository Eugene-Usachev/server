# Use norelation DBMS only to scale it horizontally
Aerospike (with pack in one index a lot of messages) or ScyllaDB for messages;
Aerospike or Tarantool or Redis (persistent mode) for users profiles;
Aerospike (with pack a lot of posts) or ScyllaDB for posts;
Aerospike (with pack a lot of comments) or ScyllaDB for comments;
Aerospike or Tarantool or Redis (persistent mode) for chats;

### Why I didn't do it? Because I was poor schoolboy who has no money for production server, and I use poor VDS that can use only Redis and Postgresql.

# Use microservices
Microservices can improve readability and speed of writing tests and a code.

### Why I didn't do it? Because I was a poor schoolboy and couldn't afford an expensive VDS, so I had to get the most out of my VDS, and the microservice architecture can hurt performance.

# Use Rust instead of Go for new services
Rust don't have GC and faster than Go, so using Rust can improve the performance

# Use nginx for sharing static files
nginx is better for sharing static files, because it can improve the performance

### Why I didn't do it? Because I was a poor schoolboy and couldn't afford an expensive VDS, so I use VDS that does not provide nginx.