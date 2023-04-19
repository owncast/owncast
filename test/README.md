# Tests

## Load Tests

1. Install [artillery](https://artillery.io/) from NPM/Yarn/Whatever Javascript package manager is popular this week.
1. Start an instance of the server on localhost.
1. `artillery run httpGetTest.yaml` for endpoint load tests.
1. `artillery run websocketTest.yaml` for websocket load tests.


## Chat test

This will send automated fake chat messages to your localhost instance.
Edit the messages, usernames or point to a different instance.

1. `npm install`
1. `node fakeChat.js`

## Public Testing

Run `./test-local.sh` and it'll create a public URL that you can access your local Owncast instance from. This is particularly useful for testing mobile and other external devices, as well as webhooks. Make sure Owncast is running under port 8080.

If you'd like your own custom hostname that uses your username follow the instructions printed, otherwise use auto-generated name printed to the console for testing.

```
$ ./test/test-local.sh
Please wait. Making the local server port 8080 available at https://oc-gabek-develop.serveo.net
Forwarding HTTP traffic from https://oc-gabek-develop.serveo.net
```
