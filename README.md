<br />
<p align="center">
  <a href="https://github.com/owncast/owncast" alt="Owncast">
    <img src="https://owncast.online/images/logo.png" alt="Logo" width="200">
  </a>

<br/>

  <p align="center">
    <h2>Take control over your content and stream it yourself.</h2>
    <br />
    <a href="http://owncast.online"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://watch.owncast.online/">View Demo</a>
    ·
    <a href="https://broadcast.owncast.online/">Use Our Server for Testing</a>
    ·
    <a href="https://owncast.online/docs/faq/">FAQ</a>
    ·
    <a href="https://github.com/owncast/owncast/issues">Report Bug</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->

## Table of Contents

- [About the Project](#about-the-project)
- [Getting Started](#getting-started)
  - [Getting Started](#getting-started)
  - [Configuration](#configuration)
  - [Web Interface & Chat](#web-interface--chat)
- [Use with your broadcasting software](#use-with-your-existing-broadcasting-software)
- [Video storage and distribution options](#video-storage-options)
- [Building from source](#building-from-source)
- [License](#license)
- [Contact](#contact)

<!-- ABOUT THE PROJECT -->

## About The Project

<p align="center">
  <a href="https://owncast.online/images/owncast-splash.png">
    <img src="https://owncast.online/images/owncast-splash.png" width="70%">
  </a>
</p>

Owncast is an open source, self-hosted, decentralized, single user live streaming and chat server for running your own live streams similar in style to the large mainstream options.  It offers complete ownership over your content, interface, moderation and audience. <a href="https://watch.owncast.online">Visit the demo</a> for an example.

<div>
    <img alt="GitHub all releases" src="https://img.shields.io/github/downloads/owncast/owncast/total?style=for-the-badge">
    <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/gabekangas/owncast?style=for-the-badge">
    <a href="https://github.com/owncast/owncast/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22">
      <img alt="GitHub issues by-label" src="https://img.shields.io/github/issues-raw/owncast/owncast/good%20first%20issue?style=for-the-badge">
    </a>
    <a href="https://opencollective.com/owncast">
      <img alt="Open Collective backers and sponsors" src="https://img.shields.io/opencollective/all/owncast?style=for-the-badge">
    </a>
</div>


---

<!-- GETTING STARTED -->

## Getting Started

The goal is to have a single service that you can run and it works out of the box. **Visit the [Quickstart](https://owncast.online/docs/quickstart/) to get up and running.**

## Configuration

Many aspects can be adjusted and customized to your preferences. [Read more about Configuration](https://owncast.online/docs/configuration/) to update the web UI, video settings, and more.

## Web interface + chat

Owncast includes a web interface to your video with built-in chat that is available once you start the server.

The web interface was specifically built to be editable by anybody comfortable tweaking a web page. It's not bundled or transpiled into anything, it's just HTML + Javascript + CSS that you can start editing.

Read more about the features provided and how to configure them in the [web documentation](https://owncast.online/docs/website/).

## Use with your existing broadcasting software

In general Owncast is compatible with any software that uses `RTMP` to broadcast to a remote server. `RTMP` is what all the major live streaming services use, so if you’re currently using one of those it’s likely that you can point your existing software at your Owncast instance instead.

OBS, Streamlabs, Restream and many others have been used with Owncast. [Read more about compatibility with existing software](https://owncast.online/docs/broadcasting/).

## Video storage options

Two ways of storing and distributing the video are supported.

1. Locally via the Owncast server.
2. [S3-compatible storage](https://owncast.online/docs/s3/).

### Local file distribution

This is the simplest and works out of the box. In this scenario video will be served to the public from the computer that is running the server. If you have a fast internet connection, enough bandwidth alotted to you, and a small audience this may be fine for many people.

### S3-Compatible Storage

Instead of serving video directly from your personal server you can use a S3 compatible storage provider to offload the bandwidth and storage requirements elsewhere.

Read [more detailed documentation about configuration of S3-compatible services](https://owncast.online/docs/s3/).

## Building from Source

1. Ensure you have the gcc compiler configured.
1. Install the [Go toolchain](https://golang.org/dl/).
1. Clone the repo. `git clone https://github.com/owncast/owncast`
1. `go run main.go pkged.go` will run from source.
1. Point your [broadcasting software](https://owncast.online/docs/broadcasting/) at your new server and start streaming.

There is also a supplied `Dockerfile` so you can spin it up from source with little effort. [Read more about running from source](https://owncast.online/docs/building/).

### Bundling in latest admin from source

The admin ui is built at: https://github.com/owncast/owncast-admin it is bundled into the final binary using pkger.

To bundle in the latest admin UI:

1. Install pkger. `go install github.com/markbates/pkger/cmd/...`
1. From the owncast directory run the packager script: `./build/admin/bundleAdmin.sh`
1. Compile or run like above. `go run main.go pkged.go`

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE` for more information.

<!-- CONTACT -->

## Contact

Project chat: [Join us on Rocket.Chat](https://owncast.rocket.chat/home) if you want to contribute, follow along, or if you have questions.

Gabe Kangas - [@gabek@mastodon.social](https://mastodon.social/@gabek) - email [gabek@real-ity.com](mailto:gabek@real-ity.com)

Project Link: [https://github.com/owncast/owncast](https://github.com/owncast/owncast)
