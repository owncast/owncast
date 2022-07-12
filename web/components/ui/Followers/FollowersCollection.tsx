import { Col, Pagination, Row } from 'antd';
import { Follower } from '../../../interfaces/follower';
import SingleFollower from './Follower';
import s from './Followers.module.scss';

interface Props {
  total: number;
  followers: Follower[];
}

export default function FollowerCollection(props: Props) {
  const ITEMS_PER_PAGE = 24;

  const { followers, total } = props;
  const pages = Math.ceil(total / ITEMS_PER_PAGE);

  const noFollowers = (
    <div>A message explaining how to follow goes here since there are no followers.</div>
  );

  if (followers.length === 0) {
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

      <Pagination current={1} pageSize={ITEMS_PER_PAGE} total={pages || 1} hideOnSinglePage />
    </div>
  );
}
