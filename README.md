# Toy Reverse Proxy

## Running
For development, start the postgres server and then source the test environment and then run the program.

```
docker compose up -d
source test.env
go run cmd/main.go
```

## Management Server

There is a single endpoint listening on port 9080 (`MANAGEMENT_PORT`) which listens for registration events.

```bash
curl http://localhost:9080 -Lk -d hostname=customer.com -d page_data_url=example
Status page registered with ID: 3
```

This will kick off the certificate issuance in the background.

## Testing

If you don't have a real dns behind this you need to either add to your hosts file or use `--resolve` in curl to get SNI to work:

```bash
curl --resolve customer.com:8443:127.0.0.1 https://customer.com:8443 -k -v    
```


## Let's Encrypt Integration
I haven't actually tested this, but the skeleton is in place to use `NewAcmeHttp01Issuer` with a letsencrypt account instead of the `SelfSignedIssuer`
