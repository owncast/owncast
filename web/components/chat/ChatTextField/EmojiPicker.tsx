import React, { useEffect, useRef, useState } from 'react';
import { createPicker } from 'picmo';

const CUSTOM_EMOJI_URL = '/api/emoji';
interface Props {
  // eslint-disable-next-line react/no-unused-prop-types
  onEmojiSelect: (emoji: string) => void;
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const EmojiPicker = (props: Props) => {
  const [customEmoji, setCustomEmoji] = useState([]);
  const { onEmojiSelect } = props;
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
    const picker = createPicker({ rootElement: ref.current, custom: customEmoji });
    picker.addEventListener('emoji:select', event => {
      console.log('Emoji selected:', event.emoji);
      onEmojiSelect(event);
    });
  }, [customEmoji]);

  return <div ref={ref} />;
};
export default EmojiPicker;
