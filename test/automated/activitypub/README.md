# ActivityPub Federation Test

This test verifies Owncast's ActivityPub federation by having snac2 users follow the Owncast instance and confirming message delivery.

All test infrastructure (snac2, Caddy, mkcert, Go) runs inside a Docker container so you don't need to install anything on the host besides Docker.

## Prerequisites

- Docker installed and running

## Running the Tests

```bash
# Run the federation test with default 100 users
./run.sh

# Run with fewer users for quick testing
USER_COUNT=10 ./run.sh

# Run the follower validation test
./run.sh test-follower-validation.sh

# Run the featured streams test (two Owncast instances, no snac2)
./run.sh test-featured-streams.sh

# Run the user-authentication tests
./run.sh test-fediverse-otp.sh   # Fediverse OTP login (uses snac2)
./run.sh test-indieauth.sh       # IndieAuth login (uses a fake provider, no snac2)

# Keep servers running after test for debugging
KEEP_RUNNING=true ./run.sh

# Adjust follow request throttling (default 0.1s)
FOLLOW_DELAY=0.2 ./run.sh
```

## Configuration Options

| Variable | Default | Description |
|----------|---------|-------------|
| `USER_COUNT` | 100 | Number of test users to create |
| `FOLLOW_DELAY` | 0.1 | Delay in seconds between follow requests |
| `KEEP_RUNNING` | false | Keep servers running after test for debugging |
| `CI` | false | Always true inside the container |
| `PROXY_PORT` | 8443 | HTTPS proxy port |
| `SNAC_PORT` | 9080 | snac2 HTTP port |
| `OWNCAST_PORT` | 8080 | Owncast HTTP port |
| `OWNCAST2_PORT` | 8081 | Second Owncast HTTP port (featured streams test only) |

## What the Test Does

### `test-federation.sh` / `test-following.sh` (default, snac2-based)

1. Creates a temporary snac2 instance with test users
2. Starts an HTTPS reverse proxy for TLS termination
3. Starts Owncast configured for federation
4. Has all snac2 users follow Owncast
5. Sends a message from Owncast
6. Verifies all followers received the message

### `test-featured-streams.sh` (two Owncast instances)

The featured-streams "mini-directory" is an Owncast-to-Owncast feature, so
this test runs two Owncast instances (`owncast.local` and `owncast2.local`)
behind the shared proxy. It:

1. Starts both instances with federation enabled and public
2. Has instance 1 add instance 2 via `POST /api/admin/federation/servers`
3. Verifies instance 2 is persisted in instance 1's directory immediately
   (a pending follow record)
4. Verifies the follow transitions to `accepted` once instance 2 returns its
   ActivityPub `Accept`
5. Verifies the listing is readable on the public, unauthenticated endpoint
6. Verifies the reverse direction (instance 2 adds instance 1) also works

### `test-fediverse-otp.sh` (Fediverse authentication, snac2-based)

End-to-end test of logging in with a Fediverse account, with no backend stubs.
It:

1. Registers an anonymous chat user on Owncast (gets an access token)
2. Requests a Fediverse OTP for a snac2 account; Owncast delivers the code as a
   real ActivityPub direct message (webfinger + signed inbox POST)
3. Reads the delivered code back out of snac2's stored DM
4. Submits the code to Owncast's verify endpoint
5. Verifies the admin users API now reports the user as authenticated with a
   `Fediverse` provider

### `test-indieauth.sh` (IndieAuth authentication, fake provider, no snac2)

End-to-end test of logging in with IndieAuth against a fake provider hosted by
the Caddy proxy at `indieauth.local` (see `Caddyfile.indieauth`). The provider
serves a discovery document, auto-approves the authorization request, and
returns a canonical `me` on code exchange. The test:

1. Registers an anonymous chat user on Owncast
2. Starts the IndieAuth flow pointed at the fake provider; Owncast discovers its
   authorization endpoint and returns a redirect URL
3. Follows the redirect like a browser would, through the provider and back to
   Owncast's callback, where Owncast exchanges the code and links the identity
4. Verifies the admin users API now reports the user as authenticated with an
   `IndieAuth` provider

Neither auth test changes the backend: `OWNCAST_ALLOW_INTERNAL_FEDERATION` lets
Owncast accept the local test hosts, and the mkcert CA makes its outbound HTTPS
calls trust the test certificates.

## Test Results

The test reports:
- **Followers Registered**: Number of successful follow requests
- **Messages Delivered**: Number of users who received the message
- **Delivery Time**: Time to deliver to all followers
- **Follow Success Rate**: Percentage of follow requests that succeeded
- **Delivery Rate**: Percentage of registered followers who received the message

## Docker Image Details

The Docker image (`owncast-ap-test`) bundles all dependencies:
- Go (for building Owncast)
- snac2 (built from source)
- Caddy (HTTPS reverse proxy)
- mkcert (TLS certificates trusted by the container)
- sqlite3, jq, curl

Go module and build caches are stored in named Docker volumes (`owncast-ap-test-gomod`, `owncast-ap-test-gobuild`) so repeated runs are faster.

## Troubleshooting

### Docker build fails

Make sure Docker is running. On macOS, Docker Desktop or a compatible runtime (colima, OrbStack, etc.) is required.

### Port already in use

If a previous container didn't shut down cleanly:
```bash
docker ps -a | grep owncast-ap-test
docker rm -f <container_id>
```

### Cleaning up Docker resources

```bash
# Remove the image
docker rmi owncast-ap-test

# Remove Go caches
docker volume rm owncast-ap-test-gomod owncast-ap-test-gobuild
```
