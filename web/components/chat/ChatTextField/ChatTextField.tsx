import React, { FC, useMemo, useState } from 'react';
import dynamic from 'next/dynamic';
import { useRecoilValue } from 'recoil';
import { Transforms, createEditor, BaseEditor, Text, Descendant, Editor, Node, Path } from 'slate';
import { Slate, Editable, withReact, ReactEditor, useSelected, useFocused } from 'slate-react';
import { MessageType } from '~/interfaces/socket-events';
import { Popover } from 'antd';
import { SendOutlined, SmileOutlined } from '@ant-design/icons';
import WebsocketService from '~/services/websocket-service';
import { websocketServiceAtom } from '~/components/stores/ClientConfigStore';
import styles from './ChatTextField.module.scss';

// Lazy loaded components

const EmojiPicker = dynamic(() => import('./EmojiPicker').then(mod => mod.EmojiPicker));

type CustomElement = { type: 'paragraph' | 'span'; children: CustomText[] } | ImageNode;
type CustomText = { text: string };

type EmptyText = {
  text: string;
};

type ImageNode = {
  type: 'image';
  alt: string;
  src: string;
  name: string;
  children: EmptyText[];
};

declare module 'slate' {
  interface CustomTypes {
    Editor: BaseEditor & ReactEditor;
    Element: CustomElement;
    Text: CustomText;
  }
}

const Image = p => {
  const { attributes, element, children } = p;

  const selected = useSelected();
  const focused = useFocused();
  return (
    <span {...attributes} contentEditable={false}>
      <img
        alt={element.alt}
        src={element.src}
        title={element.name}
        style={{
          display: 'inline',
          maxWidth: '50px',
          maxHeight: '20px',
          boxShadow: `${selected && focused ? '0 0 0 3px #B4D5FF' : 'none'}`,
        }}
      />
      {children}
    </span>
  );
};

const withImages = editor => {
  const { isVoid } = editor;

  // eslint-disable-next-line no-param-reassign
  editor.isVoid = element => (element.type === 'image' ? true : isVoid(element));
  // eslint-disable-next-line no-param-reassign
  editor.isInline = element => element.type === 'image';

  return editor;
};

const serialize = node => {
  if (Text.isText(node)) {
    const string = node.text;
    return string;
  }

  let children;
  if (node.children.length === 0) {
    children = [{ text: '' }];
  } else {
    children = node.children?.map(n => serialize(n)).join('');
  }

  switch (node.type) {
    case 'paragraph':
      return `<p>${children}</p>`;
    case 'image':
      return `<img src="${node.src}" alt="${node.alt}" title="${node.name}" class="emoji"/>`;
    default:
      return children;
  }
};

export type ChatTextFieldProps = {};

export const ChatTextField: FC<ChatTextFieldProps> = () => {
  const [showEmojis, setShowEmojis] = useState(false);
  const websocketService = useRecoilValue<WebsocketService>(websocketServiceAtom);
  const editor = useMemo(() => withReact(withImages(createEditor())), []);

  const defaultEditorValue: Descendant[] = [
    {
      type: 'paragraph',
      children: [{ text: '' }],
    },
  ];

  const sendMessage = () => {
    if (!websocketService) {
      console.log('websocketService is not defined');
      return;
    }

    const message = serialize(editor);
    websocketService.send({ type: MessageType.CHAT, body: message });

    // Clear the editor.
    Transforms.delete(editor, {
      at: {
        anchor: Editor.start(editor, []),
        focus: Editor.end(editor, []),
      },
    });
  };

  const createImageNode = (alt, src, name): ImageNode => ({
    type: 'image',
    alt,
    src,
    name,
    children: [{ text: '' }],
  });

  const insertImage = (url, name) => {
    if (!url) return;

    const { selection } = editor;
    const image = createImageNode(name, url, name);

    Transforms.insertNodes(editor, image, { select: true });

    if (selection) {
      const [parentNode, parentPath] = Editor.parent(editor, selection.focus?.path);

      if (editor.isVoid(parentNode) || Node.string(parentNode).length) {
        // Insert the new image node after the void node or a node with content
        Transforms.insertNodes(editor, image, {
          at: Path.next(parentPath),
          select: true,
        });
      } else {
        // If the node is empty, replace it instead
        // Transforms.removeNodes(editor, { at: parentPath });
        Transforms.insertNodes(editor, image, { at: parentPath, select: true });
        Editor.normalize(editor, { force: true });
      }
    } else {
      // Insert the new image node at the bottom of the Editor when selection
      // is falsey
      Transforms.insertNodes(editor, image, { select: true });
    }
  };

  const onEmojiSelect = (e: any) => {
    ReactEditor.focus(editor);

    if (e.url) {
      // Custom emoji
      const { url } = e;
      insertImage(url, url);
    } else {
      // Native emoji
      const { emoji } = e;
      Transforms.insertText(editor, emoji);
    }
  };

  const onCustomEmojiSelect = (e: any) => {
    ReactEditor.focus(editor);
    const { url } = e;
    insertImage(url, url);
  };

  const onKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      sendMessage();
    }
  };

  const renderElement = p => {
    switch (p.element.type) {
      case 'image':
        return <Image {...p} />;
      default:
        return <p {...p} />;
    }
  };

  return (
    <div className={styles.root}>
      <div className={styles.inputWrap}>
        <Slate editor={editor} value={defaultEditorValue}>
          <Editable
            onKeyDown={onKeyDown}
            renderElement={renderElement}
            placeholder="Chat message goes here..."
            style={{ width: '100%' }}
            autoFocus
          />
          <Popover
            content={
              <EmojiPicker
                onEmojiSelect={onEmojiSelect}
                onCustomEmojiSelect={onCustomEmojiSelect}
              />
            }
            trigger="click"
            onVisibleChange={visible => setShowEmojis(visible)}
            visible={showEmojis}
          />
        </Slate>

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
      </div>
    </div>
  );
};
