import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import { FeatureRequests } from '../components/admin/FederatedServers/FeatureRequests';
import { FeatureRequest } from '../hooks/useFeatureRequests';

jest.mock('next/router', () => ({
  useRouter: () => ({ query: {}, pathname: '/', asPath: '/', push: jest.fn(), replace: jest.fn() }),
}));

jest.mock('antd', () => ({
  ...jest.requireActual('antd'),
  message: { success: jest.fn(), error: jest.fn() },
}));

describe('FeatureRequests', () => {
  const requests: FeatureRequest[] = [
    {
      link: 'https://peer.example.com/federation/user/streamer',
      name: 'Peer Server',
      username: 'streamer@peer.example.com',
      image: 'https://peer.example.com/logo.png',
    },
  ];

  const onApprove = jest.fn();
  const onReject = jest.fn();

  beforeEach(() => {
    onApprove.mockReset();
    onReject.mockReset();
  });

  it('renders nothing when there are no requests', () => {
    const { container } = render(
      <FeatureRequests requests={[]} onApprove={onApprove} onReject={onReject} />,
    );
    expect(container).toBeEmptyDOMElement();
  });

  it('lists pending requests with approve and reject actions', () => {
    render(<FeatureRequests requests={requests} onApprove={onApprove} onReject={onReject} />);

    expect(screen.getByText('Peer Server')).toBeInTheDocument();
    expect(screen.getByText('Approve')).toBeInTheDocument();
    expect(screen.getByText('Reject')).toBeInTheDocument();
  });

  it('calls onReject with the actor IRI immediately', async () => {
    onReject.mockResolvedValue(undefined);
    render(<FeatureRequests requests={requests} onApprove={onApprove} onReject={onReject} />);

    fireEvent.click(screen.getByText('Reject'));

    await waitFor(() => {
      expect(onReject).toHaveBeenCalledWith(requests[0].link);
    });
    expect(onApprove).not.toHaveBeenCalled();
  });

  it('does not render a non-http link as a clickable anchor', () => {
    // A deliberately hostile non-http(s) link is the whole point of this test.
    // eslint-disable-next-line no-script-url
    const maliciousLink = 'javascript:alert(1)';
    const malicious: FeatureRequest[] = [
      { link: maliciousLink, name: 'Evil Server', username: 'evil' },
    ];
    const { container } = render(
      <FeatureRequests requests={malicious} onApprove={onApprove} onReject={onReject} />,
    );

    // No anchor should ever carry a javascript: href.
    expect(container.querySelector('a[href^="javascript:"]')).toBeNull();
    // The name is still shown, just as plain text.
    expect(screen.getByText('Evil Server')).toBeInTheDocument();
  });

  it('calls onApprove with the actor IRI after confirmation', async () => {
    onApprove.mockResolvedValue(undefined);
    render(<FeatureRequests requests={requests} onApprove={onApprove} onReject={onReject} />);

    fireEvent.click(screen.getByText('Approve'));
    // Popconfirm shows a confirm button ("Yes").
    const confirmButtons = await screen.findAllByText('Yes');
    fireEvent.click(confirmButtons[confirmButtons.length - 1]);

    await waitFor(() => {
      expect(onApprove).toHaveBeenCalledWith(requests[0].link);
    });
  });
});
