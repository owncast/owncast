import { Pagination } from 'antd';
import { Follower } from '../interfaces/follower';
import SingleFollower from './Follower';

interface Props {
  total: number;
  followers: Follower[];
}

export default function FollowerCollection(props: Props) {
  const ITEMS_PER_PAGE = 24;

  const { followers, total } = props;
  const pages = Math.ceil(total / ITEMS_PER_PAGE);

  return (
    <div>
      {followers.map(follower => (
        <SingleFollower key={follower.link} follower={follower} />
      ))}
      <Pagination current={1} pageSize={ITEMS_PER_PAGE} total={pages || 1} />
    </div>
  );
}
