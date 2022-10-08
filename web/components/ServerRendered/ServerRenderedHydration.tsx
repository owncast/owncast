import { FC } from 'react';

export type ServerRenderedHydrationProps = {
  hydrationScript: string;
};

export const ServerRenderedHydration: FC<ServerRenderedHydrationProps> = ({ hydrationScript }) => (
  // eslint-disable-next-line react/no-danger
  <script dangerouslySetInnerHTML={{ __html: hydrationScript }} />
);
