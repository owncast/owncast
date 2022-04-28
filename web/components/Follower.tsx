import { Follower } from '../interfaces/follower';

interface Props {
  follower: Follower;
}

export default function FollowerCollection(props: Props) {
  return <div>This is a single follower</div>;
}
