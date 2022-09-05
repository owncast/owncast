import { Avatar, Col, Row } from 'antd';
import React from 'react';
import { Follower } from '../../../../interfaces/follower';
import s from './SingleFollower.module.scss';

interface Props {
  follower: Follower;
}

export const SingleFollower = (props: Props) => {
  const { follower } = props;

  return (
    <div className={s.follower}>
      <a href={follower.link} target="_blank" rel="noreferrer">
        <Row wrap={false}>
          <Col span={6}>
            <Avatar src={follower.image} alt="Avatar" className={s.avatar}>
              <img src="/logo" alt="Logo" className={s.placeholder} />
            </Avatar>
          </Col>
          <Col>
            <Row>{follower.name}</Row>
            <Row className={s.account}>{follower.username}</Row>
          </Col>
        </Row>
      </a>
    </div>
  );
};
export default SingleFollower;
