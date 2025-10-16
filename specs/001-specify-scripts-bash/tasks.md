# Tasks: ChromLLM - Smart Chrome Extension with LLM Integration

**Input**: Design documents from `/specs/001-specify-scripts-bash/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: The examples below include test tasks. Tests are OPTIONAL - only include them if explicitly requested in the feature specification.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create project structure per implementation plan
- [ ] T002 Initialize Chrome extension project with manifest.json in chrome-extension/
- [ ] T003 Initialize Go backend project with go.mod in backend/
- [ ] T004 [P] Configure Jest testing framework for Chrome extension
- [ ] T005 [P] Configure Go testing framework for backend
- [ ] T006 [P] Set up ESLint and Prettier for Chrome extension
- [ ] T007 [P] Set up Go linting and formatting tools
- [ ] T008 Create documentation structure in docs/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T009 Setup Redis configuration and connection in backend/pkg/cache/
- [ ] T010 [P] Implement JWT authentication middleware in backend/internal/middleware/auth.go
- [ ] T011 [P] Setup Gin framework with basic routing in backend/cmd/server/main.go
- [ ] T012 [P] Create configuration management in backend/config/config.go
- [ ] T013 [P] Implement error handling middleware in backend/internal/middleware/error.go
- [ ] T014 [P] Setup logging infrastructure in backend/internal/middleware/logging.go
- [ ] T015 [P] Create base models for UserSession, Conversation, Message, WebpageContent, LLMResponse in backend/internal/models/
- [ ] T016 [P] Implement Redis session storage in backend/pkg/cache/session.go
- [ ] T017 [P] Setup WebSocket server infrastructure in backend/internal/websocket/

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Basic Content Extraction and Question Answering (Priority: P1) 🎯 MVP

**Goal**: Enable users to extract webpage content and ask questions about it

**Independent Test**: Can be fully tested by installing the Chrome extension, visiting a webpage, extracting content, and asking questions about that content

### Tests for User Story 1 (OPTIONAL - only if tests requested) ⚠️

**NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T018 [P] [US1] Contract test for session creation endpoint in tests/contract/test_sessions.go
- [ ] T019 [P] [US1] Contract test for conversation creation endpoint in tests/contract/test_conversations.go
- [ ] T020 [P] [US1] Integration test for content extraction and question flow in tests/integration/test_content_extraction.go
- [ ] T021 [P] [US1] Unit test for content extraction logic in chrome-extension/content/content.test.js

### Implementation for User Story 1

- [ ] T022 [P] [US1] Create Chrome extension manifest.json in chrome-extension/manifest.json
- [ ] T023 [P] [US1] Implement content script for webpage extraction in chrome-extension/content/content.js
- [ ] T024 [P] [US1] Create popup UI interface in chrome-extension/popup/popup.html
- [ ] T025 [P] [US1] Implement popup JavaScript logic in chrome-extension/popup/popup.js
- [ ] T026 [P] [US1] Style popup interface in chrome-extension/popup/popup.css
- [ ] T027 [P] [US1] Implement session creation endpoint in backend/internal/handlers/sessions.go
- [ ] T028 [P] [US1] Implement conversation creation endpoint in backend/internal/handlers/conversations.go
- [ ] T029 [P] [US1] Implement message sending endpoint in backend/internal/handlers/messages.go
- [ ] T030 [P] [US1] Create OpenAI integration package in backend/pkg/llm/openai.go
- [ ] T031 [US1] Implement content extraction service in backend/internal/services/extraction.go
- [ ] T032 [US1] Implement question answering service in backend/internal/services/qa.go
- [ ] T033 [US1] Add validation and error handling for content extraction
- [ ] T034 [US1] Add logging for user story 1 operations

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Real-time Streaming Responses (Priority: P2)

**Goal**: Provide real-time streaming of LLM responses for immediate user feedback

**Independent Test**: Can be tested by asking a complex question and verifying that responses appear incrementally rather than all at once

### Tests for User Story 2 (OPTIONAL - only if tests requested) ⚠️

- [ ] T035 [P] [US2] Contract test for WebSocket streaming endpoint in tests/contract/test_stream.go
- [ ] T036 [P] [US2] Integration test for real-time streaming in tests/integration/test_streaming.go
- [ ] T037 [P] [US2] Unit test for streaming logic in backend/pkg/llm/streaming.test.go

### Implementation for User Story 2

- [ ] T038 [P] [US2] Implement WebSocket streaming endpoint in backend/internal/websocket/stream.go
- [ ] T039 [P] [US2] Create WebSocket client in Chrome extension in chrome-extension/popup/websocket.js
- [ ] T040 [P] [US2] Implement OpenAI streaming API integration in backend/pkg/llm/streaming.go
- [ ] T041 [P] [US2] Create LLMResponse model for streaming state in backend/internal/models/llm_response.go
- [ ] T042 [US2] Implement streaming response handling in popup UI in chrome-extension/popup/streaming.js
- [ ] T043 [US2] Add streaming progress indicators in popup UI
- [ ] T044 [US2] Implement connection management and reconnection logic
- [ ] T045 [US2] Add error handling for streaming interruptions
- [ ] T046 [US2] Integrate with User Story 1 components for seamless streaming

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Session Management and History (Priority: P3)

**Goal**: Maintain conversation history and manage multiple sessions across browsing sessions

**Independent Test**: Can be tested by having multiple conversations, closing/reopening the extension, and verifying history persistence

### Tests for User Story 3 (OPTIONAL - only if tests requested) ⚠️

- [ ] T047 [P] [US3] Contract test for conversation history endpoint in tests/contract/test_history.go
- [ ] T048 [P] [US3] Integration test for session persistence in tests/integration/test_sessions.go
- [ ] T049 [P] [US3] Unit test for localStorage management in chrome-extension/popup/storage.test.js

### Implementation for User Story 3

- [ ] T050 [P] [US3] Implement conversation history endpoint in backend/internal/handlers/history.go
- [ ] T051 [P] [US3] Create localStorage management for extension data in chrome-extension/popup/storage.js
- [ ] T052 [P] [US3] Implement conversation list UI in popup in chrome-extension/popup/history.html
- [ ] T053 [P] [US3] Add conversation history JavaScript logic in chrome-extension/popup/history.js
- [ ] T054 [US3] Implement session persistence and renewal logic in backend/internal/services/session.go
- [ ] T055 [US3] Add conversation continuation across page navigation
- [ ] T056 [US3] Implement conversation archiving and cleanup in backend/internal/services/conversation.go
- [ ] T057 [US3] Add settings page for extension preferences in chrome-extension/options/options.html
- [ ] T058 [US3] Implement settings management in chrome-extension/options/options.js

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T059 [P] Documentation updates in docs/api/ and docs/extension/
- [ ] T060 [P] Code cleanup and refactoring across both extension and backend
- [ ] T061 [P] Performance optimization across all stories
- [ ] T062 [P] Additional unit tests in tests/unit/
- [ ] T063 [P] Security hardening and input validation
- [ ] T064 [P] Accessibility improvements for extension UI
- [ ] T065 [P] Error handling and user feedback improvements
- [ ] T066 [P] Rate limiting implementation in backend/internal/middleware/rate_limit.go
- [ ] T067 [P] Health check endpoint implementation in backend/internal/handlers/health.go
- [ ] T068 [P] Deployment configuration and scripts
- [ ] T069 [P] Run quickstart.md validation scenarios

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Depends on US1 for basic conversation flow
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Depends on US1 for conversation structure

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together (if tests requested):
Task: "Contract test for session creation endpoint in tests/contract/test_sessions.go"
Task: "Contract test for conversation creation endpoint in tests/contract/test_conversations.go"
Task: "Integration test for content extraction and question flow in tests/integration/test_content_extraction.go"
Task: "Unit test for content extraction logic in chrome-extension/content/content.test.js"

# Launch all models for User Story 1 together:
Task: "Create Chrome extension manifest.json in chrome-extension/manifest.json"
Task: "Implement content script for webpage extraction in chrome-extension/content/content.js"
Task: "Create popup UI interface in chrome-extension/popup/popup.html"
Task: "Implement popup JavaScript logic in chrome-extension/popup/popup.js"
Task: "Style popup interface in chrome-extension/popup/popup.css"
Task: "Implement session creation endpoint in backend/internal/handlers/sessions.go"
Task: "Implement conversation creation endpoint in backend/internal/handlers/conversations.go"
Task: "Implement message sending endpoint in backend/internal/handlers/messages.go"
Task: "Create OpenAI integration package in backend/pkg/llm/openai.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (Chrome extension + basic backend)
   - Developer B: User Story 2 (Streaming + WebSocket)
   - Developer C: User Story 3 (Session management + history)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence