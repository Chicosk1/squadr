# mobile — Kotlin Multiplatform + Compose Multiplatform

App do Squadr para Android e iOS, com **UI compartilhada** em Compose
Multiplatform. Desenvolvimento primário no **Android Studio**.

## Estrutura

```
mobile/
├── androidApp/     # casca fina: Activity, manifest, ícone
├── iosApp/         # casca fina: projeto Xcode, AppDelegate
└── shared/src/
    ├── commonMain/kotlin/com/squadr/
    │   ├── data/
    │   │   ├── remote/       # cliente HTTP e WebSocket contra o backend Go
    │   │   ├── local/        # cache e preferências
    │   │   └── repository/   # implementações das interfaces de domain
    │   ├── domain/
    │   │   ├── model/        # tipos do domínio, sem framework
    │   │   ├── repository/   # interfaces
    │   │   └── usecase/      # regra que faz sentido no cliente
    │   ├── ui/               # telas em Compose Multiplatform
    │   │   ├── login/  matching/  chat/  profile/  groups/  theme/
    │   ├── di/               # composição de dependências
    │   └── platform/         # declarações expect
    ├── androidMain/kotlin/com/squadr/platform/   # actual (Android)
    ├── iosMain/kotlin/com/squadr/platform/       # actual (iOS)
    └── commonMain/composeResources/              # imagens, fontes, strings
```

Critério de "o que entra em cada pasta":
[`docs/context/arquitetura.md`](../docs/context/arquitetura.md), seção 4.

## Três regras que mantêm o projeto multiplataforma de verdade

1. **`androidApp/` e `iosApp/` são cascas finas.** Se você está escrevendo uma
   tela dentro de uma delas, pare — ela deveria estar em `shared/.../ui/`. Nas
   cascas só entra o que o sistema operacional exige (manifest, `Info.plist`,
   ciclo de vida, permissões).
2. **`expect`/`actual` só sob `platform/`.** Espalhar isso pelo código é o
   caminho mais curto para o projeto virar dois projetos. Casos legítimos:
   armazenamento seguro do token, KMPNotifier, abrir navegador para o OAuth.
3. **Regra de negócio não mora aqui.** Matching, bloqueio, reputação e limites
   são do backend Go. `domain/usecase/` é só para o que é genuinamente de
   cliente: validação de formulário, formatação, ordenação local.

## Rodando

> ⚠️ O projeto Gradle ainda não foi gerado (Fase 2 do roadmap). Até então, os
> comandos abaixo não funcionam.

- **Android:** abra `mobile/` no Android Studio e rode a configuração
  `androidApp` no emulador.
- **iOS:** exige macOS com Xcode. Sem Mac, o build sai pelo Codemagic (Fase 6) e
  o teste local fica limitado a dispositivo físico.

```bash
./gradlew :shared:build
```

```bash
./gradlew :shared:allTests
```

## Comunicação com o backend

Só `data/remote/` conhece URL de API. O token de acesso é injetado
automaticamente no header `Authorization`, e o cliente renova o token uma vez em
caso de 401 (sem entrar em loop).

O contrato dos endpoints é [`contracts/openapi.yaml`](../contracts/openapi.yaml).
Se o app precisa de um campo que não está lá, o contrato muda **primeiro**.
