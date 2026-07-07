import React, { ReactElement, useCallback, useEffect, useState } from 'react';
import { Alert, Spin } from 'antd';
import { useRouter } from 'next/router';
import { useTranslation } from 'next-export-i18n';
import { AdminLayout } from '../../../components/layouts/AdminLayout';
import { fetchData, PLUGINS_LIST } from '../../../utils/apis';
import { Plugin } from '../../../interfaces/plugin';
import { PluginDetail } from '../../../components/admin/plugins/PluginDetail';
import { Localization } from '../../../types/localization';
import s from './configure.module.scss';

// Per-plugin admin view at /admin/plugins/configure/?id=<plugin-id>.
// The `id` query parameter is the plugin's slug (canonical identifier).
// The URL param is named `id` to keep author/operator-facing surfaces
// free of the slug jargon; internal code and the manifest still call
// it slug.
//
// Why a query-string route and not a dynamic /admin/plugins/[id]/ route:
// the admin UI is statically exported (next.config.js sets output: 'export'
// for non-dev), and plugin identifiers are discovered at runtime by the
// Owncast server. A dynamic Next.js route would either fail the static
// export or require build-time enumeration we can't provide. A static
// page with a runtime query param sidesteps both: one HTML file is
// generated, all plugin-selection logic runs client-side off router.query.
//
// The plugin list fetch itself primes the admin session cookie as a side
// effect of the RequireAdminAuth middleware on the host, so by the time
// we render the iframe the browser has a cookie that authenticates the
// iframe's top-level GET (which can't carry a custom header).
const PluginConfigurePage = () => {
  const { t } = useTranslation();
  const router = useRouter();
  const idParam = router.query.id;
  const slug = typeof idParam === 'string' ? idParam : undefined;

  const [plugin, setPlugin] = useState<Plugin | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [notFound, setNotFound] = useState(false);

  const load = useCallback(async () => {
    if (!slug) {
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      const result = await fetchData(PLUGINS_LIST);
      const list: Plugin[] = Array.isArray(result) ? result : [];
      const found = list.find(p => p.slug === slug) ?? null;
      setPlugin(found);
      setNotFound(found === null);
      setError(null);
    } catch (e) {
      setError(e instanceof Error ? e.message : String(e));
    } finally {
      setLoading(false);
    }
  }, [slug]);

  useEffect(() => {
    if (router.isReady) load();
  }, [router.isReady, load]);

  if (!router.isReady || (loading && !plugin)) {
    return (
      <div className={s.loading}>
        <Spin />
      </div>
    );
  }

  if (error) {
    return (
      <Alert
        type="error"
        showIcon
        message={t(Localization.Admin.Plugins.errorTitle)}
        description={error}
      />
    );
  }

  if (notFound || !plugin) {
    return (
      <Alert
        type="warning"
        showIcon
        message={t(Localization.Admin.Plugins.notFoundTitle)}
        description={t(Localization.Admin.Plugins.notFoundDescription, { name: slug ?? '' })}
      />
    );
  }

  return <PluginDetail plugin={plugin} />;
};

PluginConfigurePage.getLayout = function getLayout(page: ReactElement) {
  return <AdminLayout page={page} />;
};

export default PluginConfigurePage;
