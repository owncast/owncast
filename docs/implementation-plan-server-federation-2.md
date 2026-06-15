# Implementation Plan: Server Federation Phase 2 - ActivityPub Status Updates

## Overview
This plan outlines the implementation of ActivityPub "Arrive" and "Leave" activity support for Owncast server federation, allowing servers to share online/offline status and stream metadata.

## Current State Analysis
- ✅ Existing ActivityPub infrastructure (inbox, outbox, workers, crypto, resolvers)
- ✅ Basic federation enabled with follower support
- ✅ Frontend components partially exist (FederatedServersTable, useFederatedServers hook)
- ✅ Existing models.Status struct with StreamTitle field (reusable pattern)
- ✅ Stream metadata handling patterns in config repository
- ❌ No database table for federated servers
- ❌ No Arrive/Leave activity handlers
- ❌ No API endpoints for federated servers management

## Implementation Plan

### Phase 1: Database Layer
#### 1.1 Database Schema (Priority: High)
- **File**: `db/schema.sql`
- **Action**: Add new table `federated_servers` with fields:
  - `id` (PRIMARY KEY, auto-increment)
  - `iri` (TEXT, UNIQUE) - Server IRI/URL
  - `name` (TEXT) - Server display name
  - `logo_url` (TEXT) - Server logo image URL
  - `is_online` (BOOLEAN, DEFAULT FALSE) - Current online status
  - `stream_title` (TEXT) - Current stream title (when online)
  - `stream_description` (TEXT) - Current stream description
  - `stream_tags` (TEXT) - JSON array of tags
  - `thumbnail_url` (TEXT) - Stream thumbnail URL
  - `last_seen_online` (TIMESTAMP) - Last time server was seen online
  - `last_status_update` (TIMESTAMP) - Last activity received
  - `added_at` (TIMESTAMP DEFAULT CURRENT_TIMESTAMP)
  - `followed_at` (TIMESTAMP) - When we started following this server

#### 1.2 Database Queries (Priority: High)
- **File**: `db/query.sql`
- **Actions**: Add SQLC queries for:
  - `GetFederatedServers` - List all federated servers
  - `GetFederatedServer` - Get server by IRI
  - `AddFederatedServer` - Add new federated server
  - `UpdateFederatedServerStatus` - Update server online/offline status
  - `UpdateFederatedServerMetadata` - Update stream metadata
  - `RemoveFederatedServer` - Remove server from federation

#### 1.3 Generated Go Models (Priority: High)
- **File**: Auto-generated via sqlc
- **Action**: Run `sqlc generate` to create Go structs

### Phase 2: ActivityPub Activity Handlers
#### 2.1 Arrive Activity Handler (Priority: High)
- **File**: `activitypub/inbox/arrive.go` (new)
- **Function**: `handleArriveInboxRequest(context.Context, vocab.ActivityStreamsArrive) error`
- **Logic**:
  - Extract actor IRI (the federated Owncast server)
  - Extract object URL (should be the server URL)
  - Validate actor is a known federated server
  - Update server status to online
  - Extract stream metadata from activity (title, description, tags, thumbnail)
  - Update `last_seen_online` and `last_status_update` timestamps
  - Follow existing patterns from `models.Status` for stream metadata handling

#### 2.2 Leave Activity Handler (Priority: High)
- **File**: `activitypub/inbox/leave.go` (new)
- **Function**: `handleLeaveInboxRequest(context.Context, vocab.ActivityStreamsLeave) error`
- **Logic**:
  - Extract actor IRI (the federated Owncast server)
  - Extract object URL (should be the server URL)
  - Validate actor is a known federated server
  - Update server status to offline
  - Clear current stream metadata
  - Update `last_status_update` timestamp

#### 2.3 Worker Integration (Priority: High)
- **File**: `activitypub/inbox/worker.go`
- **Action**: Update resolver call to include new activity handlers:
  ```go
  resolvers.Resolve(context.Background(), request.Body,
    handleUpdateRequest, handleFollowInboxRequest, handleLikeRequest,
    handleAnnounceRequest, handleUndoInboxRequest, handleCreateRequest,
    handleArriveInboxRequest, handleLeaveInboxRequest)
  ```

### Phase 3: Database Persistence Layer
#### 3.1 Federated Servers Repository (Priority: High)
- **File**: `persistence/federatedservers/federatedservers.go` (new)
- **Functions**:
  - `GetFederatedServers() ([]models.FederatedServer, error)`
  - `GetFederatedServer(iri string) (*models.FederatedServer, error)`
  - `AddFederatedServer(iri, name, logoURL string) error`
  - `UpdateServerStatus(iri string, isOnline bool, metadata *FederatedStreamUpdate) error`
  - `RemoveFederatedServer(iri string) error`

#### 3.2 Models (Priority: High)
- **File**: `models/federatedserver.go` (new)
- **Analysis**: Existing `models.Status` already contains `StreamTitle`, so we'll align with existing patterns
- **Structs**:
  ```go
  type FederatedServer struct {
    ID                 int       `json:"id"`
    IRI                string    `json:"iri"`
    Name               string    `json:"name"`
    LogoURL            *string   `json:"logoUrl,omitempty"`
    IsOnline           bool      `json:"isOnline"`
    StreamTitle        *string   `json:"streamTitle,omitempty"`
    StreamDescription  *string   `json:"streamDescription,omitempty"`
    StreamTags         []string  `json:"tags,omitempty"`
    ThumbnailURL       *string   `json:"thumbnailUrl,omitempty"`
    LastSeenOnline     *time.Time `json:"lastSeenOnline,omitempty"`
    LastStatusUpdate   *time.Time `json:"lastStatusUpdate,omitempty"`
    AddedAt            time.Time `json:"addedAt"`
    FollowedAt         *time.Time `json:"followedAt,omitempty"`
  }

  // Lightweight struct for passing stream metadata in ActivityPub handlers
  // Reuses existing patterns from models.Status
  type FederatedStreamUpdate struct {
    Title       *string  `json:"title,omitempty"`
    Description *string  `json:"description,omitempty"`
    Tags        []string `json:"tags,omitempty"`
    ThumbnailURL *string  `json:"thumbnailUrl,omitempty"`
  }
  ```

