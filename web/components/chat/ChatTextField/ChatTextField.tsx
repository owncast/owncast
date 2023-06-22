import { Popover } from 'antd';
import React, { FC, useReducer, useRef, useState } from 'react';
import { useRecoilValue } from 'recoil';
import ContentEditable from 'react-contenteditable';
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

function convertToText(str = '') {
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
    .map((line = '') => line.trim())
    .join('\n');

  // No more than 2x newline, per "paragraph".
  value = value.replace(/\n\n+/g, '\n\n');

  // Clean up spaces.
  value = value.replace(/[ ]+/g, ' ');
  value = value.trim();

  // Expose string.
  return value;
}

export const ChatTextField: FC<ChatTextFieldProps> = ({ defaultText, enabled, focusInput }) => {
  const [showEmojis, setShowEmojis] = useState(false);
  const [characterCount, setCharacterCount] = useState(defaultText?.length);
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const text = useRef(defaultText || '');
  const [savedCursorLocation, setSavedCursorLocation] = useState(0);

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

  const insertTextAtCursor = (textToInsert: string) => {
    const output = [
      text.current.slice(0, savedCursorLocation),
      textToInsert,
      text.current.slice(savedCursorLocation),
    ].join('');
    text.current = output;
    forceUpdate();
  };

  const convertOnPaste = (event: React.ClipboardEvent) => {
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
    insertTextAtCursor(value);
  };

  // Native emoji
  const onEmojiSelect = (emoji: string) => {
    insertTextAtCursor(emoji);
  };

  // Custom emoji images
  const onCustomEmojiSelect = (name: string, emoji: string) => {
    const html = `<img src="${emoji}" alt="${name}" title=${name} class="emoji" />`;
    insertTextAtCursor(html);
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    const charCount = getCharacterCount() + 1;

    // Send the message when hitting enter.
    if (e.key === 'Enter') {
      e.preventDefault();
      sendMessage();
      return;
    }
    // Always allow backspace.
    if (e.key === 'Backspace') {
      setCharacterCount(charCount - 1);
      return;
    }

    // Always allow delete.
    if (e.key === 'Delete') {
      setCharacterCount(charCount - 1);
      return;
    }

    // Always allow ctrl + a.
    if (e.key === 'a' && e.ctrlKey) {
      return;
    }

    // Limit the number of characters.
    if (charCount + 1 > characterLimit) {
      e.preventDefault();
    }
    setCharacterCount(charCount + 1);
  };

  const handleChange = evt => {
    text.current = evt.target.value;
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

  return (
    <div id="chat-input" className={styles.root}>
      <div
        className={classNames(
          styles.inputWrap,
          characterCount >= characterLimit && styles.maxCharacters,
        )}
      >
        <Popover
          content={
            <EmojiPicker onEmojiSelect={onEmojiSelect} onCustomEmojiSelect={onCustomEmojiSelect} />
          }
          trigger="click"
          placement="topRight"
          onOpenChange={open => setShowEmojis(open)}
          open={showEmojis}
        />
        <ContentEditable
          id="chat-input-content-editable"
          html={text.current}
          placeholder={enabled ? 'Send a message to chat' : 'Chat is disabled'}
          disabled={!enabled}
          onKeyDown={onKeyDown}
          onPaste={convertOnPaste}
          onChange={handleChange}
          onBlur={handleBlur}
          onFocus={handleFocus}
          autoFocus={focusInput}
          style={{ width: '100%' }}
          role="textbox"
          aria-label="Chat text input"
        />
        {enabled && (
          <div style={{ display: 'flex', paddingLeft: '5px' }}>
            <button
              type="button"
              className={styles.emojiButton}
              title="Emoji picker button"
              onClick={() => setShowEmojis(!showEmojis)}
            >
              <SmileOutlined />
            </button>
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
