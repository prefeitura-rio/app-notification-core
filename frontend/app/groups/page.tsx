'use client';

import { useState, useEffect } from 'react';
import { Card } from '@/components/Card';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { Modal } from '@/components/Modal';
import { api } from '@/lib/api';
import type { Group, Member } from '@/types';

export default function GroupsPage() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [loading, setLoading] = useState(true);
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [selectedGroup, setSelectedGroup] = useState<Group | null>(null);
  const [editingMember, setEditingMember] = useState<Member | null>(null);
  const [editingGroup, setEditingGroup] = useState<Group | null>(null);

  const [formData, setFormData] = useState({
    name: '',
    description: '',
  });

  const [memberData, setMemberData] = useState({
    cpf: '',
    phone: '',
    email: '',
    name: '',
  });

  const [editMemberData, setEditMemberData] = useState({
    cpf: '',
    phone: '',
    email: '',
    name: '',
  });

  useEffect(() => {
    loadGroups();
  }, []);

  const loadGroups = async () => {
    try {
      setLoading(true);
      const data = await api.get<Group[]>('/api/v1/groups');
      setGroups(data || []);
    } catch (error) {
      console.error('Failed to load groups:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreateGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/api/v1/groups', formData);
      setFormData({ name: '', description: '' });
      setShowCreateForm(false);
      loadGroups();
    } catch (error) {
      console.error('Failed to create group:', error);
      alert('Erro ao criar grupo');
    }
  };

  const handleDeleteGroup = async (id: string) => {
    if (!confirm('Tem certeza que deseja deletar este grupo?')) return;

    try {
      await api.delete(`/api/v1/groups/${id}`);
      loadGroups();
    } catch (error) {
      console.error('Failed to delete group:', error);
      alert('Erro ao deletar grupo');
    }
  };

  const handleAddMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedGroup) return;

    try {
      await api.post(`/api/v1/groups/${selectedGroup.id}/members`, memberData);
      setMemberData({ cpf: '', phone: '', email: '', name: '' });
      loadGroupDetails(selectedGroup.id);
    } catch (error) {
      console.error('Failed to add member:', error);
      alert('Erro ao adicionar membro');
    }
  };

  const loadGroupDetails = async (id: string) => {
    try {
      const group = await api.get<Group>(`/api/v1/groups/${id}`);
      setSelectedGroup(group);
    } catch (error) {
      console.error('Failed to load group details:', error);
    }
  };

  const handleRemoveMember = async (groupId: string, memberId: string) => {
    if (!confirm('Tem certeza que deseja remover este membro?')) return;

    try {
      await api.delete(`/api/v1/groups/${groupId}/members/${memberId}`);
      loadGroupDetails(groupId);
    } catch (error) {
      console.error('Failed to remove member:', error);
      alert('Erro ao remover membro');
    }
  };

  const handleEditMember = (member: Member) => {
    setEditingMember(member);
    setEditMemberData({
      cpf: member.cpf || '',
      phone: member.phone || '',
      email: member.email || '',
      name: member.name || '',
    });
  };

  const handleUpdateMember = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingMember || !selectedGroup) return;

    try {
      await api.put(`/api/v1/groups/${selectedGroup.id}/members/${editingMember.id}`, editMemberData);
      setEditingMember(null);
      setEditMemberData({ cpf: '', phone: '', email: '', name: '' });
      loadGroupDetails(selectedGroup.id);
    } catch (error) {
      console.error('Failed to update member:', error);
      alert('Erro ao atualizar membro');
    }
  };

  const handleEditGroup = (group: Group) => {
    setEditingGroup(group);
    setFormData({
      name: group.name,
      description: group.description,
    });
  };

  const handleUpdateGroup = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingGroup) return;

    try {
      await api.put(`/api/v1/groups/${editingGroup.id}`, formData);
      setEditingGroup(null);
      setFormData({ name: '', description: '' });
      loadGroups();
    } catch (error) {
      console.error('Failed to update group:', error);
      alert('Erro ao atualizar grupo');
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-800 mb-2">Grupos</h1>
          <p className="text-gray-600">Gerencie grupos de usuários</p>
        </div>
        <Button onClick={() => setShowCreateForm(!showCreateForm)}>
          {showCreateForm ? 'Cancelar' : 'Novo Grupo'}
        </Button>
      </div>

      {showCreateForm && (
        <Card title="Criar Novo Grupo">
          <form onSubmit={handleCreateGroup} className="space-y-4">
            <Input
              label="Nome do Grupo"
              placeholder="Ex: Usuários Premium"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              required
            />
            <Input
              label="Descrição"
              placeholder="Descrição do grupo"
              value={formData.description}
              onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            />
            <Button type="submit">Criar Grupo</Button>
          </form>
        </Card>
      )}

      {loading ? (
        <Card>
          <p className="text-center text-gray-500">Carregando...</p>
        </Card>
      ) : groups.length === 0 ? (
        <Card>
          <p className="text-center text-gray-500">Nenhum grupo criado ainda</p>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {groups.map((group) => (
            <Card key={group.id}>
              <div className="flex justify-between items-start mb-3">
                <div>
                  <h3 className="font-bold text-lg text-gray-800">{group.name}</h3>
                  <p className="text-sm text-gray-600">{group.description}</p>
                </div>
                <div className="flex gap-2">
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => handleEditGroup(group)}
                  >
                    Editar
                  </Button>
                  <Button
                    variant="danger"
                    size="sm"
                    onClick={() => handleDeleteGroup(group.id)}
                  >
                    Deletar
                  </Button>
                </div>
              </div>
              <Button
                variant="secondary"
                size="sm"
                onClick={() => loadGroupDetails(group.id)}
              >
                Ver Membros ({group.members?.length || 0})
              </Button>
            </Card>
          ))}
        </div>
      )}

      {selectedGroup && (
        <Card title={`Membros de "${selectedGroup.name}"`}>
          <div className="space-y-4">
            <form onSubmit={handleAddMember} className="grid grid-cols-2 gap-3">
              <Input
                placeholder="CPF"
                value={memberData.cpf}
                onChange={(e) => setMemberData({ ...memberData, cpf: e.target.value })}
              />
              <Input
                placeholder="Telefone"
                value={memberData.phone}
                onChange={(e) => setMemberData({ ...memberData, phone: e.target.value })}
              />
              <Input
                placeholder="Email"
                value={memberData.email}
                onChange={(e) => setMemberData({ ...memberData, email: e.target.value })}
              />
              <Input
                placeholder="Nome"
                value={memberData.name}
                onChange={(e) => setMemberData({ ...memberData, name: e.target.value })}
              />
              <div className="col-span-2">
                <Button type="submit" size="sm">Adicionar Membro</Button>
              </div>
            </form>

            {selectedGroup.members && selectedGroup.members.length > 0 ? (
              <div className="space-y-2">
                {selectedGroup.members.map((member) => (
                  <div key={member.id} className="flex justify-between items-center p-3 bg-gray-50 rounded">
                    <div className="flex-1">
                      <p className="font-medium">{member.name || 'Sem nome'}</p>
                      <div className="text-sm text-gray-600 space-y-1">
                        {member.cpf && <p>CPF: {member.cpf}</p>}
                        {member.phone && <p>Tel: {member.phone}</p>}
                        {member.email && <p>Email: {member.email}</p>}
                      </div>
                    </div>
                    <div className="flex gap-2">
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => handleEditMember(member)}
                      >
                        Editar
                      </Button>
                      <Button
                        variant="danger"
                        size="sm"
                        onClick={() => handleRemoveMember(selectedGroup.id, member.id)}
                      >
                        Remover
                      </Button>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-center text-gray-500">Nenhum membro neste grupo</p>
            )}
          </div>
        </Card>
      )}

      {/* Modal de Edição de Grupo */}
      <Modal
        isOpen={editingGroup !== null}
        onClose={() => {
          setEditingGroup(null);
          setFormData({ name: '', description: '' });
        }}
        title="Editar Grupo"
      >
        <form onSubmit={handleUpdateGroup} className="space-y-4">
          <Input
            label="Nome do Grupo"
            placeholder="Ex: Usuários Premium"
            value={formData.name}
            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            required
          />
          <Input
            label="Descrição"
            placeholder="Descrição do grupo"
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
          />
          <div className="flex gap-2">
            <Button type="submit">Salvar</Button>
            <Button
              type="button"
              variant="secondary"
              onClick={() => {
                setEditingGroup(null);
                setFormData({ name: '', description: '' });
              }}
            >
              Cancelar
            </Button>
          </div>
        </form>
      </Modal>

      {/* Modal de Edição de Membro */}
      <Modal
        isOpen={editingMember !== null}
        onClose={() => {
          setEditingMember(null);
          setEditMemberData({ cpf: '', phone: '', email: '', name: '' });
        }}
        title="Editar Membro"
      >
        <form onSubmit={handleUpdateMember} className="space-y-4">
          <Input
            label="Nome"
            placeholder="Nome do membro"
            value={editMemberData.name}
            onChange={(e) => setEditMemberData({ ...editMemberData, name: e.target.value })}
          />
          <Input
            label="CPF"
            placeholder="12345678901"
            value={editMemberData.cpf}
            onChange={(e) => setEditMemberData({ ...editMemberData, cpf: e.target.value })}
          />
          <Input
            label="Telefone"
            placeholder="11999999999"
            value={editMemberData.phone}
            onChange={(e) => setEditMemberData({ ...editMemberData, phone: e.target.value })}
          />
          <Input
            label="Email"
            type="email"
            placeholder="usuario@exemplo.com"
            value={editMemberData.email}
            onChange={(e) => setEditMemberData({ ...editMemberData, email: e.target.value })}
          />
          <div className="flex gap-2">
            <Button type="submit">Salvar</Button>
            <Button
              type="button"
              variant="secondary"
              onClick={() => {
                setEditingMember(null);
                setEditMemberData({ cpf: '', phone: '', email: '', name: '' });
              }}
            >
              Cancelar
            </Button>
          </div>
        </form>
      </Modal>
    </div>
  );
}
