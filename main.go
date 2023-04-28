package main

import (
	"bufio"
	"fmt"
	"github.com/kasugamirai/my-cli-app/chatgpt"
	"os"
	"time"
)

func handleInput(prompt string) chan string {
	output := make(chan string)
	go chatgpt.ChatWithGPT(prompt, output)
	return output
}

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("error: OPENAI_API_KEY environment variable not set. use \"export OPENAI_API_KEY=\"your openAI API KEY\" to set\"")
		os.Exit(1)
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter your question and press Enter to submit. Type 'exit' to quit the program.")

InputLoop:
	for {
		fmt.Print("> ")
		scanner.Scan()
		prompt := scanner.Text()

		if prompt == "exit" {
			break
		}

		output := handleInput(prompt)

		// Loop to print continuous updates from the output channel
		for {
			select {
			case value, ok := <-output:
				if ok {
					fmt.Print(value)
				} else {
					// The channel is closed, break the loop
					fmt.Println()
					goto InputLoop
				}
			case <-time.After(60 * time.Second):
				// Timeout: no new values received for 1 second, break the loop
				goto InputLoop
			}
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
		os.Exit(1)
	}
}
