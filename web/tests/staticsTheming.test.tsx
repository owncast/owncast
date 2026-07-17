import React, { useEffect } from 'react';
import { act, render } from '@testing-library/react';
import { message } from 'antd';
import { Provider, useSetAtom } from 'jotai';
import { AntdProvider } from '../components/theme/AntdProvider';
import { clientConfigStateAtom } from '../components/stores/ClientConfigStore';
import { makeEmptyClientConfig } from '../interfaces/client-config.model';
import antdDefaultTheme from '../components/theme/antd-default-theme.json';

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
  const setClientConfig = useSetAtom(clientConfigStateAtom);
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
// scope class its provider assigns: the stable cssVar key from
// antd-default-theme.json for the default theme, and a distinct "-custom"
// key once appearance customizations exist (so the pre-extracted default
// values in antd.css cannot win the cascade over the customized ones).
const messageHolderCss = (expectedKey: string) => {
  const holder = document.querySelector('.ant-message');
  expect(holder).not.toBeNull();
  expect(holder.classList).toContain(expectedKey);
  return Array.from(document.querySelectorAll('style'))
    .map(s => s.textContent || '')
    .filter(cssText => cssText.includes(`.${expectedKey}`))
    .join('\n');
};

describe('static API theming', () => {
  it('applies appearance changes to statics opened afterwards', async () => {
    const { rerender } = render(
      <Provider>
        <AntdProvider>
          <ApplyTheme apply={false} />
        </AntdProvider>
      </Provider>,
    );

    // Open a message with the default theme so the holder exists before the
    // appearance change: the regression this guards is a holder that keeps
    // its mount-time theme forever.
    await act(async () => {
      message.open({ content: 'before theme change' });
    });
    expect(messageHolderCss(antdDefaultTheme.cssVar.key)).not.toContain(ACTION_COLOR);

    // Admin changes the appearance.
    await act(async () => {
      rerender(
        <Provider>
          <AntdProvider>
            <ApplyTheme apply />
          </AntdProvider>
        </Provider>,
      );
    });

    // A static opened after the change renders with the new primary color,
    // under the customized scope key.
    await act(async () => {
      message.open({ content: 'after theme change' });
    });
    expect(messageHolderCss(`${antdDefaultTheme.cssVar.key}-custom`)).toContain(ACTION_COLOR);
  });
});
