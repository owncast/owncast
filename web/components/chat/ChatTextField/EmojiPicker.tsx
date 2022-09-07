import React, { useEffect, useRef, useState } from 'react';
import { createPicker } from 'picmo';

const CUSTOM_EMOJI_URL = '/api/emoji';
interface Props {
  // eslint-disable-next-line react/no-unused-prop-types
  onEmojiSelect: (emoji: string) => void;
  onCustomEmojiSelect: (emoji: string) => void;
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const EmojiPicker = (props: Props) => {
  const [customEmoji, setCustomEmoji] = useState([]);
  const { onEmojiSelect, onCustomEmojiSelect } = props;
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
        onCustomEmojiSelect(event);
      } else {
        onEmojiSelect(event);
      }
    });
  }, [customEmoji]);

  return <div ref={ref} />;
};
