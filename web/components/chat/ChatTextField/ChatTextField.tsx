import { Popover } from 'antd';
import React, { FC, useEffect, useReducer, useRef, useState } from 'react';
import { useRecoilValue } from 'recoil';
import ContentEditable from 'react-contenteditable';
import sanitizeHtml from 'sanitize-html';

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

function getCaretPosition(node) {
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

function setCaretPosition(editableDiv, position) {
  try {
    const range = document.createRange();
    const sel = window.getSelection();
    range.selectNode(editableDiv);
    range.setStart(editableDiv.childNodes[0], position);
    range.collapse(true);

    sel.removeAllRanges();
    sel.addRange(range);
  } catch (e) {
    console.debug(e);
  }
}

export const ChatTextField: FC<ChatTextFieldProps> = ({ defaultText, enabled, focusInput }) => {
  const [characterCount, setCharacterCount] = useState(defaultText?.length);
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const text = useRef(defaultText || '');
  const [savedCursorLocation, setSavedCursorLocation] = useState(0);
  const [customEmoji, setCustomEmoji] = useState([]);

  // This is a bit of a hack to force the component to re-render when the text changes.
  // By default when updating a ref the component doesn't re-render.
  const [, forceUpdate] = useReducer(x => x + 1, 0);

  const getCharacterCount = () => text.current.length;

  const sendMessage = () => {
    if (!websocketService) {
      console.log('websocketService is not defined');
      return;
    }

    let message = text.current;
    // Strip the opening and closing <p> tags.
    message = message.replace(/^<p>|<\/p>$/g, '');
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
    insertTextAtEnd(emoji);
  };

  // Custom emoji images
  const onCustomEmojiSelect = (name: string, emoji: string) => {
    const html = `<img src="${emoji}" alt="${name}" title="${name}" class="emoji" />`;
    insertTextAtEnd(html);
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    // Allow native line breaks
    if (e.key === 'Enter' && e.shiftKey) {
      return;
    }

    const charCount = getCharacterCount() + 1;

    // Always allow backspace.
    if (e.key === 'Backspace') {
      return;
    }

    // Always allow delete.
    if (e.key === 'Delete') {
      return;
    }

    // Always allow ctrl + a.
    if (e.key === 'a' && e.ctrlKey) {
      return;
    }

    // Limit the number of characters.
    if (charCount + 1 > characterLimit) {
      e.preventDefault();
      return;
    }

    // Send the message when hitting enter.
    if (e.key === 'Enter') {
      e.preventDefault();
      sendMessage();
      return;
    }

    setCharacterCount(charCount + 1);
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
    text.current = sanitized;

    const charCountString = sanitized.replace(/<\/?[^>]+(>|$)/g, '');
    setCharacterCount(charCountString.length);

    setSavedCursorLocation(
      getCaretPosition(document.getElementById('chat-input-content-editable')),
    );
  };

  const handleBlur = () => {
    // Save the cursor location.
    setSavedCursorLocation(
      getCaretPosition(document.getElementById('chat-input-content-editable')),
    );
  };

  const handleFocus = () => {
    if (!savedCursorLocation) {
      return;
    }

    // Restore the cursor location.
    setCaretPosition(document.getElementById('chat-input-content-editable'), savedCursorLocation);
    setSavedCursorLocation(0);
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
          characterCount >= characterLimit && styles.maxCharacters,
        )}
      >
        <ContentEditable
          id="chat-input-content-editable"
          html={text.current}
          placeholder={enabled ? 'Send a message to chat' : 'Chat is disabled'}
          disabled={!enabled}
          onKeyDown={onKeyDown}
          onChange={handleChange}
          onBlur={handleBlur}
          onFocus={handleFocus}
          style={{ width: '100%' }}
          role="textbox"
          aria-label="Chat text input"
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
