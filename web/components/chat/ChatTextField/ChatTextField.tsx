import { Popover } from 'antd';
import React, { FC, useEffect, useReducer, useRef, useState } from 'react';
import { useRecoilValue } from 'recoil';
import ContentEditable from 'react-contenteditable';
import sanitizeHtml from 'sanitize-html';
import GraphemeSplitter from 'grapheme-splitter';

import dynamic from 'next/dynamic';
import classNames from 'classnames';
import WebsocketService from '../../../services/websocket-service';
import { websocketServiceAtom } from '../../stores/ClientConfigStore';
import { MessageType } from '../../../interfaces/socket-events';
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
};

const characterLimit = 300;
const maxNodeDepth = 10;
const graphemeSplitter = new GraphemeSplitter();

const getNodeTextContent = (node, depth) => {
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
            text += getNodeTextContent(node.childNodes[i], depth + 1);
          }
          text += '**';
          break;
        }
        case 'em':
        case 'i': {
          /* markdown representation of italic/emphasis */
          text = '*';
          for (let i = 0; i < node.childNodes.length; i += 1) {
            text += getNodeTextContent(node.childNodes[i], depth + 1);
          }
          text += '*';
          break;
        }
        case 'p': {
          text = '\n';
          for (let i = 0; i < node.childNodes.length; i += 1) {
            text += getNodeTextContent(node.childNodes[i], depth + 1);
          }
          break;
        }
        case 'a':
        case 'span':
        case 'div': {
          for (let i = 0; i < node.childNodes.length; i += 1) {
            text += getNodeTextContent(node.childNodes[i], depth + 1);
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

const getTextContent = node => {
  const text = getNodeTextContent(node, 0)
    .replace(/^\s+/, '') /* remove leading whitespace */
    .replace(/\s+$/, '') /* remove trailing whitespace */
    .replace(/\n([^\n])/g, '  \n$1'); /* single line break to markdown break */
  return text;
};

export const ChatTextField: FC<ChatTextFieldProps> = ({ defaultText, enabled, focusInput }) => {
  const [characterCount, setCharacterCount] = useState(defaultText?.length);
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const text = useRef(defaultText || '');
  const contentEditable = React.createRef<HTMLElement>();
  const [customEmoji, setCustomEmoji] = useState([]);

  // This is a bit of a hack to force the component to re-render when the text changes.
  // By default when updating a ref the component doesn't re-render.
  const [, forceUpdate] = useReducer(x => x + 1, 0);

  const getCharacterCount = () => {
    const message = getTextContent(contentEditable.current);
    return graphemeSplitter.countGraphemes(message);
  };

  const sendMessage = () => {
    if (!websocketService) {
      console.log('websocketService is not defined');
      return;
    }

    const message = getTextContent(contentEditable.current);
    const count = graphemeSplitter.countGraphemes(message);
    if (count === 0 || count > characterLimit) return;

    websocketService.send({ type: MessageType.CHAT, body: message });

    // Clear the input.
    text.current = '';
    setCharacterCount(0);
    forceUpdate();
  };

  const insertTextAtEnd = (textToInsert: string) => {
    const output = text.current + textToInsert;
    text.current = output;

    forceUpdate();
  };

  // Native emoji
  const onEmojiSelect = (emoji: string) => {
    setCharacterCount(getCharacterCount() + 1);
    insertTextAtEnd(emoji);
  };

  // Custom emoji images
  const onCustomEmojiSelect = (name: string, emoji: string) => {
    const html = `<img src="${emoji}" alt=":${name}:" title=":${name}:" class="emoji" />`;
    setCharacterCount(getCharacterCount() + name.length + 2);
    insertTextAtEnd(html);
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !(e.shiftKey || e.metaKey || e.ctrlKey || e.altKey)) {
      e.preventDefault();
      sendMessage();
    }
  };

  const handleChange = evt => {
    const sanitized = sanitizeHtml(evt.target.value, {
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

    if (text.current !== sanitized) text.current = sanitized;

    setCharacterCount(getCharacterCount());
  };

  // Focus the input when the component mounts.
  useEffect(() => {
    if (!focusInput) {
      return;
    }
    document.getElementById('chat-input-content-editable').focus();
  }, []);

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
          html={text.current}
          placeholder={enabled ? 'Send a message to chat' : 'Chat is disabled'}
          disabled={!enabled}
          onKeyDown={onKeyDown}
          onChange={handleChange}
          style={{ width: '100%' }}
          role="textbox"
          aria-label="Chat text input"
          innerRef={contentEditable}
        />
        {enabled && (
          <div style={{ display: 'flex', paddingLeft: '5px' }}>
            <Popover
              content={
                <EmojiPicker
                  customEmoji={customEmoji}
                  onEmojiSelect={onEmojiSelect}
                  onCustomEmojiSelect={onCustomEmojiSelect}
                />
              }
              trigger="click"
              placement="topRight"
            >
              <button type="button" className={styles.emojiButton} title="Emoji picker button">
                <SmileOutlined />
              </button>
            </Popover>
            <button
              type="button"
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
