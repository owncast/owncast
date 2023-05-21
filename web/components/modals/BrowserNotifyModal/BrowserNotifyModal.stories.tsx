import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { BrowserNotifyModal } from './BrowserNotifyModal';
import BrowserNotifyModalMock from '../../../stories/assets/mocks/notify-modal.png';

const Example = () => (
  <div>
    <BrowserNotifyModal />
  </div>
);

export default {
  title: 'owncast/Modals/Browser Notifications',
  component: BrowserNotifyModal,
  parameters: {
    chromatic: { diffThreshold: 0.9 },
    design: {
      type: 'image',
      url: BrowserNotifyModalMock,
      scale: 0.5,
    },
    docs: {
      description: {
        component: `The notify modal allows an end user to get notified when the stream goes live via [Browser Push Notifications](https://developers.google.com/web/ilt/pwa/introduction-to-push-notifications) It must:

- Verify the browser supports notifications.
- Handle errors that come back from the server.
- Have an enabled and disabled state with accurate information about each.
`,
      },
    },
  },
} as ComponentMeta<typeof BrowserNotifyModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof BrowserNotifyModal> = () => (
  <RecoilRoot>
    <Example />
  </RecoilRoot>
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
