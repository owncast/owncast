import { Popover } from 'antd';
import React, { FC, useEffect, useId, useState } from 'react';
import { useAtom, useAtomValue } from 'jotai';
import sanitizeHtml from 'sanitize-html';
import Graphemer from 'graphemer';

import dynamic from 'next/dynamic';
import classNames from 'classnames';
import { useTranslation } from 'next-export-i18n';
import ContentEditable from './ContentEditable';
import { websocketServiceAtom, chatInputDraftAtom } from '../../stores/ClientConfigStore';
import { MessageType } from '../../../interfaces/socket-events';
import { User } from '../../../interfaces/user.model';
import { Localization } from '../../../types/localization';
import styles from './ChatTextField.module.scss';

// Lazy loaded components

const EmojiPicker = dynamic(() => import('./EmojiPicker').then(mod => mod.EmojiPicker), {
  ssr: false,
});

const SendOutlined = dynamic(() => import('@ant-design/icons/SendOutlined'), {
  ssr: false,
});

const SmileOutlined = dynamic(() => import('@ant-design/icons/SmileOutlined'), {
  ssr: false,
});

export type ChatTextFieldProps = {
  defaultText?: string;
  enabled: boolean;
  focusInput: boolean;
  disabledPlaceholder?: string;
  mentionCandidates?: Pick<User, 'id' | 'displayName'>[];
};

const characterLimit = 300;
const maxMentions = 16;
const maxSocketPayloadBytes = 2048;
const maxNodeDepth = 10;
const graphemer = new Graphemer();

const getNodeTextContent = (node, depth, serializeMention?: (element: HTMLElement) => string) => {
  let text = '';

  if (depth > maxNodeDepth) return text;
  if (node === null) return text;

  switch (node.nodeType) {
    case Node.CDATA_SECTION_NODE: // unlikely
    case Node.TEXT_NODE: {
      text = node.nodeValue;
      break;
    }
    case Node.ELEMENT_NODE: {
      if (serializeMention && node.tagName === 'SPAN' && node.dataset.chatMention) {
        return serializeMention(node);
      }
      switch (node.tagName.toLowerCase()) {
        case 'img': {
          text = node.getAttribute('alt') || '';
          break;
        }
        case 'br': {
          text = '\n';
          break;
        }
        case 'strong':
        case 'b': {
          /* markdown representation of bold/strong */
          text = '**';
          for (let i = 0; i < node.childNodes.length; i += 1) {
            text += getNodeTextContent(node.childNodes[i], depth + 1, serializeMention);
          }
          text += '**';
          break;
        }
        case 'em':
        case 'i': {
          /* markdown representation of italic/emphasis */
          text = '*';
          for (let i = 0; i < node.childNodes.length; i += 1) {
            text += getNodeTextContent(node.childNodes[i], depth + 1, serializeMention);
          }
          text += '*';
          break;
        }
        case 'p': {
          text = '\n';
          for (let i = 0; i < node.childNodes.length; i += 1) {
            text += getNodeTextContent(node.childNodes[i], depth + 1, serializeMention);
          }
          break;
        }
        case 'a':
        case 'span':
        case 'div': {
          for (let i = 0; i < node.childNodes.length; i += 1) {
            text += getNodeTextContent(node.childNodes[i], depth + 1, serializeMention);
          }
          break;
        }
        /* nodes which should specifically not be parsed */
        case 'script':
        case 'style': {
          break;
        }
        default: {
          text = node.textContent;
          break;
        }
      }
      break;
    }
    default:
      break;
  }
  return text;
};

const getTextContent = (node, serializeMention?: (element: HTMLElement) => string) => {
  const text = getNodeTextContent(node, 0, serializeMention)
    .replace(/^\s+/, '') /* remove leading whitespace */
    .replace(/\s+$/, '') /* remove trailing whitespace */
    .replace(/\n([^\n])/g, '  \n$1'); /* single line break to markdown break */
  return text;
};

function getChatPayload(editor: HTMLElement, body: string) {
  const message = { type: MessageType.CHAT, body };
  if (!editor.querySelector('[data-chat-mention]')) return message;

  const encoder = new TextEncoder();
  const selected: { userId: string; label: string }[] = [];
  let prefix = 'owncastdraftmention';
  while (body.includes(prefix)) prefix += 'x';
  const markedBody = getTextContent(editor, element => {
    const { textContent: label, dataset } = element;
    if (selected.length >= maxMentions) return label;
    selected.push({ userId: dataset.chatMention, label });
    return `${prefix}${selected.length - 1}z`;
  });

  let text = '';
  let cursor = 0;
  const mentions = selected.map(({ userId, label }, index) => {
    const marker = `${prefix}${index}z`;
    const position = markedBody.indexOf(marker, cursor);
    text += markedBody.slice(cursor, position);
    const start = encoder.encode(text).length;
    text += label;
    cursor = position + marker.length;
    return { userId, start, end: encoder.encode(text).length };
  });

  const payload = { ...message, mentions };
  if (text + markedBody.slice(cursor) !== body) return message;
  while (
    mentions.length &&
    encoder.encode(JSON.stringify(payload)).length > maxSocketPayloadBytes
  ) {
    mentions.pop();
  }
  return mentions.length ? payload : message;
}

