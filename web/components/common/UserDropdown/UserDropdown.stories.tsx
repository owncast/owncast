import { useEffect } from 'react';
import { StoryFn, Meta } from '@storybook/nextjs';
import { RecoilRoot, useSetRecoilState } from 'recoil';
import { UserDropdown } from './UserDropdown';
import { CurrentUser } from '../../../interfaces/current-user';
import {
  currentUserAtom,
  appStateAtom,
  chatStateAtom,
  ChatState,
} from '../../stores/ClientConfigStore';

const meta = {
  title: 'owncast/Components/User settings menu',
  component: UserDropdown,
  parameters: {},
} satisfies Meta<typeof UserDropdown>;

export default meta;

// This component uses Recoil internally so wrap it in a RecoilRoot.
const Example = args => {
  const setCurrentUser = useSetRecoilState<CurrentUser>(currentUserAtom);
  const setAppState = useSetRecoilState(appStateAtom);
  const setChatState = useSetRecoilState(chatStateAtom);

  useEffect(() => {
    setCurrentUser({
      id: '1',
      displayName: 'Test User',
      displayColor: 3,
      isModerator: false,
    });

    setAppState({
      chatAvailable: true,
      chatLoading: false,
      videoAvailable: true,
      appLoading: false,
    });

    setChatState(ChatState.VISIBLE);
  }, []);

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
