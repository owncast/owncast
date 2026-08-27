import { FC, useEffect } from 'react';
import Picker from '@emoji-mart/react';
import data from '@emoji-mart/data';

export type EmojiPickerProps = {
  onEmojiSelect: (emoji) => void;
  customEmoji: any[];
};

export const EmojiPicker: FC<EmojiPickerProps> = ({ onEmojiSelect, customEmoji }) => {
  const custom = [
    {
      id: 'custom',
      name: 'Custom',
      emojis: customEmoji.map(emoji => ({
        id: emoji.name,
        name: emoji.name,
        skins: [{ src: emoji.url }],
      })),
    },
  ];
  useEffect(() => {
    // hack to make the picker work with viewbox only svgs, 24px is default size
    const shadow = document.querySelector('em-emoji-picker').shadowRoot;
    const pickerStyles = new CSSStyleSheet();
    pickerStyles.replaceSync('.emoji-mart-emoji {width: 24px;}');
    shadow.adoptedStyleSheets = [pickerStyles];
  }, []);

  return <Picker data={data} custom={custom} onEmojiSelect={onEmojiSelect} dynamicWidth />;
};
