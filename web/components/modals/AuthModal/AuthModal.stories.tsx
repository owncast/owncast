import { useEffect } from 'react';
import { StoryFn, Meta } from '@storybook/nextjs';
import { Provider, useSetAtom } from 'jotai';
import { AuthModal } from './AuthModal';
import { currentUserAtom } from '../../stores/ClientConfigStore';

const Example = () => {
  const setCurrentUser = useSetAtom(currentUserAtom);

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
      <AuthModal open handleClose={null} forceTabs />
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
  <Provider>
    <Example />
  </Provider>
);

export const Basic = {
  render: Template,
};
