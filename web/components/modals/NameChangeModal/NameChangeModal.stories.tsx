import { useEffect } from 'react';
import { StoryFn, Meta } from '@storybook/nextjs';
import { Provider, useSetAtom } from 'jotai';
import { NameChangeModal } from './NameChangeModal';
import { currentUserAtom } from '../../stores/ClientConfigStore';

const meta = {
  title: 'owncast/Modals/Name Change',
  component: NameChangeModal,
  parameters: {},
} satisfies Meta<typeof NameChangeModal>;

export default meta;

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
      <NameChangeModal closeModal={() => {}} />
    </div>
  );
};

const Template: StoryFn<typeof NameChangeModal> = () => (
  <Provider>
    <Example />
  </Provider>
);

export const Basic = {
  render: Template,
};
