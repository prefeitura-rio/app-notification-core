'use client';

import { useState, useEffect } from 'react';
import { Card } from '@/components/Card';
import { Button } from '@/components/Button';

interface VAPIDKeys {
  public_key: string;
  private_key: string;
  subject: string;
}

interface IntegrationConfig {
  backend_url: string;
  websocket_url: string;
  current_vapid: VAPIDKeys;
  api_endpoints: string[];
  swagger_url: string;
}

interface EnvTemplates {
  backend: string;
  frontend: string;
}

export default function IntegrationsPage() {
  const [config, setConfig] = useState<IntegrationConfig | null>(null);
  const [envTemplates, setEnvTemplates] = useState<EnvTemplates | null>(null);
  const [loading, setLoading] = useState(true);
  const [generating, setGenerating] = useState(false);
  const [generatedKeys, setGeneratedKeys] = useState<VAPIDKeys | null>(null);
  const [showKeysAlert, setShowKeysAlert] = useState(false);

  useEffect(() => {
    loadConfig();
    loadEnvTemplates();
  }, []);

  const loadConfig = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/integration/config');
      const data = await response.json();
      setConfig(data);
    } catch (error) {
      console.error('Erro ao carregar config:', error);
    } finally {
      setLoading(false);
    }
  };

  const loadEnvTemplates = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/v1/integration/env-template');
      const data = await response.json();
      setEnvTemplates(data);
    } catch (error) {
      console.error('Erro ao carregar templates:', error);
    }
  };

  const generateVAPIDKeys = async () => {
    setGenerating(true);
    try {
      const response = await fetch('http://localhost:8080/api/v1/integration/vapid/generate', {
        method: 'POST',
      });
      const keys = await response.json();
      setGeneratedKeys(keys);
      setShowKeysAlert(true);
    } catch (error) {
      console.error('Erro ao gerar chaves:', error);
      alert('Erro ao gerar chaves VAPID');
    } finally {
      setGenerating(false);
    }
  };

  const copyToClipboard = (text: string, label: string) => {
    navigator.clipboard.writeText(text);
    alert(`${label} copiado para a área de transferência!`);
  };

  const hasValidVAPIDKeys = () => {
    if (!config) return false;
    const pubKey = config.current_vapid.public_key;
    return pubKey && pubKey !== 'your_vapid_public_key_here' && pubKey.length > 20;
  };

  if (loading) {
    return <div className="text-center py-8">Carregando...</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-800 mb-2">Integrações</h1>
        <p className="text-gray-600">
          Configure e gerencie integrações com aplicações frontend
        </p>
      </div>

      {/* Status das Chaves VAPID */}
      <Card title="🔑 Chaves VAPID">
        <div className="space-y-4">
          {!hasValidVAPIDKeys() && (
            <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
              <div className="flex items-start gap-3">
                <span className="text-2xl">⚠️</span>
                <div>
                  <h3 className="font-bold text-yellow-800">Chaves VAPID não configuradas</h3>
                  <p className="text-sm text-yellow-700 mt-1">
                    Para habilitar push notifications, você precisa gerar e configurar chaves VAPID.
                  </p>
                </div>
              </div>
            </div>
          )}

          {hasValidVAPIDKeys() && (
            <div className="bg-green-50 border border-green-200 rounded-lg p-4">
              <div className="flex items-start gap-3">
                <span className="text-2xl">✅</span>
                <div>
                  <h3 className="font-bold text-green-800">Chaves VAPID configuradas</h3>
                  <p className="text-sm text-green-700 mt-1">
                    Push notifications estão habilitadas e prontas para uso.
                  </p>
                </div>
              </div>
            </div>
          )}

          {config && (
            <div className="space-y-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Chave Pública
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={config.current_vapid.public_key}
                    readOnly
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 font-mono text-sm"
                  />
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => copyToClipboard(config.current_vapid.public_key, 'Chave pública')}
                  >
                    Copiar
                  </Button>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Chave Privada
                </label>
                <div className="flex gap-2">
                  <input
                    type="password"
                    value={config.current_vapid.private_key}
                    readOnly
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 font-mono text-sm"
                  />
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => copyToClipboard(config.current_vapid.private_key, 'Chave privada')}
                  >
                    Copiar
                  </Button>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Subject
                </label>
                <input
                  type="text"
                  value={config.current_vapid.subject}
                  readOnly
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50"
                />
              </div>
            </div>
          )}

          <div className="pt-4 border-t">
            <Button onClick={generateVAPIDKeys} disabled={generating}>
              {generating ? 'Gerando...' : '🔄 Gerar Novas Chaves VAPID'}
            </Button>
            <p className="text-sm text-gray-600 mt-2">
              Gera um novo par de chaves VAPID para push notifications. Após gerar, você precisará
              atualizar o arquivo .env do backend.
            </p>
          </div>
        </div>
      </Card>

      {/* Chaves Geradas */}
      {generatedKeys && showKeysAlert && (
        <Card title="✨ Novas Chaves Geradas">
          <div className="space-y-4">
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
              <div className="flex items-start gap-3">
                <span className="text-2xl">💡</span>
                <div>
                  <h3 className="font-bold text-blue-800">Importante!</h3>
                  <p className="text-sm text-blue-700 mt-1">
                    Copie estas chaves e atualize o arquivo <code className="bg-blue-100 px-1 rounded">.env</code> do backend.
                    Depois reinicie o servidor para aplicar as mudanças.
                  </p>
                </div>
              </div>
            </div>

            <div className="space-y-3">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  VAPID_PUBLIC_KEY
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={generatedKeys.public_key}
                    readOnly
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 font-mono text-sm"
                  />
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => copyToClipboard(generatedKeys.public_key, 'Chave pública')}
                  >
                    Copiar
                  </Button>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  VAPID_PRIVATE_KEY
                </label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={generatedKeys.private_key}
                    readOnly
                    className="flex-1 px-3 py-2 border border-gray-300 rounded-lg bg-gray-50 font-mono text-sm"
                  />
                  <Button
                    variant="secondary"
                    size="sm"
                    onClick={() => copyToClipboard(generatedKeys.private_key, 'Chave privada')}
                  >
                    Copiar
                  </Button>
                </div>
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  VAPID_SUBJECT
                </label>
                <input
                  type="text"
                  value={generatedKeys.subject}
                  readOnly
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg bg-gray-50"
                  placeholder="mailto:your-email@example.com"
                />
                <p className="text-xs text-gray-500 mt-1">
                  Substitua por seu email, ex: mailto:seu-email@example.com
                </p>
              </div>
            </div>

            <Button variant="secondary" onClick={() => setShowKeysAlert(false)}>
              Fechar
            </Button>
          </div>
        </Card>
      )}

      {/* Templates .env */}
      {envTemplates && (
        <>
          <Card title="📄 Template .env Backend">
            <div className="space-y-3">
              <p className="text-sm text-gray-600">
                Copie e cole no arquivo <code className="bg-gray-100 px-1 rounded">.env</code> do backend
              </p>
              <div className="relative">
                <pre className="bg-gray-900 text-green-400 p-4 rounded-lg overflow-x-auto text-sm font-mono">
                  {envTemplates.backend}
                </pre>
                <Button
                  variant="secondary"
                  size="sm"
                  className="absolute top-2 right-2"
                  onClick={() => copyToClipboard(envTemplates.backend, 'Configuração backend')}
                >
                  Copiar
                </Button>
              </div>
            </div>
          </Card>

          <Card title="📄 Template .env.local Frontend">
            <div className="space-y-3">
              <p className="text-sm text-gray-600">
                Copie e cole no arquivo <code className="bg-gray-100 px-1 rounded">frontend/.env.local</code>
              </p>
              <div className="relative">
                <pre className="bg-gray-900 text-green-400 p-4 rounded-lg overflow-x-auto text-sm font-mono">
                  {envTemplates.frontend}
                </pre>
                <Button
                  variant="secondary"
                  size="sm"
                  className="absolute top-2 right-2"
                  onClick={() => copyToClipboard(envTemplates.frontend, 'Configuração frontend')}
                >
                  Copiar
                </Button>
              </div>
            </div>
          </Card>
        </>
      )}

      {/* Endpoints da API */}
      {config && (
        <Card title="📡 Endpoints da API">
          <div className="space-y-3">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Base URL
              </label>
              <code className="block bg-gray-100 px-3 py-2 rounded-lg text-sm">
                {config.backend_url}
              </code>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                WebSocket URL
              </label>
              <code className="block bg-gray-100 px-3 py-2 rounded-lg text-sm">
                {config.websocket_url}
              </code>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Principais Endpoints
              </label>
              <div className="space-y-1">
                {config.api_endpoints.map((endpoint, index) => (
                  <code key={index} className="block bg-gray-50 px-3 py-1 rounded text-sm border">
                    {endpoint}
                  </code>
                ))}
              </div>
            </div>
          </div>
        </Card>
      )}

      {/* Links Úteis */}
      <Card title="📚 Documentação">
        <div className="space-y-3">
          <a
            href="/INTEGRATION.md"
            target="_blank"
            className="block p-3 bg-blue-50 border border-blue-200 rounded-lg hover:bg-blue-100 transition-colors"
          >
            <h3 className="font-bold text-blue-800">📖 Guia de Integração</h3>
            <p className="text-sm text-blue-700 mt-1">
              Tutorial completo de como integrar com aplicações Next.js
            </p>
          </a>

          {config && (
            <a
              href={config.swagger_url}
              target="_blank"
              className="block p-3 bg-green-50 border border-green-200 rounded-lg hover:bg-green-100 transition-colors"
            >
              <h3 className="font-bold text-green-800">📋 Swagger API</h3>
              <p className="text-sm text-green-700 mt-1">
                Documentação interativa da API
              </p>
            </a>
          )}

          <a
            href="/test"
            className="block p-3 bg-purple-50 border border-purple-200 rounded-lg hover:bg-purple-100 transition-colors"
          >
            <h3 className="font-bold text-purple-800">🧪 Modo de Teste</h3>
            <p className="text-sm text-purple-700 mt-1">
              Teste notificações in-app e push diretamente no admin
            </p>
          </a>
        </div>
      </Card>
    </div>
  );
}
