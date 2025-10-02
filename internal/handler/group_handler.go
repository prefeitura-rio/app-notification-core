package handler

import (
	"net/http"
	"strconv"

	"github.com/prefeitura-rio/app-notification-core/internal/entity"
	"github.com/prefeitura-rio/app-notification-core/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GroupHandler struct {
	service service.GroupService
}

func NewGroupHandler(service service.GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

// Create godoc
// @Summary Criar novo grupo
// @Description Cria um novo grupo de usuários
// @Tags groups
// @Accept json
// @Produce json
// @Param group body entity.Group true "Dados do grupo"
// @Success 201 {object} entity.Group
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /groups [post]
func (h *GroupHandler) Create(c *gin.Context) {
	var group entity.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateGroup(&group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, group)
}

// Get godoc
// @Summary Buscar grupo por ID
// @Description Retorna um grupo específico pelo ID
// @Tags groups
// @Produce json
// @Param id path string true "ID do grupo"
// @Success 200 {object} entity.Group
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /groups/{id} [get]
func (h *GroupHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	group, err := h.service.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}

	c.JSON(http.StatusOK, group)
}

// List godoc
// @Summary Listar grupos
// @Description Retorna lista de grupos com paginação
// @Tags groups
// @Produce json
// @Param limit query int false "Limite de resultados" default(20)
// @Param offset query int false "Offset para paginação" default(0)
// @Success 200 {array} entity.Group
// @Failure 500 {object} map[string]string
// @Router /groups [get]
func (h *GroupHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	groups, err := h.service.ListGroups(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, groups)
}

// Update godoc
// @Summary Atualizar grupo
// @Description Atualiza dados de um grupo existente
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "ID do grupo"
// @Param group body entity.Group true "Dados atualizados do grupo"
// @Success 200 {object} entity.Group
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /groups/{id} [put]
func (h *GroupHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var group entity.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group.ID = id
	if err := h.service.UpdateGroup(&group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, group)
}

// Delete godoc
// @Summary Deletar grupo
// @Description Remove um grupo pelo ID
// @Tags groups
// @Param id path string true "ID do grupo"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /groups/{id} [delete]
func (h *GroupHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	if err := h.service.DeleteGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddMember godoc
// @Summary Adicionar membro ao grupo
// @Description Adiciona um novo membro a um grupo específico
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "ID do grupo"
// @Param member body entity.Member true "Dados do membro"
// @Success 201 {object} entity.Member
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /groups/{id}/members [post]
func (h *GroupHandler) AddMember(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	var member entity.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member.GroupID = groupID
	if err := h.service.AddMemberToGroup(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, member)
}

// GetMembers godoc
// @Summary Listar membros do grupo
// @Description Retorna todos os membros de um grupo específico
// @Tags groups
// @Produce json
// @Param id path string true "ID do grupo"
// @Success 200 {array} entity.Member
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /groups/{id}/members [get]
func (h *GroupHandler) GetMembers(c *gin.Context) {
	groupID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	members, err := h.service.GetGroupMembers(groupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, members)
}

// RemoveMember godoc
// @Summary Remover membro do grupo
// @Description Remove um membro específico de um grupo
// @Tags groups
// @Param id path string true "ID do grupo"
// @Param memberId path string true "ID do membro"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /groups/{id}/members/{memberId} [delete]
func (h *GroupHandler) RemoveMember(c *gin.Context) {
	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member ID"})
		return
	}

	if err := h.service.RemoveMemberFromGroup(memberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GetMember godoc
// @Summary Buscar membro específico
// @Description Retorna dados de um membro específico do grupo
// @Tags groups
// @Produce json
// @Param id path string true "ID do grupo"
// @Param memberId path string true "ID do membro"
// @Success 200 {object} entity.Member
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /groups/{id}/members/{memberId} [get]
func (h *GroupHandler) GetMember(c *gin.Context) {
	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member ID"})
		return
	}

	member, err := h.service.GetMember(memberID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "member not found"})
		return
	}

	c.JSON(http.StatusOK, member)
}

// UpdateMember godoc
// @Summary Atualizar membro
// @Description Atualiza dados de um membro do grupo
// @Tags groups
// @Accept json
// @Produce json
// @Param id path string true "ID do grupo"
// @Param memberId path string true "ID do membro"
// @Param member body entity.Member true "Dados atualizados do membro"
// @Success 200 {object} entity.Member
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /groups/{id}/members/{memberId} [put]
func (h *GroupHandler) UpdateMember(c *gin.Context) {
	memberID, err := uuid.Parse(c.Param("memberId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid member ID"})
		return
	}

	var member entity.Member
	if err := c.ShouldBindJSON(&member); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	member.ID = memberID
	if err := h.service.UpdateMember(&member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, member)
}
