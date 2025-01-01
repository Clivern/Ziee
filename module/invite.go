// Copyright 2026 Clivern. All rights reserved.
// License can be found in the LICENSE file.

package module

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/clivern/actx0/conf"
	"github.com/clivern/actx0/db"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

const (
	InviteStatusPending  = "pending"
	InviteStatusAccepted = "accepted"
	InviteStatusRejected = "rejected"
	InviteStatusExpired  = "expired"
)

var (
	ErrInviteNotFound          = errors.New("invite not found")
	ErrInviteNoLongerValid     = errors.New("invite no longer valid")
	ErrInviteExpired           = errors.New("invite expired")
	ErrInviteEmailMismatch     = errors.New("invite email mismatch")
	ErrPendingInviteExists     = errors.New("pending invite already exists")
	ErrUserWithEmailRegistered = errors.New("user with email already registered")
	ErrFailedCreateInvite      = errors.New("failed create invite")
	ErrFailedListInvites       = errors.New("failed list invites")
	ErrFailedGetInvite         = errors.New("failed get invite")
	ErrFailedDeleteInvite      = errors.New("failed delete invite")
	ErrFailedAcceptInvite      = errors.New("failed accept invite")
	ErrFailedRejectInvite      = errors.New("failed reject invite")
)

// InviteMailer sends invite emails.
type InviteMailer interface {
	SendInviteEmail(to, inviteLink, pname string) error
}

// Invite is the module for workspace invite CRUD and accept/reject flows.
type Invite struct {
	UserInviteRepository    db.UserInviteRepository
	UserRepository          db.UserRepository
	ConfigRepository        db.ConfigRepository
	WorkspaceRepository     db.WorkspaceRepository
	WorkspaceUserRepository db.WorkspaceUserRepository
	Mailer                  InviteMailer
}

// NewInvite creates an invite module with the given repositories.
func NewInvite(
	invites db.UserInviteRepository,
	users db.UserRepository,
	config db.ConfigRepository,
	workspaces db.WorkspaceRepository,
	workspaceUsers db.WorkspaceUserRepository,
	mailer InviteMailer,
) *Invite {
	return &Invite{
		UserInviteRepository:    invites,
		UserRepository:          users,
		ConfigRepository:        config,
		WorkspaceRepository:     workspaces,
		WorkspaceUserRepository: workspaceUsers,
		Mailer:                  mailer,
	}
}

// CreateInviteRequest is what you pass when creating an invite.
type CreateInviteRequest struct {
	Email string `json:"email" validate:"required,email,min=4,max=255" label:"Email"`
	Role  string `json:"role" validate:"required,oneof=admin regular readonly" label:"Role"`
}

// InviteResponse is an invite shaped for API responses.
type InviteResponse struct {
	Id            db.Id   `json:"id"`
	Email         string  `json:"email"`
	Role          string  `json:"role"`
	Token         string  `json:"token,omitempty"`
	Status        string  `json:"status"`
	InviterUserId db.Id   `json:"inviterUserId"`
	WorkspaceId   db.Id   `json:"workspaceId"`
	WorkspaceName string  `json:"workspaceName,omitempty"`
	ExpiresAt     string  `json:"expiresAt"`
	AcceptedAt    *string `json:"acceptedAt,omitempty"`
	CreatedAt     string  `json:"createdAt"`
	UpdatedAt     string  `json:"updatedAt,omitempty"`
}

// ListInvitesResponse is returned when listing invites.
type ListInvitesResponse struct {
	Invites []*InviteResponse
	Total   int64
}

