import {
  CHAT_INITIAL_PLACEHOLDER_TEXT,
  CHAT_PLACEHOLDER_TEXT,
  CHAT_PLACEHOLDER_OFFLINE,
} from './constants.js';

// Taken from https://stackoverflow.com/a/46902361
export function getCaretPosition(node) {
  const selection = window.getSelection();

  if (selection.rangeCount === 0) {
    return 0;
  }

  const range = selection.getRangeAt(0);
  const preCaretRange = range.cloneRange();
  const tempElement = document.createElement('div');

  preCaretRange.selectNodeContents(node);
  preCaretRange.setEnd(range.endContainer, range.endOffset);
  tempElement.appendChild(preCaretRange.cloneContents());

  return tempElement.innerHTML.length;
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
    return hasSentFirstChatMessage
      ? CHAT_PLACEHOLDER_TEXT
      : CHAT_INITIAL_PLACEHOLDER_TEXT;
  }
  return CHAT_PLACEHOLDER_OFFLINE;
}

export function extraUserNamesFromMessageHistory(messages) {
  const list = [];
  if (messages) {
    messages
      .filter((m) => m.user && m.user.displayName)
      .forEach(function (message) {
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

  // Cleanup the emoji titles.
  value = value.replace(/\u200C{2}/gi, '');

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
export function convertOnPaste(event = { preventDefault() {} }, emojiList) {
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

  const HTML = emojify(value, emojiList);

  // Insert text.
  if (typeof document.execCommand === 'function') {
    document.execCommand('insertHTML', false, HTML);
  }
}

export function createEmojiMarkup(data, isCustom) {
  const emojiUrl = isCustom ? data.emoji : data.url;
  const emojiName = (
    isCustom
      ? data.name
      : data.url.split('\\').pop().split('/').pop().split('.').shift()
  ).toLowerCase();
  return (
    '<img class="emoji" alt=":‌‌' +
    emojiName +
    '‌‌:" title=":‌‌' +
    emojiName +
    '‌‌:" src="' +
    emojiUrl +
    '"/>'
  );
}

// trim html white space characters from ends of messages for more accurate counting
export function trimNbsp(html) {
  return html.replace(/^(?:&nbsp;|\s)+|(?:&nbsp;|\s)+$/gi, '');
}

export function emojify(HTML, emojiList) {
  const textValue = convertToText(HTML);

  for (var lastPos = textValue.length; lastPos >= 0; lastPos--) {
    const endPos = textValue.lastIndexOf(':', lastPos);
    if (endPos <= 0) {
      break;
    }
    const startPos = textValue.lastIndexOf(':', endPos - 1);
    if (startPos === -1) {
      break;
    }
    const typedEmoji = textValue.substring(startPos + 1, endPos).trim();
    const emojiIndex = emojiList.findIndex(function (emojiItem) {
      return emojiItem.name.toLowerCase() === typedEmoji.toLowerCase();
    });

    if (emojiIndex != -1) {
      const emojiImgElement = createEmojiMarkup(emojiList[emojiIndex], true);
      HTML = HTML.replace(':' + typedEmoji + ':', emojiImgElement);
    }
  }
  return HTML;
}

// MODERATOR UTILS
export function checkIsModerator(message) {
  const { user } = message;
  const { scopes } = user;

  if (!scopes || scopes.length === 0) {
    return false;
  }

  return scopes.includes('MODERATOR');
}
