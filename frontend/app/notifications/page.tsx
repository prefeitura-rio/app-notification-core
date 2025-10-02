'use client';

import { useState, useEffect } from 'react';
import { Card } from '@/components/Card';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { NotificationBadge } from '@/components/NotificationBadge';
import { api } from '@/lib/api';
import type { Notification } from '@/types';

export default function NotificationsPage() {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchCPF, setSearchCPF] = useState('');
  const [searchPhone, setSearchPhone] = useState('');
  const [searchEmail, setSearchEmail] = useState('');

  useEffect(() => {
    loadNotifications();
  }, []);

  const loadNotifications = async () => {
    try {
      setLoading(true);
      const data = await api.get<Notification[]>('/api/v1/notifications');
      setNotifications(data || []);
    } catch (error) {
      console.error('Failed to load notifications:', error);
    } finally {
      setLoading(false);
    }
  };

  const searchByCPF = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchCPF) return;

    try {
      setLoading(true);
      const data = await api.get<Notification[]>(`/api/v1/notifications/cpf/${searchCPF}`);
      setNotifications(data || []);
    } catch (error) {
      console.error('Failed to search by CPF:', error);
    } finally {
      setLoading(false);
    }
  };

  const searchByPhone = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchPhone) return;

    try {
      setLoading(true);
      const data = await api.get<Notification[]>(`/api/v1/notifications/phone/${searchPhone}`);
      setNotifications(data || []);
    } catch (error) {
      console.error('Failed to search by phone:', error);
    } finally {
      setLoading(false);
    }
  };

  const searchByEmail = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!searchEmail) return;

    try {
      setLoading(true);
      const data = await api.get<Notification[]>(`/api/v1/notifications/email/${searchEmail}`);
      setNotifications(data || []);
    } catch (error) {
      console.error('Failed to search by email:', error);
    } finally {
      setLoading(false);
    }
  };

  const markAsRead = async (id: string) => {
    try {
      await api.post(`/api/v1/notifications/${id}/read`, {});
      loadNotifications();
    } catch (error) {
      console.error('Failed to mark as read:', error);
    }
  };

  const deleteNotification = async (id: string) => {
    if (!confirm('Tem certeza que deseja deletar esta notificação?')) return;

    try {
      await api.delete(`/api/v1/notifications/${id}`);
      loadNotifications();
    } catch (error) {
      console.error('Failed to delete notification:', error);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-800 mb-2">Notificações</h1>
        <p className="text-gray-600">Histórico de todas as notificações enviadas</p>
      </div>

      <Card title="Buscar Notificações">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <form onSubmit={searchByCPF} className="space-y-3">
            <Input
              label="Buscar por CPF"
              placeholder="Digite o CPF"
              value={searchCPF}
              onChange={(e) => setSearchCPF(e.target.value)}
            />
            <Button type="submit" size="sm">Buscar</Button>
          </form>

          <form onSubmit={searchByPhone} className="space-y-3">
            <Input
              label="Buscar por Telefone"
              placeholder="Digite o telefone"
              value={searchPhone}
              onChange={(e) => setSearchPhone(e.target.value)}
            />
            <Button type="submit" size="sm">Buscar</Button>
          </form>

          <form onSubmit={searchByEmail} className="space-y-3">
            <Input
              label="Buscar por Email"
              type="email"
              placeholder="Digite o email"
              value={searchEmail}
              onChange={(e) => setSearchEmail(e.target.value)}
            />
            <Button type="submit" size="sm">Buscar</Button>
          </form>
        </div>

        <div className="mt-4 pt-4 border-t">
          <Button variant="secondary" size="sm" onClick={loadNotifications}>
            Ver Todas
          </Button>
        </div>
      </Card>

      <Card title={`Notificações (${notifications.length})`}>
        {loading ? (
          <p className="text-center text-gray-500 py-8">Carregando...</p>
        ) : notifications.length === 0 ? (
          <p className="text-center text-gray-500 py-8">Nenhuma notificação encontrada</p>
        ) : (
          <div className="space-y-3">
            {notifications.map((notification) => (
              <div
                key={notification.id}
                className="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <div className="flex items-start justify-between mb-3">
                  <div className="flex-1">
                    <h3 className="font-semibold text-gray-800 mb-1">{notification.title}</h3>
                    <p className="text-gray-600 text-sm mb-2">{notification.message}</p>
                  </div>
                  <div className="flex gap-2 ml-4">
                    <NotificationBadge type={notification.type} />
                    <NotificationBadge status={notification.status} />
                  </div>
                </div>

                <div className="flex items-center justify-between text-xs text-gray-500">
                  <div className="flex gap-3">
                    <span>{new Date(notification.created_at).toLocaleString('pt-BR')}</span>
                    {notification.broadcast && (
                      <span className="px-2 py-0.5 bg-orange-100 text-orange-800 rounded">
                        Broadcast
                      </span>
                    )}
                    {notification.user_cpf && (
                      <span className="px-2 py-0.5 bg-blue-100 text-blue-800 rounded">
                        CPF: {notification.user_cpf}
                      </span>
                    )}
                    {notification.user_phone && (
                      <span className="px-2 py-0.5 bg-purple-100 text-purple-800 rounded">
                        Tel: {notification.user_phone}
                      </span>
                    )}
                    {notification.user_email && (
                      <span className="px-2 py-0.5 bg-green-100 text-green-800 rounded">
                        Email: {notification.user_email}
                      </span>
                    )}
                    {notification.is_html && (
                      <span className="px-2 py-0.5 bg-indigo-100 text-indigo-800 rounded">
                        HTML
                      </span>
                    )}
                  </div>

                  <div className="flex gap-2">
                    {notification.status !== 'read' && (
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => markAsRead(notification.id)}
                      >
                        Marcar Lida
                      </Button>
                    )}
                    <Button
                      variant="danger"
                      size="sm"
                      onClick={() => deleteNotification(notification.id)}
                    >
                      Deletar
                    </Button>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </Card>
    </div>
  );
}
