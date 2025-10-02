'use client';

import { useState, useEffect, useCallback } from 'react';
import { Card } from '@/components/Card';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { NotificationBadge } from '@/components/NotificationBadge';

interface Notification {
  id: string;
  title: string;
  message: string;
  type: string;
  status: string;
  created_at: string;
  is_html: boolean;
}

export default function TestPage() {
  const [userIdentifier, setUserIdentifier] = useState('');
  const [identifierType, setIdentifierType] = useState<'cpf' | 'phone' | 'email'>('cpf');
  const [isConnected, setIsConnected] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [pushEnabled, setPushEnabled] = useState(false);
  const [logs, setLogs] = useState<string[]>([]);

  const addLog = useCallback((message: string) => {
    const timestamp = new Date().toLocaleTimeString('pt-BR');
    setLogs((prev) => [`[${timestamp}] ${message}`, ...prev].slice(0, 50));
  }, []);

  const connectWebSocket = useCallback(() => {
    if (!userIdentifier) {
      alert('Digite um identificador de usuário');
      return;
    }

    if (ws) {
      ws.close();
    }

    addLog(`🔄 Conectando ao WebSocket com ${identifierType}: ${userIdentifier}...`);

    const wsUrl = `ws://localhost:8080/api/v1/ws?user_id=${encodeURIComponent(userIdentifier)}`;
    const websocket = new WebSocket(wsUrl);

    websocket.onopen = () => {
      addLog('✅ Conectado ao WebSocket');
      setIsConnected(true);
    };

    websocket.onmessage = (event) => {
      try {
        const notification = JSON.parse(event.data);
        addLog(`📨 Notificação recebida: ${notification.title}`);
        setNotifications((prev) => [notification, ...prev]);

        // Tentar mostrar notificação nativa do navegador
        if ('Notification' in window && Notification.permission === 'granted') {
          new Notification(notification.title, {
            body: notification.message,
            icon: '/favicon.ico',
          });
        }
      } catch (error) {
        addLog(`❌ Erro ao processar mensagem: ${error}`);
      }
    };

    websocket.onerror = (error) => {
      addLog(`❌ Erro no WebSocket: ${error}`);
      setIsConnected(false);
    };

    websocket.onclose = () => {
      addLog('🔌 WebSocket desconectado');
      setIsConnected(false);
    };

    setWs(websocket);
  }, [userIdentifier, identifierType, ws, addLog]);

  const disconnect = () => {
    if (ws) {
      ws.close();
      setWs(null);
      setIsConnected(false);
      addLog('🔌 Desconectado manualmente');
    }
  };

  const requestNotificationPermission = async () => {
    if (!('Notification' in window)) {
      addLog('❌ Este navegador não suporta notificações');
      return;
    }

    const permission = await Notification.requestPermission();
    if (permission === 'granted') {
      addLog('✅ Permissão de notificação concedida');
    } else {
      addLog('❌ Permissão de notificação negada');
    }
  };

  const registerPush = async () => {
    if (!userIdentifier) {
      alert('Digite um identificador de usuário');
      return;
    }

    if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
      addLog('❌ Push notifications não suportadas neste navegador');
      return;
    }

    try {
      addLog('🔄 Solicitando permissão...');
      const permission = await Notification.requestPermission();
      if (permission !== 'granted') {
        addLog('❌ Permissão negada');
        return;
      }

      addLog('🔄 Registrando Service Worker...');
      const registration = await navigator.serviceWorker.register('/sw.js');
      await navigator.serviceWorker.ready;
      addLog('✅ Service Worker registrado');

      // VAPID key (configure no .env.local)
      const vapidPublicKey = process.env.NEXT_PUBLIC_VAPID_PUBLIC_KEY;

      if (!vapidPublicKey || vapidPublicKey === 'your_vapid_public_key_here') {
        addLog('❌ VAPID_PUBLIC_KEY não configurada. Configure NEXT_PUBLIC_VAPID_PUBLIC_KEY no .env.local');
        return;
      }

      addLog('🔄 Criando subscrição push...');
      const subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
      });

      // Enviar para o servidor
      const subscriptionData = {
        endpoint: subscription.endpoint,
        p256dh: arrayBufferToBase64(subscription.getKey('p256dh')!),
        auth: arrayBufferToBase64(subscription.getKey('auth')!),
      };

      if (identifierType === 'cpf') {
        Object.assign(subscriptionData, { user_cpf: userIdentifier });
      } else if (identifierType === 'phone') {
        Object.assign(subscriptionData, { user_phone: userIdentifier });
      }

      addLog('🔄 Enviando subscrição para o servidor...');
      const response = await fetch('http://localhost:8080/api/v1/subscriptions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(subscriptionData),
      });

      if (response.ok) {
        addLog('✅ Push notifications registradas com sucesso');
        setPushEnabled(true);
      } else {
        addLog(`❌ Erro ao registrar: ${response.status}`);
      }
    } catch (error) {
      addLog(`❌ Erro ao registrar push: ${error}`);
    }
  };

  const clearNotifications = () => {
    setNotifications([]);
    addLog('🗑️ Notificações limpas');
  };

  const clearLogs = () => {
    setLogs([]);
  };

  useEffect(() => {
    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, [ws]);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-800 mb-2">Modo de Teste</h1>
        <p className="text-gray-600">
          Teste notificações in-app (WebSocket) e push notifications
        </p>
      </div>

      {/* Configuração */}
      <Card title="Configuração do Teste">
        <div className="space-y-4">
          <div className="flex gap-3">
            <div className="flex gap-2">
              <label className="flex items-center gap-1">
                <input
                  type="radio"
                  checked={identifierType === 'cpf'}
                  onChange={() => setIdentifierType('cpf')}
                />
                <span className="text-sm">CPF</span>
              </label>
              <label className="flex items-center gap-1">
                <input
                  type="radio"
                  checked={identifierType === 'phone'}
                  onChange={() => setIdentifierType('phone')}
                />
                <span className="text-sm">Telefone</span>
              </label>
              <label className="flex items-center gap-1">
                <input
                  type="radio"
                  checked={identifierType === 'email'}
                  onChange={() => setIdentifierType('email')}
                />
                <span className="text-sm">Email</span>
              </label>
            </div>
          </div>

          <Input
            label={`Identificador (${identifierType.toUpperCase()})`}
            placeholder={
              identifierType === 'cpf'
                ? '12345678901'
                : identifierType === 'phone'
                ? '11999999999'
                : 'teste@exemplo.com'
            }
            value={userIdentifier}
            onChange={(e) => setUserIdentifier(e.target.value)}
          />

          <div className="flex gap-2">
            {!isConnected ? (
              <Button onClick={connectWebSocket}>Conectar WebSocket</Button>
            ) : (
              <Button variant="danger" onClick={disconnect}>
                Desconectar
              </Button>
            )}

            <Button variant="secondary" onClick={requestNotificationPermission}>
              Solicitar Permissão
            </Button>

            <Button
              variant={pushEnabled ? 'secondary' : 'primary'}
              onClick={registerPush}
              disabled={pushEnabled}
            >
              {pushEnabled ? '✓ Push Habilitado' : 'Habilitar Push'}
            </Button>
          </div>

          {/* Status */}
          <div className="flex items-center gap-4 pt-4 border-t">
            <div className="flex items-center gap-2">
              <span
                className={`w-3 h-3 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`}
              />
              <span className="text-sm font-medium">
                WebSocket: {isConnected ? 'Conectado' : 'Desconectado'}
              </span>
            </div>
            <div className="flex items-center gap-2">
              <span className={`w-3 h-3 rounded-full ${pushEnabled ? 'bg-green-500' : 'bg-gray-300'}`} />
              <span className="text-sm font-medium">
                Push: {pushEnabled ? 'Habilitado' : 'Desabilitado'}
              </span>
            </div>
          </div>
        </div>
      </Card>

      {/* Notificações Recebidas */}
      <Card title={`Notificações Recebidas (${notifications.length})`}>
        <div className="flex justify-end mb-4">
          <Button variant="secondary" size="sm" onClick={clearNotifications}>
            Limpar Notificações
          </Button>
        </div>

        <div className="space-y-3">
          {notifications.length === 0 ? (
            <p className="text-center text-gray-500 py-8">
              Aguardando notificações... Envie uma notificação pela página "Enviar" para testar.
            </p>
          ) : (
            notifications.map((notification) => (
              <div
                key={notification.id}
                className="p-4 bg-blue-50 border border-blue-200 rounded-lg"
              >
                <div className="flex justify-between items-start mb-2">
                  <h3 className="font-bold text-gray-800">{notification.title}</h3>
                  <NotificationBadge type={notification.type as any} />
                </div>

                {notification.is_html ? (
                  <div
                    className="text-gray-700 text-sm"
                    dangerouslySetInnerHTML={{ __html: notification.message }}
                  />
                ) : (
                  <p className="text-gray-700 text-sm">{notification.message}</p>
                )}

                <p className="text-xs text-gray-500 mt-2">
                  {new Date(notification.created_at).toLocaleString('pt-BR')}
                </p>
              </div>
            ))
          )}
        </div>
      </Card>

      {/* Logs */}
      <Card title="Logs de Debug">
        <div className="flex justify-end mb-4">
          <Button variant="secondary" size="sm" onClick={clearLogs}>
            Limpar Logs
          </Button>
        </div>

        <div className="bg-black text-green-400 p-4 rounded-lg font-mono text-xs max-h-96 overflow-y-auto">
          {logs.length === 0 ? (
            <p className="text-gray-500">Nenhum log ainda...</p>
          ) : (
            logs.map((log, index) => (
              <div key={index} className="mb-1">
                {log}
              </div>
            ))
          )}
        </div>
      </Card>

      {/* Instruções */}
      <Card title="📖 Como Testar">
        <div className="prose prose-sm">
          <ol className="space-y-2">
            <li>
              <strong>Digite um identificador:</strong> CPF, telefone ou email de teste (ex:
              12345678901)
            </li>
            <li>
              <strong>Conecte ao WebSocket:</strong> Clique em "Conectar WebSocket"
            </li>
            <li>
              <strong>Solicite permissão:</strong> Clique em "Solicitar Permissão" para notificações
              do navegador
            </li>
            <li>
              <strong>(Opcional) Habilite Push:</strong> Clique em "Habilitar Push" para receber
              push notifications
            </li>
            <li>
              <strong>Envie uma notificação:</strong> Vá até a página "Enviar" e envie uma
              notificação para o identificador que você configurou
            </li>
            <li>
              <strong>Observe:</strong> A notificação aparecerá aqui em tempo real!
            </li>
          </ol>

          <div className="mt-4 p-3 bg-blue-50 border border-blue-200 rounded">
            <p className="text-sm text-blue-800">
              💡 <strong>Dica:</strong> Abra a página "Enviar" em outra aba para testar o envio de
              notificações enquanto monitora esta página.
            </p>
          </div>
        </div>
      </Card>
    </div>
  );
}

// Funções auxiliares para push notifications
function urlBase64ToUint8Array(base64String: string) {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);
  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}

function arrayBufferToBase64(buffer: ArrayBuffer) {
  const bytes = new Uint8Array(buffer);
  let binary = '';
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return window.btoa(binary);
}
