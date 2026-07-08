import { useEffect } from 'react';
import { StoryFn, Meta } from '@storybook/nextjs';
import { Provider, useSetAtom } from 'jotai';
import { UserDropdown } from './UserDropdown';
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

// This component reads jotai atoms internally so wrap it in a Provider.
const Example = args => {
  const setCurrentUser = useSetAtom(currentUserAtom);
  const setAppState = useSetAtom(appStateAtom);
  const setChatState = useSetAtom(chatStateAtom);

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
  <Provider>
    <Example {...args} />
  </Provider>
);

export const ChatEnabled = {
  render: Template,

  args: {
    username: 'test-user',
  },
};
