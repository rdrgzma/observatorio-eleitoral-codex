# Observatório Eleitoral Brasileiro

Plataforma **em desenvolvimento** para pesquisa quantitativa, exploração de **cenários eleitorais condicionais**, validação retrospectiva e acompanhamento dos dados oficiais de totalização de eleições brasileiras.

> **Aviso metodológico:** cenários simulados não são pesquisas de intenção de voto, resultados oficiais, recomendações eleitorais nem previsões certificadas. Dados observados, hipóteses e resultados simulados devem ser identificados separadamente.

## Escopo

- **Eleições municipais de 2028:** preparação da plataforma para estudos de prefeitos e vereadores.
- **Eleições gerais seguintes:** arquitetura extensível para presidente, governadores, senadores e deputados.
- **Pesquisa metodológica:** estudos sobre mudança eleitoral (`swing`), regressão, modelos territoriais, Monte Carlo, replays e calibração.
- **Apuração:** módulo próprio, separado das simulações pré-eleitorais, condicionado à disponibilidade dos arquivos oficiais.

## Estado atual

**MVP-00 — Fundação técnica e documental.** Não há modelo validado, dados oficiais importados, projeção pública nem ambiente de produção.

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
