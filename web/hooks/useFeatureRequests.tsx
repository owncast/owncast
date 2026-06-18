import { useState, useEffect, useCallback } from 'react';
import { message } from 'antd';
import { useTranslation } from 'next-export-i18n';
import { Localization } from '../types/localization';

// FeatureRequest is a pending request from another Owncast server asking to
// feature this server's stream in its directory. It matches the Follower
// shape the backend serializes (link is the remote actor IRI).
export interface FeatureRequest {
  link: string;
  name?: string;
  username?: string;
  image?: string;
  timestamp?: string;
}

export interface UseFeatureRequestsResult {
  requests: FeatureRequest[];
  loading: boolean;
  approve: (actorIRI: string) => Promise<void>;
  reject: (actorIRI: string) => Promise<void>;
  refetch: () => void;
}

const API_FEATURE_REQUESTS = '/api/admin/federation/feature-requests';
const API_APPROVE_FOLLOWER = '/api/admin/followers/approve';

// useFeatureRequests fetches pending requests from other servers to feature
// this stream and lets the admin approve or reject them. Approval reuses the
// follower-approval endpoint, which records the approval and returns the
// ActivityPub Accept that completes the featured-streams handshake.
export function useFeatureRequests(): UseFeatureRequestsResult {
  const { t } = useTranslation();
  const [requests, setRequests] = useState<FeatureRequest[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchRequests = useCallback(async () => {
    setLoading(true);
    try {
      const response = await fetch(API_FEATURE_REQUESTS, { credentials: 'include' });
      if (!response.ok) {
        throw new Error(`Failed to fetch feature requests: ${response.statusText}`);
      }
      const data = await response.json();
      setRequests(data.requests || []);
    } catch (err: any) {
      message.error(err.message || t(Localization.Admin.FeaturedStreams.failedToApprove));
    } finally {
      setLoading(false);
    }
  }, [t]);

  const respond = async (actorIRI: string, approved: boolean) => {
    const response = await fetch(API_APPROVE_FOLLOWER, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ actorIRI, approved }),
    });

    if (!response.ok) {
      const errorKey = approved
        ? Localization.Admin.FeaturedStreams.failedToApprove
        : Localization.Admin.FeaturedStreams.failedToReject;
      throw new Error(t(errorKey));
    }

    await fetchRequests();
    message.success(
      t(
        approved
          ? Localization.Admin.FeaturedStreams.featureRequestApproved
          : Localization.Admin.FeaturedStreams.featureRequestRejected,
      ),
    );
  };

  const approve = (actorIRI: string) => respond(actorIRI, true);
  const reject = (actorIRI: string) => respond(actorIRI, false);

  // Fetch once on mount. fetchRequests must NOT be a dependency here: it is a
  // useCallback keyed on `t`, and next-export-i18n returns a fresh `t` on every
  // render, so fetchRequests changes identity every render. Depending on it
  // would re-run this effect on each render -> setState -> re-render -> refetch,
  // an infinite fetch loop. Consumers refetch explicitly via approve/reject.
  useEffect(() => {
    fetchRequests();
  }, []);

  return { requests, loading, approve, reject, refetch: fetchRequests };
}