// CreateInvite creates a workspace invite and sends the invite email.
func (i *Invite) CreateInvite(workspaceId db.Id, req *CreateInviteRequest, inviter *db.User) (*InviteResponse, error) {
	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	existingUser, err := i.UserRepository.GetByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateInvite, err)
	}
	if existingUser != nil {
		return nil, ErrUserWithEmailRegistered
	}

	count, err := i.UserInviteRepository.CountPendingByEmailInWorkspace(workspaceId, req.Email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateInvite, err)
	}
	if count > 0 {
		return nil, ErrPendingInviteExists
	}

	token, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateInvite, err)
	}

	expiresAt := time.Now().UTC().Add(conf.InviteExpiry)
	invite := &db.UserInvite{
		Email:         req.Email,
		Role:          req.Role,
		Token:         token.String(),
		Status:        InviteStatusPending,
		InviterUserId: inviter.Id,
		WorkspaceId:   workspaceId,
		ExpiresAt:     expiresAt,
	}

	err = i.UserInviteRepository.Create(invite)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedCreateInvite, err)
	}

	purl := viper.GetString("app.url")
	pname := viper.GetString("app.name")

	err = i.Mailer.SendInviteEmail(req.Email, fmt.Sprintf("%s/invite/%s", purl, invite.Token), pname)
	if err != nil {
		log.Warn().
			Err(err).
			Str("userId", inviter.Id.String()).
			Str("email", req.Email).
			Msg("Failed to send invite email; invite was created, link available in UI")
	}

	return &InviteResponse{
		Id:            invite.Id,
		Email:         invite.Email,
		Role:          invite.Role,
		Token:         invite.Token,
		Status:        invite.Status,
		InviterUserId: invite.InviterUserId,
		WorkspaceId:   invite.WorkspaceId,
		ExpiresAt:     invite.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt:     invite.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     invite.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}

// ListInvites returns paginated invites for a workspace.
func (i *Invite) ListInvites(workspaceId db.Id, limit, offset int) (*ListInvitesResponse, error) {
	count, err := i.UserInviteRepository.MarkExpiredAsExpired()
	if err == nil && count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired invites marked")
	}

	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	total, err := i.UserInviteRepository.CountByWorkspaceId(workspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListInvites, err)
	}

	invites, err := i.UserInviteRepository.ListByWorkspaceId(workspaceId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListInvites, err)
	}

	list := make([]*InviteResponse, 0, len(invites))
	for _, inv := range invites {
		resp := &InviteResponse{
			Id:            inv.Id,
			Email:         inv.Email,
			Role:          inv.Role,
			Token:         inv.Token,
			Status:        inv.Status,
			InviterUserId: inv.InviterUserId,
			WorkspaceId:   inv.WorkspaceId,
			ExpiresAt:     inv.ExpiresAt.UTC().Format(time.RFC3339),
			CreatedAt:     inv.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:     inv.UpdatedAt.UTC().Format(time.RFC3339),
			AcceptedAt: lo.TernaryF(
				inv.AcceptedAt != nil,
				func() *string { return new(inv.AcceptedAt.UTC().Format(time.RFC3339)) },
				func() *string { return nil },
			),
		}
		list = append(list, resp)
	}

	return &ListInvitesResponse{Invites: list, Total: total}, nil
}

// ListUserInvites returns paginated invites addressed to the current user's email.
func (i *Invite) ListUserInvites(user *db.User, limit, offset int) (*ListInvitesResponse, error) {
	count, err := i.UserInviteRepository.MarkExpiredAsExpired()
	if err == nil && count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired invites marked")
	}

	email := strings.TrimSpace(user.Email)
	total, err := i.UserInviteRepository.CountByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListInvites, err)
	}

	invites, err := i.UserInviteRepository.ListByEmail(email, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedListInvites, err)
	}

	list := make([]*InviteResponse, 0, len(invites))
	for _, inv := range invites {
		resp := &InviteResponse{
			Id:            inv.Id,
			Email:         inv.Email,
			Role:          inv.Role,
			Status:        inv.Status,
			InviterUserId: inv.InviterUserId,
			WorkspaceId:   inv.WorkspaceId,
			ExpiresAt:     inv.ExpiresAt.UTC().Format(time.RFC3339),
			CreatedAt:     inv.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:     inv.UpdatedAt.UTC().Format(time.RFC3339),
			AcceptedAt: lo.TernaryF(
				inv.AcceptedAt != nil,
				func() *string { return new(inv.AcceptedAt.UTC().Format(time.RFC3339)) },
				func() *string { return nil },
			),
		}
		list = append(list, resp)
	}

	return &ListInvitesResponse{Invites: list, Total: total}, nil
}

