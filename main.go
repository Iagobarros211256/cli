package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Função para cores por comando
func getCommandColor(cmd string) string {
	switch cmd {
	case "play", "s":
		return "\033[32m" // verde
	case "rec", "r":
		return "\033[35m" // magenta
	case "stream", "st":
		return "\033[34m" // azul
	case "dev", "d":
		return "\033[36m" // ciano
	case "music", "m":
		return "\033[33m" // amarelo
	default:
		return "\033[0m" // padrão
	}
}

// Função para enviar notificações de desktop
func sendNotification(title, message string) {
	exec.Command("notify-send", title, message).Start()
}

// Função para obter CPU e RAM (simplificado)
func getSystemStats() (cpuUsage string, ramUsage string) {
	// CPU
	cpuCmd := exec.Command("bash", "-c", "top -bn1 | grep 'Cpu(s)' | awk '{print $2+$4}'")
	cpuOut, err := cpuCmd.Output()
	if err != nil {
		cpuUsage = "N/A"
	} else {
		cpuUsage = strings.TrimSpace(string(cpuOut)) + "%"
	}

	// RAM
	ramCmd := exec.Command("bash", "-c", "free -h | awk '/Mem:/ {print $3\"/\"$2}'")
	ramOut, err := ramCmd.Output()
	if err != nil {
		ramUsage = "N/A"
	} else {
		ramUsage = strings.TrimSpace(string(ramOut))
	}

	return
}

// Função para obter uso de GPU (Nvidia) (amd) (intel-caso de falha)
func getGPUUsage() string {
	// Detectar GPU
	out, err := exec.Command("bash", "-c", "lspci | grep -E 'VGA|3D'").Output()
	if err != nil {
		return "N/A"
	}

	gpuInfo := strings.ToLower(string(out))

	// Nvidia
	if strings.Contains(gpuInfo, "nvidia") {
		cmd := exec.Command("bash", "-c", "nvidia-smi --query-gpu=utilization.gpu --format=csv,noheader,nounits")
		out, err := cmd.Output()
		if err != nil {
			return "N/A"
		}
		return strings.TrimSpace(string(out)) + "%"
	}

	// AMD
	if strings.Contains(gpuInfo, "amd") || strings.Contains(gpuInfo, "radeon") {
		// Usando radeontop em modo batch
		cmd := exec.Command("bash", "-c", "radeontop -d - | head -n1 | awk '{print $2}'")
		out, err := cmd.Output()
		if err != nil {
			return "N/A"
		}
		return strings.TrimSpace(string(out)) + "%"
	}

	// Se Intel ou outra GPU
	return "N/A"
}

// Função para gerar prompt dinâmico com dashboard
func dynamicPrompt(lastCmd string) string {
	now := time.Now()
	timeStr := now.Format("15:04")

	dir, err := os.Getwd()
	if err != nil {
		dir = "unknown"
	}

	color := getCommandColor(lastCmd)
	cpu, ram := getSystemStats()
	gpu := getGPUUsage()

	return fmt.Sprintf("%s[%s | %s | CPU:%s | RAM:%s | GPU:%s] MyCLI> \033[0m",
		color, timeStr, dir, cpu, ram, gpu)
}

