import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { message } from 'antd';
import Plugins from '../pages/admin/plugins';
import {
  fetchData,
  PLUGINS_LIST,
  PLUGIN_REGISTRY_INSTALL,
  PLUGIN_REGISTRY_LIST,
  PLUGIN_UPLOAD,
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

jest.mock('antd', () => {
  const actual = jest.requireActual('antd');
  return {
    ...actual,
    message: {
      ...actual.message,
      error: jest.fn(),
      success: jest.fn(),
    },
    Upload: ({
      children,
      customRequest,
    }: {
      children: React.ReactNode;
      customRequest: (options: { file: File; onError: (error: Error) => void }) => void;
    }) => (
      <>
        {children}
        <input
          data-testid="plugin-upload"
          type="file"
          onChange={() =>
            customRequest({ file: new File(['plugin'], 'duplicate.ocpkg'), onError: jest.fn() })
          }
        />
      </>
    ),
  };
});

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

const originalFetch = global.fetch;

describe('Plugins admin page', () => {
  beforeEach(() => {
    mockedFetchData.mockReset();
    jest.clearAllMocks();
  });
  afterEach(() => {
    global.fetch = originalFetch;
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

    expect(message.error).toHaveBeenCalledWith(
      'invalid manifest: manifest.scripts requires the "http.serve" permission',
    );
  });

  test('shows a duplicate slug upload rejection in the admin alert', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 400,
      text: async () =>
        JSON.stringify({ error: 'plugin upload rejected: slug "hello-world" is already used' }),
    } as Response);

    mockedFetchData.mockImplementation(async url => {
      if (url === PLUGINS_LIST || url === PLUGIN_REGISTRY_LIST) {
        return [];
      }
      throw new Error(`unexpected fetchData call: ${url}`);
    });

    render(<Plugins />);
    fireEvent.change(await screen.findByTestId('plugin-upload'), {
      target: { files: [new File(['plugin'], 'duplicate.ocpkg')] },
    });

    await waitFor(() => {
      expect(
        screen.getByText('plugin upload rejected: slug "hello-world" is already used'),
      ).toBeInTheDocument();
    });
    expect(global.fetch).toHaveBeenCalledWith(
      PLUGIN_UPLOAD,
      expect.objectContaining({ method: 'POST' }),
    );
    expect(message.error).toHaveBeenCalledWith(
      'plugin upload rejected: slug "hello-world" is already used',
    );
  });
});
