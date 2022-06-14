<br />
<p align="center">
  <a href="https://github.com/owncast/owncast" alt="Owncast">
    <img src="https://owncast.online/images/logo.png" alt="Logo" width="200">
  </a>
</p>

<br/>

  <p align="center">
    <strong>Take control over your content and stream it yourself.</strong>
    <br />
    <a href="http://owncast.online"><strong>Explore the docs »</strong></a>
    <br />
    <a href="https://watch.owncast.online/">View Demo</a>
    ·
    <a href="https://broadcast.owncast.online/">Use Our Server for Testing</a>
    ·
    <a href="https://owncast.online/faq/">FAQ</a>
    ·
    <a href="https://github.com/owncast/owncast/issues">Report Bug</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->

## Table of Contents

- [About the Project](#about-the-project)
- [Getting Started](#getting-started)
- [Use with your broadcasting software](#use-with-your-existing-broadcasting-software)
- [Building from source](#building-from-source)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

<!-- ABOUT THE PROJECT -->

## About The Project

<p align="center">
  <a href="https://owncast.online/images/owncast-splash.png">
    <img src="https://owncast.online/images/owncast-splash.png" width="70%">
  </a>
</p>

Owncast is an open source, self-hosted, decentralized, single user live video streaming and chat server for running your own live streams similar in style to the large mainstream options.  It offers complete ownership over your content, interface, moderation and audience. <a href="https://watch.owncast.online">Visit the demo</a> for an example.

<div>
    <img alt="GitHub all releases" src="https://img.shields.io/github/downloads/owncast/owncast/total?style=for-the-badge">
	  <a href="https://hub.docker.com/r/gabekangas/owncast">
      <img alt="Docker Pulls" src="https://img.shields.io/docker/pulls/gabekangas/owncast?style=for-the-badge">
	  </a>
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

## Use with your existing broadcasting software

In general Owncast is compatible with any software that uses `RTMP` to broadcast to a remote server. `RTMP` is what all the major live streaming services use, so if you’re currently using one of those it’s likely that you can point your existing software at your Owncast instance instead.

OBS, Streamlabs, Restream and many others have been used with Owncast. [Read more about compatibility with existing software](https://owncast.online/docs/broadcasting/).

## Building from Source

1. Ensure you have the gcc compiler installed.
1. Install the [Go toolchain](https://golang.org/dl/) (1.16 or above).
1. Clone the repo. `git clone https://github.com/owncast/owncast`
1. `go run main.go` will run from source.
1. Visit `http://yourserver:8080` to access the web interface or `http://yourserver:8080/admin` to access the admin.
1. Point your [broadcasting software](https://owncast.online/docs/broadcasting/) at your new server and start streaming.

There is also a supplied `Dockerfile` so you can spin it up from source with little effort. [Read more about running from source](https://owncast.online/docs/building/).

### Bundling in latest admin from source

The admin ui is built at: https://github.com/owncast/owncast-admin it is bundled into the final binary using pkger.

To bundle in the latest admin UI:

1. From the owncast directory run the packager script: `./build/admin/bundleAdmin.sh`
1. Compile or run like above. `go run main.go`

## Contributing

Owncast is a growing open source project that is giving freedom, flexibility and fun to live streamers.
And while we have a small team of kind, talented and thoughtful volunteers, we have gaps in our skillset that we’d love to fill so we can get even better at building tools that make a difference for people.

We abide by our [Code of Conduct](https://owncast.online/contribute/) and feel strongly about open, appreciative, and empathetic people joining us.
We’ve been very lucky to have this so far, so maybe you can help us with your skills and passion, too!

There is a larger, more detailed, and more up-to-date [guide for helping contribute to Owncast on our website](https://owncast.online/help/).

### Architecture

Owncast consists of two repositories with two standalone projects. [The repo you're looking at now](https://github.com/owncast/owncast) is the core repository with the backend and frontend.  [owncast/owncast-admin](https://github.com/owncast/owncast-admin) is an additional web project that is built separately and used for configuration and management of an Owncast server.

### Suggestions when working with the Owncast codebase

1. Install [golangci-lint](https://golangci-lint.run/usage/install/) for helpful warnings and suggestions [directly in your editor](https://golangci-lint.run/usage/integrations/) when writing Go.
1. If using VSCode install the [lit-html](https://marketplace.visualstudio.com/items?itemName=bierner.lit-html) extension to aid in syntax highlighting of our frontend HTML + Preact.
1. Run the project with `go run main.go`.


<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE` for more information.

<!-- CONTACT -->

## Contact

Project chat: [Join us on Rocket.Chat](https://owncast.rocket.chat/home) if you want to contribute, follow along, or if you have questions.

Gabe Kangas - [@gabek@fosstodon.org](https://fosstodon.org/@gabek) - email [gabek@real-ity.com](mailto:gabek@real-ity.com)

Project Link: [https://github.com/owncast/owncast](https://github.com/owncast/owncast)
