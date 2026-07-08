import { Meta, StoryObj } from '@storybook/nextjs';
import { MainLayout } from './MainLayout';
import { ServerStatusContext } from '../../utils/server-status-context';

// The admin chrome itself: sider menu, header with stream title editor and
// status indicator, footer. Whole-page smoke stories render without the
// layout, so this is the only story that covers admin navigation.
//
// The Consumer-inside-Provider trick reuses the context's built-in default
// value as the fixture base without exporting the initial state objects.
const placeholderContent = (
  <div style={{ padding: '24px' }}>
    <h1>Page content</h1>
    <p>The active admin page renders here.</p>
  </div>
);

const OnlineFixture = () => (
  <ServerStatusContext.Consumer>
    {defaults => (
      <ServerStatusContext.Provider
        value={{
          ...defaults,
          online: true,
          viewerCount: 12,
          serverConfig: {
            ...defaults.serverConfig,
            instanceDetails: {
              ...defaults.serverConfig.instanceDetails,
              name: 'Demo stream',
              streamTitle: 'Weekly community show',
            },
            federation: {
              ...defaults.serverConfig.federation,
              enabled: true,
            },
          },
        }}
      >
        <MainLayout>{placeholderContent}</MainLayout>
      </ServerStatusContext.Provider>
    )}
  </ServerStatusContext.Consumer>
);

const meta = {
  title: 'owncast/Admin/Main layout',
  component: MainLayout,
  decorators: [
    Story => (
      <>
        {/* In the app AdminLayout owns the admin stylesheets and renders
            them as head links. This story renders MainLayout directly, so
            pull in the chrome stylesheet itself (served from web/public
            via the storybook staticDirs). */}
        <link rel="stylesheet" href="/styles/admin/main-layout.css" />
        <Story />
      </>
    ),
  ],
} satisfies Meta<typeof MainLayout>;

export default meta;
type Story = StoryObj<typeof meta>;

// Default context: offline server, empty config.
export const Offline: Story = {
  args: { children: placeholderContent },
};

// Online with federation enabled, which adds the Fediverse menu section and
// the compose button.
export const OnlineWithFederation: Story = {
  args: { children: placeholderContent },
  render: () => <OnlineFixture />,
};
