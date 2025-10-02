'use client';

import { useEffect, useState, useCallback } from 'react';
import type { Notification } from '@/types';

const WS_URL = process.env.NEXT_PUBLIC_API_URL?.replace('http', 'ws') || 'ws://localhost:8080';

export function useWebSocket(userId: string) {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [isConnected, setIsConnected] = useState(false);
  const [ws, setWs] = useState<WebSocket | null>(null);

  useEffect(() => {
    if (!userId) return;

    const websocket = new WebSocket(`${WS_URL}/api/v1/ws?user_id=${userId}`);

    websocket.onopen = () => {
      console.log('WebSocket connected');
      setIsConnected(true);
    };

    websocket.onmessage = (event) => {
      try {
        const notification = JSON.parse(event.data);
        setNotifications(prev => [notification, ...prev]);
      } catch (error) {
        console.error('Failed to parse notification:', error);
      }
    };

    websocket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    websocket.onclose = () => {
      console.log('WebSocket disconnected');
      setIsConnected(false);
    };

    setWs(websocket);

    return () => {
      websocket.close();
    };
  }, [userId]);

  const clearNotifications = useCallback(() => {
    setNotifications([]);
  }, []);

  return { notifications, isConnected, clearNotifications };
}
