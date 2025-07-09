# edgegrid-proxy

`edgegrid-proxy` is a tool written in Go to proxy EdgeGrid API requests.  
This proxy adds an Authentication header for the Akamai OPEN EdgeGrid API.

## Usage:

```
Usage:
  edgegrid-proxy [OPTIONS]

Application Options:
  -a, --addr=          Proxy host address (default: 127.0.0.1:8080)
  -f, --file=          Location of EdgeGrid file (default: ~/.edgerc)
  -s, --section=       Section of EdgeGrid file (default: default)
      --key=           Account switch key [$EDGEGRID_ACCOUNT_KEY]
      --host=          EdgeGrid Host [$EDGEGRID_HOST]
      --client-token=  EdgeGrid ClientToken [$EDGEGRID_CLIENT_TOKEN]
      --client-secret= EdgeGrid ClientSecret [$EDGEGRID_CLIENT_SECRET]
      --access-token=  EdgeGrid AccessToken [$EDGEGRID_ACCESS_TOKEN]
      --tls-crt=       Proxy TLS/SSL certificate file path
      --tls-key=       Proxy TLS/SSL key file path

Help Options:
  -h, --help           Show this help message
```

```
$ edgegrid-proxy
2025/07/01 00:00:00 Starting EdgeGrid proxy on http://127.0.0.1:8080
```

## Send EdgeGrid API request

```
$ curl -s http://127.0.0.1:8080/identity-management/v3/user-profile | jq
{
  "uiIdentityId": "X-W-ABCDEF",
  "firstName": "John",
  "lastName": "Doe",
  "uiUserName": "jdoe",
  "email": "jdoe@akamai.com",
  "accountId": "1-2345",
  "phone": "+12345678910",
  "timeZone": "Asia/Tokyo",
  "lastLoginDate": "2025-07-01T00:00:00.000Z",
  "tfaEnabled": true,
  "additionalAuthentication": "TFA",
  "preferredLanguage": "English",
  "sessionTimeOut": 64800,
  "passwordExpiryDate": "2026-07-01T00:00:00.000Z",
  "mobilePhone": "+123456789101",
  "address": "1-1 Chiyoda",
  "city": "Tokyo",
  "zipCode": "100-0000",
  "country": "Japan",
  "jobTitle": "Engineer",
  "tfaConfigured": false,
  "additionalAuthenticationConfigured": false,
  "userStatus": "ACTIVE",
  "isLocked": false
}
```
