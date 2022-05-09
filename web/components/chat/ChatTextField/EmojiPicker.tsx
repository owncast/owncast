import data from '@emoji-mart/data';
import React, { useRef, useEffect } from 'react';

export default function EmojiPicker(props) {
  const ref = useRef();

  // TODO: Pull this custom emoji data in from the emoji API.
  const custom = [
    {
      emojis: [
        {
          id: 'party_parrot',
          name: 'Party Parrot',
          keywords: ['dance', 'dancing'],
          skins: [{ src: 'https://watch.owncast.online/img/emoji/bluntparrot.gif' }],
        },
      ],
    },
  ];

  // TODO: Fix the emoji picker from throwing errors.
  // useEffect(() => {
  //   import('emoji-mart').then(EmojiMart => {
  //     new EmojiMart.Picker({ ...props, data, ref });
  //   });
  // }, []);

  return <div ref={ref} />;
}
