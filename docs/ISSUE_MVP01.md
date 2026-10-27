# MVP-01 — Dados históricos: contratos, proveniência e integração oficial

Rascunho local de issue: a criação via GitHub CLI retornou 403 em
`api.github.com/graphql`. Nenhum número de issue foi inventado.

Objetivo: cumprir o MVP-01 de `docs/ROADMAP.md`, seguindo `AGENTS.md` e os
protocolos. Dependência: MVP-00 incorporado pelo PR #2.

Primeiro incremento: contrato canônico proposto para totais de prefeito,
manifestos, importação exclusivamente sintética, preservação de originais,
hashes e testes de integridade/temporalidade. Inventário oficial, adaptador TSE,
persistência analítica e pareamento histórico continuam pendentes enquanto
não houver acesso às fontes primárias. Não implementar fórmulas eleitorais.

Aceite integral do marco:

- [ ] Inventário por pleito/cargo/turno, cobertura e licença verificados no TSE.
- [ ] Esquema/dicionário e adaptador versionados conferidos com fonte primária.
- [ ] Download controlado mediante manifesto aprovado, sem grandes bases automáticas.
- [ ] Importação real reexecutável, com hashes, originais e relatório de qualidade.
- [ ] Pareamento auditável, ambiguidades e cobertura por eleitorado documentados.
- [ ] Testes de consistência, duplicidade e disponibilidade temporal executados.
- [ ] PR revisado; nenhum merge automático.

A entrega sintética isolada não encerra este marco nem substitui nenhum desses
critérios que exige fontes oficiais. Evidências do incremento em
`docs/VALIDACAO_MVP01.md`.
