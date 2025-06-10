package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/TylerBrock/colorjson"
	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/google/shlex"
	"github.com/parashmaity/fleare-cli/comm"
)

func HandleCommand(conn *Connection) error {

	// Create a channel to listen for system signals
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	home, _ := os.UserHomeDir()

	historyFile := home + "/.fleare_history" // File to store command history
	underline := color.New(color.FgCyan, color.Underline).SprintFunc()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      fmt.Sprint(underline(conn.Conn.RemoteAddr().String(), ">"), " "),
		HistoryFile: historyFile, // Enable history
	})
	if err != nil {
		fmt.Println("Error setting up readline:", err)
		return err
	}

	defer rl.Close()

	for {

		line, err := rl.Readline()
		if err != nil {
			if err.Error() == "Interrupt" {
				fmt.Println("\nUse 'exit' to close the client.")
				// continue
			}
			fmt.Println("Error reading line:", err)
			break
		}
		line = strings.TrimSpace(line)
		if len(line) == 0 || line == "\n" {
			continue
		}
		line = strings.TrimSpace(line)

		// Check if the user wants to exit
		if line == "exit" {
			fmt.Println("Exiting client...")
			conn.Close() // Gracefully close the connection
			os.Exit(0)
		}

		args, err := shlex.Split(line)
		if err != nil {
			fmt.Printf("Error parsing input: %v\n", err)
			continue
		}

		cmd := &comm.Command{
			Command: args[0],
			Args:    args[1:],
		}

		if err := conn.Write(cmd); err != nil {
			fmt.Println("Failed to send command: %w", err)
			continue
		}

		var resp comm.Response
		if err = conn.Read(&resp); err != nil {
			continue
		}

		if resp.Status == "Ok" {
			printSuccess(&resp)
			continue
		} else {
			printError(&resp)
		}

	}

	return nil
}

func printError(resp *comm.Response) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("\n%s %s\n\n", red(resp.Status), string(resp.Result))
}

func printSuccess(resp *comm.Response) {
	green := color.New(color.FgGreen).SprintFunc()

	var obj map[string]any
	if err := json.Unmarshal([]byte(string(resp.Result)), &obj); err != nil {
		fmt.Printf("\n%s %s\n\n", green(resp.Status), string(resp.Result))
		return
	}

	f := colorjson.NewFormatter()
	f.Indent = 2
	s, _ := f.Marshal(obj)

	fmt.Printf("\n%s %s\n\n", green(resp.Status), string(s))
}
