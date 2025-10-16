# Data Model: ChromLLM

**Feature**: Smart Chrome Extension with LLM Integration
**Date**: 2025-10-16
**Branch**: 001-specify-scripts-bash

## Core Entities

### UserSession
Represents an authenticated user session with conversation history and preferences.

**Fields**:
- `id` (string): Unique session identifier
- `user_id` (string): User identifier (optional, for future multi-user support)
- `created_at` (timestamp): Session creation time
- `last_active` (timestamp): Last activity timestamp
- `preferences` (object): User preferences (theme, model settings, etc.)
- `jwt_token` (string): Authentication token

**Relationships**:
- Has many `Conversation` entities
- Belongs to optional `User` (future extension)

**Validation Rules**:
- `id` must be unique
- `created_at` must be before or equal to `last_active`
- `jwt_token` must be valid and not expired

### Conversation
Represents a series of questions and answers within a specific webpage context.

**Fields**:
- `id` (string): Unique conversation identifier
- `session_id` (string): Reference to UserSession
- `url` (string): Webpage URL where conversation started
- `title` (string): Webpage title
- `created_at` (timestamp): Conversation start time
- `updated_at` (timestamp): Last message timestamp
- `message_count` (integer): Total messages in conversation
- `is_active` (boolean): Whether conversation is currently active

**Relationships**:
- Belongs to `UserSession`
- Has many `Message` entities
- Has one `WebpageContent` (optional)

**Validation Rules**:
- `url` must be valid URL format
- `message_count` must be >= 0
- `updated_at` must be >= `created_at`

### Message
Represents a single question or answer in a conversation.

**Fields**:
- `id` (string): Unique message identifier
- `conversation_id` (string): Reference to Conversation
- `type` (enum): "user_question" or "llm_response"
- `content` (string): Message text content
- `timestamp` (timestamp): When message was sent/received
- `sequence_number` (integer): Order in conversation
- `metadata` (object): Additional context (model used, tokens, etc.)

**Relationships**:
- Belongs to `Conversation`

**Validation Rules**:
- `type` must be one of allowed values
- `content` must not be empty
- `sequence_number` must be unique within conversation
- `timestamp` must be in chronological order

### WebpageContent
Represents extracted and processed content from webpages for LLM context.

**Fields**:
- `id` (string): Unique content identifier
- `url` (string): Source webpage URL
- `title` (string): Webpage title
- `extracted_text` (string): Main content text
- `extraction_timestamp` (timestamp): When content was extracted
- `content_hash` (string): Hash of extracted content for deduplication
- `metadata` (object): Extraction metadata (word count, language, etc.)

**Relationships**:
- Associated with one `Conversation`

**Validation Rules**:
- `url` must be valid URL format
- `extracted_text` must not be empty
- `content_hash` must be unique for identical content

### LLMResponse
Represents streaming responses from large language model APIs.

**Fields**:
- `id` (string): Unique response identifier
- `message_id` (string): Reference to Message
- `stream_id` (string): Streaming session identifier
- `content` (string): Accumulated response content
- `is_complete` (boolean): Whether streaming is finished
- `tokens_used` (integer): Total tokens consumed
- `model_used` (string): LLM model identifier
- `created_at` (timestamp): Response start time
- `completed_at` (timestamp): Response completion time

**Relationships**:
- Belongs to `Message`

**Validation Rules**:
- `stream_id` must be unique for active streams
- `tokens_used` must be >= 0
- `completed_at` must be >= `created_at` if set

## State Transitions

### Conversation Lifecycle
```
created → active → archived
```

**Rules**:
- New conversations start in "active" state
- Conversations become "archived" after 30 days of inactivity
- Archived conversations can be reactivated by new messages

### Message Flow
```
user_question → llm_response (streaming) → complete
```

**Rules**:
- Each user question must be followed by an LLM response
- LLM responses can be streamed in chunks
- Response is marked complete when streaming finishes

### Session Management
```
created → active → expired → renewed
```

**Rules**:
- Sessions expire after 24 hours of inactivity
- JWT tokens have 1-hour validity
- Sessions can be renewed with valid refresh tokens

## Data Storage Strategy

### Redis (Cache)
- User sessions with TTL
- Active conversations
- LLM response streaming state
- Rate limiting counters

**Key Patterns**:
- `session:{session_id}` - User session data
- `conversation:{conversation_id}` - Active conversation
- `stream:{stream_id}` - LLM streaming state
- `rate_limit:{user_id}` - API rate limiting

### Browser localStorage (Extension)
- User preferences
- Recent conversation history
- Authentication tokens
- Extension settings

**Key Patterns**:
- `user_preferences` - User settings
- `recent_conversations` - Conversation metadata
- `auth_token` - JWT token
- `extension_settings` - Extension configuration

### Database (Future Extension)
- Persistent conversation history
- User accounts and profiles
- Analytics and usage data

## Security Considerations

- JWT tokens contain minimal user information
- Session data encrypted in transit and at rest
- Content extraction respects user privacy
- API rate limiting to prevent abuse
- Input validation for all user-provided data

## Performance Optimizations

- Redis caching for frequently accessed data
- Content deduplication using hashing
- Pagination for large conversation histories
- Efficient indexing for search operations
- Background processing for non-critical tasks