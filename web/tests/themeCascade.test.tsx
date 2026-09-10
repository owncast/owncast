import { act, fireEvent, render } from '@testing-library/react';
import { createStore, Provider } from 'jotai';
import { ReactElement } from 'react';
import { clientConfigStateAtom } from '../components/stores/ClientConfigStore';
import { Theme } from '../components/theme/Theme';
import { PluginTabFrame } from '../components/ui/PluginTabFrame/PluginTabFrame';
import { ClientConfig, makeEmptyClientConfig } from '../interfaces/client-config.model';

const CASCADE_VARIABLE = '--cascade-probe';
const pluginStyles = `:root { ${CASCADE_VARIABLE}: plugin; }`;
const appearanceVariables = { 'cascade-probe': 'appearance' };
const customStyles = `:root { ${CASCADE_VARIABLE}: custom; }`;

const renderWithConfig = (component: ReactElement, config: Partial<ClientConfig>) => {
  const store = createStore();
  store.set(clientConfigStateAtom, { ...makeEmptyClientConfig(), ...config });
  return { store, ...render(<Provider store={store}>{component}</Provider>) };
};

const cascadeValue = (doc: Document) =>
  doc.defaultView?.getComputedStyle(doc.documentElement).getPropertyValue(CASCADE_VARIABLE).trim();

describe('theme cascade', () => {
  it('orders plugin styles, appearance variables, and custom CSS on the viewer', () => {
    const config = {
      ...makeEmptyClientConfig(),
      pluginStyles,
      appearanceVariables,
      customStyles,
    };
    const { store } = renderWithConfig(<Theme />, config);

    expect(cascadeValue(document)).toBe('custom');

    act(() => store.set(clientConfigStateAtom, { ...config, customStyles: '' }));
    expect(cascadeValue(document)).toBe('appearance');

    act(() =>
      store.set(clientConfigStateAtom, {
        ...config,
        customStyles: '',
        appearanceVariables: {},
      }),
    );
    expect(cascadeValue(document)).toBe('plugin');
  });

  it.each([
    ['plugin styles over author styles', { pluginStyles }, 'plugin'],
    [
      'appearance variables over plugin styles',
      { pluginStyles, appearanceVariables },
      'appearance',
    ],
    [
      'custom CSS over appearance variables',
      { pluginStyles, appearanceVariables, customStyles },
      'custom',
    ],
  ])('applies %s inside plugin tabs', (_name, config, expected) => {
    const { getByTitle } = renderWithConfig(
      <PluginTabFrame content={`<style>:root { ${CASCADE_VARIABLE}: author; }</style>`} />,
      config,
    );
    const frame = getByTitle('Plugin tab') as HTMLIFrameElement;
    const doc = frame.contentDocument;
    const win = frame.contentWindow;
    expect(doc).not.toBeNull();
    expect(win).not.toBeNull();

    doc!.head.innerHTML = `<style>:root { ${CASCADE_VARIABLE}: author; }</style>`;
    Object.defineProperty(win, 'ResizeObserver', {
      configurable: true,
      value: class ResizeObserver {
        observe() {}
        disconnect() {}
      },
    });
    fireEvent.load(frame);

    expect(cascadeValue(doc!)).toBe(expected);
  });
});
