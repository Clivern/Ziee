// Copyright 2026 Ziee. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"fmt"

	"github.com/clivern/ziee/db"

	"github.com/samber/lo"
)

// ConversationMessage is one assistant or user message to append to a conversation.
type ConversationMessage struct {
	TargetType   string
	TargetNumber int
	Role         string
	Author       string
	CommentId    *int64
	Body         string
}

// Conversation is the module for GitHub assistant conversation history.
type Conversation struct {
	ConversationRepository db.RepositoryConversationRepository
	MetaRepository         db.RepositoryConversationMetaRepository
	RepoRepository         db.RepositoriesRepository
}

// NewConversation creates a conversation module with the given stores.
func NewConversation(
	conversations db.RepositoryConversationRepository,
	meta db.RepositoryConversationMetaRepository,
	repos db.RepositoriesRepository,
) *Conversation {
	return &Conversation{
		ConversationRepository: conversations,
		MetaRepository:         meta,
		RepoRepository:         repos,
	}
}

// Record appends a message to the history of an issue or pull request.
func (c *Conversation) Record(githubRepoId int64, msg ConversationMessage) (*db.RepositoryConversation, error) {
	repo, err := c.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		return nil, fmt.Errorf("get repo: %w", err)
	}

	item := &db.RepositoryConversation{
		RepositoryId: repo.Id,
		TargetType:   msg.TargetType,
		TargetNumber: msg.TargetNumber,
		Role:         msg.Role,
		Author:       lo.EmptyableToPtr(msg.Author),
		CommentId:    msg.CommentId,
		Body:         msg.Body,
	}

	err = c.ConversationRepository.Create(item)
	if err != nil {
		return nil, fmt.Errorf("record conversation message: %w", err)
	}

	return item, nil
}

// History returns the messages of an issue or pull request, oldest first.
func (c *Conversation) History(githubRepoId int64, targetType string, targetNumber int) ([]*db.RepositoryConversation, error) {
	repo, err := c.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		return nil, fmt.Errorf("get repo: %w", err)
	}

	items, err := c.ConversationRepository.ListByTarget(repo.Id, targetType, targetNumber)
	if err != nil {
		return nil, fmt.Errorf("list conversation messages: %w", err)
	}

	return items, nil
}

// Clear deletes the message history of an issue or pull request.
func (c *Conversation) Clear(githubRepoId int64, targetType string, targetNumber int) error {
	repo, err := c.RepoRepository.GetByGitHubId(githubRepoId)
	if err != nil {
		return fmt.Errorf("get repo: %w", err)
	}

	err = c.ConversationRepository.DeleteByTarget(repo.Id, targetType, targetNumber)
	if err != nil {
		return fmt.Errorf("clear conversation: %w", err)
	}

	return nil
}

// UpsertMeta stores a repository_conversations_meta value by message id.
func (c *Conversation) UpsertMeta(conversationId db.Id, key, value string) error {
	err := c.MetaRepository.Upsert(conversationId, key, value)
	if err != nil {
		return fmt.Errorf("upsert conversation meta: %w", err)
	}

	return nil
}

// GetMeta returns a repository_conversations_meta value by message id.
func (c *Conversation) GetMeta(conversationId db.Id, key string) (string, error) {
	meta, err := c.MetaRepository.Get(conversationId, key)
	if err != nil {
		return "", fmt.Errorf("get conversation meta: %w", err)
	}

	return lo.FromPtr(meta).Value, nil
}
