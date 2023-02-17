import React, { useEffect } from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot, useSetRecoilState } from 'recoil';
import { AuthModal } from './AuthModal';
import { currentUserAtom } from '../../stores/ClientConfigStore';
import { CurrentUser } from '../../../interfaces/current-user';

const Example = () => {
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

  return (
    <div>
      <AuthModal forceTabs />
    </div>
  );
};

export default {
  title: 'owncast/Modals/Auth',
  component: AuthModal,
  parameters: {},
} as ComponentMeta<typeof AuthModal>;

const Template: ComponentStory<typeof AuthModal> = () => (
  <RecoilRoot>
    <Example />
  </RecoilRoot>
);

// eslint-disable-next-line @typescript-eslint/no-unused-vars
export const Basic = Template.bind({});
