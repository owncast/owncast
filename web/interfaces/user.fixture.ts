import { User } from './user.model';

export const createUser = (name: string, color: number, createdAt: Date): User => ({
  id: name,
  displayName: name,
  displayColor: color,
  createdAt,
  authenticated: true,
  nameChangedAt: createdAt,
  previousNames: [],
  scopes: [],
  isBot: false,
  isModerator: false,
});

export const spidermanUser = createUser('Spiderman', 1, new Date(2020, 1, 2));
export const grootUser = createUser('Groot', 1, new Date(2020, 2, 3));
