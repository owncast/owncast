import {
  CHAT_INITIAL_PLACEHOLDER_TEXT,
  CHAT_PLACEHOLDER_TEXT,
  CHAT_PLACEHOLDER_OFFLINE,
} from './constants.js';


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
      if (!list.includes(message.user.displayName)) {
        list.push(message.user.displayName);
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

export function createEmojiMarkup(data, isCustom) {
  const emojiUrl = isCustom ? data.emoji : data.url;
  const emojiName = (isCustom ? data.name : data.url.split('\\').pop().split('/').pop().split('.').shift()).toLowerCase();
  return '<img class="emoji" alt=":' + emojiName + ':" title=":' + emojiName + ':" src="' + emojiUrl + '"/>';
}
