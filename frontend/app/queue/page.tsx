'use client';

import { useState, useEffect } from 'react';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';

interface QueueStats {
  queue_name: string;
  messages: number;
  consumers: number;
  dlq_messages: number;
  last_checked: string;
}

export default function QueuePage() {
  const [stats, setStats] = useState<QueueStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshInterval, setRefreshInterval] = useState(5000);
  const [autoRefresh, setAutoRefresh] = useState(true);

  const loadStats = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/queue/stats');
      const data = await response.json();
      setStats(data);
    } catch (error) {
      console.error('Erro ao carregar estatísticas:', error);
    } finally {
      setLoading(false);
    }
  };

  const purgeQueue = async () => {
    if (!confirm('Tem certeza que deseja limpar toda a fila? Esta ação não pode ser desfeita!')) {
      return;
    }

    try {
      const response = await fetch('http://localhost:8080/api/v1/queue/purge', {
        method: 'POST',
      });

      if (response.ok) {
        alert('Fila limpa com sucesso!');
        loadStats();
      } else {
        alert('Erro ao limpar a fila');
      }
    } catch (error) {
      console.error('Erro ao limpar fila:', error);
      alert('Erro ao limpar a fila');
    }
  };

  useEffect(() => {
    loadStats();
  }, []);

  useEffect(() => {
    if (!autoRefresh) return;

    const interval = setInterval(loadStats, refreshInterval);
    return () => clearInterval(interval);
  }, [autoRefresh, refreshInterval]);

  if (loading) {
    return <div className="text-center py-8">Carregando...</div>;
  }

  const getStatusColor = (count: number) => {
    if (count === 0) return 'text-green-600';
    if (count < 100) return 'text-yellow-600';
    if (count < 1000) return 'text-orange-600';
    return 'text-red-600';
  };

  const getStatusBg = (count: number) => {
    if (count === 0) return 'bg-green-50 border-green-200';
    if (count < 100) return 'bg-yellow-50 border-yellow-200';
    if (count < 1000) return 'bg-orange-50 border-orange-200';
    return 'bg-red-50 border-red-200';
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-800 mb-2">Gerenciamento de Filas</h1>
          <p className="text-gray-600">
            Monitore e gerencie a fila de notificações RabbitMQ
          </p>
        </div>

        <div className="flex gap-2">
          <Button
            variant={autoRefresh ? 'primary' : 'secondary'}
            onClick={() => setAutoRefresh(!autoRefresh)}
          >
            {autoRefresh ? '⏸️ Pausar' : '▶️ Retomar'} Auto-Refresh
          </Button>
          <Button variant="secondary" onClick={loadStats}>
            🔄 Atualizar
          </Button>
        </div>
      </div>

      {/* Estatísticas Principais */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {/* Mensagens Pendentes */}
        <div className={`p-6 rounded-lg border-2 ${getStatusBg(stats?.messages || 0)}`}>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Mensagens Pendentes</p>
              <p className={`text-3xl font-bold mt-2 ${getStatusColor(stats?.messages || 0)}`}>
                {stats?.messages || 0}
              </p>
            </div>
            <div className="text-4xl">📬</div>
          </div>
        </div>

        {/* Workers Ativos */}
        <div className="p-6 bg-blue-50 border-2 border-blue-200 rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Workers Ativos</p>
              <p className="text-3xl font-bold text-blue-600 mt-2">{stats?.consumers || 0}</p>
            </div>
            <div className="text-4xl">👷</div>
          </div>
        </div>

        {/* Dead Letter Queue */}
        <div className={`p-6 rounded-lg border-2 ${getStatusBg(stats?.dlq_messages || 0)}`}>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">DLQ (Falhas)</p>
              <p className={`text-3xl font-bold mt-2 ${getStatusColor(stats?.dlq_messages || 0)}`}>
                {stats?.dlq_messages || 0}
              </p>
            </div>
            <div className="text-4xl">💀</div>
          </div>
        </div>

        {/* Taxa de Processamento */}
        <div className="p-6 bg-purple-50 border-2 border-purple-200 rounded-lg">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium text-gray-600">Status</p>
              <p className="text-lg font-bold text-purple-600 mt-2">
                {stats?.messages === 0 ? '✅ Limpo' : stats?.messages! < 100 ? '⚡ Normal' : '⚠️ Alto'}
              </p>
            </div>
            <div className="text-4xl">📊</div>
          </div>
        </div>
      </div>

      {/* Informações da Fila */}
      {stats && (
        <Card title="ℹ️ Informações da Fila">
          <div className="space-y-3">
            <div className="flex justify-between py-2 border-b">
              <span className="font-medium text-gray-700">Nome da Fila:</span>
              <code className="bg-gray-100 px-2 py-1 rounded">{stats.queue_name}</code>
            </div>

            <div className="flex justify-between py-2 border-b">
              <span className="font-medium text-gray-700">Última Verificação:</span>
              <span className="text-gray-600">
                {new Date(stats.last_checked).toLocaleString('pt-BR')}
              </span>
            </div>

            <div className="flex justify-between py-2 border-b">
              <span className="font-medium text-gray-700">Auto-Refresh:</span>
              <div className="flex items-center gap-2">
                <select
                  value={refreshInterval}
                  onChange={(e) => setRefreshInterval(Number(e.target.value))}
                  className="px-2 py-1 border border-gray-300 rounded text-sm"
                  disabled={!autoRefresh}
                >
                  <option value={1000}>1 segundo</option>
                  <option value={2000}>2 segundos</option>
                  <option value={5000}>5 segundos</option>
                  <option value={10000}>10 segundos</option>
                  <option value={30000}>30 segundos</option>
                </select>
              </div>
            </div>
          </div>
        </Card>
      )}

      {/* Alertas */}
      {stats && stats.messages > 1000 && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <span className="text-2xl">⚠️</span>
            <div>
              <h3 className="font-bold text-red-800">Alto volume de mensagens!</h3>
              <p className="text-sm text-red-700 mt-1">
                A fila possui mais de 1000 mensagens pendentes. Considere aumentar o número de workers ou verificar se há problemas no processamento.
              </p>
            </div>
          </div>
        </div>
      )}

      {stats && stats.dlq_messages > 0 && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <span className="text-2xl">⚠️</span>
            <div>
              <h3 className="font-bold text-yellow-800">Mensagens na Dead Letter Queue</h3>
              <p className="text-sm text-yellow-700 mt-1">
                Existem {stats.dlq_messages} mensagem(ns) que falharam após 3 tentativas. Verifique os logs para identificar o problema.
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Controles */}
      <Card title="🛠️ Controles da Fila">
        <div className="space-y-4">
          <div className="p-4 bg-gray-50 rounded-lg border">
            <h3 className="font-bold text-gray-800 mb-2">Limpar Fila</h3>
            <p className="text-sm text-gray-600 mb-3">
              Remove todas as mensagens pendentes da fila. Esta ação não pode ser desfeita!
            </p>
            <Button
              variant="danger"
              onClick={purgeQueue}
              disabled={stats?.messages === 0}
            >
              🗑️ Limpar Fila ({stats?.messages || 0} mensagens)
            </Button>
          </div>

          <div className="p-4 bg-blue-50 rounded-lg border border-blue-200">
            <h3 className="font-bold text-blue-800 mb-2">RabbitMQ Management UI</h3>
            <p className="text-sm text-blue-700 mb-3">
              Acesse a interface de gerenciamento do RabbitMQ para controles avançados.
            </p>
            <a
              href="http://localhost:15672"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-block px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              🐰 Abrir RabbitMQ Management
            </a>
            <p className="text-xs text-blue-600 mt-2">
              Credenciais: admin / admin123
            </p>
          </div>
        </div>
      </Card>

      {/* Métricas Visuais */}
      <Card title="📈 Métricas Visuais">
        <div className="space-y-4">
          <div>
            <div className="flex justify-between mb-2">
              <span className="text-sm font-medium text-gray-700">Capacidade da Fila</span>
              <span className="text-sm text-gray-600">
                {stats?.messages || 0} / 100,000 mensagens
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-4">
              <div
                className={`h-4 rounded-full transition-all ${
                  (stats?.messages || 0) < 1000
                    ? 'bg-green-500'
                    : (stats?.messages || 0) < 10000
                    ? 'bg-yellow-500'
                    : 'bg-red-500'
                }`}
                style={{ width: `${Math.min(((stats?.messages || 0) / 100000) * 100, 100)}%` }}
              />
            </div>
          </div>

          <div>
            <div className="flex justify-between mb-2">
              <span className="text-sm font-medium text-gray-700">Workers Ativos</span>
              <span className="text-sm text-gray-600">
                {stats?.consumers || 0} / 3 configurados
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-4">
              <div
                className="bg-blue-500 h-4 rounded-full transition-all"
                style={{ width: `${((stats?.consumers || 0) / 3) * 100}%` }}
              />
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}
