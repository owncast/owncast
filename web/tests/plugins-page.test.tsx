import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import Plugins from '../pages/admin/plugins';
import {
  fetchData,
  PLUGINS_LIST,
  PLUGIN_REGISTRY_INSTALL,
  PLUGIN_REGISTRY_LIST,
} from '../utils/apis';

jest.mock('../pages/admin/plugins.module.scss', () => ({
  errorAlert: 'errorAlert',
}));

jest.mock('next/router', () => ({
  useRouter: () => ({
    push: jest.fn(),
  }),
}));

jest.mock('next-export-i18n', () => ({
  useTranslation: () => ({
    t: (key: string) => key,
  }),
}));

jest.mock('../components/layouts/AdminLayout', () => ({
  AdminLayout: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

jest.mock('../components/admin/plugins/PluginsList', () => ({
  PluginsList: ({
    onUpdate,
  }: {
    onUpdate?: (plugin: any, version: string) => Promise<void> | void;
  }) => (
    <button
      type="button"
      onClick={() =>
        onUpdate?.(
          { slug: 'scripts-demo', name: 'Example Scripts Demo', version: '0.1.0' },
          '0.2.0',
        )
      }
    >
      trigger update
    </button>
  ),
}));

jest.mock('../components/admin/plugins/BrowseRegistry', () => ({
  BrowseRegistry: () => <div>browse registry</div>,
}));

jest.mock('../components/admin/plugins/InstallConfirmModal', () => ({
  InstallConfirmModal: () => null,
}));

jest.mock('../utils/apis', () => {
  const actual = jest.requireActual('../utils/apis');
  return {
    ...actual,
    fetchData: jest.fn(),
  };
});

const mockedFetchData = fetchData as jest.MockedFunction<typeof fetchData>;

describe('Plugins admin page', () => {
  beforeEach(() => {
    mockedFetchData.mockReset();
  });

  test('surfaces registry install errors to the admin page', async () => {
    mockedFetchData.mockImplementation(async url => {
      if (url === PLUGINS_LIST) {
        return [];
      }
      if (url === PLUGIN_REGISTRY_LIST) {
        return [];
      }
      if (url === PLUGIN_REGISTRY_INSTALL) {
        throw new Error('invalid manifest: manifest.scripts requires the "http.serve" permission');
      }
      throw new Error(`unexpected fetchData call: ${url}`);
    });

    render(<Plugins />);

    fireEvent.click(await screen.findByText('trigger update'));

    await waitFor(() => {
      expect(
        screen.getByText('invalid manifest: manifest.scripts requires the "http.serve" permission'),
      ).toBeInTheDocument();
    });
  });
});
