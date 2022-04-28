import { Follower } from '../interfaces/follower';

interface Props {
  followers: Follower[];
}

export default function FollowerCollection(props: Props) {
  return <div>List of followers go here</div>;
}
