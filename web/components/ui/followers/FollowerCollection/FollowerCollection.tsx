import { FC, useEffect, useState } from 'react';
import { Col, Pagination, Row } from 'antd';
import { Follower } from '../../../../interfaces/follower';
import { SingleFollower } from '../SingleFollower/SingleFollower';
import styles from './FollowerCollection.module.scss';
import { FollowButton } from '../../../action-buttons/FollowButton';

export type FollowerCollectionProps = {
  name: string;
};

export const FollowerCollection: FC<FollowerCollectionProps> = ({ name }) => {
  const ENDPOINT = '/api/followers';
  const ITEMS_PER_PAGE = 24;

  const [followers, setFollowers] = useState<Follower[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const pages = Math.ceil(total / ITEMS_PER_PAGE);

  const getFollowers = async () => {
    try {
      const response = await fetch(`${ENDPOINT}?page=${page}`);
      const data = await response.json();
      const { results, total: totalResults } = data;

      setFollowers(results);
      setTotal(totalResults);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    getFollowers();
  }, []);

  useEffect(() => {
    getFollowers();
  }, [page]);

  const noFollowers = (
    <div className={styles.noFollowers}>
      <h2>Be the first follower!</h2>
      <p>
        {name !== 'Owncast' ? name : 'This server'} is a part of the{' '}
        <a href="https://owncast.online/join-fediverse">Fediverse</a>, an interconnected network of
        independent users and servers.
      </p>
      <p>
        By following {name !== 'Owncast' ? name : 'this server'} you&apos;ll be able to get updates
        from the stream, share it with others, and and show your appreciation when it goes live, all
        from your own Fediverse account.
      </p>
      <FollowButton />
    </div>
  );

  if (!followers?.length) {
    return noFollowers;
  }

  return (
    <div className={styles.followers}>
      <Row wrap gutter={[10, 10]}>
        {followers.map(follower => (
          <Col>
            <SingleFollower key={follower.link} follower={follower} />
          </Col>
        ))}
      </Row>

      <Pagination
        current={page}
        pageSize={ITEMS_PER_PAGE}
        total={pages || 1}
        onChange={p => {
          setPage(p);
        }}
        hideOnSinglePage
      />
    </div>
  );
};
