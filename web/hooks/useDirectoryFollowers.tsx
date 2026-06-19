import { useState, useEffect, useCallback } from 'react';
import { message } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../types/localization';

// DirectoryFollower is a directory that is currently listing this server: it
// followed this server with the directory marker and the operator approved it.
// It matches the Follower shape the backend serializes (link is the remote
// actor IRI).
export interface DirectoryFollower {
  link: string;
  name?: string;
  username?: string;
  image?: string;
  timestamp?: string;
}

export interface UseDirectoryFollowersResult {
  directories: DirectoryFollower[];
  loading: boolean;
  remove: (actorIRI: string) => Promise<void>;
  refetch: () => void;
}

const API_DIRECTORY_FOLLOWERS = '/api/admin/followers/directory';
const API_REMOVE_FOLLOWER = '/api/admin/followers/remove';

// useDirectoryFollowers fetches the directories that are currently listing this
// server and lets the operator remove one. Removing sends the directory a
// Reject of its follow, so it drops this server from its listing rather than
// leaving it showing offline forever.
export function useDirectoryFollowers(): UseDirectoryFollowersResult {
  const { t } = useTranslation();
  const [directories, setDirectories] = useState<DirectoryFollower[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchDirectories = useCallback(async () => {
    setLoading(true);
    try {
      const response = await fetch(API_DIRECTORY_FOLLOWERS, { credentials: 'include' });
      if (!response.ok) {
        throw new Error(`Failed to fetch directory listings: ${response.statusText}`);
      }
      const data = await response.json();
      setDirectories(data || []);
    } catch (err: any) {
      message.error(err.message || t(Localization.Admin.FeaturedStreams.failedToRemoveDirectory));
    } finally {
      setLoading(false);
    }
  }, [t]);

  const remove = async (actorIRI: string) => {
    const response = await fetch(API_REMOVE_FOLLOWER, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ actorIRI }),
    });

    if (!response.ok) {
      throw new Error(t(Localization.Admin.FeaturedStreams.failedToRemoveDirectory));
    }

    await fetchDirectories();
    message.success(t(Localization.Admin.FeaturedStreams.directoryRemoved));
  };

  // Fetch once on mount. fetchDirectories must NOT be a dependency: it is keyed
  // on `t`, which next-export-i18n returns fresh each render, so depending on it
  // would loop. Consumers refetch explicitly via remove.
  useEffect(() => {
    fetchDirectories();
  }, []);

  return { directories, loading, remove, refetch: fetchDirectories };
}
