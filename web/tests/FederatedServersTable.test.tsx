import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import {
  FederatedServersTable,
  FederatedServerData,
} from '../components/admin/FederatedServers/FederatedServersTable';

// Mock next/router for next-export-i18n
jest.mock('next/router', () => ({
  useRouter: () => ({
    query: {},
    pathname: '/',
    asPath: '/',
    push: jest.fn(),
    replace: jest.fn(),
  }),
}));

// Mock antd message
jest.mock('antd', () => ({
  ...jest.requireActual('antd'),
  message: {
    success: jest.fn(),
    error: jest.fn(),
  },
}));

describe('FederatedServersTable', () => {
  const mockServers: FederatedServerData[] = [
    {
      id: 1,
      iri: 'https://server1.example.com',
      name: 'Server 1',
      isOnline: true,
      lastStatusUpdate: '2024-01-01 12:00',
      addedAt: '2023-12-01',
    },
    {
      id: 2,
      iri: 'https://server2.example.com',
      name: 'Server 2',
      isOnline: false,
      addedAt: '2023-12-02',
    },
  ];

  const mockOnRemove = jest.fn();

  beforeEach(() => {
    mockOnRemove.mockClear();
  });

  it('renders servers correctly', () => {
    render(<FederatedServersTable servers={mockServers} onRemove={mockOnRemove} />);

    expect(screen.getByText('Server 1')).toBeInTheDocument();
    expect(screen.getByText('Server 2')).toBeInTheDocument();
    expect(screen.getByText('https://server1.example.com')).toBeInTheDocument();
    expect(screen.getByText('https://server2.example.com')).toBeInTheDocument();
  });

  it('displays online/offline status correctly', () => {
    render(<FederatedServersTable servers={mockServers} onRemove={mockOnRemove} />);

    const onlineTags = screen.getAllByText('Online');
    const offlineTags = screen.getAllByText('Offline');

    expect(onlineTags).toHaveLength(1);
    expect(offlineTags).toHaveLength(1);
  });

  it('displays last checked time or "Never"', () => {
    render(<FederatedServersTable servers={mockServers} onRemove={mockOnRemove} />);

    expect(screen.getByText('2024-01-01 12:00')).toBeInTheDocument();
    expect(screen.getByText('Never')).toBeInTheDocument();
  });

  it('shows confirmation dialog when removing server', async () => {
    render(<FederatedServersTable servers={mockServers} onRemove={mockOnRemove} />);

    const removeButtons = screen.getAllByText('Unfeature');
    fireEvent.click(removeButtons[0]);

    // Popconfirm opens with Yes/No buttons
    expect(screen.getByText('Yes')).toBeInTheDocument();
    expect(screen.getByText('No')).toBeInTheDocument();
  });

  it('calls onRemove when confirmed', async () => {
    mockOnRemove.mockResolvedValue(undefined);

    render(<FederatedServersTable servers={mockServers} onRemove={mockOnRemove} />);

    const removeButtons = screen.getAllByText('Unfeature');
    fireEvent.click(removeButtons[0]);

    const confirmButton = screen.getByText('Yes');
    fireEvent.click(confirmButton);

    await waitFor(() => {
      expect(mockOnRemove).toHaveBeenCalledWith(1);
    });
  });

  it('cancels removal when declined', async () => {
    render(<FederatedServersTable servers={mockServers} onRemove={mockOnRemove} />);

    const removeButtons = screen.getAllByText('Unfeature');
    fireEvent.click(removeButtons[0]);

    const cancelButton = screen.getByText('No');
    fireEvent.click(cancelButton);

    expect(mockOnRemove).not.toHaveBeenCalled();
  });

  it('shows loading state correctly', () => {
    const { container } = render(
      <FederatedServersTable servers={[]} loading onRemove={mockOnRemove} />,
    );

    // Ant Design Table renders a spinner when loading
    expect(container.querySelector('.ant-spin')).toBeInTheDocument();
  });

  it('handles removal error gracefully', async () => {
    // eslint-disable-next-line @typescript-eslint/no-var-requires, global-require
    const { message: antdMessage } = require('antd');
    mockOnRemove.mockRejectedValue(new Error('Removal failed'));

    render(<FederatedServersTable servers={mockServers} onRemove={mockOnRemove} />);

    const removeButtons = screen.getAllByText('Unfeature');
    fireEvent.click(removeButtons[0]);

    const confirmButton = screen.getByText('Yes');
    fireEvent.click(confirmButton);

    await waitFor(() => {
      expect(antdMessage.error).toHaveBeenCalled();
    });
  });

  it('renders external links for servers', () => {
    const { container } = render(
      <FederatedServersTable servers={mockServers} onRemove={mockOnRemove} />,
    );

    const links = container.querySelectorAll('a[target="_blank"]');
    expect(links).toHaveLength(2);
    expect(links[0]).toHaveAttribute('href', 'https://server1.example.com');
    expect(links[1]).toHaveAttribute('href', 'https://server2.example.com');
  });
});
