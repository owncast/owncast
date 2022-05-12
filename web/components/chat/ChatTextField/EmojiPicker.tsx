// import data from '@emoji-mart/data';
import React, { useRef } from 'react';

interface Props {
  // eslint-disable-next-line react/no-unused-prop-types
  onEmojiSelect: (emoji: string) => void;
}

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export default function EmojiPicker(props: Props) {
  const ref = useRef();

  // TODO: Pull this custom emoji data in from the emoji API.
  // const custom = [
  //   {
  //     emojis: [
  //       {
  //         id: 'party_parrot',
  //         name: 'Party Parrot',
  //         keywords: ['dance', 'dancing'],
  //         skins: [{ src: 'https://watch.owncast.online/img/emoji/bluntparrot.gif' }],
  //       },
  //     ],
  //   },
  // ];

  // TODO: Fix the emoji picker from throwing errors.
  // useEffect(() => {
  //   import('emoji-mart').then(EmojiMart => {
  //     new EmojiMart.Picker({ ...props, data, ref });
  //   });
  // }, []);

  return <div ref={ref} />;
}
