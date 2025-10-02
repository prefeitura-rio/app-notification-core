package service

import (
	"errors"
	"github.com/prefeitura-rio/app-notification-core/internal/entity"
	"github.com/prefeitura-rio/app-notification-core/internal/repository"
	"github.com/google/uuid"
)

type GroupService interface {
	CreateGroup(group *entity.Group) error
	GetGroup(id uuid.UUID) (*entity.Group, error)
	ListGroups(limit, offset int) ([]entity.Group, error)
	UpdateGroup(group *entity.Group) error
	DeleteGroup(id uuid.UUID) error
	AddMemberToGroup(member *entity.Member) error
	RemoveMemberFromGroup(memberID uuid.UUID) error
	GetGroupMembers(groupID uuid.UUID) ([]entity.Member, error)
	GetMember(id uuid.UUID) (*entity.Member, error)
	UpdateMember(member *entity.Member) error
}

type groupService struct {
	repo repository.GroupRepository
}

func NewGroupService(repo repository.GroupRepository) GroupService {
	return &groupService{repo: repo}
}

func (s *groupService) CreateGroup(group *entity.Group) error {
	if group.Name == "" {
		return errors.New("group name is required")
	}
	return s.repo.Create(group)
}

func (s *groupService) GetGroup(id uuid.UUID) (*entity.Group, error) {
	return s.repo.FindByID(id)
}

func (s *groupService) ListGroups(limit, offset int) ([]entity.Group, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.FindAll(limit, offset)
}

func (s *groupService) UpdateGroup(group *entity.Group) error {
	if group.Name == "" {
		return errors.New("group name is required")
	}
	return s.repo.Update(group)
}

func (s *groupService) DeleteGroup(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *groupService) AddMemberToGroup(member *entity.Member) error {
	if member.CPF == "" && member.Phone == "" {
		return errors.New("either CPF or phone is required")
	}
	return s.repo.AddMember(member)
}

func (s *groupService) RemoveMemberFromGroup(memberID uuid.UUID) error {
	return s.repo.RemoveMember(memberID)
}

func (s *groupService) GetGroupMembers(groupID uuid.UUID) ([]entity.Member, error) {
	return s.repo.FindMembers(groupID)
}

func (s *groupService) GetMember(id uuid.UUID) (*entity.Member, error) {
	return s.repo.FindMemberByID(id)
}

func (s *groupService) UpdateMember(member *entity.Member) error {
	if member.CPF == "" && member.Phone == "" && member.Email == "" {
		return errors.New("at least one of CPF, phone, or email is required")
	}
	return s.repo.UpdateMember(member)
}
