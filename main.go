package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\033[36mBem-vindo ao MyCLI! Digite 'exit' para sair.\033[0m") // ciano

	for {
		// Prompt colorido em verde
		fmt.Print("\033[32mMyCLI> \033[0m")
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

		// Comandos pré-definidos
		switch cmdName {
		case "play":
			fmt.Println("\033[34mAbrindo Steam...\033[0m")
			exec.Command("steam").Start()
		case "rec":
			fmt.Println("\033[34mAbrindo OBS Studio...\033[0m")
			exec.Command("obs").Start()
		case "stream":
			apps := []string{"steam", "obs"}
			for _, app := range apps {
				fmt.Println("\033[34mAbrindo", app, "...\033[0m")
				exec.Command(app).Start()
			}
		case "dev":
			apps := []string{"code", "docker"}
			for _, app := range apps {
				fmt.Println("\033[34mAbrindo", app, "...\033[0m")
				exec.Command(app).Start()
			}
		case "music":
			fmt.Println("\033[34mAbrindo Spotify...\033[0m")
			// Flatpak Spotify
			err := exec.Command("flatpak", "run", "com.spotify.Client").Start()
			if err != nil {
				fmt.Println("\033[31mErro ao abrir Spotify:", err, "\033[0m")
			}
		case "open":
			if len(args) == 0 {
				fmt.Println("\033[33mUse: open <nome_do_app>\033[0m") // amarelo
				continue
			}
			app := args[0]
			fmt.Println("\033[34mAbrindo", app, "...\033[0m")
			err := exec.Command(app).Start()
			if err != nil {
				fmt.Println("\033[31mErro ao abrir", app+":", err, "\033[0m")
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
			}
		}
	}
}
