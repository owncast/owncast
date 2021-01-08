<br />
<p align="center">
  <a href="https://github.com/owncast/owncast" alt="Owncast">
    <img src="https://owncast.online/images/logo.png" alt="Logo" width="200">
  </a>


  <p align="center">
    Take control over your content and stream it yourself.
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

* [About the Project](#about-the-project)
* [Getting Started](#getting-started)
  * [Getting Started](#getting-started)
  * [Configuration](#configuration)
  * [Web Interface & Chat](#web-interface--chat)
* [Use with your broadcasting software](#use-with-your-existing-broadcasting-software)
* [Video storage and distribution options](#video-storage-options)
* [Building from source](#building-from-source)
* [License](#license)
* [Contact](#contact)


<!-- ABOUT THE PROJECT -->
## About The Project

<p align="center">
  <a href="https://owncast.online/images/owncast-screenshot.png">
    <img src="https://owncast.online/images/owncast-screenshot.png" width="70%">
  </a>
</p>

In 2020 the world changed when everyone become stuck in their homes, looking for creative outlets to share their art, skills and themselves from inside their bedroom.

This created an explosion of live streaming on Facebook Live, YouTube Live, Instagram, and Twitch.  These services provided everything they needed, an easy way to live stream to the world, and a chat for users to be a part of their community.

That's when I wanted a better option for people. Something you could run yourself and get all the functionality of these services, where you could live stream to an audience and and allow them to take part in the chat, just like they've been used to on all the other services. **There should be a independent, standalone _Twitch in a Box_.**

**Keep in mind that while streaming to the big social companies is always free, you pay for it with your identity and your data, as well as the identity and data of every person that tunes in.  When you self-host anything you'll have to pay with your money instead.  But running a self-hosted live stream server can be done for as cheap as $5/mo, and that's a much better deal than selling your soul to Facebook, Google or Amazon.**

---

<!-- GETTING STARTED -->

## Getting Started

The goal is to have a single service that you can run and it works out of the box. **Visit the [Quickstart](https://owncast.online/docs/quickstart/) to get up and running.**

## Configuration

Many aspects can be adjusted and customized to your preferences.  [Read more about Configuration](https://owncast.online/docs/configuration/) to update the web UI, video settings, and more.

## Web interface + chat

Owncast includes a web interface to your video with built-in chat that is available once you start the server.

The web interface was specifically built to be editable by anybody comfortable tweaking a web page.  It's not bundled or transpiled into anything, it's just HTML + Javascript + CSS that you can start editing.

Read more about the features provided and how to configure them in the [web documentation](https://owncast.online/docs/website/).

## Use with your existing broadcasting software

In general Owncast is compatible with any software that uses `RTMP` to broadcast to a remote server. `RTMP` is what all the major live streaming services use, so if you’re currently using one of those it’s likely that you can point your existing software at your Owncast instance instead.

OBS, Streamlabs, Restream and many others have been used with Owncast.  [Read more about compatibility with existing software](https://owncast.online/docs/broadcasting/).

## Video storage options

Two ways of storing and distributing the video are supported.

1. Locally via the Owncast server.
2. [S3-compatible storage](https://owncast.online/docs/s3/).

### Local file distribution

This is the simplest and works out of the box.  In this scenario video will be served to the public from the computer that is running the server.  If you have a fast internet connection, enough bandwidth alotted to you, and a small audience this may be fine for many people.

### S3-Compatible Storage

Instead of serving video directly from your personal server you can use a S3 compatible storage provider to offload the bandwidth and storage requirements elsewhere.

Read [more detailed documentation about configuration of S3-compatible services](https://owncast.online/docs/s3/).


## Building from Source

1. Ensure you have the gcc compiler configured.
1. Install the [Go toolchain](https://golang.org/dl/).
1. Clone the repo.  `git clone https://github.com/owncast/owncast`
1. `go run main.go pkged.go` will run from source.
1. Point your [broadcasting software](https://owncast.online/docs/broadcasting/) at your new server and start streaming.

There is also a supplied `Dockerfile` so you can spin it up from source with little effort.  [Read more about running from source](https://owncast.online/docs/building/).


### Bundling in latest admin from source

The admin ui is built at: https://github.com/owncast/owncast-admin it is bundled into the final binary using pkger.

To bundle in the latest admin UI:

1. Install pkger. `go install github.com/markbates/pkger/cmd/...
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