// Função para log
func logCommand(cmd string, success bool) {
	f, err := os.OpenFile(os.Getenv("HOME")+"/.mycli_history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	status := "OK"
	if !success {
		status = "ERROR"
	}
	logLine := fmt.Sprintf("[%s] %s : %s\n", time.Now().Format("2006-01-02 15:04:05"), cmd, status)
	f.WriteString(logLine)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	lastCmd := ""

	fmt.Println("\033[36mBem-vindo ao MyCLI! Digite 'exit' para sair.\033[0m") // ciano

	for {
		fmt.Print(dynamicPrompt(lastCmd))
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			fmt.Println("\033[31mSaindo...\033[0m") // vermelho
			break
		}

		if input == "" {
			continue
		}

		parts := strings.Fields(input)
		cmdName := parts[0]
		args := parts[1:]

		// Aliases
		switch cmdName {
		case "s":
			cmdName = "play"
		case "r":
			cmdName = "rec"
		case "st":
			cmdName = "stream"
		case "d":
			cmdName = "dev"
		case "m":
			cmdName = "music"
		}

		lastCmd = cmdName // para o prompt dinâmico e cores

		success := true // para logs

		// Comandos pré-definidos
		switch cmdName {
		case "play":
			fmt.Println(getCommandColor(cmdName), "Abrindo Steam...\033[0m")
			err := exec.Command("setsid", "steam").Start()
			if err != nil {
				fmt.Println("\033[31mErro ao abrir Steam:", err, "\033[0m")
				sendNotification("MyCLI - Erro", "Falha ao abrir Steam!")
				success = false
			} else {
				sendNotification("MyCLI", "Steam abriu com sucesso!")
			}

		case "rec":
			fmt.Println(getCommandColor(cmdName), "Abrindo OBS Studio...\033[0m")
			err := exec.Command("setsid", "obs").Start()
			if err != nil {
				fmt.Println("\033[31mErro ao abrir OBS:", err, "\033[0m")
				sendNotification("MyCLI - Erro", "Falha ao abrir OBS Studio!")
				success = false
			} else {
				sendNotification("MyCLI", "OBS Studio abriu com sucesso!")
			}

		case "stream":
			apps := []string{"steam", "obs"}
			for _, app := range apps {
				fmt.Println(getCommandColor(cmdName), "Abrindo", app, "...\033[0m")
				err := exec.Command(app).Start()
				if err != nil {
					fmt.Println("\033[31mErro ao abrir", app+":", err, "\033[0m")
					sendNotification("MyCLI - Erro", "Falha ao abrir "+app)
					success = false
				} else {
					sendNotification("MyCLI", app+" abriu com sucesso!")
				}
			}

		case "dev":
			apps := []string{"code", "docker"}
			for _, app := range apps {
				fmt.Println(getCommandColor(cmdName), "Abrindo", app, "...\033[0m")
				err := exec.Command(app).Start()
				if err != nil {
					fmt.Println("\033[31mErro ao abrir", app+":", err, "\033[0m")
					sendNotification("MyCLI - Erro", "Falha ao abrir "+app)
					success = false
				} else {
					sendNotification("MyCLI", app+" abriu com sucesso!")
				}
			}

		case "music":
			fmt.Println(getCommandColor(cmdName), "Abrindo Spotify...\033[0m")
			err := exec.Command("setsid", "flatpak", "run", "com.spotify.Client").Start()
			if err != nil {
				fmt.Println("\033[31mErro ao abrir Spotify:", err, "\033[0m")
				sendNotification("MyCLI - Erro", "Falha ao abrir Spotify")
				success = false
			} else {
				sendNotification("MyCLI", "Spotify abriu com sucesso!")
			}

		case "open":
			if len(args) == 0 {
				fmt.Println("\033[33mUse: open <nome_do_app>\033[0m") // amarelo
				success = false
				continue
			}
			app := args[0]
			fmt.Println(getCommandColor(cmdName), "Abrindo", app, "...\033[0m")
			err := exec.Command(app).Start()
			if err != nil {
				fmt.Println("\033[31mErro ao abrir", app+":", err, "\033[0m")
				sendNotification("MyCLI - Erro", "Falha ao abrir "+app)
				success = false
			} else {
				sendNotification("MyCLI", app+" abriu com sucesso!")
			}
		case "log":
			logPath := os.Getenv("HOME") + "/.mycli_history.log"

			fmt.Println("\033[36m--- MyCLI Logs ---\033[0m")

			cmd := exec.Command("cat", logPath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Println("\033[31mErro ao ler o log\033[0m")
				success = false
			}

		default:
			// Qualquer outro comando Linux
			cmd := exec.Command("bash", "-c", input)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin
			err := cmd.Run()
			if err != nil {
				fmt.Println("\033[31mErro ao executar comando:", err, "\033[0m")
				success = false
			}
		}

		// Log
		logCommand(input, success)
	}
}
