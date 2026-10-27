# Evidências MVP-00 — issue #1

Execução real em 2026-10-27 no ambiente Codex Linux x86_64. Classificação:
`METHOD_DOCUMENTED` (verificação de infraestrutura, sem resultados eleitorais).

## Resultados finais

| Verificação executada | Resultado |
| --- | --- |
| `python3 scripts/init-env.py` | `.env` local criado; valores omitidos dos logs |
| `docker compose config --quiet` | configuração válida |
| `./scripts/smoke.sh` com override local de rede descrito abaixo | builds concluídos; três serviços saudáveis; `PASS /healthz`, `PASS /`, `PASS /readyz` |
| `./scripts/check.sh` | exit 0; todas as verificações abaixo executadas |
| unittest em `tests/` | 5 testes passaram: chave, permissões, preservação de configuração e integridade/idempotência dos arquivos Composer |
| Ruff 0.14.2 | passou para scripts, testes e motor Python |
| gofmt / `go vet ./...` / `go test -race ./...` | passaram; um teste com três subtestes GET/POST/DELETE; pacote main sem testes unitários, exercitado pelo smoke test |
| unittest do motor Python | 2 testes passaram, incluindo execução real do comando |
| `composer validate --no-check-publish` | válido; advertência de licença indefinida, já prevista no README |
| `composer check-platform-reqs` | todos os requisitos instalados atendidos |
| Pint | 20 arquivos passaram |
| `vendor/bin/phpunit --fail-on-warning` | 5 testes passaram, 10 asserções, sem warnings |
| PostgreSQL `SELECT 1 AS readiness` via `psql` | uma linha, valor 1 |
| `sh -n scripts/check.sh scripts/smoke.sh` / `git diff --check` | passaram |

PHP executado: 8.4.27; Laravel instalado: 13.39.0; Filament: 5.11.0.
Versões exatas PHP estão em `painel/composer.lock`. O painel roda como
`www-data`; a API roda como UID 65534. Os healthchecks e a requisição `/readyz`
foram executados nos contêineres reais, não somente em mocks.

## Diagnóstico e limites do ambiente

- PHP/Composer não existem no host; foram executados em imagens Docker.
  O `/usr/bin/go` do host não é o compilador Go esperado; a imagem oficial
  `golang:1.26-bookworm` foi usada para build, lint e testes.
- O build sem override falhou na resolução DNS externa. Para este ambiente,
  foi criado somente em `/tmp/observatorio-compose-cloud.json` um override de
  Compose com `build.extra_hosts` para o proxy existente, encaminhamento das
  variáveis `http_proxy`/`https_proxy` e BuildKit secret `custom_ca` apontando
  ao bundle confiável `/etc/ssl/certs/ca-certificates.crt`. Host e IP do proxy
  foram derivados da configuração existente, sem gravar credenciais.
  Comandos efetivamente usados:

  ```sh
  COMPOSE_FILE="$PWD/compose.yaml:/tmp/observatorio-compose-cloud.json" \
    DOCKER_CONFIG=/tmp/observatorio-docker ./scripts/smoke.sh
  DOCKER_CONFIG=/tmp/observatorio-docker ./scripts/check.sh
  ```

- A API do GitHub respondeu `403 Forbidden`. Pacotes foram instalados pelos
  arquivos oficiais de `codeload.github.com`, preservando commits fixados e
  verificações TLS. O normalizador não altera checksums existentes.
- Uma tentativa de instalação por Git source foi interrompida porque copiava
  históricos desnecessários. A instalação final foi por arquivos ZIP oficiais.
- O teste inicial encontrou permissões restritas herdadas dos arquivos locais.
  O Dockerfile foi corrigido para permitir leitura pelo usuário não privilegiado,
  mantendo escrita em `storage` e `bootstrap/cache`. O smoke test passou após a
  correção. Caches temporários de builds foram limpos após esgotarem o disco;
  nenhum volume de banco ou arquivo de usuário foi removido.
- `php artisan test` inicialmente gerou avisos do wrapper sobre ausência do
  `.env` dentro da imagem (variáveis são injetadas pelo Compose). O runner salvo
  usa PHPUnit diretamente, com `--fail-on-warning`; os cinco testes executaram
  sem avisos. Não foi necessário copiar segredos para a imagem.
- A CI foi adicionada, mas seu resultado remoto não foi verificado neste
  ambiente. Os mesmos scripts foram executados localmente conforme acima.
- Não há testes temporais/eleitorais neste marco, porque não existe ingestão,
  dataset ou modelo. Nenhuma formulação estatística foi implementada; futuras
  escolhas matemáticas seguem `A_CONFIRMAR` e exigem documentação e revisão.

Filament está instalado e descoberto pelo Laravel, mas o painel administrativo
permanece futuro: não há rotas administrativas, autenticação ou usuários de
demonstração. Não houve implantação pública nem publicação de simulações.
