chatear-backend/
│
├── application/              # Casos de uso genéricos / nível cross-domain
├── cmd/                      # Pontos de entrada (main.go, workers, CLI)
├── config/                   # Leitura e binding de configs/env
├── domain/                   # Entidades de domínio base (opcional se já dentro de internal/*)
├── docs/                     # Documentação
├── graph/                    # Schema e resolvers GraphQL
├── infrastructure/           # Implementações técnicas globais (DB, Redis, NATS, SMTP, etc.)
├── internal/
│   ├── user/
│   │   ├── application/      # Casos de uso do domínio user
│   │   ├── domain/           # Entidades + regras de negócio puras
│   │   ├── infrastructure/   # Repositórios / adaptadores (Supabase, Redis, etc)
│   │   └── presentation/     # GraphQL/HTTP/CLI handlers específicos
│   └── notification/
│       ├── application/
│       ├── domain/
│       ├── infrastructure/
│       ├── worker/           # Event consumers (NATS)
│       └── example_usage.go  # Teste de integração ou demonstração
├── pkg/                      # Utilidades compartilhadas e genéricas
├── shared/                   # Tipos comuns (ex: DTOs, erros, eventos)
├── presentation/              # Interfaces globais (GraphQL root)
└── migrations/                # Banco

