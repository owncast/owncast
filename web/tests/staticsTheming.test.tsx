import React, { useEffect } from 'react';
import { act, render } from '@testing-library/react';
import { message } from 'antd';
import { RecoilRoot, useSetRecoilState } from 'recoil';
import { AntdProvider } from '../components/theme/AntdProvider';
import { clientConfigStateAtom } from '../components/stores/ClientConfigStore';
import { makeEmptyClientConfig } from '../interfaces/client-config.model';

jest.mock('next/router', () => ({
  useRouter: () => ({ query: {}, pathname: '/', asPath: '/', push: jest.fn(), replace: jest.fn() }),
}));

// The static message/notification/Modal.confirm APIs render into a detached
// React root outside AntdProvider. The provider keeps a module-level theme
// reference current, and antd re-invokes the configured holderRender on
// every static call, so statics opened after an admin appearance change
// must pick up the new tokens. This guards that contract end to end.

const ACTION_COLOR = '#123456';

const ApplyTheme: React.FC<{ apply: boolean }> = ({ apply }) => {
  const setClientConfig = useSetRecoilState(clientConfigStateAtom);
  useEffect(() => {
    if (apply) {
      setClientConfig({
        ...makeEmptyClientConfig(),
        appearanceVariables: { 'theme-color-action': ACTION_COLOR },
      });
    }
  }, [apply, setClientConfig]);
  return null;
};

// The css variable VALUES for a component land in a style rule keyed by the
// css-var-* scope class its provider assigns, so the message holder's rule
// is distinguishable from the main provider tree's.
const messageHolderCss = () => {
  const holder = document.querySelector('.ant-message');
  expect(holder).not.toBeNull();
  const varClass = Array.from(holder.classList).find(c => /^css-var-/.test(c));
  expect(varClass).toBeDefined();
  return Array.from(document.querySelectorAll('style'))
    .map(s => s.textContent || '')
    .filter(cssText => cssText.includes(`.${varClass}`))
    .join('\n');
};

describe('static API theming', () => {
  it('applies appearance changes to statics opened afterwards', async () => {
    const { rerender } = render(
      <RecoilRoot>
        <AntdProvider>
          <ApplyTheme apply={false} />
        </AntdProvider>
      </RecoilRoot>,
    );

    // Open a message with the default theme so the holder exists before the
    // appearance change: the regression this guards is a holder that keeps
    // its mount-time theme forever.
    await act(async () => {
      message.open({ content: 'before theme change' });
    });
    expect(messageHolderCss()).not.toContain(ACTION_COLOR);

    // Admin changes the appearance.
    await act(async () => {
      rerender(
        <RecoilRoot>
          <AntdProvider>
            <ApplyTheme apply />
          </AntdProvider>
        </RecoilRoot>,
      );
    });

    // A static opened after the change renders with the new primary color.
    await act(async () => {
      message.open({ content: 'after theme change' });
    });
    expect(messageHolderCss()).toContain(ACTION_COLOR);
  });
});
