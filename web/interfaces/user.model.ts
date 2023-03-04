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
  }

  id: string;

  displayName: string;

  displayColor: number;

  createdAt: Date;

  previousNames: string[];

  nameChangedAt: Date;

  scopes: string[];

  authenticated: boolean;

  public isModerator = (): boolean => {
    if (!this.scopes || this.scopes.length === 0) {
      return false;
    }

    return this.scopes.includes('moderator');
  };
}
