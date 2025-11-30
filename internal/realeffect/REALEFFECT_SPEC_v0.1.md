# RealEffect DSL — Especificação v0.1

## 1. Visão geral

RealEffect DSL (`.reff`) é uma linguagem declarativa baseada em YAML para descrever **missões de impacto real** (meio ambiente, saúde, educação, etc.) com regras imutáveis de validação:

- **Humana e legível** (YAML).
- **Determinística**: mesma entrada → mesmo resultado.
- **Auditável**: regras explícitas, revisáveis e documentadas.
- **Neutra**: sem ideologia; foca em provas concretas.

O núcleo da linguagem é definido por duas camadas:

- **RE-0 (Entrega)**: 100% das provas exigidas devem ser entregues antes de qualquer avaliação.
- **RE-1 (Aceitação)**: pelo menos **80% do peso total das provas** deve ser aceito para que a missão seja considerada validada.

---

## 2. Estrutura básica de um arquivo `.reff`

```yaml
spec_version: "0.1"
mission_id: "string-única"

context:
  title: "Título da missão"
  description: "Descrição legível"
  location: "Local ou região"
  timeframe_start: "YYYY-MM-DD"
  timeframe_end: "YYYY-MM-DD"
  owner: "id_da_DAO_ou_organização"

participants:
  - role: "nome_do_papel"
    min: 1
    max: 10        # opcional

evidence_slots:
  - id: "slug_da_prova"
    description: "Descrição da prova exigida"
    category: "tipo_de_prova"
    weight: 0.2
    required: true
