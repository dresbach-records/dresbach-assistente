 licença 
---

## 🔐 Regras de Negócio Fundamentais

- Nunca informar preço de projeto sem diagnóstico
- Diagnóstico Tech Ops é sempre pago
- Apenas uma pergunta por mensagem
- Linguagem humana, profissional e clara
- Fluxos Tech Ops e Área do Cliente são totalmente independentes
- Pagamentos ocorrem apenas no fluxo Tech Ops

---

## ⚙️ Stack Tecnológica

- Go (Golang)
- WhatsApp Cloud API (Meta)
- Webhooks HTTP
- Integração com:
  - Gateways de pagamento
  - Agenda (calendário)
  - WHM / cPanel
  - Sistemas internos

---

## 📁 Estrutura do Projeto (sugerida)

# Dresbach Assistente

Assistente oficial da **Dresbach Hosting do Brasil LTDA**, desenvolvido em **Go**, operando via **WhatsApp Cloud API**, com fluxos inteligentes de atendimento, Tech Ops (consultoria especializada) e Área do Cliente de Hospedagem.

---

## 📌 Visão Geral

O **Dresbach Assistente** é um backend conversacional orientado a **máquina de estados**, projetado para operar **100% via WhatsApp**, sem dependência de frontend web.

Ele centraliza:
- Atendimento automatizado
- Qualificação de leads
- Consultoria técnica e jurídica (Tech Ops)
- Cobrança e validação de pagamentos
- Agendamento automático
- Transferência para operadores humanos
- Área do Cliente de Hospedagem no WhatsApp

---

## 🧠 Principais Módulos

### 🔹 1. Tech Ops (Consultoria Especializada)
Fluxo premium voltado a:
- Diagnóstico técnico
- Segurança da informação
- Arquitetura de sistemas
- LGPD e governança digital

Funciona como funil:
Qualificação → Diagnóstico pago → Pagamento → Agendamento → Humano

---

### 🔹 2. Área do Cliente (Hospedagem)
Área de autoatendimento via WhatsApp:
- Login seguro (CPF/CNPJ + senha + 2FA)
- Serviços ativos
- Acesso a cPanel / Webmail / WHM
- Domínios e DNS
- Tickets de suporte
- Faturas (visualização)

> ⚠️ Importante:  
> **Nenhum pagamento é realizado dentro da Área do Cliente.**

---

## 🧱 Arquitetura

- Linguagem: **Go (Golang)**
- Modelo: **State Machine**
- Comunicação: **Webhooks (WhatsApp Cloud API)**
- Backend: **Stateless**
- Integrações externas desacopladas
