import { useEffect } from 'react';
import { StoryFn, Meta } from '@storybook/react';
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

const meta = {
  title: 'owncast/Modals/Auth',
  component: AuthModal,
  parameters: {},
} satisfies Meta<typeof AuthModal>;

export default meta;

const Template: StoryFn<typeof AuthModal> = () => (
  <RecoilRoot>
    <Example />
  </RecoilRoot>
);

export const Basic = {
  render: Template,
};
