import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom';
import { StreamCard, StreamCardProps } from '../components/ui/StreamCard/StreamCard';

// Mock window.open
const mockOpen = jest.fn();
window.open = mockOpen;

describe('StreamCard', () => {
  const defaultProps: StreamCardProps = {
    serverName: 'Test Server',
    serverUrl: 'https://test.example.com',
    isOnline: false,
  };

  beforeEach(() => {
    mockOpen.mockClear();
  });

  it('renders offline server correctly', () => {
    render(<StreamCard {...defaultProps} />);

    expect(screen.getByText('Test Server')).toBeInTheDocument();
    expect(screen.getByText('OFFLINE')).toBeInTheDocument();
  });

  it('shows the real server hostname so a spoofed name cannot fully masquerade', () => {
    render(<StreamCard {...defaultProps} serverName="Totally Legit News" />);

    // The display name is shown, but so is the actual hostname of the link.
    expect(screen.getByText('Totally Legit News')).toBeInTheDocument();
    expect(screen.getByText('test.example.com')).toBeInTheDocument();
  });

  it('renders online server with stream info', () => {
    const props: StreamCardProps = {
      ...defaultProps,
      isOnline: true,
      streamTitle: 'Live Stream Title',
      streamDescription: 'This is a test stream description',
      tags: ['test', 'stream'],
      thumbnail: 'https://test.example.com/thumb.jpg',
    };

    render(<StreamCard {...props} />);

    expect(screen.getByText('Test Server')).toBeInTheDocument();
    expect(screen.getByText('LIVE')).toBeInTheDocument();
    expect(screen.getByText('Live Stream Title')).toBeInTheDocument();
    expect(screen.getByText('This is a test stream description')).toBeInTheDocument();
    expect(screen.getByText('test')).toBeInTheDocument();
    expect(screen.getByText('stream')).toBeInTheDocument();
  });

  it('displays server logo when provided', () => {
    const props: StreamCardProps = {
      ...defaultProps,
      serverLogo: 'https://test.example.com/logo.png',
    };

    render(<StreamCard {...props} />);

    const logos = screen.getAllByAltText('Test Server');
    expect(logos.length).toBeGreaterThan(0);
    expect(logos[0]).toHaveAttribute('src', 'https://test.example.com/logo.png');
  });

  it('opens server URL in new tab when clicked', () => {
    render(<StreamCard {...defaultProps} />);

    const card = screen.getByRole('article');
    fireEvent.click(card);

    expect(mockOpen).toHaveBeenCalledWith(
      'https://test.example.com',
      '_blank',
      'noopener,noreferrer',
    );
  });

  it('calls custom onClick handler when provided', () => {
    const mockOnClick = jest.fn();
    const props: StreamCardProps = {
      ...defaultProps,
      onClick: mockOnClick,
    };

    render(<StreamCard {...props} />);

    const card = screen.getByRole('article');
    fireEvent.click(card);

    expect(mockOnClick).toHaveBeenCalled();
    expect(mockOpen).not.toHaveBeenCalled();
  });

  it('truncates long tags list to 3 items', () => {
    const props: StreamCardProps = {
      ...defaultProps,
      tags: ['tag1', 'tag2', 'tag3', 'tag4', 'tag5'],
    };

    render(<StreamCard {...props} />);

    expect(screen.getByText('tag1')).toBeInTheDocument();
    expect(screen.getByText('tag2')).toBeInTheDocument();
    expect(screen.getByText('tag3')).toBeInTheDocument();
    expect(screen.queryByText('tag4')).not.toBeInTheDocument();
    expect(screen.queryByText('tag5')).not.toBeInTheDocument();
  });

  it('applies correct CSS classes for online/offline states', () => {
    const { container: offlineContainer } = render(<StreamCard {...defaultProps} />);
    expect(offlineContainer.querySelector('.offline')).toBeInTheDocument();

    const { container: onlineContainer } = render(<StreamCard {...defaultProps} isOnline />);
    expect(onlineContainer.querySelector('.online')).toBeInTheDocument();
  });
});
