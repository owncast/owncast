# Automated browser tests

The tests currently address the following interfaces:

1. The main web frontend of Owncast
1. The embeddable video player
1. The embeddable read-only chat
1. the embeddable read-write chat

Each have a set of test to make sure they load, have the expected elements on the screen, that API requests are successful, and that there are no errors being thrown in the console.

The main web frontend additionally iterates its tests over a set of different device characteristics to verify mobile and tablet usage and goes through some interactive usage of the page such as changing their name and sending a chat message by clicking and typing.

While it emulates the user agent, screen size, and touch features of different devices, they're still just a copy of Chromium running and not a true emulation of these other devices. So any "it breaks only on Safari" type bugs will not get caught.

It can't actually play video, so anything specific about video playback cannot be verified with these tests.

## Setup

`npm install`

## Run

`./run.sh`
## Screenshots

After the tests finish a set of screenshots will be saved into the `screenshots` directory to aid in troubleshooting or sanity checking different viewport sizes. three