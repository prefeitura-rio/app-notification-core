'use client';

import { useState } from 'react';
import { Card } from '@/components/Card';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { useWebSocket } from '@/hooks/useWebSocket';
import { NotificationBadge } from '@/components/NotificationBadge';

export default function Home() {
  const [userId, setUserId] = useState('');
  const [activeUserId, setActiveUserId] = useState('');
  const { notifications, isConnected, clearNotifications } = useWebSocket(activeUserId);

  const handleConnect = (e: React.FormEvent) => {
    e.preventDefault();
    if (userId.trim()) {
      setActiveUserId(userId);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-800 mb-2">Dashboard</h1>
        <p className="text-gray-600">Monitore notificações em tempo real</p>
      </div>

      <Card title="Monitor WebSocket">
        <form onSubmit={handleConnect} className="space-y-4">
          <Input
            label="CPF ou Telefone do Usuário"
            placeholder="Ex: 12345678901 ou 11999999999"
            value={userId}
            onChange={(e) => setUserId(e.target.value)}
          />
          <div className="flex gap-3">
            <Button type="submit">Conectar</Button>
            {activeUserId && (
              <Button
                type="button"
                variant="secondary"
                onClick={() => setActiveUserId('')}
              >
                Desconectar
              </Button>
            )}
          </div>
        </form>

        {activeUserId && (
          <div className="mt-4 pt-4 border-t">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <div className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
                <span className="text-sm font-medium">
                  {isConnected ? 'Conectado' : 'Desconectado'} ({activeUserId})
                </span>
              </div>
              <Button variant="secondary" size="sm" onClick={clearNotifications}>
                Limpar
              </Button>
            </div>
          </div>
        )}
      </Card>

      {activeUserId && (
        <Card title={`Notificações Recebidas (${notifications.length})`}>
          {notifications.length === 0 ? (
            <p className="text-gray-500 text-center py-8">
              Nenhuma notificação recebida ainda...
            </p>
          ) : (
            <div className="space-y-3">
              {notifications.map((notification, index) => (
                <div
                  key={notification.id || index}
                  className="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-start justify-between mb-2">
                    <h3 className="font-semibold text-gray-800">{notification.title}</h3>
                    <div className="flex gap-2">
                      <NotificationBadge type={notification.type} />
                      <NotificationBadge status={notification.status} />
                    </div>
                  </div>
                  <p className="text-gray-600 text-sm mb-2">{notification.message}</p>
                  <div className="flex items-center gap-4 text-xs text-gray-500">
                    <span>{new Date(notification.created_at).toLocaleString('pt-BR')}</span>
                    {notification.broadcast && (
                      <span className="px-2 py-0.5 bg-orange-100 text-orange-800 rounded">
                        Broadcast
                      </span>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      )}
    </div>
  );
}
