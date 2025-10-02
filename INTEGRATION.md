# Guia de Integração - Sistema de Notificações

Este guia mostra como integrar o sistema de notificações em sua aplicação Next.js para receber notificações in-app (WebSocket) e push notifications.

## 📋 Índice

1. [Configuração Inicial](#configuração-inicial)
2. [Notificações In-App (WebSocket)](#notificações-in-app-websocket)
3. [Push Notifications](#push-notifications)
4. [API REST](#api-rest)
5. [Exemplos Completos](#exemplos-completos)

---

## 🚀 Configuração Inicial

### 1. Variáveis de Ambiente

Adicione ao seu `.env.local`:

```env
NEXT_PUBLIC_NOTIFICATION_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_NOTIFICATION_WS_URL=ws://localhost:8080/api/v1/ws
NEXT_PUBLIC_VAPID_PUBLIC_KEY=sua_chave_vapid_publica
```

### 2. Cliente API

Crie um cliente API em `lib/notificationApi.ts`:

```typescript
const API_URL = process.env.NEXT_PUBLIC_NOTIFICATION_API_URL || 'http://localhost:8080/api/v1';

export const notificationApi = {
  async getNotifications(cpf?: string, phone?: string, email?: string) {
    let url = `${API_URL}/notifications`;

    if (cpf) url = `${API_URL}/notifications/cpf/${cpf}`;
    else if (phone) url = `${API_URL}/notifications/phone/${phone}`;
    else if (email) url = `${API_URL}/notifications/email/${email}`;

    const response = await fetch(url);
    return response.json();
  },

  async markAsRead(id: string) {
    await fetch(`${API_URL}/notifications/${id}/read`, {
      method: 'POST',
    });
  },
};
```

---

## 🔔 Notificações In-App (WebSocket)

### Hook Customizado

Crie `hooks/useNotifications.ts`:

```typescript
'use client';

import { useEffect, useState, useCallback } from 'react';

interface Notification {
  id: string;
  title: string;
  message: string;
  type: string;
  status: string;
  created_at: string;
  read_at?: string;
  is_html: boolean;
}

export function useNotifications(userIdentifier?: string) {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [ws, setWs] = useState<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);

  const connect = useCallback(() => {
    const WS_URL = process.env.NEXT_PUBLIC_NOTIFICATION_WS_URL || 'ws://localhost:8080/api/v1/ws';
    let wsUrl = WS_URL;

    // Adicionar identificador do usuário como query parameter
    if (userIdentifier) {
      wsUrl += `?user_id=${encodeURIComponent(userIdentifier)}`;
    }

    const websocket = new WebSocket(wsUrl);

    websocket.onopen = () => {
      console.log('✅ Conectado ao servidor de notificações');
      setIsConnected(true);
    };

    websocket.onmessage = (event) => {
      try {
        const notification = JSON.parse(event.data);
        console.log('📨 Nova notificação recebida:', notification);

        setNotifications((prev) => [notification, ...prev]);

        // Mostrar notificação do navegador (opcional)
        if ('Notification' in window && Notification.permission === 'granted') {
          new Notification(notification.title, {
            body: notification.message,
            icon: '/notification-icon.png',
          });
        }
      } catch (error) {
        console.error('Erro ao processar notificação:', error);
      }
    };

    websocket.onerror = (error) => {
      console.error('❌ Erro no WebSocket:', error);
      setIsConnected(false);
    };

    websocket.onclose = () => {
      console.log('🔌 Desconectado do servidor de notificações');
      setIsConnected(false);

      // Reconectar após 5 segundos
      setTimeout(() => {
        console.log('🔄 Tentando reconectar...');
        connect();
      }, 5000);
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, [userIdentifier]);

  useEffect(() => {
    const cleanup = connect();
    return cleanup;
  }, [connect]);

  const markAsRead = async (id: string) => {
    try {
      await notificationApi.markAsRead(id);
      setNotifications((prev) =>
        prev.map((n) =>
          n.id === id ? { ...n, status: 'read', read_at: new Date().toISOString() } : n
        )
      );
    } catch (error) {
      console.error('Erro ao marcar como lida:', error);
    }
  };

  return {
    notifications,
    isConnected,
    markAsRead,
  };
}
```

### Uso no Componente

```typescript
'use client';

import { useNotifications } from '@/hooks/useNotifications';

export default function NotificationPanel() {
  const userCpf = '12345678901'; // Obter do contexto de autenticação
  const { notifications, isConnected, markAsRead } = useNotifications(userCpf);

  return (
    <div>
      <div className="flex items-center gap-2">
        <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
        <span>{isConnected ? 'Conectado' : 'Desconectado'}</span>
      </div>

      <div className="space-y-2">
        {notifications.map((notification) => (
          <div
            key={notification.id}
            className={`p-4 rounded-lg border ${
              notification.status === 'read' ? 'bg-gray-50' : 'bg-blue-50'
            }`}
          >
            <h3 className="font-bold">{notification.title}</h3>
            <p>{notification.message}</p>
            <p className="text-xs text-gray-500">
              {new Date(notification.created_at).toLocaleString('pt-BR')}
            </p>
            {notification.status !== 'read' && (
              <button
                onClick={() => markAsRead(notification.id)}
                className="text-sm text-blue-600 hover:underline mt-2"
              >
                Marcar como lida
              </button>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

## 📲 Push Notifications

### 1. Solicitar Permissão e Registrar

Crie `lib/pushNotifications.ts`:

```typescript
const VAPID_PUBLIC_KEY = process.env.NEXT_PUBLIC_VAPID_PUBLIC_KEY;
const API_URL = process.env.NEXT_PUBLIC_NOTIFICATION_API_URL;

export async function registerPushNotifications(userCpf?: string, userPhone?: string) {
  if (!('serviceWorker' in navigator) || !('PushManager' in window)) {
    console.error('Push notifications não suportadas neste navegador');
    return;
  }

  try {
    // 1. Solicitar permissão
    const permission = await Notification.requestPermission();
    if (permission !== 'granted') {
      console.log('Permissão de notificação negada');
      return;
    }

    // 2. Validar VAPID key
    if (!VAPID_PUBLIC_KEY || VAPID_PUBLIC_KEY === 'your_vapid_public_key_here') {
      console.error('VAPID_PUBLIC_KEY não configurada. Configure NEXT_PUBLIC_VAPID_PUBLIC_KEY no .env.local');
      return;
    }

    // 3. Registrar Service Worker
    const registration = await navigator.serviceWorker.register('/sw.js');
    await navigator.serviceWorker.ready;

    // 4. Criar subscrição
    const subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(VAPID_PUBLIC_KEY),
    });

    // 5. Enviar subscrição para o servidor
    await fetch(`${API_URL}/subscriptions`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        user_cpf: userCpf,
        user_phone: userPhone,
        endpoint: subscription.endpoint,
        p256dh: arrayBufferToBase64(subscription.getKey('p256dh')!),
        auth: arrayBufferToBase64(subscription.getKey('auth')!),
      }),
    });

    console.log('✅ Push notifications registradas com sucesso');
  } catch (error) {
    console.error('Erro ao registrar push notifications:', error);
  }
}

