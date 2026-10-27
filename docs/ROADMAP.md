# Roadmap — Observatório Eleitoral

## MVP-00: Fundação
- [ ] Incorporar a documentação científica e técnica e revisar as regras do Codex.
- [ ] Registrar ADR da estrutura inicial e convenções de versionamento.
- [ ] Criar estrutura mínima Go/Python/Laravel sem implementar algoritmos eleitorais.
- [ ] Preparar Docker Compose, `.env.example`, `.gitignore`, lint e CI básico.
- [ ] Criar testes de saúde e documentar comandos de execução.
**Aceite:** ambiente reproduzível; testes executados com resultados reais, sem dados ou resultados inventados.

## MVP-01: Dados históricos
- [ ] Inventário oficial de fontes TSE por pleito/cargo/turno.
- [ ] Esquema canônico, dicionário e parser versionado.
- [ ] Download controlado, manifestos, hash e licenças.
- [ ] Pareamento histórico auditável, sem imputação silenciosa.
- [ ] Testes de consistência, cobertura e duplicidade.
**Aceite:** importação reexecutável com qualidade medida.

## MVP-02: Laboratório retrospectivo
- [ ] Baseline de agregação e modelo territorial documentados.
- [ ] Replays temporais (somente quando houver snapshots) e sintéticos, rotulados.
- [ ] Métricas de erro, cobertura, viés e análise por estágio e território.
**Aceite:** evidências reexecutáveis de desempenho, sem vazamento de dados futuros.

## MVP-03: Simulações condicionais
- [ ] Definir cenários e distribuições justificadas em ADR.
- [ ] Monte Carlo reproduzível, com seed, sensibilidade e convergência.
- [ ] Separar comparecimento, votos válidos e participações.
**Aceite:** testes estatísticos e documentação de hipóteses.

## MVP-04: Painel e operação
- [ ] Interface que distingue dados oficiais, cenários simulados e experimentos.
- [ ] Integrar coleta da totalização sem misturar dados de versões diferentes.
- [ ] Validar desempenho e revisão legal antes de disponibilizar publicamente.
**Aceite:** auditoria de cada execução e apresentação sem falsa precisão.

## MVP-05: Outros cargos
- [ ] Modelo majoritário e modalidades do Senado.
- [ ] Modelos proporcionais com regras legais versionadas por pleito.
- [ ] Validação retrospectiva por circunscrição e cargo.

**Horizonte:** municipais de 2028 e gerais de 2030, sem cronograma fixado nesta etapa.
