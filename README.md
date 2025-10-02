# Notification Service API

Sistema completo de gerenciamento de notificações em Golang com Gin Framework, suportando notificações em tempo real via WebSocket e Push Notifications para aplicações web.

## Tecnologias

- **Golang 1.21+**
- **Gin** - Framework web
- **GORM** - ORM para PostgreSQL
- **Viper** - Gerenciamento de variáveis de ambiente
- **Gorilla WebSocket** - WebSocket para notificações em tempo real
- **Swagger** - Documentação da API
- **Docker & Docker Compose** - Containerização
- **Just** - Task runner (alternativa ao Make)

## Estrutura do Projeto

```
.
├── cmd/
│   └── server/          # Ponto de entrada da aplicação
├── internal/
│   ├── config/          # Configurações e conexão com banco
│   ├── entity/          # Modelos de dados
│   ├── handler/         # Controllers/Handlers HTTP
│   ├── middleware/      # Middlewares customizados
│   ├── repository/      # Camada de acesso a dados
│   ├── service/         # Lógica de negócio
│   └── websocket/       # Hub e cliente WebSocket
├── pkg/
│   ├── auth/            # Autenticação JWT (parse de tokens)
│   ├── queue/           # Cliente RabbitMQ
│   └── utils/           # Utilitários reutilizáveis
├── docs/                # Documentação Swagger (gerada)
├── .env.example         # Exemplo de variáveis de ambiente
├── docker-compose.yml   # Configuração Docker
├── Dockerfile           # Imagem Docker da aplicação
├── Justfile             # Comandos reutilizáveis
└── README.md

```

## Funcionalidades

### Notificações

- ✅ CRUD completo de notificações
- ✅ Envio para usuário específico (CPF ou telefone)
- ✅ Envio para grupos de usuários
- ✅ Broadcast (todos os usuários)
- ✅ Notificações em tempo real via WebSocket
- ✅ Suporte a Push Notifications
- ✅ Marcação de leitura
- ✅ Histórico de notificações

### Grupos

- ✅ CRUD completo de grupos
- ✅ Gerenciamento de membros (adicionar/remover)
- ✅ Listagem de membros por grupo
- ✅ Envio de notificações para grupos

### WebSocket

- ✅ Conexão em tempo real
- ✅ Múltiplas sessões por usuário
- ✅ Ping/Pong para manter conexão
- ✅ Broadcast de notificações

### Autenticação

- ✅ Middleware JWT para extração de informações do usuário
- ✅ Parse de tokens JWT (sem validação de assinatura)
- ✅ Extração automática de CPF, email, nome, telefone e roles
- ✅ Rotas protegidas com autenticação obrigatória
- ✅ Rotas com autenticação opcional
- ✅ Compatível com tokens IDRio (Keycloak)

### Push Notifications

- ✅ Subscrição de dispositivos
- ✅ Gerenciamento de subscriptions
- ✅ Suporte a VAPID keys

### Integrações

- ✅ Gerador de chaves VAPID integrado
- ✅ Templates .env para backend e frontend
- ✅ Visualização de endpoints da API
- ✅ Copiar configurações com um clique
- ✅ Status de configuração em tempo real

### Filas RabbitMQ

- ✅ Processamento assíncrono de notificações
- ✅ Sistema de retry automático (até 3 tentativas)
- ✅ Dead Letter Queue para mensagens com falha
- ✅ Workers configuráveis para escalabilidade
- ✅ Dashboard de monitoramento em tempo real
- ✅ Controles de gerenciamento (pausar, limpar, purgar)
- ✅ Métricas visuais de capacidade e throughput
- ✅ Alertas automáticos para alto volume

## Instalação

### Pré-requisitos

- Go 1.21+
- PostgreSQL 15+
- Node.js 18+ (para o frontend)
- Docker e Docker Compose (opcional)
- Just (opcional, mas recomendado)

### Instalar Just

```bash
# macOS
brew install just

# Linux
curl --proto '=https' --tlsv1.2 -sSf https://just.systems/install.sh | bash -s -- --to /usr/local/bin

# Windows (via Chocolatey)
choco install just
```

### Configuração

1. Clone o repositório:

```bash
git clone <repository-url>
cd app-notification-core
```

