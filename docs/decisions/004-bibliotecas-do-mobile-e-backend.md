# 004 — Bibliotecas do mobile e do backend

## Contexto

Seguindo o item 5.3 do roadmap, decidindo cliente HTTP e DI do mobile, e roteador
HTTP e biblioteca de WebSocket do backend.

## Decisões

- **Cliente HTTP (mobile): Ktor Client.** Fuel foi descartado por ter suporte a iOS
  ainda em estágio alfa (3.0.0-alpha04) — risco alto para a camada de rede do app.
  Ktor é mantido pela JetBrains, nativo em Kotlin Multiplatform, com engines maduros
  para Android (OkHttp) e iOS (Darwin).

- **Injeção de dependência (mobile): Koin.** Existe uma opção mais nova, Metro (DI
  compilado), mas com comunidade e documentação ainda pequenas — escolhida Koin pela
  maturidade, integração direta com Compose Multiplatform (`koinViewModel`,
  `koinInject`) e menor risco combinado com o resto da stack, que já é nova.
  Dificuldade esperada: Koin usa resolução em runtime (service locator), exige
  disciplina para não virar um "container global" desorganizado conforme o app cresce.

- **Roteador HTTP (backend): chi.** Desde o Go 1.22 o chi passou a usar o roteamento
  nativo do `net/http` por baixo dos panos — a escolha aqui é por ergonomia de
  agrupamento de rotas e middleware, não por performance. Mantém 100% de
  compatibilidade com `http.Handler`, sem lock-in de framework.

- **Biblioteca de WebSocket (backend): coder/websocket.** gorilla/websocket está
  arquivado. coder/websocket é o substituto recomendado pela comunidade Go: suporte
  nativo a `context.Context`, escrita concorrente seguro por padrão, mantido
  ativamente.

## Dificuldades esperadas

- Koin exige atenção para não deixar o grafo de dependências virar acoplamento
  implícito difícil de rastrear — revisar periodicamente se algo deveria ser
  passado explicitamente em vez de injetado.
- Reavaliar Metro como possível migração futura, quando o time (ou você) já tiver
  base sólida em KMP.