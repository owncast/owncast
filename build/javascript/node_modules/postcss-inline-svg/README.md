# postcss-inline-svg [![Build Status][travis-img]][travis]

[PostCSS] plugin to reference an SVG file and control its attributes with CSS syntax.

[PostCSS]: https://github.com/postcss/postcss
[travis-img]: https://travis-ci.org/TrySound/postcss-inline-svg.svg
[travis]: https://travis-ci.org/TrySound/postcss-inline-svg

```postcss
@svg-load nav url(img/nav.svg) {
    fill: #cfc;
    path:nth-child(2) {
        fill: #ff0;
    }
}
.nav {
    background: svg-inline(nav);
}
.up {
    background: svg-load('img/arrow-up.svg', fill=#000, stroke=#fff);
}
```

```css
.nav {
    background: url("data:image/svg+xml;charset=utf-8,%3Csvg fill='%23cfc'%3E%3Cpath d='...'/%3E%3Cpath d='...' fill='%23ff0'/%3E%3Cpath d='...'/%3E%3C/svg%3E");
}
.up {
    background: url("data:image/svg+xml;charset=utf-8,%3Csvg fill='%23000' stroke='%23fff'%3E...%3C/svg%3E");
}
```

PostCSS parsers allow to use different syntax (but only one syntax in one svg-load() definition):

```postcss
.up {
    background: svg-load('img/arrow-up.svg', fill: #000, stroke: #fff);
}
.down {
    background: svg-load('img/arrow-down.svg', fill=#000, stroke=#fff);
}
```

## Usage

```js
postcss([ require('postcss-inline-svg')(options) ])
```

See [PostCSS] docs for examples for your environment.

### Options

#### options.paths

Array of paths to resolve URL. Paths are tried in order, until an existing file is found.

Default: `false` - path will be relative to source file if it was specified.

#### options.removeFill

Default: `false` - with `true` removes all `fill` attributes before applying specified.
Passed RegExp filters files by ID.

#### options.removeStroke

Default: `false` - with `true` removes all `stroke` attributes before applying specified.
Passed RegExp filters files by ID.


#### options.encode(svg)

Processes SVG after applying parameters. Default will be ommited if passed `false`.

Default:

```js
function encode(code) {
    return code
        .replace(/%/g, '%25')
        .replace(/</g, '%3C')
        .replace(/>/g, '%3E')
        .replace(/&/g, '%26')
        .replace(/#/g, "%23")
        .replace(/{/g, "%7B")
        .replace(/}/g, "%7D");
}
```

#### options.transform(svg, path)

Transforms SVG after `encode` function and generates URL.

#### options.xmlns

type: boolean
default: true

Adds `xmlns` attribute to SVG if not present.


## Frequently asked questions

### Why svg-load() doesn't work and the color still black (or red or whatever)?

That's because `svg-load()` overrides attributes only in `<svg>` element and children inherit that color.
But if there is already color on children nothing will be inherited.

For example

```xml
<svg>
    <path fill="#ff0000" d="..." />
</svg>
```

after inline-svg (fill: #000) will looks like

```xml
<svg fill="#000">
    <path fill="#ff0000" d="..." />
</svg>
```

There are three solutions: to remove that attribute (preferable), to use extended `@svg-load` notation or to use removeFill option.

### How to optimize svg on build step?

> There is a plugin. :)

You are able to add [postcss-svgo](https://github.com/cssnano/cssnano/tree/master/packages/postcss-svgo) in your PostCSS plugins list which will detect every URL which contains data SVG URI and minify via [svgo](https://github.com/svg/svgo).

```js
postcss([
    require('postcss-inline-svg'),
    require('postcss-svgo')
])
```

Or if you use [cssnano](https://cssnano.co/) your SVG URLs already minified
as cssnano includes postcss-svgo.

```js
postcss([
    require('postcss-inline-svg'),
    require('cssnano')
])
```


# License

MIT © [Bogdan Chadkin](mailto:trysound@yandex.ru)

## Support on Beerpay
Hey dude! Help me out for a couple of :beers:!

[![Beerpay](https://beerpay.io/TrySound/postcss-inline-svg/badge.svg?style=beer-square)](https://beerpay.io/TrySound/postcss-inline-svg)  [![Beerpay](https://beerpay.io/TrySound/postcss-inline-svg/make-wish.svg?style=flat-square)](https://beerpay.io/TrySound/postcss-inline-svg?focus=wish)