// GetInvite returns one invite by Id within a workspace.
func (i *Invite) GetInvite(workspaceId, inviteId db.Id) (*InviteResponse, error) {
	count, err := i.UserInviteRepository.MarkExpiredAsExpired()
	if err == nil && count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired invites marked")
	}

	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, ErrWorkspaceNotFound
	}

	invite, err := i.UserInviteRepository.GetById(inviteId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetInvite, err)
	}
	if invite == nil || invite.WorkspaceId != workspaceId {
		return nil, ErrInviteNotFound
	}

	resp := &InviteResponse{
		Id:            invite.Id,
		Email:         invite.Email,
		Role:          invite.Role,
		Status:        invite.Status,
		InviterUserId: invite.InviterUserId,
		WorkspaceId:   invite.WorkspaceId,
		ExpiresAt:     invite.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt:     invite.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     invite.UpdatedAt.UTC().Format(time.RFC3339),
		AcceptedAt: lo.TernaryF(
			invite.AcceptedAt != nil,
			func() *string { return new(invite.AcceptedAt.UTC().Format(time.RFC3339)) },
			func() *string { return nil },
		),
	}

	return resp, nil
}

// DeleteInvite deletes an invite by Id within a workspace.
func (i *Invite) DeleteInvite(workspaceId, inviteId db.Id) error {
	workspace, err := i.WorkspaceRepository.GetById(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return ErrWorkspaceNotFound
	}

	invite, err := i.UserInviteRepository.GetById(inviteId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteInvite, err)
	}
	if invite == nil || invite.WorkspaceId != workspaceId {
		return ErrInviteNotFound
	}

	err = i.UserInviteRepository.Delete(inviteId)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrFailedDeleteInvite, err)
	}

	return nil
}

// GetUserInviteByToken returns invite details for the signup page.
func (i *Invite) GetUserInviteByToken(token string) (*InviteResponse, error) {
	count, err := i.UserInviteRepository.MarkExpiredAsExpired()
	if err == nil && count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired invites marked")
	}

	invite, err := i.UserInviteRepository.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetInvite, err)
	}
	if invite == nil {
		return nil, ErrInviteNotFound
	}
	if invite.Status != InviteStatusPending {
		return nil, ErrInviteNoLongerValid
	}

	workspace, err := i.WorkspaceRepository.GetById(invite.WorkspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetInvite, err)
	}

	resp := &InviteResponse{
		Id:            invite.Id,
		Email:         invite.Email,
		Role:          invite.Role,
		Status:        invite.Status,
		InviterUserId: invite.InviterUserId,
		WorkspaceId:   invite.WorkspaceId,
		WorkspaceName: workspace.Name,
		ExpiresAt:     invite.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt:     invite.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     invite.UpdatedAt.UTC().Format(time.RFC3339),
		AcceptedAt: lo.TernaryF(
			invite.AcceptedAt != nil,
			func() *string { return new(invite.AcceptedAt.UTC().Format(time.RFC3339)) },
			func() *string { return nil },
		),
	}

	return resp, nil
}

// GetAuthenticatedUserInviteByToken returns invite details when the invite belongs to the user.
func (i *Invite) GetAuthenticatedUserInviteByToken(token string, user *db.User) (*InviteResponse, error) {
	count, err := i.UserInviteRepository.MarkExpiredAsExpired()
	if err == nil && count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired invites marked")
	}

	invite, err := i.UserInviteRepository.GetByToken(token)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetInvite, err)
	}
	if invite == nil {
		return nil, ErrInviteNotFound
	}
	if invite.Status != InviteStatusPending {
		return nil, ErrInviteNoLongerValid
	}
	if !strings.EqualFold(user.Email, invite.Email) {
		return nil, ErrInviteNotFound
	}

	workspace, err := i.WorkspaceRepository.GetById(invite.WorkspaceId)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedGetInvite, err)
	}

	resp := &InviteResponse{
		Id:            invite.Id,
		Email:         invite.Email,
		Role:          invite.Role,
		Status:        invite.Status,
		InviterUserId: invite.InviterUserId,
		WorkspaceId:   invite.WorkspaceId,
		WorkspaceName: workspace.Name,
		ExpiresAt:     invite.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt:     invite.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     invite.UpdatedAt.UTC().Format(time.RFC3339),
		AcceptedAt: lo.TernaryF(
			invite.AcceptedAt != nil,
			func() *string { return new(invite.AcceptedAt.UTC().Format(time.RFC3339)) },
			func() *string { return nil },
		),
	}

	return resp, nil
}

