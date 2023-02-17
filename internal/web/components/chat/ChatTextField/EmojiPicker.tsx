import React, { FC, useEffect, useRef, useState } from 'react';
import { createPicker } from 'picmo';

const CUSTOM_EMOJI_URL = '/api/emoji';
export type EmojiPickerProps = {
  onEmojiSelect: (emoji: string) => void;
  onCustomEmojiSelect: (name: string, url: string) => void;
};

export const EmojiPicker: FC<EmojiPickerProps> = ({ onEmojiSelect, onCustomEmojiSelect }) => {
  const [customEmoji, setCustomEmoji] = useState([]);
  const ref = useRef();

  const getCustomEmoji = async () => {
    try {
      const response = await fetch(CUSTOM_EMOJI_URL);
      const emoji = await response.json();
      setCustomEmoji(emoji);
    } catch (e) {
      console.error('cannot fetch custom emoji', e);
    }
  };

  // Fetch the custom emoji on component mount.
  useEffect(() => {
    getCustomEmoji();
  }, []);

  // Recreate the emoji picker when the custom emoji changes.
  useEffect(() => {
    const e = customEmoji.map(emoji => ({
      emoji: emoji.name,
      label: emoji.name,
      url: emoji.url,
    }));

    const picker = createPicker({
      rootElement: ref.current,
      custom: e,
      initialCategory: 'custom',
      showPreview: false,
      showRecents: true,
    });
    picker.addEventListener('emoji:select', event => {
      if (event.url) {
        onCustomEmojiSelect(event.name, event.url);
      } else {
        onEmojiSelect(event.emoji);
      }
    });
  }, [customEmoji]);

  return <div ref={ref} />;
};