2. Copie o arquivo de exemplo de variáveis de ambiente:

```bash
cp .env.example .env
```

3. Configure as variáveis no arquivo `.env`:

```env
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
SERVER_MODE=debug

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=notification_db
DB_SSLMODE=disable

VAPID_PUBLIC_KEY=your_vapid_public_key_here
VAPID_PRIVATE_KEY=your_vapid_private_key_here
VAPID_SUBJECT=mailto:your-email@example.com
```

4. **Gerar chaves VAPID** para push notifications:

```bash
# Instalar web-push CLI globalmente
npm install -g web-push

# Gerar chaves VAPID
web-push generate-vapid-keys
```

Copie as chaves geradas e atualize o `.env` com os valores:
- `VAPID_PUBLIC_KEY`: Chave pública
- `VAPID_PRIVATE_KEY`: Chave privada
- `VAPID_SUBJECT`: Seu email (ex: mailto:seu-email@example.com)

5. Configure as variáveis do frontend em `frontend/.env.local`:

```bash
cd frontend
cp .env.local.example .env.local
```

Edite o `frontend/.env.local` e adicione a **mesma chave pública VAPID** do backend:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_VAPID_PUBLIC_KEY=<sua_chave_publica_vapid_aqui>
```

### Executar com Docker

```bash
just docker-up
```

Ou manualmente:

```bash
docker-compose up -d
```

### Executar localmente

1. Instale as dependências do backend e frontend:

```bash
just install
just frontend-install
```

2. Inicie o PostgreSQL (se não estiver usando Docker):

```bash
docker run -d \
  --name notification-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=notification_db \
  -p 5432:5432 \
  postgres:15-alpine
