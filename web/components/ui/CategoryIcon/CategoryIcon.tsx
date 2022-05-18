import React from 'react';
import s from './CategoryIcon.module.scss';

import GAME from './icons/games.svg';
import CHAT from './icons/chat.svg';
import CONFERENCE from './icons/conference.svg';
import GOVERNMENT from './icons/gov.svg';
import RELIGION from './icons/religion.svg';
import TECH from './icons/tech.svg';
import MUSIC from './icons/music.svg';

interface Props {
  tags: string[];
}

const ICON_MAPPING = {
  game: GAME,
  gamer: GAME,
  games: GAME,
  playing: GAME,

  chat: CHAT,
  chatting: CHAT,

  conference: CONFERENCE,
  event: CONFERENCE,
  convention: CONFERENCE,

  government: GOVERNMENT,

  religion: RELIGION,
  god: RELIGION,
  church: RELIGION,
  bible: RELIGION,

  tech: TECH,
  technology: TECH,
  code: TECH,
  software: TECH,
  linux: TECH,
  development: TECH,

  music: MUSIC,
  listening: MUSIC,
};

function getCategoryIconFromTags(tags: string[]): any {
  const categories = tags.map(tag => ICON_MAPPING[tag]).filter(Boolean);
  if (categories.length > 0) {
    return categories[0];
  }
  return null;
}

export default function CategoryIcon({ tags }: Props) {
  const Icon = getCategoryIconFromTags(tags);
  if (!Icon) {
    return null;
  }

  return <Icon className={s.icon} />;
}
