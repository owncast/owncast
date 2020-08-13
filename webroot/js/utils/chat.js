export const KEY_USERNAME = 'owncast_username';
export const KEY_AVATAR = 'owncast_avatar';
export const KEY_CHAT_DISPLAYED = 'owncast_chat';
export const KEY_CHAT_FIRST_MESSAGE_SENT = 'owncast_first_message_sent';
export const CHAT_INITIAL_PLACEHOLDER_TEXT = 'Type here to chat, no account necessary.';
export const CHAT_PLACEHOLDER_TEXT = 'Message';
export const CHAT_PLACEHOLDER_OFFLINE = 'Chat is offline.';

export function formatMessageText(message) {
  showdown.setFlavor('github');
  var markdownToHTML = new showdown.Converter({
    emoji: true,
    openLinksInNewWindow: true,
    tables: false,
    simplifiedAutoLink: false,
    literalMidWordUnderscores: true,
    strikethrough: true,
    ghMentions: false,
  }).makeHtml(this.body);
  const linked = autoLink(markdownToHTML, {
    embed: true,
    removeHTTP: true,
    linkAttr: {
      target: '_blank'
    }
  });
  const highlighted = highlightUsername(linked);
  return addNewlines(highlighted);
}

function highlightUsername(message) {
  const username = document.getElementById('self-message-author').value;
  const pattern = new RegExp('@?' + username.replace(/[-\/\\^$*+?.()|[\]{}]/g, '\\$&'), 'gi');
  return message.replace(pattern, '<span class="highlighted">$&</span>');
}
