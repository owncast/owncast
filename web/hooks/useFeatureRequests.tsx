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

// How often the sidebar badge re-checks for pending requests.
const PENDING_REQUESTS_POLL_INTERVAL = 60_000;

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

// usePendingFeatureRequestCount returns how many feature requests are waiting
// for the admin to approve, polling so the sidebar badge stays current. It is
// deliberately silent: a transient failure should never raise a toast from the
// layout, and it does nothing while federation is disabled (the endpoint
// requires it). Pass federationEnabled so the count clears and the polling
// stops when federation is off.
export function usePendingFeatureRequestCount(enabled: boolean): number {
  const [count, setCount] = useState(0);

  useEffect(() => {
    if (!enabled) {
      setCount(0);
      return undefined;
    }

    let cancelled = false;
    const load = async () => {
      try {
        const response = await fetch(API_FEATURE_REQUESTS, { credentials: 'include' });
        if (!response.ok) {
          return;
        }
        const data = await response.json();
        if (!cancelled) {
          setCount((data.requests || []).length);
        }
      } catch {
        // Silent: the badge is a hint, not a place to surface fetch errors.
      }
    };

    load();
    const intervalId = setInterval(load, PENDING_REQUESTS_POLL_INTERVAL);
    return () => {
      cancelled = true;
      clearInterval(intervalId);
    };
  }, [enabled]);

  return count;
}
