import React, { useEffect } from 'react';
import { ComponentStory, ComponentMeta } from '@storybook/react';
import { RecoilRoot, useSetRecoilState } from 'recoil';
import { NameChangeModal } from './NameChangeModal';
import { CurrentUser } from '../../../interfaces/current-user';
import { currentUserAtom } from '../../stores/ClientConfigStore';

export default {
  title: 'owncast/Modals/Name Change',
  component: NameChangeModal,
  parameters: {},
} as ComponentMeta<typeof NameChangeModal>;

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
      <NameChangeModal />
    </div>
  );
};

const Template: ComponentStory<typeof NameChangeModal> = () => (
  <RecoilRoot>
    <Example />
  </RecoilRoot>
);

export const Basic = Template.bind({});
