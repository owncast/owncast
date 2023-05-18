import React from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { FollowModal } from './FollowModal';
import FollowModalMock from '../../../stories/assets/mocks/follow-modal.png';

const Example = () => (
  <div>
    <FollowModal handleClose={null} account="@fake@server.name" name="Fake Owncast Server" />
  </div>
);

export default {
  title: 'owncast/Modals/Follow',
  component: FollowModal,
  parameters: {
    design: {
      type: 'image',
      url: FollowModalMock,
      scale: 0.5,
    },
    docs: {
      description: {
        component: `The Follow modal allows an end user to type in their Fediverse account information to follow this Owncast instance. It must:

- Validate the input to make sure it's a valid looking account.
- Handle errors that come back from the server.
- Perform the redirect to the remote server when the backend response is received.
`,
      },
    },
  },
} as ComponentMeta<typeof FollowModal>;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
const Template: ComponentStory<typeof FollowModal> = () => <Example />;

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
