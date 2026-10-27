# AGENTS.md — Observatório Eleitoral Brasileiro

## Missão
Construir plataforma pública, neutra, reproduzível e auditável para análise histórica e simulações condicionais de eleições brasileiras. Operação pré-eleitoral e projeção durante apuração são fluxos **distintos**.

## Regras inegociáveis
1. Não inventar dados do TSE, pesquisas, fontes, resultados de testes ou interpretações legais.
2. Classificar as informações em `DATA_FACT`, `METHOD_DOCUMENTED`, `METHOD_EXPERIMENTAL`, `HYPOTHESIS` e `OUTPUT_SIMULATED`.
3. Não atribuir a Rodrigo Mahlmeister fórmulas ou parâmetros não verificados em fontes primárias.
4. Bloquear vazamento temporal: cada experimento registra um `cutoff` e apenas utiliza dados disponíveis até ele.
5. Preservar fonte, URL, hash, horário de coleta, versão do parser, versão do modelo, parâmetros e semente.
6. Dados históricos brutos são imutáveis; operações de normalização produzem novos artefatos.
7. Verificar invariantes eleitorais e registrar exceções justificadas: votos não negativos, votos válidos <= comparecimento <= eleitorado, e proporções válidas.
8. Cada mudança matemática exige documentação (ADR), testes e revisão.
9. Não apresentar simulação como pesquisa, apuração oficial ou previsão certificada; não produzir recomendação eleitoral.
10. Nunca enviar dados pessoais de eleitores, credenciais ou segredos a prompts.
11. Exigir revisão estatística e jurídica antes de qualquer publicação externa de simulações.
12. Começar por baselines e dados sintéticos rotulados; não afirmar que infraestrutura ou testes funcionam sem efetiva execução.

## Stack pretendida
- Go: coleta e API;
- Python: estatística, replays e Monte Carlo;
- PostgreSQL: entidades e metadados;
- Parquet: bases analíticas;
- Laravel 13 + Filament 5: painel;
- Docker Compose: desenvolvimento;
- GitHub + Codex: implementação assistida e revisão.

## Fluxo de trabalho do Codex
- Ler este arquivo, README e protocolos pertinentes antes de modificar.
- Trabalhar por issue e branch específica, com PR revisável; não alterar `main` diretamente nas tarefas seguintes.
- Fornecer plano sucinto, arquivos modificados, justificativas, premissas e riscos.
- Incluir testes de unidade, bordas, integridade e temporalidade; relatar os comandos realmente executados.
- Não implementar fórmulas ainda não aprovadas; documentar pendências como `A_CONFIRMAR`.
- Não ativar publicação pública de projeções por padrão.
- Se uma fonte estiver indisponível, parar a implementação dependente e registrar a limitação.
