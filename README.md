# Observatório Eleitoral Brasileiro

Plataforma **em desenvolvimento** para pesquisa quantitativa, exploração de **cenários eleitorais condicionais**, validação retrospectiva e acompanhamento dos dados oficiais de totalização de eleições brasileiras.

> **Aviso metodológico:** cenários simulados não são pesquisas de intenção de voto, resultados oficiais, recomendações eleitorais nem previsões certificadas. Dados observados, hipóteses e resultados simulados devem ser identificados separadamente.

## Escopo

- **Eleições municipais de 2028:** preparação da plataforma para estudos de prefeitos e vereadores.
- **Eleições gerais seguintes:** arquitetura extensível para presidente, governadores, senadores e deputados.
- **Pesquisa metodológica:** estudos sobre mudança eleitoral (`swing`), regressão, modelos territoriais, Monte Carlo, replays e calibração.
- **Apuração:** módulo próprio, separado das simulações pré-eleitorais, condicionado à disponibilidade dos arquivos oficiais.

## Estado atual

**MVP-01 — primeiro incremento sintético, integração oficial bloqueada.** A fundação MVP-00 foi incorporada pelo PR #2. Não há modelo validado, dados oficiais importados, projeção pública nem ambiente de produção.

## Arquitetura planejada

| Camada | Tecnologias previstas | Responsabilidade |
|---|---|---|
| Ingestão/API | Go | Aquisição e normalização das fontes oficiais |
| Estatística | Python | Baselines, regressões, simulações e validação |
| Painel | Laravel + Filament | Administração e consultas |
| Banco | PostgreSQL | Metadados, entidades e execuções |
| Arquivos analíticos | Parquet | Dados históricos e resultados derivados |
| Desenvolvimento | Docker Compose, GitHub, Codex | Reprodutibilidade e revisão |

As ferramentas são uma proposta inicial, sujeita a decisões de arquitetura registradas.

## Princípios

1. **Fontes verificáveis:** priorizar dados publicados pelo TSE; preservar URI, data de obtenção, licença, versão e hash.
2. **Sem vazamento temporal:** em replays, não utilizar informações posteriores ao ponto reconstruído.
3. **Separação documental:** distinguir `DATA_FACT`, `METHOD_DOCUMENTED`, `METHOD_EXPERIMENTAL`, `HYPOTHESIS` e `OUTPUT_SIMULATED`.
4. **Integridade:** validar conservação dos votos, comparecimento, eleitorado e consistência de percentuais.
5. **Pesquisa auditável:** versionar parâmetros, sementes, ambientes, artefatos e métricas.
6. **Neutralidade:** não oferecer orientação de voto ou promover candidaturas; representar dados e limitações fielmente.
7. **Revisão humana:** o Codex apoia o desenvolvimento; decisões estatísticas e publicações requerem aprovação de responsáveis.

## Estrutura de trabalho

Consulte [AGENTS.md](AGENTS.md) para regras do Codex, [docs/](docs/) para protocolos e [docs/ROADMAP.md](docs/ROADMAP.md) para as entregas.

## Fontes oficiais

