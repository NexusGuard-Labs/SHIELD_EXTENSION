# 🛡️ SHIELD_EXTENSION

<div align="center">

![Shield Badge](https://img.shields.io/badge/status-em%20desenvolvimento-orange)
![Security](https://img.shields.io/badge/seguran%C3%A7a-an%C3%A1lise%20de%20extens%C3%B5es-blue)
![Focus](https://img.shields.io/badge/foco-malware%20%26%20dados-red)

</div>

> Repositório dedicado à análise de extensões e à verificação de comportamentos potencialmente maliciosos.

## 📌 Visão geral

O SHIELD_EXTENSION foi pensado como um ambiente para investigar extensões de navegador e verificar se o código-fonte inclui componentes suspeitos, coleta indevida de dados ou ações que possam comprometer a privacidade e a segurança do usuário.

A proposta central do projeto é detectar sinais de malware, abuso de permissões e exfiltração de informação, com foco em análise técnica e rastreio de risco.

## 🎯 Objetivo

Este repositório busca:

- avaliar o código-fonte de extensões;
- verificar a presença de scripts ou comportamentos suspeitos;
- identificar coleta de dados sem necessidade clara;
- detectar comunicação com servidores externos;
- registrar evidências de risco e comportamento anômalo;
- apoiar investigações de segurança digital.

## 🔍 O que a análise procura

- permissões excessivas ou desnecessárias;
- chamadas a APIs e domínios suspeitos;
- scripts obfuscados ou pouco transparentes;
- acesso a dados do navegador ou do usuário;
- comportamento anômalo em execução;
- transporte de informação para servidores externos sem justificativa clara.

## 🧱 Estrutura do repositório

```text
SHIELD_EXTENSION/
├── README.md
├── docs/
│   ├── analise-estatica.md
│   ├── analise-dinamica.md
│   └── checklist-seguranca.md
├── evidencias/
│   ├── logs/
│   ├── capturas/
│   └── relatorios/
├── src/
│   └── codigo-analisado/
└── reports/
    └── relatorio-final.md
```

## 🛠️ Fluxo de trabalho

1. Definir a extensão a ser analisada;
2. revisar o código-fonte e os manifestos;
3. verificar permissões e chamadas externas;
4. observar comportamento em execução;
5. registrar evidências e riscos;
6. produzir um relatório final com conclusão técnica.

## 📚 Documentação

- [Análise estática](docs/analise-estatica.md)
- [Análise dinâmica](docs/analise-dinamica.md)
- [Checklist de segurança](docs/checklist-seguranca.md)

## 🚀 Como usar o analisador automático

O projeto agora conta com uma ferramenta otimizada em **Go** para automatizar a análise estática das extensões. O uso de Go permite realizar a verredura concorrente de milhares de arquivos JavaScript de forma extremamente rápida.

1. Certifique-se de ter o [Go](https://golang.org/) instalado.
2. Compile a ferramenta:
   ```bash
   go build -o shield_analyzer main.go
   ```
3. Execute passando o link do repositório da extensão ou o caminho local:
   ```bash
   ./shield_analyzer https://github.com/usuario/extensao-exemplo
   ```
4. A ferramenta irá baixar o código para `src/codigo-analisado/`, realizar a análise concorrente e gerar um relatório automático na pasta `reports/`.

## 🚨 Importante

Este projeto é voltado para análise e segurança, e deve ser usado apenas em contextos legítimos, éticos e autorizados, com foco em prevenção, investigação e redução de risco.

## ✅ Status

Projeto em desenvolvimento, com foco em organização, documentação e análise técnica de extensões.

---

<p align="center">
  <strong>SHIELD_EXTENSION</strong> — proteção, análise e investigação de extensões.
</p>

## Desenvolvimento local após o Codespaces

Abra a pasta `SHIELD_EXTENSION` no VS Code. Os caminhos do projeto são relativos à raiz e não dependem do Codespaces. No momento, o repositório contém guias e modelos de análise; ainda não há uma extensão executável, dependências para instalar ou testes automatizados.

Coloque o código a analisar em `src/codigo-analisado/`. A configuração de depuração fica vazia até existir um programa com ponto de entrada definido. Extensões de navegador devem ser inspecionadas no navegador. Os guias pessoais de criação do repositório e de iniciação ficam apenas no computador e são ignorados pelo Git.

### Enviar alterações locais para o GitHub

Execute no terminal, dentro desta pasta:

```bash
git status
git diff
git add .
git diff --cached
git commit -m "Atualiza arquivos do projeto local"
git push origin HEAD:refs/heads/main
```

Esse fluxo publica os arquivos locais. No VS Code, use **Push / Enviar**; **Sincronizar alterações** também pode trazer alterações do remoto. Se o push for rejeitado por divergência, preserve o trabalho local em um commit e examine o histórico remoto antes de decidir como reconciliar. Não use `git reset --hard` nem `git push --force` como rotina.

