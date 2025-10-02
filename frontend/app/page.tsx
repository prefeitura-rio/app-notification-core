'use client';

import { useState, useEffect } from 'react';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';
import Link from 'next/link';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface NotificationStats {
  total: number;
  by_type: {
    'in-app': number;
    push: number;
    email: number;
    both: number;
    all: number;
  };
  by_status: {
    pending: number;
    sent: number;
    failed: number;
  };
  recent: Array<{
    id: string;
    title: string;
    type: string;
    status: string;
    created_at: string;
  }>;
}

interface QueueStats {
  queue_name: string;
  messages: number;
  consumers: number;
  dlq_messages: number;
}

export default function Home() {
  const [stats, setStats] = useState<NotificationStats | null>(null);
  const [queueStats, setQueueStats] = useState<QueueStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
    const interval = setInterval(loadDashboardData, 30000); // Atualiza a cada 30s
    return () => clearInterval(interval);
  }, []);

  const loadDashboardData = async () => {
    try {
      // Buscar notificações recentes e calcular estatísticas
      const notificationsRes = await fetch(`${API_URL}/api/v1/notifications?limit=100`);
      const notifications = await notificationsRes.json();

      // Calcular estatísticas
      const byType: any = { 'in-app': 0, push: 0, email: 0, both: 0, all: 0 };
      const byStatus: any = { pending: 0, sent: 0, failed: 0 };

      notifications.forEach((n: any) => {
        if (byType[n.type] !== undefined) byType[n.type]++;
        if (byStatus[n.status] !== undefined) byStatus[n.status]++;
      });

      setStats({
        total: notifications.length,
        by_type: byType,
        by_status: byStatus,
        recent: notifications.slice(0, 5),
      });

      // Buscar estatísticas da fila
      try {
        const queueRes = await fetch(`${API_URL}/api/v1/queue/stats`);
        const queue = await queueRes.json();
        setQueueStats(queue);
      } catch (err) {
        console.log('Queue stats not available');
      }

      setLoading(false);
    } catch (error) {
      console.error('Erro ao carregar dados:', error);
      setLoading(false);
    }
  };

  const getTypeColor = (type: string) => {
    const colors: Record<string, string> = {
      'in-app': 'bg-blue-100 text-blue-800',
      'push': 'bg-purple-100 text-purple-800',
      'email': 'bg-green-100 text-green-800',
      'both': 'bg-indigo-100 text-indigo-800',
      'all': 'bg-pink-100 text-pink-800',
    };
    return colors[type] || 'bg-gray-100 text-gray-800';
  };

  const getStatusColor = (status: string) => {
    const colors: Record<string, string> = {
      pending: 'bg-yellow-100 text-yellow-800',
      sent: 'bg-green-100 text-green-800',
      failed: 'bg-red-100 text-red-800',
    };
    return colors[status] || 'bg-gray-100 text-gray-800';
  };

  const getQueueStatusColor = (count: number) => {
    if (count === 0) return 'text-green-600';
    if (count < 100) return 'text-yellow-600';
    if (count < 1000) return 'text-orange-600';
    return 'text-red-600';
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="text-center">
          <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
          <p className="text-gray-600">Carregando dashboard...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-800 mb-2">Dashboard</h1>
          <p className="text-gray-600">Visão geral do sistema de notificações</p>
        </div>
        <Button onClick={loadDashboardData} variant="secondary">
          🔄 Atualizar
        </Button>
      </div>

      {/* Cards de Métricas Principais */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {/* Total de Notificações */}
        <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg p-6 text-white shadow-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-blue-100 text-sm font-medium">Total de Notificações</p>
              <p className="text-4xl font-bold mt-2">{stats?.total || 0}</p>
            </div>
            <div className="text-5xl opacity-20">📬</div>
          </div>
        </div>

        {/* Enviadas */}
        <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-lg p-6 text-white shadow-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-green-100 text-sm font-medium">Enviadas</p>
              <p className="text-4xl font-bold mt-2">{stats?.by_status.sent || 0}</p>
            </div>
            <div className="text-5xl opacity-20">✅</div>
          </div>
        </div>

        {/* Pendentes */}
        <div className="bg-gradient-to-br from-yellow-500 to-yellow-600 rounded-lg p-6 text-white shadow-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-yellow-100 text-sm font-medium">Pendentes</p>
              <p className="text-4xl font-bold mt-2">{stats?.by_status.pending || 0}</p>
            </div>
            <div className="text-5xl opacity-20">⏳</div>
          </div>
        </div>

        {/* Falhas */}
        <div className="bg-gradient-to-br from-red-500 to-red-600 rounded-lg p-6 text-white shadow-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-red-100 text-sm font-medium">Falhas</p>
              <p className="text-4xl font-bold mt-2">{stats?.by_status.failed || 0}</p>
            </div>
            <div className="text-5xl opacity-20">❌</div>
          </div>
        </div>
      </div>

      {/* Estatísticas da Fila RabbitMQ */}
      {queueStats && (
        <Card title="📊 Status da Fila RabbitMQ">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div className="p-4 bg-gray-50 rounded-lg border">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-600">Mensagens Pendentes</span>
                <span className="text-2xl">📬</span>
              </div>
              <p className={`text-3xl font-bold ${getQueueStatusColor(queueStats.messages)}`}>
                {queueStats.messages}
              </p>
            </div>

            <div className="p-4 bg-gray-50 rounded-lg border">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-600">Workers Ativos</span>
                <span className="text-2xl">👷</span>
              </div>
              <p className="text-3xl font-bold text-blue-600">{queueStats.consumers}</p>
            </div>

            <div className="p-4 bg-gray-50 rounded-lg border">
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm font-medium text-gray-600">DLQ (Falhas)</span>
                <span className="text-2xl">💀</span>
              </div>
              <p className={`text-3xl font-bold ${getQueueStatusColor(queueStats.dlq_messages)}`}>
                {queueStats.dlq_messages}
              </p>
            </div>
          </div>

          <div className="mt-4 pt-4 border-t">
            <Link href="/queue">
              <Button variant="secondary" size="sm">
                Ver Dashboard Completo da Fila →
              </Button>
            </Link>
          </div>
        </Card>
      )}

      {/* Notificações por Tipo */}
      <Card title="📱 Notificações por Tipo">
        <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
          {Object.entries(stats?.by_type || {}).map(([type, count]) => (
            <div key={type} className="text-center p-4 bg-gray-50 rounded-lg border">
              <div className={`inline-block px-3 py-1 rounded-full text-sm font-medium mb-2 ${getTypeColor(type)}`}>
                {type}
              </div>
              <p className="text-2xl font-bold text-gray-800">{count}</p>
            </div>
          ))}
        </div>
      </Card>

      {/* Notificações Recentes */}
      <Card title="🕐 Notificações Recentes">
        {stats?.recent && stats.recent.length > 0 ? (
          <div className="space-y-3">
            {stats.recent.map((notification) => (
              <div
                key={notification.id}
                className="flex items-center justify-between p-4 bg-gray-50 rounded-lg border hover:bg-gray-100 transition-colors"
              >
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-800">{notification.title}</h3>
                  <p className="text-sm text-gray-500">
                    {new Date(notification.created_at).toLocaleString('pt-BR')}
                  </p>
                </div>
                <div className="flex gap-2">
                  <span className={`px-3 py-1 rounded-full text-xs font-medium ${getTypeColor(notification.type)}`}>
                    {notification.type}
                  </span>
                  <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(notification.status)}`}>
                    {notification.status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-center text-gray-500 py-8">Nenhuma notificação encontrada</p>
        )}

        <div className="mt-4 pt-4 border-t">
          <Link href="/notifications">
            <Button variant="secondary" size="sm">
              Ver Todas as Notificações →
            </Button>
          </Link>
        </div>
      </Card>

      {/* Quick Actions */}
      <Card title="⚡ Ações Rápidas">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <Link href="/send">
            <div className="p-6 bg-blue-50 border-2 border-blue-200 rounded-lg hover:bg-blue-100 transition-colors cursor-pointer text-center">
              <div className="text-4xl mb-2">📤</div>
              <h3 className="font-semibold text-blue-800">Enviar Notificação</h3>
              <p className="text-xs text-blue-600 mt-1">Criar e enviar nova notificação</p>
            </div>
          </Link>

          <Link href="/groups">
            <div className="p-6 bg-purple-50 border-2 border-purple-200 rounded-lg hover:bg-purple-100 transition-colors cursor-pointer text-center">
              <div className="text-4xl mb-2">👥</div>
              <h3 className="font-semibold text-purple-800">Grupos</h3>
              <p className="text-xs text-purple-600 mt-1">Gerenciar grupos de usuários</p>
            </div>
          </Link>

          <Link href="/test">
            <div className="p-6 bg-green-50 border-2 border-green-200 rounded-lg hover:bg-green-100 transition-colors cursor-pointer text-center">
              <div className="text-4xl mb-2">🧪</div>
              <h3 className="font-semibold text-green-800">Testes</h3>
              <p className="text-xs text-green-600 mt-1">Testar WebSocket e Push</p>
            </div>
          </Link>

          <Link href="/queue">
            <div className="p-6 bg-orange-50 border-2 border-orange-200 rounded-lg hover:bg-orange-100 transition-colors cursor-pointer text-center">
              <div className="text-4xl mb-2">📬</div>
              <h3 className="font-semibold text-orange-800">Filas</h3>
              <p className="text-xs text-orange-600 mt-1">Monitorar fila RabbitMQ</p>
            </div>
          </Link>
        </div>
      </Card>

      {/* System Info */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card title="ℹ️ Informações do Sistema">
          <div className="space-y-3 text-sm">
            <div className="flex justify-between py-2 border-b">
              <span className="font-medium text-gray-700">API URL:</span>
              <code className="text-xs bg-gray-100 px-2 py-1 rounded">{API_URL}</code>
            </div>
            <div className="flex justify-between py-2 border-b">
              <span className="font-medium text-gray-700">WebSocket:</span>
              <code className="text-xs bg-gray-100 px-2 py-1 rounded">
                {API_URL.replace(/^http/, 'ws')}/api/v1/ws
              </code>
            </div>
            <div className="flex justify-between py-2 border-b">
              <span className="font-medium text-gray-700">Documentação:</span>
              <a
                href={`${API_URL}/swagger/index.html`}
                target="_blank"
                rel="noopener noreferrer"
                className="text-blue-600 hover:underline text-xs"
              >
                Swagger UI →
              </a>
            </div>
          </div>
        </Card>

        <Card title="📚 Recursos">
          <div className="space-y-2">
            <a
              href="https://github.com/prefeitura-rio/app-notification-core"
              target="_blank"
              rel="noopener noreferrer"
              className="block p-3 bg-gray-50 rounded-lg border hover:bg-gray-100 transition-colors"
            >
              <div className="flex items-center gap-2">
                <span className="text-xl">📖</span>
                <div>
                  <p className="font-medium text-gray-800">Repositório GitHub</p>
                  <p className="text-xs text-gray-600">Código fonte e documentação</p>
                </div>
              </div>
            </a>

            <Link href="/integrations">
              <div className="block p-3 bg-gray-50 rounded-lg border hover:bg-gray-100 transition-colors cursor-pointer">
                <div className="flex items-center gap-2">
                  <span className="text-xl">🔌</span>
                  <div>
                    <p className="font-medium text-gray-800">Integrações</p>
                    <p className="text-xs text-gray-600">Gerar chaves VAPID e configurações</p>
                  </div>
                </div>
              </div>
            </Link>
          </div>
        </Card>
      </div>
    </div>
  );
}
