import linkify from 'https://cdn.skypack.dev/linkifyjs';

const { tokenize, options } = linkify;
const { Options } = options;

function linkifyStr(str, opts = {}) {
  opts = new Options(opts);

  let tokens = tokenize(str);
  let result = [];

  for (let i = 0; i < tokens.length; i++) {
    let token = tokens[i];

    if (token.type === 'nl' && opts.nl2br) {
      result.push('<br>\n');
      continue;
    } else if (!token.isLink || !opts.check(token)) {
      result.push(escapeText(token.toString()));
      continue;
    }

    let {
      formatted,
      formattedHref,
      tagName,
      className,
      target,
      attributes,
    } = opts.resolve(token);

    let link = `<${tagName} href="${escapeAttr(formattedHref)}"`;

    if (className) {
      link += ` class="${escapeAttr(className)}"`;
    }

    if (target) {
      link += ` target="${escapeAttr(target)}"`;
    }

    if (attributes) {
      link += ` ${attributesToString(attributes)}`;
    }

    link += `>${escapeText(formatted)}</${tagName}>`;

    if (getYoutubeIdFromURL(formattedHref)) {
      const youtubeID = getYoutubeIdFromURL(formattedHref);
      link += '<br/>' + getYoutubeEmbedFromID(youtubeID);
    } else if (formattedHref.indexOf('instagram.com/p/') > -1) {
      link += '<br/>' + getInstagramEmbedFromURL(formattedHref);
    } else if (isImage(formattedHref)) {
      link = getImageForURL(formattedHref);
    }

    result.push(link);
  }

  return result.join('');
}


function escapeText(text) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;');
}

function escapeAttr(href) {
  return href.replace(/"/g, '&quot;');
}

function attributesToString(attributes) {
  if (!attributes) { return ''; }
  let result = [];

  for (let attr in attributes) {
    let val = attributes[attr] + '';
    result.push(`${attr}="${escapeAttr(val)}"`);
  }
  return result.join(' ');
}

function getYoutubeIdFromURL(url) {
  try {
    var regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*/;
    var match = url.match(regExp);

    if (match && match[2].length == 11) {
      return match[2];
    } else {
      return null;
    }
  } catch (e) {
    console.log(e);
    return null;
  }
}

function getYoutubeEmbedFromID(id) {
  return `<iframe class="chat-embed" src="//www.youtube.com/embed/${id}" frameborder="0" allowfullscreen></iframe>`
}

function getInstagramEmbedFromURL(url) {
  const urlObject = new URL(url.replace(/\/$/, ""));
  urlObject.pathname += "/embed"
  return `<iframe class="chat-embed" height="150px" src="${urlObject.href}" frameborder="0" allowfullscreen></iframe>`
}

function isImage(url) {
  const re = /\.(jpe?g|png|gif)$/;
  const isImage = re.test(url);
  return isImage;
}

function getImageForURL(url) {
  return `<a target="_blank" href="${url}"><img class="embedded-image" src="${url}" width="100%" height="150px"/></a>`
}

export { linkifyStr };