```

3. Execute backend e frontend:

**Opção 1: Tudo de uma vez**
```bash
just dev-all
```

**Opção 2: Separadamente**

Terminal 1 (Backend):
```bash
just run
```

Terminal 2 (Frontend):
```bash
just frontend-dev
```

4. Acesse:
- Frontend Admin: http://localhost:3000
- Gerenciamento de Filas: http://localhost:3000/queue
- Integrações: http://localhost:3000/integrations
- Modo de Teste: http://localhost:3000/test
- Backend API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- RabbitMQ Management: http://localhost:15672 (admin / admin123)

## 🧪 Modo de Teste

Para testar notificações in-app e push diretamente no painel admin:

1. Acesse http://localhost:3000/test
2. Digite um identificador de usuário (CPF, telefone ou email)
3. Clique em "Conectar WebSocket"
4. (Opcional) Clique em "Habilitar Push" para testar push notifications
5. Vá para a página "Enviar" em outra aba
6. Envie uma notificação para o identificador configurado
7. Veja a notificação aparecer em tempo real na página de teste!

**Recursos do Modo de Teste:**
- ✅ Monitor WebSocket em tempo real
- ✅ Logs detalhados de debug
- ✅ Teste de push notifications
- ✅ Histórico de notificações recebidas
- ✅ Conexão/desconexão manual

## 📬 Sistema de Filas RabbitMQ

O sistema utiliza RabbitMQ para processamento assíncrono e escalável de notificações, garantindo alta disponibilidade e throughput mesmo com grande volume.

**Acesse:** http://localhost:3000/queue

### Como Funciona

1. **Envio Assíncrono**: Quando uma notificação é criada, ela é **publicada na fila RabbitMQ** em vez de ser processada sincronamente
2. **Workers**: 3 workers (configurável) consomem mensagens da fila e processam as notificações em paralelo
3. **Retry Automático**: Se uma notificação falhar, ela é automaticamente reenfileirada (até 3 tentativas)
4. **Dead Letter Queue**: Após 3 falhas, a mensagem é movida para a DLQ para análise posterior

### Dashboard de Monitoramento

O dashboard fornece:

**📊 Métricas em Tempo Real**
- Mensagens pendentes na fila
- Workers ativos processando
- Mensagens na Dead Letter Queue (falhas)
- Status geral do sistema

**🎨 Visualizações**
- Gráficos de capacidade da fila
- Indicadores visuais de carga (verde/amarelo/vermelho)
- Taxa de workers ativos
- Auto-refresh configurável (1s a 30s)

**⚠️ Alertas Inteligentes**
- Aviso quando fila > 1000 mensagens
- Alerta de mensagens na DLQ
- Sugestões de ação

**🛠️ Controles**
- **Limpar Fila**: Remove todas as mensagens pendentes
- **Auto-Refresh**: Pausa/retoma atualização automática
- **Link direto** para RabbitMQ Management UI

### Configuração

No `.env` do backend:

```env
RABBITMQ_URL=amqp://admin:admin123@localhost:5672/
RABBITMQ_QUEUE_NOTIFICATIONS=notifications
RABBITMQ_WORKERS=3
```

**Variáveis:**
- `RABBITMQ_URL`: Conexão com RabbitMQ
- `RABBITMQ_QUEUE_NOTIFICATIONS`: Nome da fila
- `RABBITMQ_WORKERS`: Número de workers paralelos (recomendado: 3-10)

### RabbitMQ Management

Acesse a interface nativa do RabbitMQ para controles avançados:
- **URL**: http://localhost:15672
- **Usuário**: admin
- **Senha**: admin123

**Recursos avançados:**
- Visualizar mensagens na fila em tempo real
- Configurar dead letter exchanges
- Ajustar políticas de TTL
- Monitorar throughput e latência
- Gerenciar exchanges e bindings

### Vantagens

✅ **Performance**: Processamento paralelo com múltiplos workers
✅ **Confiabilidade**: Retry automático e DLQ para falhas
✅ **Escalabilidade**: Adicione mais workers conforme necessário
✅ **Observabilidade**: Dashboard completo e RabbitMQ Management
✅ **Resiliência**: Mensagens persistentes mesmo se servidor reiniciar

## 🔌 Gerenciamento de Integrações

A página de **Integrações** facilita a configuração de aplicações frontend para se conectar ao sistema de notificações:

**Acesse:** http://localhost:3000/integrations

**Funcionalidades:**

### 1. Geração de Chaves VAPID
- Gere chaves VAPID com um clique (não precisa instalar ferramentas externas!)
- Visualize chaves atuais configuradas no backend
- Status visual indica se as chaves estão configuradas corretamente
- Copie chaves individuais ou templates completos

### 2. Templates de Configuração
- **Backend (.env)**: Template completo com todas as variáveis necessárias
- **Frontend (.env.local)**: Configuração pronta para aplicações Next.js
- Botão "Copiar" em cada template para facilitar o uso
- Chaves VAPID já preenchidas automaticamente

### 3. Informações da API
- Lista de todos os endpoints disponíveis
- URLs do WebSocket e API REST
- Link direto para documentação Swagger
- Link para guia de integração completo

**Como usar:**
1. Acesse a página de Integrações
2. Se não tiver chaves VAPID, clique em "Gerar Novas Chaves VAPID"
3. Copie as chaves geradas e atualize o `.env` do backend
4. Reinicie o servidor backend
5. Copie o template frontend para `frontend/.env.local` da sua aplicação
6. Consulte o guia de integração para implementar o código

## 📚 Documentação para Desenvolvedores

Para integrar o sistema de notificações em sua aplicação Next.js, consulte:

📖 **[INTEGRATION.md](./INTEGRATION.md)** - Guia completo de integração

O guia inclui:
- 🔔 Como conectar ao WebSocket para notificações in-app
- 📲 Como implementar push notifications
- 🔌 Exemplos de código prontos para uso
- 🛠️ Hooks React customizados
- 🐛 Troubleshooting comum

## Comandos Just

### Backend

```bash
just run              # Executar aplicação
just build            # Build da aplicação
just dev              # Modo desenvolvimento com hot-reload
just test             # Executar testes
just test-coverage    # Testes com cobertura
just install          # Instalar dependências do Go
just install-tools    # Instalar ferramentas (swag, air, golangci-lint)
just docker-up        # Subir containers Docker
just docker-down      # Parar containers Docker
just swagger          # Gerar documentação Swagger (auto-instala swag se necessário)
just fmt              # Formatar código
just lint             # Executar linter
just clean            # Limpar arquivos gerados
```

### Frontend

```bash
just frontend-install # Instalar dependências do frontend
just frontend-dev     # Executar frontend em desenvolvimento
just frontend-build   # Build do frontend
just frontend-start   # Executar frontend em produção
```

### Ambos

```bash
just dev-all          # Iniciar backend e frontend juntos
just help             # Listar todos os comandos disponíveis
```

## Endpoints da API

### Health Check

```
GET /health
```

### Grupos

```
POST   /api/v1/groups                    - Criar grupo
GET    /api/v1/groups                    - Listar grupos
GET    /api/v1/groups/:id                - Obter grupo
PUT    /api/v1/groups/:id                - Atualizar grupo
DELETE /api/v1/groups/:id                - Deletar grupo
POST   /api/v1/groups/:id/members        - Adicionar membro
DELETE /api/v1/groups/:id/members/:id    - Remover membro
GET    /api/v1/groups/:id/members        - Listar membros
```

### Notificações

```
POST   /api/v1/notifications                     - Criar notificação
GET    /api/v1/notifications                     - Listar notificações
GET    /api/v1/notifications/me                  - Listar minhas notificações (autenticado) 🔐
GET    /api/v1/notifications/:id                 - Obter notificação
PUT    /api/v1/notifications/:id                 - Atualizar notificação
DELETE /api/v1/notifications/:id                 - Deletar notificação
POST   /api/v1/notifications/:id/read            - Marcar como lida
GET    /api/v1/notifications/cpf/:cpf            - Listar por CPF
GET    /api/v1/notifications/phone/:phone        - Listar por telefone
GET    /api/v1/notifications/email/:email        - Listar por email
POST   /api/v1/notifications/send/user           - Enviar para usuário
POST   /api/v1/notifications/send/group/:id      - Enviar para grupo
POST   /api/v1/notifications/send/batch          - Enviar em lote (múltiplos destinatários)
POST   /api/v1/notifications/send/broadcast      - Broadcast (todos)
```

### WebSocket

```
GET /api/v1/ws?user_id=<cpf_ou_telefone>
```

### Subscriptions (Push)

```
POST   /api/v1/subscriptions    - Criar subscription
DELETE /api/v1/subscriptions    - Deletar subscription
```

### Documentação Swagger

```
GET /swagger/index.html
```

## Exemplos de Uso

### Criar Grupo

```bash
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Premium Users",
    "description": "Usuários premium da plataforma"
  }'
