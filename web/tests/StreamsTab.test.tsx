import React from 'react';
import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { StreamsTab, FederatedServer } from '../components/ui/StreamsTab/StreamsTab';

// Mock the Translation component
jest.mock('../components/ui/Translation/Translation', () => ({
  Translation: ({ defaultText }: { defaultText: string }) => <span>{defaultText}</span>,
}));

describe('StreamsTab', () => {
  const mockServers: FederatedServer[] = [
    {
      id: '1',
      name: 'Server 1',
      url: 'https://server1.example.com',
      isOnline: true,
      streamTitle: 'Live Stream',
      streamDescription: 'Description',
      tags: ['gaming'],
    },
    {
      id: '2',
      name: 'Server 2',
      url: 'https://server2.example.com',
      isOnline: false,
    },
  ];

  it('renders loading state', () => {
    render(<StreamsTab loading />);

    expect(screen.getByText('Loading featured streams...')).toBeInTheDocument();
  });

  it('renders error state', () => {
    const errorMessage = 'Failed to load servers';
    render(<StreamsTab error={errorMessage} />);

    expect(screen.getByText('Error loading streams')).toBeInTheDocument();
    expect(screen.getByText(errorMessage)).toBeInTheDocument();
  });

  it('renders empty state when no servers', () => {
    render(<StreamsTab servers={[]} />);

    expect(screen.getByText('No featured streams available')).toBeInTheDocument();
  });

  it('renders servers list correctly', () => {
    render(<StreamsTab servers={mockServers} />);

    expect(screen.getByText('Server 1')).toBeInTheDocument();
    expect(screen.getByText('Server 2')).toBeInTheDocument();
    expect(screen.getByText('Live Stream')).toBeInTheDocument();
  });

  it('sorts servers with online servers first', () => {
    const servers: FederatedServer[] = [
      {
        id: '1',
        name: 'Offline A',
        url: 'https://a.example.com',
        isOnline: false,
      },
      {
        id: '2',
        name: 'Online B',
        url: 'https://b.example.com',
        isOnline: true,
      },
      {
        id: '3',
        name: 'Offline C',
        url: 'https://c.example.com',
        isOnline: false,
      },
      {
        id: '4',
        name: 'Online A',
        url: 'https://d.example.com',
        isOnline: true,
      },
    ];

    const { container } = render(<StreamsTab servers={servers} />);

    // Get all server name elements in order
    const serverNames = container.querySelectorAll('.serverName');
    const names = Array.from(serverNames).map(el => el.textContent);

    // Online servers should come first, then sorted alphabetically
    expect(names[0]).toBe('Online A');
    expect(names[1]).toBe('Online B');
    expect(names[2]).toBe('Offline A');
    expect(names[3]).toBe('Offline C');
  });

  it('updates when servers prop changes', () => {
    const { rerender } = render(<StreamsTab servers={[]} />);

    expect(screen.getByText('No featured streams available')).toBeInTheDocument();

    rerender(<StreamsTab servers={mockServers} />);

    expect(screen.queryByText('No featured streams available')).not.toBeInTheDocument();
    expect(screen.getByText('Server 1')).toBeInTheDocument();
  });

  it('renders correct number of StreamCard components', () => {
    const { container } = render(<StreamsTab servers={mockServers} />);

    const cards = container.querySelectorAll('.ant-card');
    expect(cards).toHaveLength(mockServers.length);
  });
});
