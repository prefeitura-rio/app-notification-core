package repository

import (
	"github.com/prefeitura-rio/app-notification-core/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRepository interface {
	Create(group *entity.Group) error
	FindByID(id uuid.UUID) (*entity.Group, error)
	FindAll(limit, offset int) ([]entity.Group, error)
	Update(group *entity.Group) error
	Delete(id uuid.UUID) error
	AddMember(member *entity.Member) error
	RemoveMember(id uuid.UUID) error
	FindMembers(groupID uuid.UUID) ([]entity.Member, error)
	FindMemberByID(id uuid.UUID) (*entity.Member, error)
	UpdateMember(member *entity.Member) error
}

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) Create(group *entity.Group) error {
	return r.db.Create(group).Error
}

func (r *groupRepository) FindByID(id uuid.UUID) (*entity.Group, error) {
	var group entity.Group
	err := r.db.Preload("Members").First(&group, "id = ?", id).Error
	return &group, err
}

func (r *groupRepository) FindAll(limit, offset int) ([]entity.Group, error) {
	var groups []entity.Group
	err := r.db.Preload("Members").Limit(limit).Offset(offset).Find(&groups).Error
	return groups, err
}

func (r *groupRepository) Update(group *entity.Group) error {
	return r.db.Save(group).Error
}

func (r *groupRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Group{}, "id = ?", id).Error
}

func (r *groupRepository) AddMember(member *entity.Member) error {
	return r.db.Create(member).Error
}

func (r *groupRepository) RemoveMember(id uuid.UUID) error {
	return r.db.Delete(&entity.Member{}, "id = ?", id).Error
}

func (r *groupRepository) FindMembers(groupID uuid.UUID) ([]entity.Member, error) {
	var members []entity.Member
	err := r.db.Where("group_id = ?", groupID).Find(&members).Error
	return members, err
}

func (r *groupRepository) FindMemberByID(id uuid.UUID) (*entity.Member, error) {
	var member entity.Member
	err := r.db.First(&member, "id = ?", id).Error
	return &member, err
}

func (r *groupRepository) UpdateMember(member *entity.Member) error {
	return r.db.Save(member).Error
}
