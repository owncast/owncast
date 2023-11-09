import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot } from 'recoil';
import { BrowserNotifyModal } from './BrowserNotifyModal';
import BrowserNotifyModalMock from '../../../stories/assets/mocks/notify-modal.png';

const Example = () => (
  <div>
    <BrowserNotifyModal />
  </div>
);

const meta = {
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
} satisfies Meta<typeof BrowserNotifyModal>;

export default meta;

const Template: StoryFn<typeof BrowserNotifyModal> = () => (
  <RecoilRoot>
    <Example />
  </RecoilRoot>
);

export const Basic = {
  render: Template,
};
