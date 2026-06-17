import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { useFederatedServers } from '../hooks/useFederatedServers';
import { StreamsTab } from '../components/ui/StreamsTab/StreamsTab';
import { FederatedServersTable } from '../components/admin/FederatedServers/FederatedServersTable';

// next-export-i18n's useTranslation reads from the router; the components and
// the hook both call it.
jest.mock('next/router', () => ({
  useRouter: () => ({
    query: {},
    pathname: '/',
    asPath: '/',
    push: jest.fn(),
    replace: jest.fn(),
  }),
}));

jest.mock('antd', () => ({
  ...jest.requireActual('antd'),
  message: { success: jest.fn(), error: jest.fn() },
}));

// Render the Translation component as its default text so assertions can match
// on stable, language-independent strings.
jest.mock('../components/ui/Translation/Translation', () => ({
  Translation: ({ defaultText }: { defaultText: string }) => <span>{defaultText}</span>,
}));

// This payload deliberately uses the EXACT field names the Go backend
// serializes for the OpenAPI FederatedServer model. The whole point of this
// test is to fail if the web ever reads a different set of names than the API
// actually returns -- which is precisely the bug that previously shipped (the
// web read url/logo/thumbnail/lastChecked while the API sent
// iri/logoUrl/thumbnailUrl/lastStatusUpdate).
const API_SERVER = {
  id: 7,
  iri: 'https://goodnight.example.com',
  name: 'goodnight', // federation username
  displayName: 'Goodnight TV', // human-friendly name
  logoUrl: 'https://goodnight.example.com/logo.png',
  isOnline: true,
  streamTitle: 'Late night coding',
  streamDescription: 'Chill coding vibes',
  tags: ['coding', 'chill'],
  thumbnailUrl: 'https://goodnight.example.com/thumb.jpg',
  lastSeenOnline: '2026-06-16T21:50:00Z',
  lastStatusUpdate: '2026-06-16T21:55:00Z',
  addedAt: '2026-06-16T21:48:50Z',
  followedAt: '2026-06-16T21:48:50Z',
  pending: false,
  username: 'goodnight',
  summary: 'Goodnight TV',
  followStatus: 'accepted',
};

function mockFetchServers(server: object = API_SERVER) {
  global.fetch = jest.fn().mockResolvedValue({
    ok: true,
    statusText: 'OK',
    json: async () => ({ servers: [server] }),
  }) as unknown as typeof fetch;
}

// Harnesses wire the real hook to the real components so the test exercises
// the full path: API JSON shape -> useFederatedServers -> rendered component.
const PublicHarness = () => {
  const { servers, loading, error } = useFederatedServers();
  return <StreamsTab servers={servers} loading={loading} error={error || undefined} />;
};

const AdminHarness = () => {
  const { servers, loading, removeServer } = useFederatedServers(true);
  return <FederatedServersTable servers={servers} loading={loading} onRemove={removeServer} />;
};

describe('useFederatedServers contract', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it('requests the public federation listing endpoint', async () => {
    mockFetchServers();
    render(<PublicHarness />);
    await waitFor(() => expect(global.fetch).toHaveBeenCalledWith('/api/federation/servers'));
  });

  describe('public StreamsTab rendering', () => {
    it('renders the friendly display name, logo, thumbnail and live status from the API fields', async () => {
      mockFetchServers();
      const { container } = render(<PublicHarness />);

      // displayName is preferred over the federation username.
      expect(await screen.findByText('Goodnight TV')).toBeInTheDocument();
      expect(screen.queryByText('goodnight')).not.toBeInTheDocument();

      // logoUrl + thumbnailUrl must reach the rendered <img> elements.
      expect(container.querySelector(`img[src="${API_SERVER.logoUrl}"]`)).toBeInTheDocument();
      expect(container.querySelector(`img[src="${API_SERVER.thumbnailUrl}"]`)).toBeInTheDocument();

      // isOnline drives the LIVE badge and reveals the stream title.
      expect(screen.getByText('LIVE')).toBeInTheDocument();
      expect(screen.getByText('Late night coding')).toBeInTheDocument();
    });

    it('falls back to the username when there is no display name', async () => {
      mockFetchServers({ ...API_SERVER, displayName: undefined });
      render(<PublicHarness />);
      expect(await screen.findByText('goodnight')).toBeInTheDocument();
    });
  });

  describe('admin FederatedServersTable rendering', () => {
    it('renders the name, the iri as the URL, the external link and the status', async () => {
      mockFetchServers();
      const { container } = render(<AdminHarness />);

      expect(await screen.findByText('Goodnight TV')).toBeInTheDocument();

      // The iri populates both the URL column text and the external link href.
      expect(screen.getByText(API_SERVER.iri)).toBeInTheDocument();
      const link = container.querySelector('a[target="_blank"]');
      expect(link).toHaveAttribute('href', API_SERVER.iri);

      // isOnline -> "Online"; lastStatusUpdate -> the Last Checked column.
      expect(screen.getByText('Online')).toBeInTheDocument();
      expect(screen.getByText(API_SERVER.lastStatusUpdate)).toBeInTheDocument();
    });

    it('shows "Never" for last checked when the server has no status update', async () => {
      mockFetchServers({ ...API_SERVER, isOnline: false, lastStatusUpdate: undefined });
      render(<AdminHarness />);
      expect(await screen.findByText('Never')).toBeInTheDocument();
    });
  });
});
