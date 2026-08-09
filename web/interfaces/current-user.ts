export interface CurrentUser {
  id: string;
  displayName: string;
  displayColor: number;
  isModerator: boolean;
  // authenticated is true when this identity proved itself through an auth
  // provider (IndieAuth/FediAuth), not merely registered for chat.
  authenticated?: boolean;
  // createdAt is when this chat identity was first registered, used to gate
  // features (like clip creation) on account age.
  createdAt?: string | Date;
}
