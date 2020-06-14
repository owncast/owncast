const re = {
  http: /.*?:\/\//g,
  url: /(\s|^)((https?|ftp):\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\w-\/\?\=\#\.])*/gi,
  image: /\.(jpe?g|png|gif)$/,
  email: /(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))/gi,
  cloudmusic: /http:\/\/music\.163\.com\/#\/song\?id=(\d+)/i,
  kickstarter: /(https?:\/\/www\.kickstarter\.com\/projects\/\d+\/[a-zA-Z0-9_-]+)(\?\w+\=\w+)?/i,
  youtube: /https?:\/\/www\.youtube\.com\/watch\?v=([a-zA-Z0-9_-]+)(\?\w+\=\w+)?/i,
  vimeo: /https?:\/\/(?:www\.|player\.)?vimeo.com\/(?:channels\/(?:\w+\/)?|groups\/[^\/]*\/videos\/|album\/\d+\/video\/|video\/|)(\d+)(?:$|\/|\?)/i,
  youku: /https?:\/\/v\.youku\.com\/v_show\/id_([a-zA-Z0-9_\=-]+).html(\?\w+\=\w+)?(\#\w+)?/i
}
/**
 * AutoLink constructor function
 *
 * @param {String} text
 * @param {Object} options
 * @constructor
 */
function AutoLink(string, options = {}) {
  this.string = options.safe ? safe_tags_replace(string) : string
  this.options = options
  this.attrs = ''
  this.linkAttr = ''
  this.imageAttr = ''
  if (this.options.sharedAttr) {
    this.attrs = getAttr(this.options.sharedAttr)
  }
  if (this.options.linkAttr) {
    this.linkAttr = getAttr(this.options.linkAttr)
  }
  if (this.options.imageAttr) {
    this.imageAttr = getAttr(this.options.imageAttr)
  }
}

AutoLink.prototype = {
  constructor: AutoLink,
  /**
   * call relative functions to parse url/email/image
   *
   * @returns {String}
   */
  parse: function() {
    var shouldReplaceImage = defaultTrue(this.options.image)
    var shouldRelaceEmail = defaultTrue(this.options.email)
    var shouldRelaceBr = defaultTrue(this.options.br)

    var result = ''
    if (!shouldReplaceImage) {
      result = this.string.replace(re.url, this.formatURLMatch.bind(this))
    } else {
      result = this.string.replace(re.url, this.formatIMGMatch.bind(this))
    }
    if (shouldRelaceEmail) {
      result = this.formatEmailMatch.call(this, result)
    }
    if (shouldRelaceBr) {
      result = result.replace(/\r?\n/g, '<br />')
    }
    return result
  },
  /**
   * @param {String} match
   * @returns {String} Offset 1
   */
  formatURLMatch: function(match, p1) {
    match = prepHTTP(match.trim())
    if (this.options.cloudmusic || this.options.embed) {
      if (match.indexOf('music.163.com/#/song?id=') > -1) {
        return match.replace(
          re.cloudmusic,
          p1 +
            '<iframe frameborder="no" border="0" marginwidth="0" marginheight="0" width=330 height=86 src="http://music.163.com/outchain/player?type=2&id=$1&auto=1&height=66"></iframe>'
        )
      }
    }
    if (this.options.kickstarter || this.options.embed) {
      if (re.kickstarter.test(match)) {
        return match.replace(
          re.kickstarter,
          p1 +
            '<iframe width="480" height="360" src="$1/widget/video.html" frameborder="0" scrolling="no"> </iframe>'
        )
      }
    }
    if (this.options.vimeo || this.options.embed) {
      if (re.vimeo.test(match)) {
        return match.replace(
          re.vimeo,
          p1 +
            '<iframe width="500" height="281" src="https://player.vimeo.com/video/$1" frameborder="0" webkitallowfullscreen mozallowfullscreen allowfullscreen></iframe>'
        )
      }
    }
    if (this.options.youtube || this.options.embed) {
      if (re.youtube.test(match)) {
        return match.replace(
          re.youtube,
          p1 +
            '<iframe width="560" height="315" src="https://www.youtube.com/embed/$1" frameborder="0" allowfullscreen></iframe>'
        )
      }
    }
    if (this.options.youku || this.options.embed) {
      if (re.youku.test(match)) {
        return match.replace(
          re.youku,
          p1 +
            '<iframe height=498 width=510 src="http://player.youku.com/embed/$1" frameborder=0 allowfullscreen></iframe>'
        )
      }
    }
    var text = this.options.removeHTTP ? removeHTTP(match) : match
    return (
      p1 +
      '<a target="_blank" href="' +
      match +
      '"' +
      this.attrs +
      this.linkAttr +
      '>' +
      text +
      '</a>'
    )
  },
  /**
   * @param {String} match
   * @param {String} Offset 1
   */
  formatIMGMatch: function(match, p1) {
    match = match.trim()
    var isIMG = re.image.test(match)
    if (isIMG) {
      return (
        p1 +
        '<img src="' +
        prepHTTP(match.trim()) +
        '"' +
        this.attrs +
        this.imageAttr +
        '/>'
      )
    }
    return this.formatURLMatch(match, p1)
  },
  /**
   * @param {String} text
   */
  formatEmailMatch: function(text) {
    return text.replace(
      re.email,
      '<a href="mailto:$&"' + this.attrs + this.linkAttr + '>$&</a>'
    )
  }
}

/**
 * return true if undefined
 * else return itself
 *
 * @param {Boolean} val
 * @returns {Boolean}
 */
function defaultTrue(val) {
  return typeof val === 'undefined' ? true : val
}

/**
 * parse attrs from object
 *
 * @param {Object} obj
 * @returns {Stirng}
 */
function getAttr(obj) {
  var attr = ''
  for (var key in obj) {
    if (key) {
      attr += ' ' + key + '="' + obj[key] + '"'
    }
  }
  return attr
}

/**
 * @param {String} url
 * @returns {String}
 */
function prepHTTP(url) {
  if (url.substring(0, 4) !== 'http' && url.substring(0, 2) !== '//') {
    return 'http://' + url
  }
  return url
}

/**
 * @param {String} url
 * @returns {String}
 */
function removeHTTP(url) {
  return url.replace(re.http, '')
}

var tagsToReplace = {
  '&': '&amp;',
  '<': '&lt;',
  '>': '&gt;'
}

/**
 * Replace tag if should be replace
 *
 * @param {String} tag
 * @returns {String}
 */
function replaceTag(tag) {
  return tagsToReplace[tag] || tag
}

/**
 * Make string safe by replacing html tag
 *
 * @param {String} str
 * @returns {String}
 */
function safe_tags_replace(str) {
  return str.replace(/[&<>]/g, replaceTag)
}

/**
 * return an instance of AutoLink
 *
 * @param {String} string
 * @param {Object} options
 * @returns {Object}
 */
function autoLink(string, options) {
  return new AutoLink(string, options).parse()
}