function getMentionQuery(editor: HTMLElement) {
  const selection = editor?.ownerDocument.getSelection();
  if (!selection?.isCollapsed || !selection.rangeCount) return null;
  const { startContainer, startOffset } = selection.getRangeAt(0);
  if (!editor.contains(startContainer) || startContainer.nodeType !== Node.TEXT_NODE) return null;
  if (startContainer.parentElement.closest('[data-chat-mention]')) return null;

  const beforeCaret = startContainer.textContent.slice(0, startOffset);
  const match = beforeCaret.match(/(?:^|\s)@([^@\r\n]*)$/u);
  if (!match) return null;
  return {
    node: startContainer,
    start: startOffset - match[1].length - 1,
    end: startOffset,
    query: match[1],
  };
}

export const ChatTextField: FC<ChatTextFieldProps> = ({
  defaultText,
  enabled,
  focusInput,
  disabledPlaceholder,
  mentionCandidates = [],
}) => {
  const [inputDraft, setInputDraft] = useAtom(chatInputDraftAtom);
  const [characterCount, setCharacterCount] = useState(defaultText?.length);
  const websocketService = useAtomValue(websocketServiceAtom);
  const [contentEditable, setContentEditable] = useState(null);
  const [customEmoji, setCustomEmoji] = useState([]);
  const [mentionQuery, setMentionQuery] = useState<string | null>(null);
  const [activeMention, setActiveMention] = useState(0);
  const mentionListId = useId();
  const { t } = useTranslation();
  const normalizedQuery = mentionQuery?.toLowerCase();
  const suggestions =
    enabled && mentionQuery !== null
      ? mentionCandidates
          .filter(user => user.displayName.toLowerCase().startsWith(normalizedQuery))
          .slice(0, 5)
      : [];
  const selectedIndex = Math.min(activeMention, suggestions.length - 1);

  const onRootRef = el => {
    setContentEditable(el);
  };

  const getCharacterCount = () => {
    const message = getTextContent(contentEditable);
    return graphemer.countGraphemes(message);
  };

  const sendMessage = () => {
    if (!websocketService) {
      console.log('websocketService is not defined');
      return;
    }

    const message = getTextContent(contentEditable);
    const count = graphemer.countGraphemes(message);
    if (count === 0 || count > characterLimit) return;

    websocketService.send(getChatPayload(contentEditable, message));
    setMentionQuery(null);
    contentEditable.innerHTML = '';
    setInputDraft('');
  };

  const insertTextAtEnd = (textToInsert: string) => {
    contentEditable.innerHTML += textToInsert;
  };

  const onEmojiSelect = emoji => {
    if (emoji.native) {
      insertTextAtEnd(emoji.native);
    } else {
      // Custom emoji images
      const html = `<img src="${emoji.src}" alt=":${emoji.name}:" title=":${emoji.name}:" class="emoji" />`;
      insertTextAtEnd(html);
    }
  };

  const updateMentionQuery = () => {
    const query = getMentionQuery(contentEditable)?.query ?? null;
    if (query !== mentionQuery) {
      setMentionQuery(query);
      setActiveMention(0);
    }
  };

  const insertMention = (user: Pick<User, 'id' | 'displayName'>) => {
    const match = getMentionQuery(contentEditable);
    if (!match) return;

    const doc = contentEditable.ownerDocument as Document;
    const selection = doc.getSelection();
    const range = doc.createRange();
    range.setStart(match.node, match.start);
    range.setEnd(match.node, match.end);
    selection.removeAllRanges();
    selection.addRange(range);

    const mention = doc.createElement('span');
    mention.setAttribute('contenteditable', 'false');
    mention.dataset.chatMention = user.id;
    mention.textContent = `@${user.displayName}`;
    // Use the same undo-preserving editor operation as formatted paste.
    doc.execCommand('insertHTML', false, `${mention.outerHTML} `);
    setMentionQuery(null);
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.nativeEvent.isComposing) return;
    if (suggestions.length && e.key === 'Escape') {
      e.preventDefault();
      setMentionQuery(null);
      return;
    }
    if (suggestions.length && (e.key === 'ArrowDown' || e.key === 'ArrowUp')) {
      e.preventDefault();
      const step = e.key === 'ArrowDown' ? 1 : -1;
      setActiveMention((selectedIndex + step + suggestions.length) % suggestions.length);
      return;
    }
    if (e.key === 'Enter' && !(e.shiftKey || e.metaKey || e.ctrlKey || e.altKey)) {
      e.preventDefault();
      if (suggestions.length) {
        insertMention(suggestions[selectedIndex]);
      } else {
        sendMessage();
      }
    }
  };

  const onPaste = evt => {
    evt.preventDefault();

    const clip = evt.clipboardData;
    const { types } = clip;
    const contentTypes = ['text/html', 'text/plain'];

    let content;

    for (let i = 0; i < contentTypes.length; i += 1) {
      const contentType = contentTypes[i];

      if (types.includes(contentType)) {
        content = clip.getData(contentType);
        break;
      }
    }

    if (!content) return;

    const sanitized = sanitizeHtml(content, {
      allowedTags: ['b', 'i', 'em', 'strong', 'a', 'br', 'p', 'img'],
      allowedAttributes: {
        img: ['class', 'alt', 'title', 'src'],
      },
      allowedClasses: {
        img: ['emoji'],
      },
      transformTags: {
        h1: 'p',
        h2: 'p',
        h3: 'p',
      },
    });

    // MDN lists this as deprecated, but it's the only way to save this paste
    // into the browser's Undo buffer. Plus it handles all the selection
    // deletion, caret positioning, etc automaticaly.
    if (sanitized) document.execCommand('insertHTML', false, sanitized);
  };

  const handleChange = () => {
    const count = getCharacterCount();
    setCharacterCount(count);

    if (count === 0 && contentEditable.children.length === 1) {
      /* if we have a single <br> element added by the browser, remove. */
      if (contentEditable.children[0].tagName.toLowerCase() === 'br') {
        contentEditable.removeChild(contentEditable.children[0]);
      }
    }

    // Persist draft to the shared draft atom so it survives mobile/desktop mode switches
    setInputDraft(contentEditable.innerHTML);
    updateMentionQuery();
  };

  // Focus the input when the component mounts.
  useEffect(() => {
    if (!focusInput) {
      return;
    }
    document.getElementById('chat-input-content-editable').focus({ preventScroll: true });
  }, []);

  // Restore draft from the shared draft atom when component mounts (e.g., after mobile/desktop switch)
  useEffect(() => {
    if (contentEditable && inputDraft) {
      contentEditable.innerHTML = inputDraft;
      setCharacterCount(graphemer.countGraphemes(getTextContent(contentEditable)));
    }
  }, [contentEditable]);

  const getCustomEmoji = async () => {
    try {
      const response = await fetch(`/api/emoji`);
      const emoji = await response.json();
      setCustomEmoji(emoji);

      emoji.forEach(e => {
        const preImg = document.createElement('link');
        preImg.href = e.url;
        preImg.rel = 'preload';
        preImg.as = 'image';
        document.head.appendChild(preImg);
      });
    } catch (e) {
      console.error('cannot fetch custom emoji', e);
    }
  };

  useEffect(() => {
    getCustomEmoji();
  }, []);

  return (
    <div id="chat-input" className={styles.root}>
      <div
        className={classNames(
          styles.inputWrap,
          characterCount > characterLimit && styles.maxCharacters,
        )}
      >
        <ContentEditable
          id="chat-input-content-editable"
          html={defaultText || ''}
          placeholder={
            enabled ? 'Send a message to chat' : disabledPlaceholder || 'Chat is disabled'
          }
          disabled={!enabled}
          onKeyDown={onKeyDown}
          onContentChange={handleChange}
          onClick={updateMentionQuery}
          onKeyUp={event => {
            if (['ArrowLeft', 'ArrowRight', 'Home', 'End'].includes(event.key))
              updateMentionQuery();
          }}
          onBlur={() => setMentionQuery(null)}
          onPaste={onPaste}
          onRootRef={onRootRef}
          style={{ whiteSpace: 'pre-wrap', width: '100%' }}
          role="combobox"
          aria-disabled={!enabled}
          aria-autocomplete="list"
          aria-expanded={suggestions.length > 0}
          aria-controls={suggestions.length ? mentionListId : undefined}
          aria-activedescendant={
            suggestions.length ? `${mentionListId}-${selectedIndex}` : undefined
          }
          aria-label="Chat text input"
        />
        {suggestions.length > 0 && (
          <ul
            id={mentionListId}
            className={styles.mentionSuggestions}
            role="listbox"
            aria-label={t(Localization.Frontend.Chat.mentionSuggestions)}
          >
            {suggestions.map((user, index) => (
              <li key={user.id} role="presentation">
                <div
                  id={`${mentionListId}-${index}`}
                  role="option"
                  tabIndex={-1}
                  aria-selected={index === selectedIndex}
                  onMouseDown={event => event.preventDefault()}
                  onClick={() => insertMention(user)}
                  onKeyDown={onKeyDown}
                >
                  @{user.displayName}
                </div>
              </li>
            ))}
          </ul>
        )}
        {enabled && (
          <div style={{ display: 'flex', paddingLeft: '5px' }}>
            <Popover
              content={
                <div className={styles.emojiPickerContainer}>
                  <EmojiPicker customEmoji={customEmoji} onEmojiSelect={onEmojiSelect} />
                </div>
              }
              trigger="click"
              placement="topRight"
            >
              <button
                type="button"
                aria-label="Emoji picker"
                id="owncast-emoji-picker-button"
                className={styles.emojiButton}
                title="Emoji picker button"
              >
                <SmileOutlined />
              </button>
            </Popover>
            <button
              type="button"
              aria-label="Send message"
              id="owncast-send-message-button"
              className={styles.sendButton}
              title="Send message Button"
              onClick={sendMessage}
            >
              <SendOutlined />
            </button>
          </div>
        )}
      </div>
    </div>
  );
};
