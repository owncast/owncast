import { Avatar, Col, Row } from 'antd';
import React, { FC } from 'react';
import { Follower } from '../../../../interfaces/follower';
import styles from './SingleFollower.module.scss';

export type SingleFollowerProps = {
  follower: Follower;
};

export const SingleFollower: FC<SingleFollowerProps> = ({ follower }) => (
  <div className={styles.follower}>
    <a href={follower.link} target="_blank" rel="noreferrer">
      <Row wrap={false}>
        <Col span={6}>
          <Avatar src={follower.image} alt="Avatar" className={styles.avatar}>
            <img src="/logo" alt="Logo" className={styles.placeholder} />
          </Avatar>
        </Col>
        <Col>
          <Row>{follower.name}</Row>
          <Row className={styles.account}>{follower.username}</Row>
        </Col>
      </Row>
    </a>
  </div>
);
export default SingleFollower;