```

### Adicionar Membro ao Grupo

```bash
curl -X POST http://localhost:8080/api/v1/groups/{group_id}/members \
  -H "Content-Type: application/json" \
  -d '{
    "cpf": "12345678901",
    "phone": "11999999999",
    "email": "usuario@exemplo.com",
    "name": "João Silva"
  }'
```

### Enviar Notificação para Usuário Específico

```bash
curl -X POST http://localhost:8080/api/v1/notifications/send/user \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Nova Mensagem",
    "message": "Você recebeu uma nova mensagem!",
    "type": "in-app",
    "cpf": "12345678901",
    "data": {
      "action": "open_message",
      "message_id": "123"
    }
  }'
```

### Buscar Minhas Notificações (Autenticado)

Este endpoint requer um token JWT válido no header `Authorization`.

```bash
curl -X GET "http://localhost:8080/api/v1/notifications/me?limit=20&offset=0" \
  -H "Authorization: Bearer <seu_token_jwt>" \
  -H "Content-Type: application/json"
```

**Resposta:**
```json
{
  "user": {
    "cpf": "12345678901",
    "email": "usuario@exemplo.com",
    "name": "João Silva",
    "email_verified": true
  },
  "notifications": [
    {
      "id": "uuid-here",
      "title": "Nova Mensagem",
      "message": "Você recebeu uma nova mensagem!",
      "type": "in-app",
      "read": false,
      "created_at": "2024-01-01T12:00:00Z"
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "count": 1
  }
}
```

### Enviar Notificação para Grupo

```bash
curl -X POST http://localhost:8080/api/v1/notifications/send/group/{group_id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Atualização",
    "message": "Nova funcionalidade disponível!",
    "type": "both"
  }'
```

### Broadcast

```bash
curl -X POST http://localhost:8080/api/v1/notifications/send/broadcast \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Manutenção Programada",
    "message": "Sistema em manutenção das 02h às 04h",
    "type": "in-app"
  }'
