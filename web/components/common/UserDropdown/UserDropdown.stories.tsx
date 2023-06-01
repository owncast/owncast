import React, { useEffect } from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot, useSetRecoilState } from 'recoil';
import { UserDropdown } from './UserDropdown';
import { CurrentUser } from '../../../interfaces/current-user';
import { currentUserAtom } from '../../stores/ClientConfigStore';

export default {
  title: 'owncast/Components/User settings menu',
  component: UserDropdown,
  parameters: {},
} as ComponentMeta<typeof UserDropdown>;

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

const Template: ComponentStory<typeof UserDropdown> = args => (
  <RecoilRoot>
    <Example {...args} />
  </RecoilRoot>
);

export const ChatEnabled = Template.bind({});
ChatEnabled.args = {
  username: 'test-user',
};
