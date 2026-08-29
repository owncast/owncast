import { fireEvent, render, screen } from '@testing-library/react';
import { BrowseRegistry, RegistryPlugin } from '../components/admin/plugins/BrowseRegistry';

// rc-select's useOpen schedules its dropdown-close through
// `new MessageChannel()` (port1.onmessage + port2.postMessage). Node's
// worker_threads-backed ports keep the jest process alive, so always
// replace the channel with a timer-based fake that covers just the
// surface rc-select uses.
class TestMessageChannel {
  port1: { onmessage: ((ev: { data: unknown }) => void) | null } = { onmessage: null };

  port2 = {
    postMessage: (data: unknown) => {
      setTimeout(() => this.port1.onmessage?.({ data }), 0);
    },
  };
}
globalThis.MessageChannel = TestMessageChannel as unknown as typeof MessageChannel;

jest.mock('next-export-i18n', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

const plugins: RegistryPlugin[] = [
  {
    slug: 'alpha',
    name: 'Alpha Plugin',
    authorName: 'Alice',
    category: 'chat-bots',
    latest: { version: '1.0.0', manifest: { permissions: ['chat.send', 'storage.kv'] } },
  },
  {
    // Category only inside the embedded manifest: exercises the
    // fallback for registries that predate the top-level field.
    slug: 'beta',
    name: 'Beta Plugin',
    authorName: 'Bob',
    latest: { version: '2.0.0', manifest: { permissions: ['chat.send'], category: 'viewer-ui' } },
  },
  {
    slug: 'gamma',
    name: 'Gamma Plugin',
    summary: 'Charts viewer counts over time.',
    latest: { version: '3.0.0' },
  },
];

const renderRegistry = (registry: RegistryPlugin[] = plugins) =>
  render(
    <BrowseRegistry
      installedVersions={new Map()}
      registry={registry}
      registryLoading={false}
      registryError={null}
      onInstall={async () => {}}
    />,
  );

// The t() mock is an identity function, so placeholders and option
// labels render as their raw Localization keys.
const PERMISSIONS_PLACEHOLDER = 'Admin.Plugins.browseFilterPermissions';
const AUTHOR_PLACEHOLDER = 'Admin.Plugins.browseFilterAuthor';
const CATEGORY_PLACEHOLDER = 'Admin.Plugins.browseFilterCategory';
const SEARCH_PLACEHOLDER = 'Admin.Plugins.browseSearch';

const openSelect = (placeholder: string) => {
  // antd v6 renders the placeholder inside .ant-select-content; the
  // dropdown opens on mousedown of the .ant-select root.
  fireEvent.mouseDown(screen.getByText(placeholder).closest('.ant-select') as HTMLElement);
};

const clickOption = (title: string) => {
  const option = document.querySelector(`.ant-select-item-option[title="${title}"]`);
  expect(option).not.toBeNull();
  fireEvent.click(option!);
};

const visiblePlugins = () =>
  ['Alpha Plugin', 'Beta Plugin', 'Gamma Plugin'].filter(name => screen.queryByText(name) !== null);

describe('BrowseRegistry filters', () => {
  it('renders all plugins and the filter bar', () => {
    renderRegistry();
    expect(visiblePlugins()).toEqual(['Alpha Plugin', 'Beta Plugin', 'Gamma Plugin']);
    expect(screen.getByText(PERMISSIONS_PLACEHOLDER)).toBeInTheDocument();
    expect(screen.getByText(AUTHOR_PLACEHOLDER)).toBeInTheDocument();
    expect(screen.getByText(CATEGORY_PLACEHOLDER)).toBeInTheDocument();
  });

  it('filters by permission requiring ALL selected permissions', () => {
    renderRegistry();
    openSelect(PERMISSIONS_PLACEHOLDER);
    clickOption('Admin.Plugins.PermissionNames.chatSend');
    expect(visiblePlugins()).toEqual(['Alpha Plugin', 'Beta Plugin']);
    clickOption('Admin.Plugins.PermissionNames.storageKv');
    expect(visiblePlugins()).toEqual(['Alpha Plugin']);
  });

  it('filters by author', () => {
    renderRegistry();
    openSelect(AUTHOR_PLACEHOLDER);
    clickOption('Bob');
    expect(visiblePlugins()).toEqual(['Beta Plugin']);
  });

  it('filters by category using taxonomy display names', () => {
    renderRegistry();
    openSelect(CATEGORY_PLACEHOLDER);
    clickOption('Admin.Plugins.Categories.chatBots');
    expect(visiblePlugins()).toEqual(['Alpha Plugin']);
  });

  it('searches by plugin name, case-insensitively', () => {
    renderRegistry();
    fireEvent.change(screen.getByPlaceholderText(SEARCH_PLACEHOLDER), {
      target: { value: 'BETA' },
    });
    expect(visiblePlugins()).toEqual(['Beta Plugin']);
  });

  it('searches by summary text and clears with the other filters', () => {
    renderRegistry();
    const search = screen.getByPlaceholderText(SEARCH_PLACEHOLDER);
    fireEvent.change(search, { target: { value: 'viewer counts' } });
    expect(visiblePlugins()).toEqual(['Gamma Plugin']);
    fireEvent.change(search, { target: { value: 'no such plugin' } });
    expect(visiblePlugins()).toEqual([]);
    fireEvent.click(screen.getByText('Admin.Plugins.browseClearFilters'));
    expect(visiblePlugins()).toEqual(['Alpha Plugin', 'Beta Plugin', 'Gamma Plugin']);
    expect((search as HTMLInputElement).value).toBe('');
  });

  it('shows the filtered-empty state and clears filters', () => {
    renderRegistry();
    // Alice + viewer UI matches nothing: Alice authored a chat-bot.
    openSelect(AUTHOR_PLACEHOLDER);
    clickOption('Alice');
    openSelect(CATEGORY_PLACEHOLDER);
    clickOption('Admin.Plugins.Categories.viewerUi');
    expect(visiblePlugins()).toEqual([]);
    expect(screen.getByText('Admin.Plugins.browseFilteredEmpty')).toBeInTheDocument();
    fireEvent.click(screen.getByText('Admin.Plugins.browseClearFilters'));
    expect(visiblePlugins()).toEqual(['Alpha Plugin', 'Beta Plugin', 'Gamma Plugin']);
  });

  it('hides the category select when no plugin carries a category', () => {
    renderRegistry(
      plugins.map(p => ({
        ...p,
        category: undefined,
        latest: p.latest && {
          ...p.latest,
          manifest: p.latest.manifest && { ...p.latest.manifest, category: undefined },
        },
      })),
    );
    expect(screen.queryByText(CATEGORY_PLACEHOLDER)).toBeNull();
    expect(screen.getByText(AUTHOR_PLACEHOLDER)).toBeInTheDocument();
  });

  it('renders a category tag on cards that have one', () => {
    renderRegistry([plugins[0]]);
    expect(screen.getByText('Admin.Plugins.Categories.chatBots')).toBeInTheDocument();
  });

  it('renders a viewer UI category tag from the embedded manifest fallback', () => {
    renderRegistry([plugins[1]]);
    expect(screen.getByText('Admin.Plugins.Categories.viewerUi')).toBeInTheDocument();
  });
});
