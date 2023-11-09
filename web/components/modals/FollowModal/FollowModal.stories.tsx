import { StoryFn, Meta } from '@storybook/react';
import { FollowModal } from './FollowModal';
import FollowModalMock from '../../../stories/assets/mocks/follow-modal.png';

const Example = () => (
  <div>
    <FollowModal handleClose={null} account="@fake@server.name" name="Fake Owncast Server" />
  </div>
);

const meta = {
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
} satisfies Meta<typeof FollowModal>;

export default meta;

const Template: StoryFn<typeof FollowModal> = () => <Example />;

export const Basic = {
  render: Template,
};
