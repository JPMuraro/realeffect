# RealEffect DSL — Core Engine v0.1

RealEffect DSL (`.reff`) é uma linguagem declarativa baseada em YAML para descrever **missões de impacto real**  
(meio ambiente, educação, saúde, etc.) com regras imutáveis de validação.

---

## Objetivos

- **Humana e legível** — arquivos `.reff` são YAML.
- **Determinística** — mesma entrada → mesmo resultado.
- **Auditável** — regras claras, explícitas, revisáveis.
- **Neutra** — foca em provas concretas, sem ideologia.

---

## Núcleo da linguagem (v0.1)

Regras imutáveis implementadas no engine:

- **RE-0 (Entrega total)**  
  100% dos slots de evidência devem ter prova entregue.  
  → se faltar qualquer slot obrigatório para qualquer participante, a missão é inválida.

- **RE-1 (80% de peso aceito)**  
  Pelo menos **80% do peso total das evidências** deve ser aceito (`ACCEPTED`).  
  → se o peso aceito for menor que 0.8, a missão falha.

Outras proteções:

- Peso máximo por slot: **≤ 0.40**
- Pelo menos **3 slots relevantes** (peso normalizado ≥ 0.10)

---

## Estrutura do projeto

```text
realeffect/
  cmd/
    realeffectc/        # CLI: avalia um arquivo .reff direto do disco
      main.go
    realeffectd/        # Daemon HTTP: expõe o engine via /evaluate
      main.go
  internal/
    realeffect/         # Núcleo da linguagem (biblioteca Go)
      spec.go           # Tipos da missão (MissionSpec, Context, EvidenceSlot, etc.)
      engine.go         # Regras RE-0, RE-1 e validações estruturais
      engine_test.go    # Testes unitários do core
      client.go         # Cliente HTTP oficial para falar com o realeffectd
  examples/
    plant_100_trees.reff      # Missão válida (plantio de 100 árvores)
    invalid_weights.reff      # Exemplo de peso inválido (> 0.40)
    invalid_relevance.reff    # Exemplo com pouca diversidade de provas
  REALEFFECT_SPEC_v0.1.md     # Especificação formal da linguagem (detalhe técnico)
  go.mod
  go.sum
