import { useEffect } from 'react';
import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot, useSetRecoilState } from 'recoil';
import { UserDropdown } from './UserDropdown';
import { CurrentUser } from '../../../interfaces/current-user';
import { currentUserAtom } from '../../stores/ClientConfigStore';

const meta = {
  title: 'owncast/Components/User settings menu',
  component: UserDropdown,
  parameters: {},
} satisfies Meta<typeof UserDropdown>;

export default meta;

// This component uses Recoil internally so wrap it in a RecoilRoot.
const Example = args => {
  const setCurrentUser = useSetRecoilState<CurrentUser>(currentUserAtom);

  useEffect(
    () =>
      setCurrentUser({
        id: '1',
        displayName: 'Test User',
        displayColor: 3,
        isModerator: false,
      }),
    [],
  );

  return <UserDropdown id="user-menu" {...args} />;
};

const Template: StoryFn<typeof UserDropdown> = args => (
  <RecoilRoot>
    <Example {...args} />
  </RecoilRoot>
);

export const ChatEnabled = {
  render: Template,

  args: {
    username: 'test-user',
  },
};
