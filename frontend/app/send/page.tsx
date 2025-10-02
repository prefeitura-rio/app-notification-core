'use client';

import { useState, useEffect } from 'react';
import { Card } from '@/components/Card';
import { Input } from '@/components/Input';
import { Button } from '@/components/Button';
import { api } from '@/lib/api';
import type { SendNotificationRequest, Group } from '@/types';

export default function SendPage() {
  const [groups, setGroups] = useState<Group[]>([]);
  const [sendType, setSendType] = useState<'user' | 'group' | 'broadcast' | 'batch'>('user');
  const [sending, setSending] = useState(false);
  const [batchRecipients, setBatchRecipients] = useState('');

  const [formData, setFormData] = useState<SendNotificationRequest>({
    title: '',
    message: '',
    type: 'in-app',
    cpf: '',
    phone: '',
    email: '',
    is_html: false,
  });

  const [channels, setChannels] = useState({
    inApp: true,
    push: false,
    email: false,
  });

  const [isScheduled, setIsScheduled] = useState(false);
  const [scheduledFor, setScheduledFor] = useState('');

  const [selectedGroupId, setSelectedGroupId] = useState('');

  useEffect(() => {
    loadGroups();
  }, []);

  const loadGroups = async () => {
    try {
      const data = await api.get<Group[]>('/api/v1/groups');
      setGroups(data || []);
    } catch (error) {
      console.error('Failed to load groups:', error);
    }
  };

  const getNotificationType = (): string => {
    const { inApp, push, email } = channels;

    if (inApp && push && email) return 'all';
    if (inApp && push) return 'both';
    if (inApp && email) return 'all'; // Não temos tipo específico para in-app + email
    if (push && email) return 'all'; // Não temos tipo específico para push + email
    if (inApp) return 'in-app';
    if (push) return 'push';
    if (email) return 'email';

    return 'in-app'; // Default
  };

  const parseBatchRecipients = () => {
    const lines = batchRecipients.trim().split('\n');
    const recipients = [];

    for (const line of lines) {
      if (!line.trim()) continue;

      const parts = line.split(',').map(p => p.trim());

      const recipient: any = {};

      // Formato: cpf,phone,email,name ou qualquer combinação
      if (parts[0]) recipient.cpf = parts[0];
      if (parts[1]) recipient.phone = parts[1];
      if (parts[2]) recipient.email = parts[2];
      if (parts[3]) recipient.name = parts[3];

      recipients.push(recipient);
    }

    return recipients;
  };

  const handleSend = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validar se ao menos um canal está selecionado
    if (!channels.inApp && !channels.push && !channels.email) {
      alert('Selecione pelo menos um canal de notificação');
      return;
    }

    setSending(true);

    const notificationType = getNotificationType();

    // Preparar dados de agendamento
    const basePayload: any = {
      type: notificationType,
      is_html: formData.is_html,
    };

    if (isScheduled && scheduledFor) {
      basePayload.is_scheduled = true;
      // Converter para RFC3339
      const date = new Date(scheduledFor);
      basePayload.scheduled_for = date.toISOString();
    }

    try {
      if (sendType === 'user') {
        await api.post('/api/v1/notifications/send/user', {
          ...formData,
          ...basePayload,
        });
        alert(isScheduled ? 'Notificação agendada com sucesso!' : 'Notificação enviada com sucesso!');
      } else if (sendType === 'group') {
        await api.post(`/api/v1/notifications/send/group/${selectedGroupId}`, {
          title: formData.title,
          message: formData.message,
          ...basePayload,
        });
        alert(isScheduled ? 'Notificação agendada para o grupo com sucesso!' : 'Notificação enviada para o grupo com sucesso!');
      } else if (sendType === 'broadcast') {
        await api.post('/api/v1/notifications/send/broadcast', {
          title: formData.title,
          message: formData.message,
          ...basePayload,
        });
        alert(isScheduled ? 'Notificação broadcast agendada com sucesso!' : 'Notificação broadcast enviada com sucesso!');
      } else if (sendType === 'batch') {
        const recipients = parseBatchRecipients();

        if (recipients.length === 0) {
          alert('Adicione pelo menos um destinatário');
          setSending(false);
          return;
        }

        const response = await api.post('/api/v1/notifications/send/batch', {
          title: formData.title,
          message: formData.message,
          ...basePayload,
          recipients: recipients,
        });

        const result = response as any;
        let message = `Envio em lote concluído!\n\nTotal: ${result.total}\nSucesso: ${result.succeeded}\nFalhas: ${result.failed}`;

        if (result.errors && result.errors.length > 0) {
          message += '\n\nErros:\n' + result.errors.join('\n');
        }

        alert(message);
        setBatchRecipients('');
      }

      setFormData({
        title: '',
        message: '',
        type: 'in-app',
        cpf: '',
        phone: '',
        email: '',
        is_html: false,
      });
      setChannels({
        inApp: true,
        push: false,
        email: false,
      });
      setIsScheduled(false);
      setScheduledFor('');
    } catch (error) {
      console.error('Failed to send notification:', error);
      alert('Erro ao enviar notificação');
    } finally {
      setSending(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-800 mb-2">Enviar Notificação</h1>
        <p className="text-gray-600">Envie notificações para usuários, grupos ou todos</p>
      </div>

      <Card>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
          <Button
            variant={sendType === 'user' ? 'primary' : 'secondary'}
            onClick={() => setSendType('user')}
          >
            Usuário Específico
          </Button>
          <Button
            variant={sendType === 'group' ? 'primary' : 'secondary'}
            onClick={() => setSendType('group')}
          >
            Grupo
          </Button>
          <Button
            variant={sendType === 'batch' ? 'primary' : 'secondary'}
            onClick={() => setSendType('batch')}
          >
            Envio em Lote
          </Button>
          <Button
            variant={sendType === 'broadcast' ? 'primary' : 'secondary'}
            onClick={() => setSendType('broadcast')}
          >
            Broadcast (Todos)
          </Button>
        </div>

        <form onSubmit={handleSend} className="space-y-4">
          <Input
            label="Título"
            placeholder="Título da notificação"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            required
          />

          <div className="flex flex-col gap-1">
            <label className="text-sm font-medium text-gray-700">Mensagem</label>
            <textarea
              className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Conteúdo da notificação"
              rows={4}
              value={formData.message}
              onChange={(e) => setFormData({ ...formData, message: e.target.value })}
              required
            />
          </div>

          <div className="flex flex-col gap-2">
            <label className="text-sm font-medium text-gray-700">Canais de Notificação *</label>
            <div className="bg-gray-50 p-4 rounded-lg border border-gray-200">
              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id="channel_inapp"
                    checked={channels.inApp}
                    onChange={(e) => setChannels({ ...channels, inApp: e.target.checked })}
                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <label htmlFor="channel_inapp" className="text-sm font-medium text-gray-700 cursor-pointer">
                    📱 In-App (WebSocket)
                  </label>
                </div>
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id="channel_push"
                    checked={channels.push}
                    onChange={(e) => setChannels({ ...channels, push: e.target.checked })}
                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <label htmlFor="channel_push" className="text-sm font-medium text-gray-700 cursor-pointer">
                    🔔 Push Notification
                  </label>
                </div>
                <div className="flex items-center gap-2">
                  <input
                    type="checkbox"
                    id="channel_email"
                    checked={channels.email}
                    onChange={(e) => setChannels({ ...channels, email: e.target.checked })}
                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <label htmlFor="channel_email" className="text-sm font-medium text-gray-700 cursor-pointer">
                    ✉️ Email
                  </label>
                </div>
              </div>
              {!channels.inApp && !channels.push && !channels.email && (
                <p className="text-xs text-red-600 mt-2">
                  Selecione pelo menos um canal
                </p>
              )}
            </div>
          </div>

          {channels.email && (
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                id="is_html"
                checked={formData.is_html}
                onChange={(e) => setFormData({ ...formData, is_html: e.target.checked })}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <label htmlFor="is_html" className="text-sm font-medium text-gray-700">
                Email em formato HTML
              </label>
            </div>
          )}

          <div className="bg-purple-50 border border-purple-200 rounded-lg p-4">
            <div className="flex items-center gap-2 mb-3">
              <input
                type="checkbox"
                id="is_scheduled"
                checked={isScheduled}
                onChange={(e) => setIsScheduled(e.target.checked)}
                className="w-4 h-4 text-purple-600 border-gray-300 rounded focus:ring-purple-500"
              />
              <label htmlFor="is_scheduled" className="text-sm font-medium text-purple-800">
                📅 Agendar envio
              </label>
            </div>

            {isScheduled && (
              <div className="flex flex-col gap-1">
                <label className="text-sm font-medium text-purple-700">
                  Data e Hora do Envio
                </label>
                <input
                  type="datetime-local"
                  className="px-3 py-2 border border-purple-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                  value={scheduledFor}
                  onChange={(e) => setScheduledFor(e.target.value)}
                  min={new Date().toISOString().slice(0, 16)}
                  required
                />
                <p className="text-xs text-purple-600 mt-1">
                  A notificação será enviada automaticamente na data e hora especificada
                </p>
              </div>
            )}
          </div>

          {sendType === 'user' && (
            <div className="grid grid-cols-3 gap-4 pt-4 border-t">
              <Input
                label="CPF (opcional)"
                placeholder="12345678901"
                value={formData.cpf}
                onChange={(e) => setFormData({ ...formData, cpf: e.target.value })}
              />
              <Input
                label="Telefone (opcional)"
                placeholder="11999999999"
                value={formData.phone}
                onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
              />
              <Input
                label="Email (opcional)"
                type="email"
                placeholder="usuario@exemplo.com"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              />
              <p className="col-span-3 text-sm text-gray-600">
                * Preencha pelo menos CPF, Telefone ou Email
              </p>
            </div>
          )}

          {sendType === 'group' && (
            <div className="pt-4 border-t">
              <label className="text-sm font-medium text-gray-700 block mb-2">
                Selecione o Grupo
              </label>
              <select
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                value={selectedGroupId}
                onChange={(e) => setSelectedGroupId(e.target.value)}
                required
              >
                <option value="">Selecione um grupo...</option>
                {groups.map((group) => (
                  <option key={group.id} value={group.id}>
                    {group.name} ({group.members?.length || 0} membros)
                  </option>
                ))}
              </select>
            </div>
          )}

          {sendType === 'batch' && (
            <div className="pt-4 border-t space-y-3">
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <p className="text-sm text-blue-800 font-medium mb-2">
                  📋 Formato CSV para destinatários (um por linha):
                </p>
                <p className="text-xs text-blue-700 font-mono">
                  cpf,telefone,email,nome
                </p>
                <p className="text-xs text-blue-600 mt-2">
                  Você pode deixar campos vazios. Exemplos:<br />
                  12345678901,11999999999,user@example.com,João Silva<br />
                  ,11988888888,,Maria Santos<br />
                  ,,admin@example.com,Admin
                </p>
              </div>

              <div className="flex flex-col gap-1">
                <label className="text-sm font-medium text-gray-700">
                  Destinatários (CSV)
                </label>
                <textarea
                  className="px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
                  placeholder="12345678901,11999999999,user@example.com,João Silva"
                  rows={8}
                  value={batchRecipients}
                  onChange={(e) => setBatchRecipients(e.target.value)}
                  required
                />
                <p className="text-xs text-gray-500">
                  {batchRecipients.trim().split('\n').filter(l => l.trim()).length} destinatário(s)
                </p>
              </div>
            </div>
          )}

          {sendType === 'broadcast' && (
            <div className="pt-4 border-t">
              <div className="bg-orange-50 border border-orange-200 rounded-lg p-4">
                <p className="text-sm text-orange-800">
                  ⚠️ Esta notificação será enviada para <strong>todos os usuários</strong> conectados.
                </p>
              </div>
            </div>
          )}

          <Button type="submit" disabled={sending}>
            {sending ? 'Enviando...' : 'Enviar Notificação'}
          </Button>
        </form>
      </Card>
    </div>
  );
}
