# Notification Admin Panel

Interface administrativa para gerenciar o sistema de notificações.

## Tecnologias

- Next.js 14 (App Router)
- TypeScript
- Tailwind CSS
- WebSocket para notificações em tempo real

## Instalação

```bash
npm install
```

## Executar em Desenvolvimento

```bash
npm run dev
```

Ou usando o Justfile na raiz do projeto:

```bash
just frontend-dev
```

O aplicativo estará disponível em [http://localhost:3000](http://localhost:3000)

## Funcionalidades

### Dashboard
- Monitor de notificações em tempo real via WebSocket
- Conectar usando CPF ou telefone do usuário
- Visualizar notificações recebidas ao vivo

### Grupos
- Criar e gerenciar grupos de usuários
- Adicionar/remover membros dos grupos
- Visualizar membros de cada grupo

### Notificações
- Histórico de todas as notificações
- Buscar por CPF ou telefone
- Marcar como lida
- Deletar notificações

### Enviar
- Enviar para usuário específico (CPF ou telefone)
- Enviar para grupo inteiro
- Broadcast para todos os usuários
- Suporte a notificações in-app, push ou ambas

## Build para Produção

```bash
npm run build
npm start
```

## Variáveis de Ambiente

Crie um arquivo `.env.local` baseado no `.env.local.example`:

```env
# API do backend de notificações
# Para desenvolvimento local: http://localhost:8080
# Para API em produção: https://sua-api-producao.com
NEXT_PUBLIC_API_URL=http://localhost:8080

# Chave VAPID pública para push notifications
# Obtenha esta chave no backend ou gere uma nova em /integrations
NEXT_PUBLIC_VAPID_PUBLIC_KEY=your_vapid_public_key_here
```

### Conectar com API em Produção

Para usar o frontend local conectando à API em produção:

1. Edite o arquivo `.env.local`
2. Altere `NEXT_PUBLIC_API_URL` para a URL da sua API em produção:
   ```env
   NEXT_PUBLIC_API_URL=https://api-producao.exemplo.com
   ```
3. Reinicie o servidor de desenvolvimento:
   ```bash
   npm run dev
   ```

**Nota**: Se a API em produção usar HTTPS, o WebSocket automaticamente usará WSS (WebSocket Secure).

## Estrutura

```
frontend/
├── app/              # Páginas (App Router)
│   ├── page.tsx     # Dashboard
│   ├── groups/      # Gerenciamento de grupos
│   ├── notifications/ # Histórico de notificações
│   └── send/        # Enviar notificações
├── components/      # Componentes reutilizáveis
├── hooks/           # Custom hooks (WebSocket)
├── lib/             # Utilitários (API client)
└── types/           # TypeScript types
```
