# Tarefa inicial do Codex — MVP-00

Você é engenheiro de software responsável pela fundação técnica do Observatório Eleitoral Brasileiro.

Leia `AGENTS.md`, `README.md` e todos os arquivos `docs/` antes de alterar código.

**Objetivo:** propor e implementar o menor ambiente local reproduzível para Go (ingestão/API), Python (motor estatístico ainda vazio), Laravel 13/Filament 5 (painel futuro), PostgreSQL e Docker Compose.

**Limites:**
- Não implementar projeção eleitoral, Monte Carlo nem inferir fórmulas de Rodrigo Mahlmeister.
- Não baixar automaticamente grandes bases sem manifesto e autorização explícita.
- Não introduzir dependências pagas, segredos ou dados pessoais.
- Não marcar algo como validado sem rodar os testes correspondentes.

**Etapas:**
1. Inspecionar o repositório e elaborar plano curto com riscos.
2. Abrir branch `feat/mvp00-foundation`.
3. Criar estrutura de diretórios, configuração de ambiente, scripts de verificação e serviços mínimos.
4. Registrar decisões relevantes em ADR.
5. Incluir testes de inicialização/healthcheck e lint onde aplicável.
6. Atualizar README com instruções verificadas e executar os comandos disponíveis.
7. Abrir PR com checklist, evidências e limitações. Não fazer merge automático.

**Aceite:** serviços que existirem iniciam localmente; os testes são reproduzíveis e sua saída é documentada. Ausência de ferramentas ou ambiente deve ser explicitada, nunca mascarada.
