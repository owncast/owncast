import { FC, useEffect, useState } from 'react';
import { Pagination, Spin } from 'antd';
import { ErrorBoundary } from 'react-error-boundary';
import { Follower } from '../../../../interfaces/follower';
import { SingleFollower } from '../SingleFollower/SingleFollower';
import styles from './FollowerCollection.module.scss';
import { FollowButton } from '../../../action-buttons/FollowButton';
import { ComponentError } from '../../ComponentError/ComponentError';

export type FollowerCollectionProps = {
  name: string;
  onFollowButtonClick: () => void;
};

export const FollowerCollection: FC<FollowerCollectionProps> = ({ name, onFollowButtonClick }) => {
  const ENDPOINT = '/api/followers';
  const ITEMS_PER_PAGE = 24;

  const [followers, setFollowers] = useState<Follower[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

  const getFollowers = async () => {
    try {
      setLoading(true);

      const response = await fetch(
        `${ENDPOINT}?offset=${(page - 1) * ITEMS_PER_PAGE}&limit=${ITEMS_PER_PAGE}`,
      );

      const data = await response.json();
      const { results, total: totalResults } = data;

      setFollowers(results);
      setTotal(totalResults);

      setLoading(false);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    getFollowers();
  }, [page]);

  const noFollowers = (
    <div className={styles.noFollowers} id="followers-collection">
      <h2>Be the first follower!</h2>
      <p>
        {name !== 'Owncast' ? name : 'This server'} is a part of the{' '}
        <a href="https://owncast.online/join-fediverse">Fediverse</a>, an interconnected network of
        independent users and servers.
      </p>
      <p>
        By following {name !== 'Owncast' ? name : 'this server'} you&apos;ll be able to get updates
        from the stream, share it with others, and show your appreciation when it goes live, all
        from your own Fediverse account.
      </p>
      <FollowButton onClick={onFollowButtonClick} />
    </div>
  );

  if (!followers?.length && !loading) {
    return noFollowers;
  }

  return (
    <ErrorBoundary
      // eslint-disable-next-line react/no-unstable-nested-components
      fallbackRender={({ error, resetErrorBoundary }) => (
        <ComponentError
          componentName="FollowerCollection"
          message={error.message}
          retryFunction={resetErrorBoundary}
        />
      )}
    >
      <Spin spinning={loading} size="large">
        <div className={styles.followers} id="followers-collection">
          <div className={styles.followerRow}>
            {followers.map(follower => (
              <SingleFollower key={follower.link} follower={follower} />
            ))}
          </div>

          <Pagination
            className={styles.pagination}
            current={page}
            pageSize={ITEMS_PER_PAGE}
            defaultPageSize={ITEMS_PER_PAGE}
            total={total}
            showSizeChanger={false}
            onChange={p => {
              setPage(p);
            }}
            hideOnSinglePage
          />
        </div>
      </Spin>
    </ErrorBoundary>
  );
};