// AcceptUserInviteByToken marks an invite as accepted using only the token.
func (i *Invite) AcceptUserInviteByToken(token string, user *db.User) (*InviteResponse, error) {
	count, err := i.UserInviteRepository.MarkExpiredAsExpired()
	if err == nil && count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired invites marked")
	}

	invite, err := i.UserInviteRepository.GetByToken(token)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, ErrInviteNotFound
	}
	if invite.Status == InviteStatusExpired || time.Now().UTC().After(invite.ExpiresAt) {
		return nil, ErrInviteExpired
	}
	if invite.Status != InviteStatusPending {
		return nil, ErrInviteNoLongerValid
	}

	if !strings.EqualFold(user.Email, invite.Email) {
		return nil, ErrInviteEmailMismatch
	}

	membership, err := i.WorkspaceUserRepository.GetByWorkspaceAndUser(invite.WorkspaceId, user.Id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedAcceptInvite, err)
	}
	if membership == nil {
		err = i.WorkspaceUserRepository.Create(&db.WorkspaceUser{
			WorkspaceId: invite.WorkspaceId,
			UserId:      user.Id,
			Role:        invite.Role,
		})
		if err != nil {
			return nil, fmt.Errorf("%w: %v", ErrFailedAcceptInvite, err)
		}
	}

	now := time.Now().UTC()
	err = i.UserInviteRepository.UpdateStatus(invite.Id, InviteStatusAccepted, &now)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedAcceptInvite, err)
	}

	invite.Status = InviteStatusAccepted
	invite.AcceptedAt = &now

	resp := &InviteResponse{
		Id:            invite.Id,
		Email:         invite.Email,
		Role:          invite.Role,
		Status:        invite.Status,
		InviterUserId: invite.InviterUserId,
		WorkspaceId:   invite.WorkspaceId,
		ExpiresAt:     invite.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt:     invite.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     invite.UpdatedAt.UTC().Format(time.RFC3339),
		AcceptedAt:    new(invite.AcceptedAt.UTC().Format(time.RFC3339)),
	}

	return resp, nil
}

// RejectUserInviteByToken marks an invite as rejected using only the token.
func (i *Invite) RejectUserInviteByToken(token string, user *db.User) (*InviteResponse, error) {
	count, err := i.UserInviteRepository.MarkExpiredAsExpired()
	if err == nil && count > 0 {
		log.Info().
			Int64("count", count).
			Msg("Expired invites marked")
	}

	invite, err := i.UserInviteRepository.GetByToken(token)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, ErrInviteNotFound
	}
	if invite.Status == InviteStatusExpired || time.Now().UTC().After(invite.ExpiresAt) {
		return nil, ErrInviteExpired
	}
	if invite.Status != InviteStatusPending {
		return nil, ErrInviteNoLongerValid
	}

	if !strings.EqualFold(user.Email, invite.Email) {
		return nil, ErrInviteEmailMismatch
	}

	err = i.UserInviteRepository.UpdateStatus(invite.Id, InviteStatusRejected, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFailedRejectInvite, err)
	}

	invite.Status = InviteStatusRejected
	invite.AcceptedAt = nil

	return &InviteResponse{
		Id:            invite.Id,
		Email:         invite.Email,
		Role:          invite.Role,
		Status:        invite.Status,
		InviterUserId: invite.InviterUserId,
		WorkspaceId:   invite.WorkspaceId,
		ExpiresAt:     invite.ExpiresAt.UTC().Format(time.RFC3339),
		CreatedAt:     invite.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     invite.UpdatedAt.UTC().Format(time.RFC3339),
	}, nil
}
