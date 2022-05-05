import data from '@emoji-mart/data';
import { Picker } from 'emoji-mart';
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

  // useEffect(() => {
  //   // eslint-disable-next-line no-new
  //   new Picker({ ...props, data, custom, ref });
  // }, []);

  return <div>emoji picker goes here</div>;
}