### Phase 4: API Endpoints
#### 4.1 OpenAPI Specification (Priority: Medium)
- **File**: `openapi.yaml`
- **Action**: Add API definitions for federated servers endpoints:
  - `GET /api/federation/servers` - List federated servers (public)
    - **Tags**: `['Internal', 'Federation']`
    - **Security**: None (public endpoint)
  - `POST /api/admin/federation/servers` - Add federated server (admin)
    - **Tags**: `['Internal', 'Admin', 'Federation']`
    - **Security**: `BasicAuth: []` (requires admin authentication)
  - `DELETE /api/admin/federation/servers/{id}` - Remove server (admin)
    - **Tags**: `['Internal', 'Admin', 'Federation']`
    - **Security**: `BasicAuth: []` (requires admin authentication)
- **Response Types**: Follow existing patterns with `401BasicAuth` responses for admin endpoints

#### 4.2 Code Generation (Priority: Medium)
- **Action**: Run `build/gen-api.sh` to generate API stubs from OpenAPI spec
- **Note**: This must be done before implementing handlers since generated code provides interfaces

#### 4.3 API Handlers (Priority: Medium)
- **File**: `webserver/handlers/admin/federatedservers.go` (new)
- **Functions**:
  - `GetFederatedServers(w http.ResponseWriter, r *http.Request)`
  - `AddFederatedServer(w http.ResponseWriter, r *http.Request)`
  - `RemoveFederatedServer(w http.ResponseWriter, r *http.Request)`
- **Note**: Implement generated interfaces from step 4.2

#### 4.4 Router Integration (Priority: Medium)
- **File**: `webserver/router/router.go`
- **Action**: Add route mappings for federated servers endpoints

### Phase 5: Frontend Updates
#### 5.1 API Integration (Priority: Medium)
- **File**: `web/hooks/useFederatedServers.tsx`
- **Action**: Update API endpoints to match backend implementation
- **Status**: Already mostly implemented, needs minor adjustments

#### 5.2 Admin Interface (Priority: Low)
- **Files**:
  - `web/components/admin/FederatedServers/FederatedServersTable.tsx`
  - `web/components/admin/FederatedServers/AddServerForm.tsx`
- **Action**: Verify compatibility with new backend API structure

#### 5.3 Public Interface (Priority: Low)
- **Files**:
  - `web/components/ui/StreamsTab/StreamsTab.tsx`
  - `web/components/ui/Content/Content.tsx`
- **Action**: Display federated servers list on main interface

### Phase 6: Testing
#### 6.1 Unit Tests (Priority: Medium)
- **Files**:
  - `activitypub/inbox/arrive_test.go` (new)
  - `activitypub/inbox/leave_test.go` (new)
  - `persistence/federatedservers/federatedservers_test.go` (new)

#### 6.2 API Integration Tests (Priority: Medium)
- **File**: `test/automated/api/010_federation_servers.test.js` (new)

#### 6.3 Browser Tests (Priority: Low)
- **File**: `test/automated/browser/cypress/e2e/federation/federated_servers.cy.js` (new)

## Implementation Order (Recommended)
1. **Database Schema & Queries** (Essential foundation)
2. **ActivityPub Handlers** (Core functionality)
3. **Persistence Layer** (Data access)
4. **API Specification & Code Generation** (Generate interfaces first)
5. **API Handler Implementation** (Implement generated interfaces)
6. **Frontend Integration** (User interface)
7. **Testing** (Quality assurance)

## Technical Considerations
### Security
- Validate all incoming ActivityPub activities using existing signature verification
- Sanitize and validate all metadata fields (titles, descriptions, URLs)
- Rate limiting for status updates from federated servers

### Performance
- Consider caching federated server list for public API
- Implement pagination for server lists if needed
- Index database appropriately for common queries

### Error Handling
- Graceful handling of malformed ActivityPub activities
- Retry logic for failed server communications
- Proper logging for debugging federation issues

### Data Integrity
- Ensure atomic updates to server status
- Handle duplicate activities gracefully
- Clean up offline servers periodically

### Code Reuse & Consistency
- **Leverage existing patterns**: Stream metadata handling should follow `models.Status` patterns
- **Reuse configuration patterns**: Stream title/description handling similar to existing config repository methods
- **Consistent JSON field naming**: Match existing API response structures where possible

## Dependencies
- Existing ActivityPub infrastructure ✅
- SQLC for query generation ✅
- Go-fed ActivityStreams library ✅
- Existing federation configuration ✅

## Rollout Strategy
1. Deploy database migrations
2. Deploy backend changes with feature flag (if desired)
3. Deploy frontend changes
4. Enable feature in production

## Success Criteria
- ✅ Owncast servers can exchange Arrive/Leave activities
- ✅ Server online/offline status is accurately tracked
- ✅ Stream metadata is properly synchronized
- ✅ Admin interface allows management of federated servers
- ✅ Public interface displays federated servers
- ✅ All tests pass
- ✅ No performance degradation