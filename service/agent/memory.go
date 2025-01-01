// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package agent

import (
	"context"
	"fmt"

	"github.com/clivern/actx0/db"
	"github.com/clivern/actx0/migration"

	"github.com/rs/zerolog/log"
)

// ConsolidateMemory summarizes and compacts session memory.
func (s *Service) ConsolidateMemory(ctx context.Context, req ConsolidateMemoryRequest) error {
	return ErrNotImplemented
}

// StoreMemory saves a memory item for a session.
func (s *Service) StoreMemory(ctx context.Context, req StoreMemoryRequest) (*db.SessionMemory, error) {
	return nil, ErrNotImplemented
}

// RecallMemory returns memories relevant to the query.
func (s *Service) RecallMemory(ctx context.Context, req RecallMemoryRequest) ([]MemoryHit, error) {
	return nil, ErrNotImplemented
}

// ListMemories returns stored memories for a session.
func (s *Service) ListMemories(ctx context.Context, sessionId db.Id, limit, offset int) ([]*db.SessionMemory, error) {
	return nil, ErrNotImplemented
}

// IndexSessionMessage embeds a session message into Qdrant.
func (s *Service) IndexSessionMessage(ctx context.Context, messageId db.Id) error {
	return ErrNotImplemented
}

// IndexSessionMemory embeds a session memory into Qdrant.
func (s *Service) IndexSessionMemory(ctx context.Context, memoryId db.Id) error {
	return ErrNotImplemented
}

// GenerateMemoryFromMessage extracts and stores memory from a session message.
func (s *Service) GenerateMemoryFromMessage(ctx context.Context, messageId db.Id) error {
	return ErrNotImplemented
}

// DeleteSessionMemory deletes vectors for one memory item.
func (s *Service) DeleteSessionMemory(ctx context.Context, memoryId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMemoriesCollection, map[string]string{
		migration.PayloadMemoryID: memoryId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("memory_id", memoryId.String()).
			Msg("Failed to delete memory vectors")
		return fmt.Errorf("%w: delete memory vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}

// DeleteSessionMemories deletes all memory vectors for a session.
func (s *Service) DeleteSessionMemories(ctx context.Context, agentId, sessionId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMemoriesCollection, map[string]string{
		migration.PayloadAgentID:   agentId.String(),
		migration.PayloadSessionID: sessionId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("agent_id", agentId.String()).
			Str("session_id", sessionId.String()).
			Msg("Failed to delete session memory vectors")
		return fmt.Errorf("%w: delete memory vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}

// DeleteSessionMessages deletes all message vectors for a session.
func (s *Service) DeleteSessionMessages(ctx context.Context, agentId, sessionId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMessagesCollection, map[string]string{
		migration.PayloadAgentID:   agentId.String(),
		migration.PayloadSessionID: sessionId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("agent_id", agentId.String()).
			Str("session_id", sessionId.String()).
			Msg("Failed to delete session message vectors")
		return fmt.Errorf("%w: delete message vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}

// DeleteSessionMessage deletes vectors for one message.
func (s *Service) DeleteSessionMessage(ctx context.Context, messageId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMessagesCollection, map[string]string{
		migration.PayloadMessageID: messageId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("message_id", messageId.String()).
			Msg("Failed to delete message vectors")
		return fmt.Errorf("%w: delete message vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}

// DeleteAgentSession deletes all message and memory vectors for a session.
func (s *Service) DeleteAgentSession(ctx context.Context, agentId, sessionId db.Id) error {
	err := s.DeleteSessionMessages(ctx, agentId, sessionId)
	if err != nil {
		return err
	}

	return s.DeleteSessionMemories(ctx, agentId, sessionId)
}

// DeleteAgent removes all session message and memory vectors for an agent.
func (s *Service) DeleteAgent(ctx context.Context, agentId db.Id) error {
	err := s.vectors.DeleteByFilter(ctx, migration.AgentSessionMessagesCollection, map[string]string{
		migration.PayloadAgentID: agentId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("agent_id", agentId.String()).
			Msg("Failed to delete agent message vectors")
		return fmt.Errorf("%w: delete message vectors: %v", ErrDeleteFailed, err)
	}

	err = s.vectors.DeleteByFilter(ctx, migration.AgentSessionMemoriesCollection, map[string]string{
		migration.PayloadAgentID: agentId.String(),
	})
	if err != nil {
		log.Warn().
			Err(err).
			Str("agent_id", agentId.String()).
			Msg("Failed to delete agent memory vectors")
		return fmt.Errorf("%w: delete memory vectors: %v", ErrDeleteFailed, err)
	}

	return nil
}
