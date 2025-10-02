export interface Group {
  id: string;
  name: string;
  description: string;
  members?: Member[];
  created_at: string;
  updated_at: string;
}

export interface Member {
  id: string;
  group_id: string;
  cpf?: string;
  phone?: string;
  email?: string;
  name?: string;
  created_at: string;
  updated_at: string;
}

export type NotificationType = 'in-app' | 'push' | 'email' | 'both' | 'all';
export type NotificationStatus = 'pending' | 'sent' | 'delivered' | 'read' | 'failed';

export interface Notification {
  id: string;
  title: string;
  message: string;
  type: NotificationType;
  status: NotificationStatus;
  data?: Record<string, any>;
  user_cpf?: string;
  user_phone?: string;
  user_email?: string;
  group_id?: string;
  broadcast: boolean;
  is_html: boolean;
  created_at: string;
  updated_at: string;
  read_at?: string;
}

export interface SendNotificationRequest {
  title: string;
  message: string;
  type: NotificationType;
  data?: Record<string, any>;
  cpf?: string;
  phone?: string;
  email?: string;
  is_html?: boolean;
}
