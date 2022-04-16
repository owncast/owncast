# Owncast Load Testing

Load testing is an important tool to surface bugs, race conditions, and to determine the overall server performance of Owncast.  The primary goal is to test the server components, not so much the front-end, as the performance of the browser and the client components may be variable.  While the test aims to push a bunch of traffic through the backend it's possible the frontend may not be able to handle it all.  Working on the performance of the frontend is a goal that should be treated separately from what the backend load tests are designed to do.

## What it will test

The test aims to reproduce the same requests and actions performed by a normal user joining an Owncast session, but at a rate that's faster than most normal environments.

## This includes

1. Downloads the configuration.
1. Registering as a brand new chat user.
1. Fetches the chat history.
1. Connects to the chat websocket and sends messages.
1. Access the viewer count `ping` endpoint to count this virtual user as a viewer.
1. Fetches the current status.

## Setup your environment

1. Install [k6](https://k6.io/open-source) by following [the instructions for your local machine](https://k6.io/docs/getting-started/installation/).
1. Start Owncast on your local machine, listening on `localhost:8080`.

## Run the tests

1. To monitor the concurrent chat users open the admin chat user's page at http://localhost:8080/admin/chat/users/.
1. To monitor the concurrent "viewers" open the admin viewers page at http://localhost:8080/admin/viewer-info/.
1. Begin the test suite by running `k6 run test.js`.


## troubleshooting

If you receive the error `ERRO[0080] dial tcp 127.0.0.1:8080: socket: too many open files` it means your OS is not configured to have enough concurrent sockets open to perform the level of testing the test is trying to accomplish.  

Using [ulimit](https://www.learnitguide.net/2015/07/how-to-increase-ulimit-values-in-linux.html) you can adjust this value.

Run `ulimit -n` and see what you're currently set at.  The default is likely `1024`, meaning your system can only open 1024 resources (files or sockets or anything) at the same time.  If you run `ulimit -Hn` you can get the _hard limit_ of your system.  So you can adjust your limit to be something between where you're at now and your hard limit.  `ulimit -n 10000` for example.

As a side note, Owncast automatically increases its own limit when you run the Owncast service to make it less likely that your Owncast server will hit this limit.