```

### Conectar via WebSocket (JavaScript)

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?user_id=12345678901');

ws.onopen = () => {
  console.log('Conectado ao servidor de notificações');
};

ws.onmessage = (event) => {
  const notification = JSON.parse(event.data);
  console.log('Nova notificação:', notification);

  // Exibir notificação
  showNotification(notification.title, notification.message);
};

ws.onerror = (error) => {
  console.error('Erro WebSocket:', error);
};

ws.onclose = () => {
  console.log('Desconectado');
};
```

### Push Notification Subscription (Next.js)

```javascript
// Subscribe to push notifications
async function subscribeToPush() {
  const registration = await navigator.serviceWorker.register('/sw.js');

  const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: 'YOUR_VAPID_PUBLIC_KEY'
  });

  await fetch('http://localhost:8080/api/v1/subscriptions', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      user_cpf: '12345678901',
      endpoint: subscription.endpoint,
      p256dh: btoa(String.fromCharCode(...new Uint8Array(subscription.getKey('p256dh')))),
      auth: btoa(String.fromCharCode(...new Uint8Array(subscription.getKey('auth'))))
    })
  });
}
```

## Tipos de Notificação

- `in-app`: Apenas notificações no aplicativo (WebSocket)
- `push`: Apenas Push Notifications
- `email`: Apenas Email (via Data Relay)
- `both`: In-app + Push
- `all`: In-app + Push + Email

## Status de Notificação

- `pending`: Pendente
- `sent`: Enviada
- `delivered`: Entregue
- `read`: Lida
- `failed`: Falha no envio

## Integração com Next.js

Para integrar com sua aplicação Next.js, consulte a documentação completa em:

📖 **[INTEGRATION.md](./INTEGRATION.md)**

Essa documentação contém:
- Guia passo a passo de integração
- Hooks React customizados (useNotifications)
- Implementação completa de push notifications
- Service Worker pronto para uso
- Exemplos de código completos
- Troubleshooting

## Testes

```bash
just test
```

Para cobertura de testes:

```bash
just test-coverage
```

## Documentação da API

Após iniciar o servidor, acesse:

```
http://localhost:8080/swagger/index.html
```

Para regenerar a documentação Swagger:

```bash
just swagger
```

## 🔐 Autenticação JWT

O sistema possui um módulo completo de autenticação JWT que **extrai informações do token sem validar a assinatura**. A validação RBAC e autenticação é feita por outra aplicação (ex: Keycloak/IDRio).

### Como Funciona

1. O cliente envia um token JWT no header `Authorization: Bearer <token>`
2. O middleware `auth.JWTMiddleware()` extrai as informações do payload
3. As informações do usuário ficam disponíveis no contexto da requisição
4. O handler pode acessar CPF, email, nome, telefone e roles do usuário

### Informações Extraídas

```go
type UserInfo struct {
    CPF           string   // Campo "preferred_username" do token
    Email         string
    Name          string
    Phone         string
    Roles         []string // Roles do realm_access
    EmailVerified bool
    Sub           string   // ID único do usuário
}
```

### Uso Básico

#### 1. Proteger uma rota

```go
// Requer autenticação
notifications.GET("/me", auth.RequireAuth(), handler.GetMyNotifications)

// Autenticação opcional
notifications.GET("/public", auth.OptionalJWTMiddleware(), handler.List)
```

#### 2. Extrair informações no handler

```go
func (h *Handler) GetMyNotifications(c *gin.Context) {
    userInfo, exists := auth.GetUserInfo(c)
    if !exists {
        c.JSON(401, gin.H{"error": "Unauthorized"})
        return
    }

    cpf := userInfo.CPF
    email := userInfo.Email
    name := userInfo.Name

    notifications, _ := h.service.GetNotificationsByCPF(cpf, 20, 0)
    c.JSON(200, notifications)
}
```

#### 3. Parse direto de token

```go
import "github.com/fzolio/app-notification-core/pkg/auth"

// Parse completo
userInfo, err := auth.ParseToken(token)

// Extrair apenas CPF
cpf, err := auth.ExtractCPF(token)

// Extrair apenas email
email, err := auth.ExtractEmail(token)
```

### Endpoints Protegidos

Atualmente, apenas o endpoint `/api/v1/notifications/me` requer autenticação:

