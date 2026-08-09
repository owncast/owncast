import { useEffect, useState } from 'react';
import { useTranslation } from 'next-export-i18n';
import { useAtomValue } from 'jotai';
import Link from 'next/link';
import Head from 'next/head';
import { Alert, Skeleton, Typography } from 'antd';
import { format } from 'date-fns';
import {
  clientConfigStateAtom,
  ClientConfigStore,
} from '../../components/stores/ClientConfigStore';
import { Theme } from '../../components/theme/Theme';
import { ClipPlayer } from '../../components/video/ClipPlayer/ClipPlayer';
import { ClipsTab } from '../../components/ui/ClipsTab/ClipsTab';
import type { Clip } from '../../interfaces/clip.model';
import { Localization } from '../../types/localization';
import styles from './ClipPage.module.scss';

const { Title, Text } = Typography;

// ClipPage is the standalone page a shared clip link opens: the clip itself
// and the few facts about it. The full viewer page, with its chat and tabs,
// is deliberately not part of this view.
export default function ClipPage() {
  const { t } = useTranslation();
  const clientConfig = useAtomValue(clientConfigStateAtom);
  const { name } = clientConfig;

  // The web app is a static export served from one document per route, so the
  // clip id comes from the path rather than a router param.
  const [clipId, setClipId] = useState<string | null>(null);
  const [clip, setClip] = useState<Clip | null>(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const match = window.location.pathname.match(/^\/clips\/([^/]+)\/?$/);
    setClipId(match ? match[1] : '');
  }, []);

  useEffect(() => {
    if (clipId === null || clipId === '') {
      setLoading(clipId === null);
      return;
    }

    const load = async () => {
      try {
        const response = await fetch(`/api/clips/${clipId}`);
        if (!response.ok) {
          throw new Error('clip not found');
        }
        setClip(await response.json());
      } catch {
        setError('This clip is no longer available.');
      } finally {
        setLoading(false);
      }
    };

    load();
  }, [clipId]);

  const createdAt = clip?.timestamp ? new Date(clip.timestamp) : null;
  const createdLabel =
    createdAt && !Number.isNaN(createdAt.getTime()) ? format(createdAt, 'PPP p') : '';

  const content = () => {
    if (loading) {
      return <Skeleton active paragraph={{ rows: 6 }} />;
    }

    // No clip in the path: this is the clips index.
    if (clipId === '') {
      return (
        <>
          <Title level={2}>Clips</Title>
          <ClipsTab />
        </>
      );
    }

    if (error || !clip) {
      return <Alert type="error" message={error || 'This clip is no longer available.'} showIcon />;
    }

    return (
      <>
        <ClipPlayer source={clip.manifest} poster={clip.thumbnail} autoplay />

        <div className={styles.details}>
          <Title level={3} className={styles.title}>
            {clip.title || 'Clip'}
          </Title>

          {clip.streamTitle && (
            <Text className={styles.streamTitle} type="secondary">
              {clip.streamTitle}
            </Text>
          )}

          {clip.clippedBy && (
            <Text className={styles.clippedBy} type="secondary">
              {t(Localization.Frontend.Clips.clippedBy).replace('{{name}}', clip.clippedBy)}
            </Text>
          )}

          {createdLabel && (
            <Text className={styles.created} type="secondary">
              {createdLabel}
            </Text>
          )}
        </div>
      </>
    );
  };

  const pageTitle = clip?.title || 'Clips';

  return (
    <>
      <Head>
        <title>{name ? `${pageTitle} - ${name}` : pageTitle}</title>
      </Head>
      <ClientConfigStore />
      <Theme />
      <main className={styles.page}>
        <Link className={styles.serverLink} href="/">
          {name}
        </Link>
        {content()}
      </main>
    </>
  );
}
