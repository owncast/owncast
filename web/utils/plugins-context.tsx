import React, {
  createContext,
  FC,
  ReactElement,
  useCallback,
  useEffect,
  useMemo,
  useState,
} from 'react';
import { fetchData, PLUGINS_LIST } from './apis';
import { Plugin } from '../interfaces/plugin';

// Shared installed-plugins state for the admin. The sidebar's Plugins
// submenu and the Plugins page both read from this single source, so an
// install/uninstall/enable from the page can call reload() and have the
// sidebar pick up the change without a full-page refresh.
export type PluginsContextValue = {
  plugins: Plugin[];
  // True until the first fetch settles. The Plugins page drives its
  // table's loading state from this.
  loading: boolean;
  // Set when the most recent fetch failed; null when it succeeded.
  error: string | null;
  // Re-fetches the installed plugin list. Awaitable so callers can
  // sequence UI (e.g. open a confirmation modal) after the refresh.
  reload: () => Promise<void>;
};

export const PluginsContext = createContext<PluginsContextValue>({
  plugins: [],
  loading: true,
  error: null,
  reload: async () => {},
});

export type PluginsProviderProps = {
  children: ReactElement;
};

const PluginsProvider: FC<PluginsProviderProps> = ({ children }) => {
  const [plugins, setPlugins] = useState<Plugin[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    try {
      const result = await fetchData(PLUGINS_LIST);
      setPlugins(Array.isArray(result) ? result : []);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    reload();
  }, [reload]);

  const value = useMemo(
    () => ({ plugins, loading, error, reload }),
    [plugins, loading, error, reload],
  );
  return <PluginsContext.Provider value={value}>{children}</PluginsContext.Provider>;
};

export default PluginsProvider;
