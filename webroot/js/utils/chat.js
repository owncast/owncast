import { addNewlines } from './helpers.js';
import {
  CHAT_INITIAL_PLACEHOLDER_TEXT,
  CHAT_PLACEHOLDER_TEXT,
  CHAT_PLACEHOLDER_OFFLINE,
} from './constants.js';

export function formatMessageText(message, username) {
  showdown.setFlavor('github');
  let formattedText = new showdown.Converter({
    emoji: true,
    openLinksInNewWindow: true,
    tables: false,
    simplifiedAutoLink: false,
    literalMidWordUnderscores: true,
    strikethrough: true,
    ghMentions: false,
  }).makeHtml(message);

  formattedText = linkify(formattedText, message);
  formattedText = highlightUsername(formattedText, username);

  return addNewlines(formattedText);
}

function highlightUsername(message, username) {
	const pattern = new RegExp('@?' + username.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&'), 'gi');
  return message.replace(pattern, '<span class="highlighted font-bold bg-orange-500">$&</span>');
}

function linkify(text, rawText) {
  const urls = getURLs(stripTags(rawText));
  if (urls) {
    urls.forEach(function (url) {
      let linkURL = url;

      // Add http prefix if none exist in the URL so it actually
      // will work in an anchor tag.
      if (linkURL.indexOf('http') === -1) {
        linkURL = 'http://' + linkURL;
      }

      // Remove the protocol prefix in the display URLs just to make
      // things look a little nicer.
      const displayURL = url.replace(/(^\w+:|^)\/\//, '');
      const link = `<a href="${linkURL}" target="_blank">${displayURL}</a>`;
      text = text.replace(url, link);

      if (getYoutubeIdFromURL(url)) {
        if (isTextJustURLs(text, [url, displayURL])) {
          text = '';
        } else {
          text += '<br/>';
        }

        const youtubeID = getYoutubeIdFromURL(url);
        text += getYoutubeEmbedFromID(youtubeID);
      } else if (url.indexOf('instagram.com/p/') > -1) {
        if (isTextJustURLs(text, [url, displayURL])) {
          text = '';
        } else {
          text += `<br/>`;
        }
        text += getInstagramEmbedFromURL(url);
      } else if (isImage(url)) {
        if (isTextJustURLs(text, [url, displayURL])) {
          text = '';
        } else {
          text += `<br/>`;
        }
        text += getImageForURL(url);
      }
    }.bind(this));
  }
  return text;
}

function isTextJustURLs(text, urls) {
  for (var i = 0; i < urls.length; i++) {
    const url = urls[i];
    if (stripTags(text) === url) {
      return true;
    }
  }
  return false;
}


function stripTags(str) {
	return str.replace(/<\/?[^>]+(>|$)/g, "");
}

function getURLs(str) {
	var exp = /((\w+:\/\/\S+)|(\w+[\.:]\w+\S+))[^\s,\.]/ig;
	return str.match(exp);
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
  return `<div class="chat-embed youtube-embed"><lite-youtube videoid="${id}" /></div>`;
}

function getInstagramEmbedFromURL(url) {
	const urlObject = new URL(url.replace(/\/$/, ""));
	urlObject.pathname += "/embed";
	return `<iframe class="chat-embed instagram-embed" src="${urlObject.href}" frameborder="0" allowfullscreen></iframe>`;
}

function isImage(url) {
	const re = /\.(jpe?g|png|gif)$/i;
	return re.test(url);
}

function getImageForURL(url) {
	return `<a target="_blank" href="${url}"><img class="chat-embed embedded-image" src="${url}" /></a>`;
}

// Taken from https://stackoverflow.com/questions/3972014/get-contenteditable-caret-index-position
export function getCaretPosition(editableDiv) {
	var caretPos = 0,
		sel, range;
	if (window.getSelection) {
		sel = window.getSelection();
		if (sel.rangeCount) {
			range = sel.getRangeAt(0);
			if (range.commonAncestorContainer.parentNode == editableDiv) {
				caretPos = range.endOffset;
			}
		}
	} else if (document.selection && document.selection.createRange) {
		range = document.selection.createRange();
		if (range.parentElement() == editableDiv) {
			var tempEl = document.createElement("span");
			editableDiv.insertBefore(tempEl, editableDiv.firstChild);
			var tempRange = range.duplicate();
			tempRange.moveToElementText(tempEl);
			tempRange.setEndPoint("EndToEnd", range);
			caretPos = tempRange.text.length;
		}
	}
	return caretPos;
}

// Might not need this anymore
// Pieced together from parts of https://stackoverflow.com/questions/6249095/how-to-set-caretcursor-position-in-contenteditable-element-div
export function setCaretPosition(editableDiv, position) {
	var range = document.createRange();
	var sel = window.getSelection();
	range.selectNode(editableDiv);
	range.setStart(editableDiv.childNodes[0], position);
	range.collapse(true);

	sel.removeAllRanges();
	sel.addRange(range);
}


export function generatePlaceholderText(isEnabled, hasSentFirstChatMessage) {
  if (isEnabled) {
    return hasSentFirstChatMessage ? CHAT_PLACEHOLDER_TEXT : CHAT_INITIAL_PLACEHOLDER_TEXT;
  }
  return CHAT_PLACEHOLDER_OFFLINE;
}

export function extraUserNamesFromMessageHistory(messages) {
  const list = [];
  if (messages) {
    messages.forEach(function(message) {
      if (!list.includes(message.author)) {
        list.push(message.author);
      }
    });
  }
  return list;
}
