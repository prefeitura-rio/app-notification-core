import type { NotificationType, NotificationStatus } from '@/types';

interface BadgeProps {
  type?: NotificationType;
  status?: NotificationStatus;
}

export function NotificationBadge({ type, status }: BadgeProps) {
  if (type) {
    const typeColors = {
      'in-app': 'bg-blue-100 text-blue-800',
      'push': 'bg-purple-100 text-purple-800',
      'both': 'bg-green-100 text-green-800',
    };

    return (
      <span className={`px-2 py-1 text-xs font-medium rounded-full ${typeColors[type]}`}>
        {type}
      </span>
    );
  }

  if (status) {
    const statusColors = {
      pending: 'bg-yellow-100 text-yellow-800',
      sent: 'bg-blue-100 text-blue-800',
      delivered: 'bg-green-100 text-green-800',
      read: 'bg-gray-100 text-gray-800',
      failed: 'bg-red-100 text-red-800',
    };

    return (
      <span className={`px-2 py-1 text-xs font-medium rounded-full ${statusColors[status]}`}>
        {status}
      </span>
    );
  }

  return null;
}
