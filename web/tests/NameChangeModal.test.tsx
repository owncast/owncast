import { render } from '@testing-library/react';
import '@testing-library/jest-dom';
import { RecoilRoot } from 'recoil';
import { NameChangeModal } from '../components/modals/NameChangeModal/NameChangeModal';
import { currentUserAtom } from '../components/stores/ClientConfigStore';

jest.mock('next/router', () => ({
  useRouter: () => ({ query: {}, pathname: '/', asPath: '/', push: jest.fn(), replace: jest.fn() }),
}));

const renderWithUser = (displayColor?: number) =>
  render(
    <RecoilRoot
      initializeState={({ set }) => {
        set(currentUserAtom, {
          id: 'user-1',
          displayName: 'tester',
          // The interface requires a number, but a user who has not yet
          // received their info over the socket has no color; the modal
          // must tolerate that (it used to crash on .toString()).
          displayColor: displayColor as number,
          isModerator: false,
        });
      }}
    >
      <NameChangeModal closeModal={() => {}} />
    </RecoilRoot>,
  );

describe('NameChangeModal', () => {
  it('shows the current user color in the closed color dropdown', () => {
    const { container } = renderWithUser(3);

    const swatch = container.querySelector<HTMLElement>(
      '.ant-select div[style*="theme-color-users-3"]',
    );
    expect(swatch).not.toBeNull();
    // antd v6's select content box sizes by line height, so a percentage
    // height collapses the swatch to zero pixels and the current color
    // appears missing. jsdom does no layout, so guard the style directly.
    expect(swatch.style.height).not.toBe('100%');
    expect(swatch.style.height).not.toBe('');
  });

  it('renders without crashing when the user has no color yet', () => {
    const { container } = renderWithUser(undefined);

    expect(container.querySelector('.ant-select')).not.toBeNull();
    expect(container.querySelector('[style*="theme-color-users"]')).toBeNull();
  });
});
