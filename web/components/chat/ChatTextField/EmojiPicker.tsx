import React, { FC, useEffect, useRef } from 'react';
import { createPicker } from 'picmo';

export type EmojiPickerProps = {
  onEmojiSelect: (emoji: string) => void;
  onCustomEmojiSelect: (name: string, url: string) => void;
  customEmoji: any[];
};

export const EmojiPicker: FC<EmojiPickerProps> = ({
  onEmojiSelect,
  onCustomEmojiSelect,
  customEmoji,
}) => {
  const ref = useRef();

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
        onCustomEmojiSelect(event.label, event.url);
      } else {
        onEmojiSelect(event.emoji);
      }
    });
  }, []);

  return <div ref={ref} />;
};
