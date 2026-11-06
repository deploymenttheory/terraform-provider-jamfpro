# Reverse-proxy validation for issue #1162

Validated against a real Jamf Pro 11.33.3 lab using a local reverse proxy. The proxy required Basic credentials in `Authorization`, accepted Jamf's current bearer token in `X-Jamf-Token`, removed `/proxypath` before forwarding, and forwarded the remaining request to Jamf over HTTPS.

Direct and proxied curl calls returned matching version and complete smart-group/Classic-category JSON responses. A real Terraform plan/apply/refreshed-plan used a Framework smart-group data source and an SDKv2 category data source. Values matched direct curl reads; the refreshed plan returned exit code 0. No managed Jamf resources were created or changed.

The proxy shortened only the OAuth response's expiry metadata to one second. Both provider implementations obtained fresh real Jamf tokens after that interval and continued reading the real version endpoint successfully. Provider tests also cover invalid/duplicate headers, token-header collisions, sensitive schema metadata, default behavior and shared SDKv2/Framework schema compatibility. HTTP library tests cover request/config immutability, custom TLS transports, URL proxy credentials, cross-origin redirect blocking and sensitive logging.

The live provider test is opt-in:

```shell
JAMFPRO_PROXY_LIVE_CONFIG=/path/to/private/config.json \
  go test -race -tags=integration ./internal/provider -run TestLiveProxyHeaders -count=1
```

The disposable test proxy must expose `/proxypath`, use `X-Jamf-Token` and shorten expiry metadata as described above. Its private configuration file supplies:

```json
{
  "proxy_url": "http://127.0.0.1:PORT",
  "client_id": "LAB_CLIENT_ID",
  "client_secret": "LAB_CLIENT_SECRET",
  "proxy_authorization": "Basic BASE64_PROXY_CREDENTIALS",
  "expected_version": "EXACT_VERSION_FROM_DIRECT_CURL"
}
```

Only sanitized comparison results are retained in `results.json`. Credentials, OAuth responses, Terraform state/plan files and the disposable proxy are cleaned up after testing.
