export interface User {
  id: string;
  displayName: string;
  displayColor: number;
  createdAt: Date;
  previousNames: string[];
  nameChangedAt: Date;
  scopes: string[];
  authenticated: boolean;
}
