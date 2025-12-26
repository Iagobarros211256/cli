# MyCLI

**MyCLI** é uma CLI personalizada para Linux que funciona como um mini-terminal.  
Ela permite abrir apps frequentes rapidamente, rodar qualquer comando Linux e ainda possui aliases e prompt colorido.

---

## **Funcionalidades**
- **Comandos pré-definidos** para abrir apps do dia a dia:
  - `play` → Steam
  - `rec` → OBS Studio
  - `stream` → Steam + OBS
  - `dev` → VS Code + Docker
  - `music` → Spotify (via Flatpak)
- **Aliases curtos** para comandos frequentes:
  - `s` → `play`
  - `r` → `rec`
  - `st` → `stream`
  - `d` → `dev`
  - `m` → `music`
- **Comando flexível** `open <app>` para abrir qualquer app instalado
- Executa **qualquer comando Linux normal**
- Prompt **colorido** e mensagens diferenciadas por tipo:
  - Azul → status / abertura de apps
  - Vermelho → erros
  - Amarelo → avisos
  - Ciano → mensagens de boas-vindas

---

## **Instalação**

1. **Clone ou baixe o código**:
```bash
git clone <repo-url>
cd mycli


 cli git:(main) ✗ go build -o mycli main.go

➜  cli git:(main) ✗ ./mycli

