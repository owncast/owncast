import { Avatar, Comment } from 'antd';
import React from 'react';
import { Follower } from '../interfaces/follower';

interface Props {
  follower: Follower;
}

export default function SingleFollower(props: Props) {
  const { follower } = props;

  return (
    <Comment
      author={follower.username}
      avatar={<Avatar src={follower.image} alt="Han Solo" />}
      content={follower.name}
    />
  );
}
