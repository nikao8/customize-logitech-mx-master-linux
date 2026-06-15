# Logitech MX Master Configuration GUI

Interface gráfica em Go com [Fyne](https://fyne.io/) para gerar e gerenciar a configuração do [logiops](https://github.com/PixlOne/logiops) (logid) para o mouse Logitech MX Master.

Gera o arquivo `/etc/logid.cfg` e instala o serviço systemd para executar o daemon `logid`.

## Dependências

### Debian / Ubuntu

```bash
sudo apt install golang-go libgl1-mesa-dev xorg-dev libxrandr-dev
```

### Arch Linux

```bash
sudo pacman -S go
```

Pacote AUR para logiops:
```bash
yay -S logiops        # ou logiops-git
```

### Fedora

```bash
sudo dnf install golang libXcursor-devel libXrandr-devel mesa-libGL-devel
```

## Build

```bash
# Com suporte gráfico (requer X11/GL dev headers)
make build-release

# Ou com fallback headless (CI tag)
make build

# Apenas verificar compilação
go build -tags ci ./...
```

O binário gerado é `logid-config-gui`.

## Instalação

```bash
# Instalar o binário em /usr/local/bin
sudo make install

# Ou copiar manualmente
sudo cp logid-config-gui /usr/local/bin/
```

## Uso

1. Descubra o nome exato do seu dispositivo:

```bash
sudo logid -v
```

2. Execute o GUI:

```bash
./logid-config-gui
```

3. Na aba **General**, ajuste:
   - **Device Name**: insira o nome exato mostrado pelo `logid -v`
   - **DPI**: sensibilidade do cursor (200-4000)
   - **SmartShift**: liga/desliga e threshold da roda
   - **Hi-Res Scrolling**: scroll de alta resolução

4. Na aba **Buttons**, configure cada botão:
   - Botão Gesture (polegar): suporta gestos direcionais
   - Botão SmartShift toggle (abaixo da roda)
   - Botões Back/Forward laterais
   - Botões esquerdo, direito e meio

5. Na aba **Thumbwheel**, configure a roda lateral:
   - **Divert**: habilitar manipulação pelo logid
   - Ações para esquerda, direita e toque

6. Na aba **Service**:
   - **Save Configuration**: salva `/etc/logid.cfg`
   - **Install & Start Service**: salva config, instala unit systemd e ativa o daemon

## Mapeamento dos botões

| CID     | Botão                    |
|---------|--------------------------|
| `0xc3`  | Gesto (polegar)          |
| `0xc4`  | Alternar SmartShift      |
| `0x53`  | Voltar (lateral cima)    |
| `0x56`  | Avançar (lateral baixo)  |
| `0x52`  | Botão do meio            |
| `0x50`  | Botão esquerdo           |
| `0x51`  | Botão direito            |

## Tipos de ação

| Tipo              | Descrição                           |
|-------------------|-------------------------------------|
| `None`            | Botão desativado                   |
| `Keypress`        | Pressiona tecla(s) ex: `KEY_LEFTCTRL, KEY_T` |
| `Gestures`        | Gestos direcionais (Up/Down/Left/Right) |
| `ToggleSmartShift` | Alterna modo SmartShift da roda    |
| `ToggleHiresScroll` | Alterna scroll de alta resolução   |
| `CycleDPI`        | Alterna entre valores de DPI       |
| `ChangeDPI`       | Incrementa/decrementa DPI          |
| `ChangeHost`      | Alterna host Bluetooth             |

## Idiomas

Interface disponível em **English** e **Português (Brasil)**, selecionável no topo da janela.

## Estrutura do projeto

```
├── app.go        # Interface Fyne
├── config.go     # Modelo e gerador de /etc/logid.cfg
├── mapping.go    # Mapas de CID, keycodes, traduções
├── service.go    # Gerenciamento do serviço systemd
├── main.go       # Ponto de entrada
├── Makefile      # Targets de build
└── README.md     # Este arquivo
```

## Referências

- [logiops - GitHub](https://github.com/PixlOne/logiops)
- [logiops Configuration Wiki](https://github.com/PixlOne/logiops/wiki/Configuration)
- [logiops CIDs](https://github.com/PixlOne/logiops/wiki/CIDs)
- [Arch Wiki - Logitech MX Master](https://wiki.archlinux.org/title/Logitech_MX_Master)
- [Linux input-event-codes.h](https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h)