- [Dados Abertos do TSE](https://dadosabertos.tse.jus.br/)
- [Especificações dos arquivos eleitorais](https://www.tse.jus.br/eleicoes/arquivos)
- [Informações técnicas sobre divulgação da totalização](https://www.tse.jus.br/eleicoes/informacoes-tecnicas-sobre-a-divulgacao-de-resultados)
- [Normas dos sistemas eleitorais](https://www.tse.jus.br/legislacao/compilada/res/2021/resolucao-no-23-677-de-16-de-dezembro-de-2021)

## Uso de IA

O Codex deve ler **AGENTS.md** e os protocolos antes de modificar o código. Fórmulas não encontradas no material de Rodrigo Mahlmeister **não** devem ser apresentadas como reprodução de seu método.

**Licença de software:** a definir pelos mantenedores antes de disponibilizar versões de distribuição.

## Desenvolvimento local — MVP-00

Pré-requisitos: Docker Engine com Compose v2 e Python 3.12+. Go, PHP e
Composer rodam nos contêineres. Os comandos abaixo partem da raiz do checkout.

```sh
python3 scripts/init-env.py
./scripts/smoke.sh
./scripts/check.sh
```

O primeiro comando cria `.env` ignorado pelo Git, com senha explicitamente local
e chave Laravel aleatória; ele se recusa a sobrescrever um arquivo existente.
Se você já tem `.env`, preserve seus valores e acrescente as variáveis ausentes
de `.env.example`, incluindo `APP_KEY` (32 bytes aleatórios em base64 com prefixo
`base64:`). Não use credenciais reais de serviços externos para este marco.

O smoke test constrói/inicia API Go, PostgreSQL e Laravel e verifica respostas
HTTP e conexão real do painel ao banco. `check.sh` executa gofmt, go vet, testes
Go com detector de corrida, Ruff, unittest Python, validação Composer, Pint e
testes Laravel. A CI repete os dois scripts. O motor Python é um comando sob o
perfil `tools`, sem serviço estatístico permanente:

```sh
docker compose --profile tools run --rm engine
# Logs e encerramento, preservando o volume do PostgreSQL:
docker compose logs --tail=100
docker compose down
```

API: porta local `8080`, rota `/healthz`. Laravel: porta local `8000`, rotas `/`
(estado da fundação), `/up` (processo) e `/readyz` (consulta PostgreSQL). Os valores
`API_PORT` e `PANEL_PORT` podem ser ajustados no `.env`. Banco não publica porta
no host. Não há painel administrativo público: Filament 5 está instalado para
trabalho futuro, sem usuários, recursos ou autenticação fictícios.

Estrutura: `services/api/`, `services/engine/`, `painel/`, `scripts/`,
`docs/adr/`. Há importação local exclusivamente sintética no MVP-01, com testes
de invariantes e cutoff; não há ingestão oficial nem modelo eleitoral. Consulte as decisões e
limites em [ADR 0001](docs/adr/0001-fundacao-local.md) e os resultados efetivamente
executados em [Validação MVP-00](docs/VALIDACAO_MVP00.md).

### Rede e ambientes corporativos

Builds precisam de Docker Hub, Debian, Packagist e GitHub. Mantenha TLS e
assinaturas APT ativos. O Dockerfile aceita um certificado adicional opcional
por BuildKit secret `custom_ca`; encaminhe os proxies já configurados usando
`--build-arg http_proxy --build-arg https_proxy`. Não grave tokens nos argumentos.
O lockfile usa arquivos oficiais de `codeload.github.com`, com os mesmos commits
fixados. Após um `composer update`, se a API do GitHub estiver bloqueada, execute
`python3 scripts/normalize_composer_dist.py` antes de instalar. O utilitário não
altera versões ou checksums e recusa referências inconsistentes. Resolução de
DNS/proxy deve funcionar dentro dos contêineres, além do host.

O servidor PHP embutido e as credenciais de exemplo são exclusivos do ambiente
local. Não há implantação nem publicação pública ativada por estes comandos.


## MVP-01 — Importação local sintética

O [inventário de fontes](docs/data/INVENTARIO_FONTES.md) registra o bloqueio HTTP
403 do TSE. O [contrato interno](docs/data/CONTRATO_HISTORICO.md) não representa
um layout oficial. O importador aceita somente CSV local `SYNTHETIC` rotulado
`HYPOTHESIS`, verifica manifesto e hash e impede normalização de dados inválidos.
Não baixa arquivos nem realiza pareamento histórico ou imputação.

Exemplo reproduzível com a fixture versionada (três seções inteiramente fictícias):

```sh
# Execute na raiz; use um novo diretório de saída em cada execução.
docker run --rm --user "$(id -u):$(id -g)" \
  -e GOCACHE=/tmp/go-cache -e GOPATH=/tmp/go -e GOTELEMETRY=off \
  -v "$PWD:/repo" -w /repo/services/api golang:1.26-bookworm sh -c '
  go build -o /tmp/import-canonical ./cmd/import-canonical &&
  /tmp/import-canonical \
    --input testdata/historical/sections.csv \
    --manifest testdata/historical/manifest.json \
    --cutoff 2026-10-27T19:22:18Z \
    --output /repo/data/processed/mvp01-synthetic-1'
```

Resultado esperado: `SYNTHETIC: rows=3 accepted=3 rejected=0`. O cutoff do
exemplo é a disponibilidade declarada **da fixture**, não de dados eleitorais de
2024. As saídas permanecem ignoradas pelo Git. Reexecutar com o mesmo destino
falha sem sobrescrever artefatos. `quality_passed` não autoriza publicação.

Os testes do novo módulo e do CLI já fazem parte de `./scripts/check.sh` e da CI.
Veja [ADR 0002](docs/adr/0002-contrato-historico-sintetico.md),
[evidências do incremento](docs/VALIDACAO_MVP01.md) e o
[rascunho da issue MVP-01](docs/ISSUE_MVP01.md). O aceite integral do marco
continua pendente de fontes oficiais, licença, adaptador verificado e
pareamento auditável com cobertura medida.
