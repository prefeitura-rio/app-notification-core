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
    icon: '/favicon.ico',
    badge: '/favicon.ico',
    vibrate: [200, 100, 200],
    tag: data.id || 'notification',
    requireInteraction: false,
    data: {
      dateOfArrival: Date.now(),
      primaryKey: data.id,
      customData: data.data || {},
      url: '/test',
    },
    actions: [
      { action: 'open', title: 'Abrir', icon: '/favicon.ico' },
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

  const url = event.notification.data.url || '/test';

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