export async function unregisterPushNotifications() {
  try {
    const registration = await navigator.serviceWorker.ready;
    const subscription = await registration.pushManager.getSubscription();

    if (subscription) {
      await fetch(`${API_URL}/subscriptions?endpoint=${encodeURIComponent(subscription.endpoint)}`, {
        method: 'DELETE',
      });
      await subscription.unsubscribe();
      console.log('✅ Push notifications canceladas');
    }
  } catch (error) {
    console.error('Erro ao cancelar push notifications:', error);
  }
}

// Funções auxiliares
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
```

### 2. Service Worker

Crie `public/sw.js`:

```javascript
self.addEventListener('push', function (event) {
  console.log('🔔 Push notification received:', event);

  let data = { title: 'Notificação', message: 'Você tem uma nova notificação' };

  if (event.data) {
    try {
      data = event.data.json();
      console.log('📨 Parsed notification data:', data);
    } catch (error) {
      console.error('❌ Error parsing push data:', error);
    }
  }

  const options = {
    body: data.message,
    icon: '/notification-icon.png',
    badge: '/badge-icon.png',
    vibrate: [200, 100, 200],
    tag: data.id || 'notification',
    requireInteraction: false,
    data: {
      dateOfArrival: Date.now(),
      primaryKey: data.id,
      customData: data.data || {},
      url: '/notifications',
    },
    actions: [
      { action: 'open', title: 'Abrir', icon: '/notification-icon.png' },
      { action: 'close', title: 'Fechar' },
    ],
  };

  console.log('✅ Showing notification:', data.title);
  event.waitUntil(self.registration.showNotification(data.title, options));
});

