import { Meta, StoryFn, StoryObj } from '@storybook/nextjs';
import { Provider, createStore } from 'jotai';
import type { Store } from 'jotai/vanilla/store';
import { makeEmptyClientConfig } from '../../../interfaces/client-config.model';
import { ServerStatus, makeEmptyServerStatus } from '../../../interfaces/server-status.model';
import {
  accessTokenAtom,
  appStateAtom,
  chatMessagesAtom,
  ChatState,
  chatStateAtom,
  clientConfigStateAtom,
  currentUserAtom,
  fatalErrorStateAtom,
  isMobileAtom,
  isVideoPlayingAtom,
  serverStatusState,
} from '../../stores/ClientConfigStore';
import { Main } from './Main';
import { ClientConfigServiceContext } from '../../../services/client-config-service';
import { ChatServiceContext } from '../../../services/chat-service';
import {
  ServerStatusServiceContext,
  ServerStatusStaticService,
} from '../../../services/status-service';
import { clientConfigServiceMockOf } from '../../../services/client-config-service.mock';
import chatServiceMockOf from '../../../services/chat-service.mock';
import serverStatusServiceMockOf from '../../../services/status-service.mock';
import { VideoSettingsServiceContext } from '../../../services/video-settings-service';
import videoSettingsServiceMockOf from '../../../services/video-settings-service.mock';
import { spidermanUser } from '../../../interfaces/user.fixture';
import { exampleChatHistory } from '../../../interfaces/chat-message.fixture';

const meta = {
  title: 'owncast/Layout/Main',
} satisfies Meta<typeof Main>;

export default meta;

// mock the Websocket to prevent ani webhook calls from being made in storybook
// @ts-ignore
window.WebSocket = class {};

type StateInitializer = (store: Store) => void;

const composeStateInitializers =
  (...fns: Array<StateInitializer>): StateInitializer =>
  state =>
    fns.forEach(fn => fn?.(state));

const defaultClientConfig = {
  ...makeEmptyClientConfig(),
  logo: 'http://localhost:8080/logo',
  name: "Spiderman's super serious stream",
  summary: 'Strong Spidey stops supervillains! Streamed Saturdays & Sundays.',
  extraPageContent: 'Spiderman is **cool**',
};

const defaultServerStatus = makeEmptyServerStatus();
const onlineServerStatus: ServerStatus = {
  ...defaultServerStatus,
  online: true,
  viewerCount: 5,
};

const initializeDefaultState = (mutableState: Store) => {
  mutableState.set(appStateAtom, {
    videoAvailable: false,
    chatAvailable: false,
  });
  mutableState.set(clientConfigStateAtom, defaultClientConfig);
  mutableState.set(chatStateAtom, ChatState.VISIBLE);
  mutableState.set(accessTokenAtom, 'token');
  mutableState.set(currentUserAtom, {
    ...spidermanUser,
    isModerator: false,
  });
  mutableState.set(serverStatusState, defaultServerStatus);
  mutableState.set(isMobileAtom, false);

  mutableState.set(chatMessagesAtom, exampleChatHistory);
  mutableState.set(isVideoPlayingAtom, false);
  mutableState.set(fatalErrorStateAtom, null);
};

const ClientConfigServiceMock = clientConfigServiceMockOf(defaultClientConfig);
const ChatServiceMock = chatServiceMockOf(exampleChatHistory, {
  ...spidermanUser,
  accessToken: 'some fake token',
});
const DefaultServerStatusServiceMock = serverStatusServiceMockOf(defaultServerStatus);
const OnlineServerStatusServiceMock = serverStatusServiceMockOf(onlineServerStatus);
const VideoSettingsServiceMock = videoSettingsServiceMockOf([]);

const Template: StoryFn<typeof Main> = ({
  initializeState,
  ServerStatusServiceMock = DefaultServerStatusServiceMock,
  ...args
}: {
  initializeState: (store: Store) => void;
  ServerStatusServiceMock: ServerStatusStaticService;
}) => {
  // A fresh store per story render replaces RecoilRoot's initializeState:
  // state is seeded before anything renders, so the mocked values never
  // flash their defaults.
  const store = createStore();
  composeStateInitializers(initializeDefaultState, initializeState)(store);
  return (
    <Provider store={store}>
      <ClientConfigServiceContext.Provider value={ClientConfigServiceMock}>
        <ChatServiceContext.Provider value={ChatServiceMock}>
          <ServerStatusServiceContext.Provider value={ServerStatusServiceMock}>
            <VideoSettingsServiceContext.Provider value={VideoSettingsServiceMock}>
              <Main {...args} />
            </VideoSettingsServiceContext.Provider>
          </ServerStatusServiceContext.Provider>
        </ChatServiceContext.Provider>
      </ClientConfigServiceContext.Provider>
    </Provider>
  );
};

export const OfflineDesktop: StoryObj<typeof Template> = {
  render: Template,

  parameters: {
    chromatic: { diffThreshold: 0.88 },
  },
};

export const OfflineMobile: StoryObj<typeof Template> = {
  render: Template,

  args: {
    initializeState: (mutableState: Store) => {
      mutableState.set(isMobileAtom, true);
    },
  },

  // Content re-derives isMobile from window.innerWidth on mount, so the
  // Chromatic capture window must actually be mobile sized or the pinned
  // atom above gets flipped back to desktop.
  parameters: {
    chromatic: { viewports: [375] },
  },

  globals: {
    viewport: {
      value: 'mobile1',
      isRotated: false,
    },
  },
};

export const OfflineTablet: StoryObj<typeof Template> = {
  render: Template,

  // Chromatic ignores the storybook viewport global, so the capture width
  // must be pinned separately or these render at desktop width.
  parameters: {
    chromatic: { viewports: [834] },
  },

  globals: {
    viewport: {
      value: 'tablet',
      isRotated: false,
    },
  },
};

export const Online: StoryObj<typeof Template> = {
  render: Template,

  args: {
    ServerStatusServiceMock: OnlineServerStatusServiceMock,
  },

  parameters: {
    chromatic: { diffThreshold: 0.88 },
  },
};

export const OnlineMobile: StoryObj<typeof Template> = {
  render: Template,

  args: {
    ServerStatusServiceMock: OnlineServerStatusServiceMock,
    initializeState: (mutableState: Store) => {
      mutableState.set(isMobileAtom, true);
    },
  },

  parameters: {
    chromatic: { viewports: [375] },
  },

  globals: {
    viewport: {
      value: 'mobile1',
      isRotated: false,
    },
  },
};

export const OnlineTablet: StoryObj<typeof Template> = {
  render: Template,

  args: {
    ServerStatusServiceMock: OnlineServerStatusServiceMock,
  },

  parameters: {
    chromatic: { viewports: [834] },
  },

  globals: {
    viewport: {
      value: 'tablet',
      isRotated: false,
    },
  },
};
