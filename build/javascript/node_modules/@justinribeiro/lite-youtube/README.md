[![npm version](https://badge.fury.io/js/%40justinribeiro%2Flite-youtube.svg)](https://badge.fury.io/js/%40justinribeiro%2Flite-youtube)

# \<lite-youtube\>

> A web component that displays render YouTube embeds faster. The shadowDom web component version of Paul's [lite-youtube-embed](https://github.com/paulirish/lite-youtube-embed).

## Features

- No dependencies; it's just a vanilla web component.
- It's fast yo.
- It's Shadow Dom encapsulated!
- It's responsive 16:9
- It's accessible via keyboard and will set ARIA via the `videotitle` attribute
- It's locale ready; you can set the `videoplay` to have a properly locale based label
- Set the `start` attribute to start at a particular place in a video
- You can set `autoload` to use Intersection Observer to load the iframe when scrolled into view.
- Loads placeholder image as WebP with a Jpeg fallback

## Install

This web component is built with ES modules in mind and is
available on NPM:

Install code-block:

```sh
npm i @justinribeiro/lite-youtube
# or
yarn add @justinribeiro/lite-youtube
```

After install, import into your project:

```js
import '@justinribeiro/lite-youtube';
```

## Install with CDN

If you want the paste-and-go version, you can simply load it via CDN:

```html
<script type="module" src="https://cdn.jsdelivr.net/npm/@justinribeiro/lite-youtube@0.6.2/lite-youtube.js">
```

## Basic Usage

```html
<lite-youtube videoid="guJLfqTFfIw"></lite-youtube>
```

## Add Video Title

```html
<lite-youtube
  videoid="guJLfqTFfIw"
  videotitle="This is a video title"
></lite-youtube>
```

## Change "Play" for Locale</h3>

```html
<lite-youtube
  videoid="guJLfqTFfIw"
  videoplay="Mirar"
  videotitle="Mis hijos se burlan de mi español"
>
</lite-youtube>
```

## Style It

Height and Width are responsive in the component.

```html
<style>
  .styleIt {
    width: 400px;
    margin: auto;
  }
</style>
<div class="styleIt">
  <lite-youtube videoid="guJLfqTFfIw"></lite-youtube>
</div>
```

## Set a video start time

```html
<!-- Start at 5 seconds -->
<lite-youtube videoid="guJLfqTFfIw" videoStartAt="5"></lite-youtube>
```

## AutoLoad with IntersectionObserver

Uses Intersection Observer if available to automatically load the YouTube iframe when scrolled into view.

```html
<lite-youtube videoid="guJLfqTFfIw" autoload> </lite-youtube>
```

## YouTube QueryParams

Use any [YouTube Embedded Players and Player Parameters](https://developers.google.com/youtube/player_parameters) you like

```html
<lite-youtube videoid="guJLfqTFfIw" params="controls=0&enablejsapi=1">
</lite-youtube>
```

## Attributes

The web component allows certain attributes to be give a little additional
flexibility.

| Name           | Description                                                      | Default |
| -------------- | ---------------------------------------------------------------- | ------- |
| `videoid`      | The YouTube videoid                                              | ``      |
| `videotitle`   | The title of the video                                           | `Video` |
| `videoplay`    | The title of the play button (for translation)                   | `Play`  |
| `videoStartAt` | Set the point at which the video should start, in seconds        | `0`     |
| `autoload`     | Use Intersection Observer to load iframe when scrolled into view | `false` |
| `params`       | Set YouTube query parameters                                     | ``      |
