# Go Rate Limiter

This project involves the development of a rate limiter in Go, which can be configured to limit the maximum number of requests per second based on a specific IP address or an access token. The goal of the limiter is to efficiently control web service traffic.

## Features

- **IP Address Limitation**: Restricts requests from a single IP address within a defined time interval.
- **Access Token Limitation**: Limits requests based on unique access tokens, allowing different expiration time limits for different tokens. Tokens must be provided in the header as API_KEY.
- **Overlapping Settings**: Access token settings take precedence over IP address settings.
- **Integration with Middleware**: Functions as a middleware injected into the web server.
- **Request Limit Configuration**: Allows setting the maximum number of requests per time.
- **Block Duration**: Option to define the duration of the block time for IP or Token after exceeding request limits.
- **Environment Variable Settings**: Limit settings can be made through environment variables or a .env file in the root folder.
- **HTTP Response When Limit Exceeded**: Responds with HTTP code 429 and a message indicating that the maximum number of requests has been exceeded.
- **Redis Database**: Uses Redis to store and query limiter information.
- **Flexible Persistence Strategy**: Strategy standard for easily switching between Redis and other persistence mechanisms.
- **Separate Limiter Logic**: The limiter logic is independent of the middleware.
- **Request Quantity Settings for IP or Token List**: Option to set an individual IP or Token its maximum number of requests per time.
- **Block Duration Settings for IP or Token List**: Option to set an individual IP or Token its block time.

## Usage Examples

- **IP Limitation Example**: If configured for a maximum of 5 requests per second per IP, the 6th request from IP 192.168.1.1 within a second should be blocked.
- **Access Token Limitation Example**: If a token `abc123` is set with a limit of 10 requests per second, the 11th request within that second should be blocked.
- **Expiration Time**: After reaching the limit, new requests from the same IP or token are only possible after the expiration time.
- **Custom IP Configuration**: Suppose `CUSTOM_MAX_REQ_PER_SEC` is set with "192.168.1.2=2", this means IP 192.168.1.2 will have a custom limit of 2 requests per time, regardless of the standard limit for other IPs.
- **Custom Token Configuration**: If `CUSTOM_MAX_REQ_PER_SEC` includes "token123=8", this indicates that the token `token123` has a custom limit of 8 requests per time, which may be different from the standard limit for other tokens.
- **Custom Block Duration**: Using `CUSTOM_BLOCK_DURATION` with "192.168.1.2=30s", IP 192.168.1.2 will be blocked for 30 seconds after exceeding its request limit, which is a specific setting different from the standard block.
- **Custom Block for Token**: For example, with "token123=1m" in `CUSTOM_BLOCK_DURATION`, the token `token123` will face a block of 1 minute after reaching its request limit.

## Environment Variables

The configuration of the rate limiter is managed through the following environment variables. Each plays a crucial role in controlling and customizing the behavior of the rate limiter:

- **REDIS_ADDRESS**: Defines the address of the Redis server used by the rate limiter.
- **REDIS_PASSWORD**: Password for authentication on the Redis server.
- **REDIS_DB**: Number of the Redis database to be used by the application.
- **DEFAULT_IP_MAX_REQ_PER_SEC**: Defines the standard limit of requests per second by IP address. This value is applied to all IPs unless a custom setting is specified.
- **DEFAULT_TOKEN_MAX_REQ_PER_SEC**: Establishes the standard limit of requests per second by access token. This limit is applied to all tokens except those with custom settings.
- **DEFAULT_IP_BLOCK_DURATION**: Standard block duration for an IP address that exceeds its request limit. Specified in a duration format, such as 10s for ten seconds.
- **DEFAULT_TOKEN_BLOCK_DURATION**: Standard block duration for a token that exceeds its request limit, specified in the same duration format as the IP block.
- **CUSTOM_MAX_REQ_PER_SEC**: Allows setting custom request limits for specific IPs or tokens. Expected format: `ip_or_token=value;another_ip_or_token=value`, for example, `127.0.0.1=2;abc123=10`.
- **CUSTOM_BLOCK_DURATION**: Sets custom block durations for specific IPs or tokens. The format is similar to CUSTOM_MAX_REQ_PER_SEC, for example, `127.0.0.1=30s;abc123=1m`.

These variables are fundamental to the flexibility and effectiveness of the limiter, allowing adaptation to specific needs. It's important to set these variables appropriately to ensure the system functions as expected.

## Rate Limit Operation

The Rate Limiter is implemented as a middleware in the HTTP server, allowing it to intercept and control requests.

### Rate Limiting Middleware

- The `RateLimiterMiddleware` middleware is applied to every request received by the server.
- It identifies each request by a unique identifier, which can be an access token (if present in the API_KEY header) or the requester's IP address.
- After identifying the request, the middleware consults the `RateLimit`, a structure that contains the limitation settings and Redis to check if the limit has been exceeded.

### Limit Check and Control

- The `IsLimitExceeded` function of the `RateLimit` structure is responsible for determining if a request exceeds the configured limit.
- It checks the current number of requests made by the identifier (IP or token) in Redis.
- Based on the identifier, the function determines the maximum number of allowed requests per second `maxReqPerSec` and the block duration `blockDuration`. These values can be customized for each IP or token or use the default values.
- If the number of requests already made is greater than or equal to the permitted limit, the function indicates that the limit has been exceeded.

### Response to Excessive Requests

- If the limit is exceeded, the middleware immediately responds with an HTTP 429 error "you have reached the maximum number of requests or actions allowed within a certain time frame", indicating to the client that the request limit has been reached.
- In case of internal errors (for example, failure to access storage), the middleware responds with an HTTP 500 error.

### Request Count Increment

- For each valid request that does not exceed the limit, the request counter in storage is incremented. This counter is used to track the number of requests made by the identifier within the specified time interval.

## Prerequisites

Make sure Docker is installed on your system. You can download and install Docker from the [official Docker website](https://www.docker.com/).

## Setup

- Clone the repository.
- Configure the `.env` file as needed.

## Running the Project

```bash
docker compose up
```

## Tests

This command performs an HTTP GET request without an access token:

```bash
curl http://localhost:8080/
```

To test with an access token, you must include the token in the request header. In the example below, we are using the token `abc123`, which is an example token. Replace `abc123` with your own access token when performing the test.

```bash
curl -H "API_KEY: abc123" http://localhost:8080/
```

To insert your own token in the tests, replace `abc123` in the curl command with the token you wish to test. The token should be passed in the request header using the API_KEY key. For example, if your token is `myToken123`, the curl command would be:

```bash
curl -H "API_KEY: myToken123" http://localhost:8080/
```