self.addEventListener('notificationclick', function (event) {
  console.log('👆 Notification clicked:', event);
  console.log('Action:', event.action);
  console.log('Notification data:', event.notification.data);

  event.notification.close();

  const url = event.notification.data.url || '/notifications';

  if (event.action === 'close') {
    // Apenas fecha a notificação
    return;
  }

  // Para ação 'open' ou clique na notificação
  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true })
      .then(function (clientList) {
        // Se já existe uma janela aberta, focar nela
        for (let i = 0; i < clientList.length; i++) {
          const client = clientList[i];
          if (client.url.indexOf(url) >= 0 && 'focus' in client) {
            return client.focus();
          }
        }
        // Caso contrário, abrir nova janela
        if (clients.openWindow) {
          return clients.openWindow(url);
        }
      })
  );
});

self.addEventListener('install', function (event) {
  console.log('Service Worker installing...');
  self.skipWaiting();
});

self.addEventListener('activate', function (event) {
  console.log('Service Worker activated');
  event.waitUntil(clients.claim());
});
```

### 3. Uso no Componente

```typescript
'use client';

import { registerPushNotifications } from '@/lib/pushNotifications';

export default function EnablePushButton() {
  const handleEnablePush = async () => {
    const userCpf = '12345678901'; // Obter do contexto
    await registerPushNotifications(userCpf);
  };

  return (
    <button
      onClick={handleEnablePush}
      className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
    >
      🔔 Ativar Notificações Push
    </button>
  );
}
```

---

## 📡 API REST

### Buscar Notificações do Usuário

```typescript
// Por CPF
const notifications = await fetch('http://localhost:8080/api/v1/notifications/cpf/12345678901');

// Por Telefone
const notifications = await fetch('http://localhost:8080/api/v1/notifications/phone/11999999999');

// Por Email
const notifications = await fetch('http://localhost:8080/api/v1/notifications/email/user@example.com');
```

### Marcar como Lida

```typescript
await fetch('http://localhost:8080/api/v1/notifications/{id}/read', {
  method: 'POST',
});
```

---

## 💡 Exemplos Completos

### Componente de Notificações Completo

```typescript
'use client';

import { useState, useEffect } from 'react';
import { useNotifications } from '@/hooks/useNotifications';
import { registerPushNotifications } from '@/lib/pushNotifications';

