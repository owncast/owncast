/* eslint-disable import/prefer-default-export */
export class User {
  constructor(u) {
    this.id = u.id;
    this.displayName = u.displayName;
    this.displayColor = u.displayColor;
    this.createdAt = u.createdAt;
    this.previousNames = u.previousNames;
    this.nameChangedAt = u.nameChangedAt;
    this.scopes = u.scopes;
    this.authenticated = u.authenticated;
    this.isBot = u.isBot;

    if (this.scopes && this.scopes.length > 0) {
      this.isModerator = this.scopes.includes('MODERATOR');
    }
  }

  id: string;

  displayName: string;

  displayColor: number;

  createdAt: Date;

  previousNames: string[];

  nameChangedAt: Date;

  scopes: string[];

  authenticated: boolean;

  isBot: boolean;

  isModerator: boolean;
}
