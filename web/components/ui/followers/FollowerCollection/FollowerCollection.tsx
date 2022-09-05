import { useEffect, useState } from 'react';
import { Col, Pagination, Row } from 'antd';
import { Follower } from '../../../../interfaces/follower';
import SingleFollower from '../SingleFollower/SingleFollower';
import s from '../SingleFollower/SingleFollower.module.scss';

export const FollowerCollection = () => {
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

      setFollowers(data.response);
      setTotal(data.total);
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
    <div>A message explaining how to follow goes here since there are no followers.</div>
  );

  if (!followers?.length) {
    return noFollowers;
  }

  return (
    <div className={s.followers}>
      <Row wrap gutter={[10, 10]} justify="space-around">
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
export default FollowerCollection;