```bash
curl -X GET "http://localhost:8080/api/v1/notifications/me" \
  -H "Authorization: Bearer <seu_token_jwt>"
```

### Formato do Token

O token esperado segue o padrão Keycloak/IDRio:

- **Header**: `Authorization: Bearer <token>`
- **CPF**: Extraído do campo `preferred_username`
- **Roles**: Extraído de `realm_access.roles`

### Documentação Completa

Veja a documentação completa em [`pkg/auth/README.md`](pkg/auth/README.md) com exemplos detalhados de uso.

## 🐛 Troubleshooting

### WebSocket não conecta (Erro 400)

Se você vê o erro `400 Bad Request` ao tentar conectar no WebSocket:

**Causa**: O parâmetro de query está incorreto.

**Solução**: Certifique-se de usar `user_id` como parâmetro:
```javascript
// ✅ Correto
ws://localhost:8080/api/v1/ws?user_id=12345678901

// ❌ Incorreto
ws://localhost:8080/api/v1/ws?user=12345678901
```

### Como Push Notifications Funcionam?

**IMPORTANTE**: Push notifications NÃO aparecem na interface do admin - elas aparecem como **notificações nativas do sistema operacional** (igual WhatsApp, Gmail, etc).

**Fluxo completo:**
1. Usuário clica em "Habilitar Push" na página de teste
2. Navegador solicita permissão para mostrar notificações
3. Service Worker é registrado no navegador
4. Subscription é criada e enviada para o backend
5. Backend armazena a subscription no banco de dados
6. Quando você envia uma notificação com tipo `push`, `both` ou `all`:
   - Backend busca todas as subscriptions do usuário
   - Envia push notification via Web Push API para cada subscription
   - Service Worker do navegador recebe o push
   - **Notificação aparece como notificação do sistema operacional**

**Onde as notificações aparecem:**
- 🪟 **Windows**: Canto inferior direito (Action Center)
- 🍎 **macOS**: Canto superior direito (Notification Center)
- 🐧 **Linux**: Depende do DE (geralmente canto superior direito)
- 📱 **Mobile**: Barra de status

**Como verificar se está funcionando:**
1. Abra os logs do backend e procure por:
   - `Found X subscription(s), sending push notifications...`
   - `Push sent successfully to subscription...`
2. No navegador, abra DevTools > Console para ver mensagens do Service Worker
3. Verifique se a permissão de notificações está concedida
4. **Importante**: Notificações NÃO aparecem se a aba do navegador estiver em foco - minimize ou troque de aba

### Push Notifications - InvalidCharacterError

Se você vê o erro `InvalidCharacterError: Failed to execute 'atob'`:

**Causa**: A chave VAPID pública não está configurada ou está inválida.

**Solução**:
1. Gere chaves VAPID se ainda não tiver:
   ```bash
   npm install -g web-push
   web-push generate-vapid-keys
   ```

2. Configure a chave pública no `.env` do backend:
   ```env
   VAPID_PUBLIC_KEY=sua_chave_publica_aqui
   ```

3. Configure a **mesma chave pública** no `frontend/.env.local`:
   ```env
   NEXT_PUBLIC_VAPID_PUBLIC_KEY=sua_chave_publica_aqui
   ```

4. Reinicie o frontend para carregar as novas variáveis:
   ```bash
   just frontend-dev
   ```

### Backend não envia emails

Se as notificações com canal `email` ou `all` não estão sendo enviadas:

**Solução**: Verifique as configurações do Data Relay no `.env`:
```env
DATA_RELAY_API_URL=https://data-relay.dados.rio/
DATA_RELAY_API_TOKEN=seu_token_aqui
```

### Service Worker não registra

Se o Service Worker não está sendo registrado:

**Causa**: Navegador requer HTTPS ou localhost.

**Solução**:
- Use `localhost` para desenvolvimento (já é o padrão)
- Em produção, certifique-se de que o site está em HTTPS

### Notificações não aparecem no modo de teste

**Soluções**:
1. Verifique se o WebSocket está conectado (indicador verde)
2. Certifique-se de que o identificador está correto
3. Envie a notificação para o mesmo identificador configurado no teste
4. Verifique o console do navegador para erros
5. Confira os logs de debug na página de teste

## Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## Licença

Este projeto está sob a licença MIT.