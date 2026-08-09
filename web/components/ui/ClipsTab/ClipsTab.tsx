import { FC, useCallback, useEffect, useState } from 'react';
import { Alert, Card, Empty, Modal, Spin, Typography } from 'antd';
import dynamic from 'next/dynamic';
import { useTranslation } from 'next-export-i18n';
import type { Clip } from '../../../interfaces/clip.model';
import { Localization } from '../../../types/localization';
import styles from './ClipsTab.module.scss';

const ClipPlayer = dynamic(
  () => import('../../video/ClipPlayer/ClipPlayer').then(mod => mod.ClipPlayer),
  { ssr: false },
);

// formatDuration renders a clip length as m:ss.
const formatDuration = (seconds: number): string => {
  const total = Math.max(Math.round(seconds), 0);
  const minutes = Math.floor(total / 60);
  const remainder = total % 60;
  return `${minutes}:${String(remainder).padStart(2, '0')}`;
};

// ClipsTab lists the clips created from this server's streams and plays the
// selected one.
export const ClipsTab: FC = () => {
  const { t } = useTranslation();
  const [clips, setClips] = useState<Clip[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedClip, setSelectedClip] = useState<Clip | null>(null);
  const clipLoadError = t(Localization.Frontend.Clips.empty);

  const loadClips = useCallback(async () => {
    try {
      const response = await fetch('/api/clips');
      if (!response.ok) {
        throw new Error(`unexpected response ${response.status}`);
      }
      const result: Clip[] = await response.json();
      setClips(result || []);
    } catch {
      setError(clipLoadError);
    } finally {
      setLoading(false);
    }
  }, [clipLoadError]);

  useEffect(() => {
    loadClips();
  }, [loadClips]);

  if (loading) {
    return (
      <div className={styles.centered}>
        <Spin size="large" />
      </div>
    );
  }

  if (error) {
    return (
      <div className={styles.centered}>
        <Alert type="error" message={error} showIcon />
      </div>
    );
  }

  if (clips.length === 0) {
    return (
      <div className={styles.centered}>
        <Empty description={t(Localization.Frontend.Clips.empty)} />
      </div>
    );
  }

  return (
    <div className={styles.clips} id="clips-tab">
      {clips.map(clip => (
        // A real anchor so the shareable clip page is one right-click away
        // (open in new tab, copy link). A plain left click stays in the page
        // and plays the clip in the modal instead.
        <a
          key={clip.id}
          className={styles.cardLink}
          href={`/clips/${clip.id}`}
          onClick={e => {
            if (e.metaKey || e.ctrlKey || e.shiftKey) {
              return;
            }
            e.preventDefault();
            setSelectedClip(clip);
          }}
        >
          <Card
            role="article"
            className={styles.clipCard}
            styles={{ body: { padding: 0 } }}
            cover={
              <div className={styles.coverContainer}>
                {clip.thumbnail ? (
                  // eslint-disable-next-line @next/next/no-img-element
                  <img
                    src={clip.thumbnail}
                    alt={clip.title || 'Clip'}
                    className={styles.thumbnail}
                  />
                ) : (
                  <div className={styles.thumbnailPlaceholder} />
                )}
                <span className={styles.durationBadge}>{formatDuration(clip.durationSeconds)}</span>
              </div>
            }
          >
            <div className={styles.cardContent}>
              <Typography.Text strong className={styles.clipTitle} ellipsis>
                {clip.title || clip.streamTitle || t(Localization.Frontend.Clips.title)}
              </Typography.Text>
              {clip.clippedBy && (
                <Typography.Text type="secondary" ellipsis className={styles.clippedBy}>
                  {t(Localization.Frontend.Clips.clippedBy).replace('{{name}}', clip.clippedBy)}
                </Typography.Text>
              )}
            </div>
          </Card>
        </a>
      ))}

      <Modal
        open={!!selectedClip}
        title={selectedClip?.title || selectedClip?.streamTitle}
        onCancel={() => setSelectedClip(null)}
        footer={null}
        width="720px"
        zIndex={999}
        centered
        destroyOnHidden
      >
        {selectedClip && (
          <ClipPlayer source={selectedClip.manifest} poster={selectedClip.thumbnail} autoplay />
        )}
      </Modal>
    </div>
  );
};
