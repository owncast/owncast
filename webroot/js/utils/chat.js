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

  return convertToMarkup(formattedText);
}

function highlightUsername(message, username) {
	const pattern = new RegExp('@?' + username.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&'), 'gi');
  return message.replace(
    pattern,
    '<span class="highlighted px-1 rounded font-bold bg-orange-500">$&</span>'
  );
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


// utils from https://gist.github.com/nathansmith/86b5d4b23ed968a92fd4
/*
  You would call this after getting an element's
  `.innerHTML` value, while the user is typing.
*/
export function convertToText(str = '') {
  // Ensure string.
  let value = String(str);

  // Convert encoding.
  value = value.replace(/&nbsp;/gi, ' ');
  value = value.replace(/&amp;/gi, '&');

  // Replace `<br>`.
  value = value.replace(/<br>/gi, '\n');

  // Replace `<div>` (from Chrome).
  value = value.replace(/<div>/gi, '\n');

  // Replace `<p>` (from IE).
  value = value.replace(/<p>/gi, '\n');

  // Trim each line.
  value = value
    .split('\n')
    .map((line = '') => {
      return line.trim();
    })
    .join('\n');

  // No more than 2x newline, per "paragraph".
  value = value.replace(/\n\n+/g, '\n\n');

  // Clean up spaces.
  value = value.replace(/[ ]+/g, ' ');
  value = value.trim();

  // Expose string.
  return value;
}

/*
  You would call this when receiving a plain text
  value back from an API, and before inserting the
  text into the `contenteditable` area on a page.
*/
export function convertToMarkup(str = '') {
  return convertToText(str).replace(/\n/g, '<br>');
}

/*
  You would call this when a user pastes from
  the clipboard into a `contenteditable` area.
*/
export function convertOnPaste( event = { preventDefault() {} }) {
  // Prevent paste.
  event.preventDefault();

  // Set later.
  let value = '';

  // Does method exist?
  const hasEventClipboard = !!(
    event.clipboardData &&
    typeof event.clipboardData === 'object' &&
    typeof event.clipboardData.getData === 'function'
  );

  // Get clipboard data?
  if (hasEventClipboard) {
    value = event.clipboardData.getData('text/plain');
  }

  // Insert into temp `<textarea>`, read back out.
  const textarea = document.createElement('textarea');
  textarea.innerHTML = value;
  value = textarea.innerText;

  // Clean up text.
  value = convertToText(value);

  // Insert text.
  if (typeof document.execCommand === 'function') {
    document.execCommand('insertText', false, value);
  }
}

export function formatTimestamp(sentAt) {
  sentAt = new Date(sentAt);

  let diffInDays = ((new Date()) - sentAt) / (24 * 3600 * 1000);
  if (diffInDays >= -1) {
    return `${sentAt.toLocaleDateString('en-US', {dateStyle: 'medium'})} at ` +
      sentAt.toLocaleTimeString();
  }

  return sentAt.toLocaleTimeString();
}