export default function NotificationCenter({ userId }: { userId: string }) {
  const [pushEnabled, setPushEnabled] = useState(false);
  const { notifications, isConnected, markAsRead } = useNotifications(userId);

  useEffect(() => {
    // Verificar se push já está habilitado
    navigator.serviceWorker.ready.then(async (registration) => {
      const subscription = await registration.pushManager.getSubscription();
      setPushEnabled(!!subscription);
    });
  }, []);

  const handleEnablePush = async () => {
    await registerPushNotifications(userId);
    setPushEnabled(true);
  };

  const unreadCount = notifications.filter((n) => n.status !== 'read').length;

  return (
    <div className="max-w-2xl mx-auto p-6">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold">Notificações</h2>
          <div className="flex items-center gap-2 mt-1">
            <span className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'}`} />
            <span className="text-sm text-gray-600">
              {isConnected ? 'Conectado' : 'Desconectado'}
            </span>
          </div>
        </div>

        {!pushEnabled && (
          <button
            onClick={handleEnablePush}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            🔔 Ativar Push
          </button>
        )}
      </div>

      {/* Badge de não lidas */}
      {unreadCount > 0 && (
        <div className="mb-4 px-4 py-2 bg-blue-100 text-blue-800 rounded">
          {unreadCount} notificação{unreadCount > 1 ? 'ões' : ''} não lida{unreadCount > 1 ? 's' : ''}
        </div>
      )}

      {/* Lista de notificações */}
      <div className="space-y-3">
        {notifications.length === 0 ? (
          <p className="text-center text-gray-500 py-8">Nenhuma notificação</p>
        ) : (
          notifications.map((notification) => (
            <div
              key={notification.id}
              className={`p-4 rounded-lg border ${
                notification.status === 'read'
                  ? 'bg-white border-gray-200'
                  : 'bg-blue-50 border-blue-200'
              }`}
            >
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <h3 className="font-semibold text-gray-900">{notification.title}</h3>
                  {notification.is_html ? (
                    <div
                      className="text-gray-700 mt-1"
                      dangerouslySetInnerHTML={{ __html: notification.message }}
                    />
                  ) : (
                    <p className="text-gray-700 mt-1">{notification.message}</p>
                  )}
                  <p className="text-xs text-gray-500 mt-2">
                    {new Date(notification.created_at).toLocaleString('pt-BR')}
                  </p>
                </div>

                {notification.status !== 'read' && (
                  <button
                    onClick={() => markAsRead(notification.id)}
                    className="ml-4 text-sm text-blue-600 hover:underline"
                  >
                    Marcar como lida
                  </button>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
```

---

## 🔐 Considerações de Segurança

1. **Autenticação**: Adicione autenticação JWT/Bearer token às chamadas da API
2. **CORS**: Configure CORS adequadamente no servidor
3. **Rate Limiting**: Implemente rate limiting para prevenir abuse
4. **Validação**: Sempre valide os dados recebidos do WebSocket

## 📚 Documentação da API

Acesse a documentação completa Swagger em:
```
http://localhost:8080/swagger/index.html
```

## 🐛 Troubleshooting

### WebSocket não conecta
- Verifique se o servidor está rodando
- Confirme a URL do WebSocket no `.env`
- Certifique-se de usar `?user_id=` como parâmetro (não `?user=`)
- Verifique o console do navegador para erros

### Push Notifications não funcionam

**IMPORTANTE**: Push notifications NÃO aparecem na interface da aplicação - elas aparecem como **notificações nativas do sistema operacional** (como WhatsApp, Gmail, etc).

**Onde aparecem:**
- 🪟 Windows: Canto inferior direito (Action Center)
- 🍎 macOS: Canto superior direito (Notification Center)
- 🐧 Linux: Geralmente canto superior direito
- 📱 Mobile: Barra de status

**Checklist de verificação:**
1. ✅ Certifique-se de que o site está em HTTPS (ou localhost)
2. ✅ Verifique se o Service Worker foi registrado (DevTools > Application > Service Workers)
3. ✅ Confirme que a permissão foi concedida (DevTools > Application > Permissions)
4. ✅ Verifique se a chave VAPID pública está configurada corretamente
5. ✅ **Minimize o navegador ou troque de aba** - notificações não aparecem se a aba estiver em foco
6. ✅ Verifique os logs do backend para: `Push sent successfully to subscription...`
7. ✅ Verifique o console do navegador para: `🔔 Push notification received`

### Notificações in-app não aparecem
- Verifique o console para erros de parsing
- Confirme que o formato da notificação está correto
- Teste a conexão WebSocket com ferramentas de debug
- Verifique se o user_id está sendo passado corretamente

### Como testar push notifications

1. **Registre a subscription:**
   - Abra sua aplicação
   - Clique no botão "Ativar Push"
   - Conceda permissão quando solicitado

2. **Envie uma notificação:**
   - Use a API para enviar uma notificação com tipo `push`, `both` ou `all`
   - Certifique-se de usar o mesmo CPF/telefone usado no registro

3. **Verifique a notificação:**
   - **IMPORTANTE**: Minimize o navegador ou troque para outra aba
   - A notificação deve aparecer como notificação do sistema
   - Se não aparecer, verifique os logs conforme checklist acima

---

## 📞 Suporte

Para mais informações ou dúvidas, consulte a documentação da API ou abra uma issue no repositório.
