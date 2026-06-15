# Server Federation Phase 4 - Implementation Plan
## ActivityPub Stream Ended Status Updates

### Overview
This implementation adds support for sending and receiving ActivityPub "Leave" activities when an Owncast server's live stream ends. This allows federated servers to update their status tracking of other servers in real-time.

### Implementation Components

## 1. ActivityPub Leave Activity Handling

### 1.1 Outbound Leave Activity (When Local Stream Ends)
**Location**: Same location where `GoLive` is called when stream starts

- **Trigger Point**: When stream transitions from online to offline (parallel to GoLive but for offline)
- **Tasks**:
  - Create ActivityPub "Leave" activity with:
    - Actor: Local Owncast server actor
    - Object: URL of the local server
    - Custom metadata using `SetOwncastMetadata`: Server logo, name, stream title, description, tags
  - Queue activity for delivery to all federated followers
  - Use existing ActivityPub outbox infrastructure

### 1.2 Inbound Leave Activity Reception
**Location**: `activitypub/inbox/worker.go` or similar inbox handler

- **Trigger**: Receive "Leave" activity in inbox
- **Tasks**:
  - Validate activity structure and signature
  - Extract server metadata using `ParseOwncastMetadata`
  - Update federated server status to offline
  - Update cached server metadata (logo, name, stream info, tags)
  - Trigger UI update notifications if needed

## 2. ActivityPub Activity Structure

### Leave Activity Format
```json
{
  "@context": "https://www.w3.org/ns/activitystreams",
  "type": "Leave",
  "id": "https://example.owncast.tv/activities/leave-123",
  "actor": "https://example.owncast.tv/federation/user/",
  "object": "https://example.owncast.tv",
  "published": "2024-01-01T00:00:00Z",
  "logo": "https://example.owncast.tv/logo.png",
  "name": "Example Server",
  "streamTitle": "My Stream",
  "streamDescription": "Stream description",
  "tags": ["gaming", "tech"]
}
```

Note: Metadata fields are set at the root level using `SetOwncastMetadata` and parsed using `ParseOwncastMetadata`, following the existing pattern.

## 3. Integration Points

### 3.1 Stream Lifecycle Integration
**Files to Modify**:
- Same file where `GoLive` is implemented for stream start
- Stream offline transition handler (parallel to GoLive)

**Implementation**:
- Create `GoOffline` function parallel to `GoLive`
- Hook into existing stream offline event
- Trigger Leave activity creation with metadata
- Ensure activity is sent after stream fully offline

### 3.2 ActivityPub Infrastructure Integration
**Existing Components to Use**:
- `activitypub/outbox/` - For sending activities
- `activitypub/inbox/` - For receiving activities
- `activitypub/persistence/` - For queue management
- `activitypub/crypto/` - For signing activities

### 3.3 Federation Server Status Management
**Files to Modify**:
- Federation server status tracker/cache
- UI notification system for status updates

## 4. Testing Strategy

### 4.1 Unit Tests
- **Leave Activity Creation**: Test proper activity structure
- **Leave Activity Processing**: Test status update logic
- **Metadata Extraction**: Test parsing of custom metadata

### 4.2 Integration Tests
**Location**: `test/automated/api/`
- Test Leave activity endpoint reception
- Test status update propagation
- Test metadata update persistence

### 4.3 End-to-End Testing
- Set up two Owncast instances
- Federate them together
- Start/stop stream on one instance
- Verify Leave activity received on other instance
- Verify status updated correctly

## 5. Implementation Steps

### Phase 1: Leave Activity Sending (2-3 days)
1. Identify where GoLive is called for stream start
2. Create GoOffline function parallel to GoLive
3. Create Leave activity builder function
4. Use SetOwncastMetadata to add metadata to activity
5. Integrate with existing outbox queue
6. Test activity creation and queuing

### Phase 2: Leave Activity Reception (2-3 days)
1. Add Leave activity handler to inbox worker
2. Use ParseOwncastMetadata for metadata extraction
3. Update federated server status
4. Persist metadata changes
5. Test activity processing

### Phase 3: Integration & Testing (2 days)
1. Full integration testing
2. Error handling and edge cases
3. Performance testing with multiple federated servers
4. Documentation updates

### Phase 4: UI Updates (Optional - 1 day)
1. Real-time status updates in UI
2. Display updated metadata
3. Visual indicators for offline servers

## 6. Technical Considerations

### Performance
- Activity sending should be async/non-blocking
- Batch updates if multiple servers go offline simultaneously
- Cache server metadata to reduce database queries

### Error Handling
- Retry failed activity deliveries
- Handle malformed activities gracefully
- Log failures for debugging

### Security
- Verify activity signatures
- Validate actor authorization
- Sanitize metadata fields

### Backwards Compatibility
- Ensure compatibility with servers not supporting Leave activities
- Graceful degradation if metadata missing
- Version checking for activity format

## 7. Dependencies

### Existing Code to Review
- `activitypub/` - ActivityPub implementation
- `core/webhooks/` - Stream event handlers
- `core/data/` - Database models
- `models/` - Data structures

### External Dependencies
- No new external dependencies required
- Uses existing ActivityPub libraries

## 8. Success Criteria

- [ ] Leave activity sent when stream goes offline
- [ ] Leave activity received and processed correctly
- [ ] Federated server status updated to offline
- [ ] Server metadata updated with latest information
- [ ] All tests passing
- [ ] No performance regression
- [ ] Error handling implemented
- [ ] Documentation updated

## 9. Risks & Mitigation

### Risk 1: Activity Delivery Failures
**Mitigation**: Implement retry mechanism with exponential backoff

### Risk 2: Metadata Size
**Mitigation**: Limit metadata fields, implement size validation

### Risk 3: Race Conditions
**Mitigation**: Use proper locking for status updates

## 10. Future Enhancements

- Support for other activity types (Join, Update)
- Bulk status updates for efficiency
- WebSocket notifications for real-time updates
- Historical status tracking

## Timeline Estimate

**Total Duration**: 7-8 days

- Days 1-3: Outbound Leave activity implementation
- Days 4-6: Inbound Leave activity handling
- Days 7-8: Integration testing and refinements

## Notes for Implementation

1. Follow existing ActivityPub patterns in codebase
2. Use OpenAPI spec generation for any new endpoints
3. Add comprehensive logging for debugging
4. Consider rate limiting for activity processing
5. Ensure all strings support localization