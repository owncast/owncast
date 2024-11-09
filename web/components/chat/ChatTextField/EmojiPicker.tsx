import React, { FC, useEffect, useState } from 'react';
import Picker from '@emoji-mart/react';
import data from '@emoji-mart/data';

export type EmojiPickerProps = {
  onEmojiSelect: (emoji) => void;
  customEmoji: any[];
};

export const EmojiPicker: FC<EmojiPickerProps> = ({ onEmojiSelect, customEmoji }) => {
  const [custom, setCustom] = useState({});
  const categories = [
    'frequent',
    'custom', // same id as in the setCustom call below
    'people',
    'nature',
    'foods',
    'activity',
    'places',
    'objects',
    'symbols',
    'flags',
  ];
  // Recreate the emoji picker when the custom emoji changes.
  useEffect(() => {
    const e = customEmoji.map(emoji => ({
      id: emoji.name,
      name: emoji.name,
      skins: [{ src: emoji.url }],
    }));

    setCustom([{ id: 'custom', name: 'Custom', emojis: e }]);

    // hack to make the picker work with viewbox only svgs, 24px is default size
    const shadow = document.querySelector('em-emoji-picker').shadowRoot;
    const pickerStyles = new CSSStyleSheet();
    pickerStyles.replaceSync('.emoji-mart-emoji {width: 24px;}');
    shadow.adoptedStyleSheets = [pickerStyles];
  }, []);

  return (
    <Picker data={data} custom={custom} onEmojiSelect={onEmojiSelect} categories={categories} />
  );
};
