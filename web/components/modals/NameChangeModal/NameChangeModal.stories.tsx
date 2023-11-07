import { useEffect } from 'react';
import { StoryFn, Meta } from '@storybook/react';
import { RecoilRoot, useSetRecoilState } from 'recoil';
import { NameChangeModal } from './NameChangeModal';
import { CurrentUser } from '../../../interfaces/current-user';
import { currentUserAtom } from '../../stores/ClientConfigStore';

const meta = {
  title: 'owncast/Modals/Name Change',
  component: NameChangeModal,
  parameters: {},
} satisfies Meta<typeof NameChangeModal>;

export default meta;

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
      <NameChangeModal closeModal={() => {}} />
    </div>
  );
};

const Template: StoryFn<typeof NameChangeModal> = () => (
  <RecoilRoot>
    <Example />
  </RecoilRoot>
);

export const Basic = {
  render: Template,
};
