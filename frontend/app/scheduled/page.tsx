'use client';

import { useState, useEffect } from 'react';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import { api } from '@/lib/api';

interface ScheduledNotification {
  id: string;
  title: string;
  message: string;
  type: string;
  status: string;
  scheduled_for: string;
  user_cpf?: string;
  user_phone?: string;
  user_email?: string;
  broadcast: boolean;
  group_id?: string;
  created_at: string;
}

export default function ScheduledPage() {
  const [notifications, setNotifications] = useState<ScheduledNotification[]>([]);
  const [loading, setLoading] = useState(true);
  const [limit] = useState(50);
  const [offset] = useState(0);

  useEffect(() => {
    loadScheduledNotifications();
    // Auto-refresh a cada 30 segundos
    const interval = setInterval(loadScheduledNotifications, 30000);
    return () => clearInterval(interval);
  }, []);

  const loadScheduledNotifications = async () => {
    try {
      const data = await api.get<ScheduledNotification[]>(
        `/api/v1/scheduled-notifications?limit=${limit}&offset=${offset}`
      );
      setNotifications(data || []);
    } catch (error) {
      console.error('Failed to load scheduled notifications:', error);
    } finally {
      setLoading(false);
    }
  };

  const cancelScheduled = async (id: string) => {
    if (!confirm('Tem certeza que deseja cancelar esta notificação agendada?')) {
      return;
    }

    try {
      await api.post(`/api/v1/scheduled-notifications/${id}/cancel`, {});
      alert('Notificação cancelada com sucesso!');
      loadScheduledNotifications();
    } catch (error) {
      console.error('Failed to cancel scheduled notification:', error);
      alert('Erro ao cancelar notificação');
    }
  };

  const formatDateTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const getRecipientInfo = (notification: ScheduledNotification) => {
    if (notification.broadcast) {
      return '📢 Broadcast (Todos)';
    }
    if (notification.group_id) {
      return `👥 Grupo: ${notification.group_id}`;
    }
    if (notification.user_cpf) {
      return `👤 CPF: ${notification.user_cpf}`;
    }
    if (notification.user_phone) {
      return `📱 Tel: ${notification.user_phone}`;
    }
    if (notification.user_email) {
      return `✉️ Email: ${notification.user_email}`;
    }
    return 'Destinatário não especificado';
  };

  const getTypeLabel = (type: string) => {
    const types: Record<string, string> = {
      'in-app': '📱 In-App',
      'push': '🔔 Push',
      'email': '✉️ Email',
      'both': '📱🔔 In-App + Push',
      'all': '📱🔔✉️ Todos',
    };
    return types[type] || type;
  };

  const getTimeUntil = (scheduledFor: string) => {
    const now = new Date();
    const scheduled = new Date(scheduledFor);
    const diff = scheduled.getTime() - now.getTime();

    if (diff < 0) {
      return '⏰ Processando...';
    }

    const hours = Math.floor(diff / (1000 * 60 * 60));
    const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));

    if (hours > 24) {
      const days = Math.floor(hours / 24);
      return `🕐 Em ${days} dia(s)`;
    }

    if (hours > 0) {
      return `🕐 Em ${hours}h ${minutes}m`;
    }

    return `🕐 Em ${minutes} minuto(s)`;
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold text-gray-800 mb-2">Notificações Agendadas</h1>
          <p className="text-gray-600">Carregando...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-800 mb-2">Notificações Agendadas</h1>
          <p className="text-gray-600">
            Gerencie notificações programadas para envio futuro
          </p>
        </div>
        <Button onClick={loadScheduledNotifications}>
          🔄 Atualizar
        </Button>
      </div>

      {notifications.length === 0 ? (
        <Card>
          <div className="text-center py-12">
            <div className="text-6xl mb-4">📅</div>
            <h3 className="text-xl font-semibold text-gray-700 mb-2">
              Nenhuma notificação agendada
            </h3>
            <p className="text-gray-600">
              As notificações agendadas aparecerão aqui quando você criar uma.
            </p>
          </div>
        </Card>
      ) : (
        <div className="space-y-4">
          <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
            <p className="text-sm text-purple-800">
              📊 <strong>{notifications.length}</strong> notificação(ões) agendada(s)
            </p>
          </div>

          {notifications.map((notification) => (
            <Card key={notification.id}>
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <div className="flex items-center gap-2 mb-2">
                    <h3 className="text-lg font-semibold text-gray-800">
                      {notification.title}
                    </h3>
                    <span className="text-xs bg-purple-100 text-purple-800 px-2 py-1 rounded">
                      {getTypeLabel(notification.type)}
                    </span>
                  </div>

                  <p className="text-gray-600 mb-3">{notification.message}</p>

                  <div className="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm">
                    <div>
                      <span className="text-gray-500">Destinatário:</span>{' '}
                      <span className="font-medium text-gray-700">
                        {getRecipientInfo(notification)}
                      </span>
                    </div>

                    <div>
                      <span className="text-gray-500">Envio programado:</span>{' '}
                      <span className="font-medium text-purple-700">
                        {formatDateTime(notification.scheduled_for)}
                      </span>
                    </div>

                    <div>
                      <span className="text-gray-500">Tempo restante:</span>{' '}
                      <span className="font-medium text-purple-700">
                        {getTimeUntil(notification.scheduled_for)}
                      </span>
                    </div>

                    <div>
                      <span className="text-gray-500">Criada em:</span>{' '}
                      <span className="text-gray-700">
                        {formatDateTime(notification.created_at)}
                      </span>
                    </div>
                  </div>
                </div>

                <Button
                  variant="danger"
                  onClick={() => cancelScheduled(notification.id)}
                  className="ml-4"
                >
                  ❌ Cancelar
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
