import { Meta, StoryFn } from '@storybook/nextjs';
import { Provider } from 'jotai';
import { MobileContent, MobileContentProps } from './MobileContent';
import { PluginTab } from '../../../interfaces/client-config.model';

const makePluginTabs = (count: number): PluginTab[] =>
  Array.from({ length: count }, (_, i) => ({
    slug: `demo-plugin/tab-${i + 1}`,
    pluginSlug: 'demo-plugin',
    title: `Plugin Tab ${i + 1}`,
    html: `<h1>Content of plugin tab ${i + 1}</h1><p>Rendered inside the plugin tab iframe.</p>`,
  }));

const meta = {
  title: 'owncast/Components/Mobile content',
  component: MobileContent,
  tags: ['screenshot-mobile'],
  globals: {
    viewport: { value: 'iphone12Portrait', isRotated: false },
  },
  parameters: {
    // Chromatic ignores the storybook viewport global, so the capture width
    // must be pinned separately or these render at desktop width.
    chromatic: { viewports: [390] },
    docs: {
      description: {
        component:
          'The tabbed lower section of the viewer page on mobile. When plugins ' +
          'contribute enough tabs to overflow the screen width, an all-tabs menu ' +
          'button appears at the right end of the tab bar so every tab stays ' +
          'reachable without swipe-scrolling. Pick a mobile viewport to see the ' +
          'overflow behavior.',
      },
    },
  },
} satisfies Meta<typeof MobileContent>;

export default meta;

// PluginTabFrame reads plugin styles from jotai atoms, so the component
// needs a Provider around it.
const Template: StoryFn<MobileContentProps> = args => (
  <Provider>
    <MobileContent {...args} />
  </Provider>
);

const defaultArgs: MobileContentProps = {
  name: 'Demo stream',
  summary: 'A live community stream with conversation, music, and the occasional guest.',
  tags: ['live music', 'community', 'weekly show'],
  socialHandles: [],
  extraPageContent: `
    <h2>Tonight’s broadcast</h2>
    <p>Join us for a relaxed live set, followed by a chat with local artists. New listeners are always welcome.</p>
    <p>Schedule, guest details, and community links are available at <a href="https://demo.example.org/schedule-and-special-guest-announcements">demo.example.org/schedule-and-special-guest-announcements</a>.</p>
    <h3>What to expect</h3>
    <p>Music, conversation, and a place to meet people who care about the same things.</p>
  `,
  pluginTabs: [],
  showFollowersTab: false,
  online: false,
  federatedServers: [],
};

// Enough plugin tabs to overflow any phone-sized viewport: the tab bar
// scrolls and the all-tabs menu button appears next to it.
export const ManyPluginTabs = Template.bind({});
ManyPluginTabs.args = {
  ...defaultArgs,
  pluginTabs: makePluginTabs(8),
};

// A single plugin tab alongside About: everything fits, so the all-tabs
// menu button stays hidden.
export const OnePluginTab = Template.bind({});
OnePluginTab.args = {
  ...defaultArgs,
  pluginTabs: makePluginTabs(1),
};

// No plugin tabs and only one built-in tab: no tab bar is rendered at all.
export const NoTabs = Template.bind({});
NoTabs.args = {
  ...defaultArgs,
};
