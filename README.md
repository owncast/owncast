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
    <a href="https://owncast.online"><strong>Explore the docs »</strong></a>
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

Owncast is an open source, self-hosted, decentralized, single user live video streaming and chat server for running your own live streams similar in style to the large mainstream options. It offers complete ownership over your content, interface, moderation and audience. <a href="https://watch.owncast.online">Visit the demo</a> for an example.

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

In general, Owncast is compatible with any software that uses `RTMP` to broadcast to a remote server. `RTMP` is what all the major live streaming services use, so if you’re currently using one of those it’s likely that you can point your existing software at your Owncast instance instead.

OBS, Streamlabs, Restream and many others have been used with Owncast. [Read more about compatibility with existing software](https://owncast.online/docs/broadcasting/).

## Building from Source

Owncast consists of two projects.

1. The Owncast backend is written in Go.
1. The frontend is written in React.

[Read more about running from source](https://owncast.online/development/).

### Important note about source code and the develop branch

The `develop` branch is always the most up-to-date state of development and this may not be what you always want. If you want to run the latest released stable version, check out the tag related to that release. For example, if you'd only like the source prior to the v0.1.0 development cycle you can check out the `v0.0.13` tag.

> Note: Currently Owncast does not natively support Windows servers. However, Windows Users can use Windows Subsystem for Linux (WSL2) to install Owncast. For details visit [this document](https://github.com/owncast/owncast/blob/develop/contrib/owncast_for_windows.md).

### Backend

The Owncast backend is a service written in Go.

1. Ensure you have prerequisites installed.
   - C compiler, such as [GCC compiler](https://gcc.gnu.org/install/download.html) or a [Musl-compatible compiler](https://musl.libc.org/)
   - [ffmpeg](https://ffmpeg.org/download.html)
1. Install the [Go toolchain](https://golang.org/dl/) (1.20 or above).
1. Clone the repo. `git clone https://github.com/owncast/owncast`
1. `go run main.go` will run from the source.
1. Visit `http://yourserver:8080` to access the web interface or `http://yourserver:8080/admin` to access the admin.
1. Point your [broadcasting software](https://owncast.online/docs/broadcasting/) at your new server and start streaming.

### Frontend

The frontend is the web interface that includes the player, chat, embed components, and other UI.

1. This project lives in the `web` directory.
1. Run `npm install` to install the Javascript dependencies.
1. Run `npm run dev`

## Contributing

Owncast is a growing open source project that is giving freedom, flexibility and fun to live streamers.
And while we have a small team of kind, talented and thoughtful volunteers, we have gaps in our skillset that we’d love to fill so we can get even better at building tools that make a difference for people.

We abide by our [Code of Conduct](https://owncast.online/contribute/) and feel strongly about open, appreciative, and empathetic people joining us.
We’ve been very lucky to have this so far, so maybe you can help us with your skills and passion, too!

There is a larger, more detailed, and more up-to-date [guide for helping contribute to Owncast on our website](https://owncast.online/help/).

<!-- LICENSE -->

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Supported by

- This project is tested with [BrowserStack](https://browserstack.com).

<!-- CONTACT -->

## Contact

Project chat: [Join us on Rocket.Chat](https://owncast.rocket.chat/home) if you want to contribute, follow along, or if you have questions.

Gabe Kangas - [@gabek@social.gabekangas.com](https://social.gabekangas.com/gabek) - email [gabek@real-ity.com](mailto:gabek@real-ity.com)

Project Link: [https://github.com/owncast/owncast](https://github.com/owncast/owncast)
