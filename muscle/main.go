package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"sort"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

// --- DATA STRUCTURES ---

type BrainRequest struct {
	SessionID string   `json:"session_id"`
	Command   string   `json:"command"`
	Cwd       string   `json:"cwd"`
	History   []string `json:"history"`
}

type BrainResponse struct {
	Output string `json:"output"`
}

type SessionState struct {
	CurrentDir string
	History    []string
}

func main() {
	// 1. SSH Server Config
	config := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
			log.Printf("⚠️  Login: User=%s Pass=%s", c.User(), string(pass))
			return nil, nil // Allow everyone
		},
	}

	// 2. Load Keys
	privateBytes, err := ioutil.ReadFile("hostkey")
	if err != nil {
		log.Fatal("Failed to load hostkey: ", err)
	}
	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	config.AddHostKey(private)

	// 3. Listen
	listener, err := net.Listen("tcp", "0.0.0.0:2222")
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}
	log.Println("🚀 GhostShell Muscle (Instant Mode) listening on 2222...")

	for {
		nConn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(nConn, config)
	}
}

func handleConnection(nConn net.Conn, config *ssh.ServerConfig) {
	_, chans, reqs, err := ssh.NewServerConn(nConn, config)
	if err != nil {
		return
	}
	go ssh.DiscardRequests(reqs)

	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				req.Reply(req.Type == "shell" || req.Type == "pty-req", nil)
			}
		}(requests)

		go runSmartShell(channel)
	}
}

func runSmartShell(channel ssh.Channel) {
	defer channel.Close()

	// Use terminal wrapper for correct line editing (backspace, arrows)
	term := term.NewTerminal(channel, "")
	
	state := &SessionState{
		CurrentDir: "/root",
		History:    []string{},
	}

	// Fake Welcome Message
	term.Write([]byte("Welcome to Ubuntu 22.04.3 LTS (GNU/Linux 5.15.0-91-generic x86_64)\r\n"))
	term.Write([]byte(" * Documentation:  https://help.ubuntu.com\r\n"))
	term.Write([]byte(" * Management:     https://landscape.canonical.com\r\n"))
	term.Write([]byte(" * Support:        https://ubuntu.com/advantage\r\n\r\n"))
	term.Write([]byte("System information as of " + state.CurrentDir + "\r\n\r\n"))

	for {
		// Set Prompt
		prompt := fmt.Sprintf("root@server:%s# ", state.CurrentDir)
		term.SetPrompt(prompt)

		line, err := term.ReadLine()
		if err != nil {
			break
		}

		rawCmd := strings.TrimSpace(line)
		if rawCmd == "" {
			continue
		}
		if rawCmd == "exit" {
			break
		}

		// Update History
		state.History = append(state.History, rawCmd)
		if len(state.History) > 10 {
			state.History = state.History[1:]
		}

		// --- INSTANT MODE: LOCAL HANDLING ---

		// 1. Handle 'pwd'
		if rawCmd == "pwd" {
			term.Write([]byte(state.CurrentDir + "\r\n"))
			continue
		}

		// 2. Handle 'ls' (Read from fakeFS)
		if rawCmd == "ls" || rawCmd == "ll" || rawCmd == "ls -la" {
			files, exists := fakeFS[state.CurrentDir]
			if !exists {
				// If dir isn't in our map, imply it's empty or fallback to AI
				term.Write([]byte("\r\n")) 
			} else {
				// Sort and print files like columns
				sort.Strings(files)
				output := strings.Join(files, "  ")
				term.Write([]byte(output + "\r\n"))
			}
			continue
		}

		// 3. Handle 'cd' (Update state instantly)
		if strings.HasPrefix(rawCmd, "cd ") {
			target := strings.TrimSpace(strings.TrimPrefix(rawCmd, "cd "))
			
			// Resolve new path
			newDir := resolvePath(state.CurrentDir, target)

			// Check if it exists in our FakeFS
			if _, ok := fakeFS[newDir]; ok {
				state.CurrentDir = newDir
			} else {
				term.Write([]byte("-bash: cd: " + target + ": No such file or directory\r\n"))
			}
			continue
		}

		// --- SLOW MODE: ASK THE BRAIN ---
		// If it wasn't ls/cd/pwd, we send it to Python/Gemini
		
		go func() { // Run in background so we don't block? (Actually we want to block for output)
			// Construct Request
			reqBody, _ := json.Marshal(BrainRequest{
				SessionID: "session-123",
				Command:   rawCmd,
				Cwd:       state.CurrentDir,
				History:   state.History,
			})

			// Call Python
			resp, err := http.Post("http://localhost:5000/hallucinate", "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				term.Write([]byte("Error: Connection to Brain failed.\r\n"))
				return
			}
			defer resp.Body.Close()

			var result BrainResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return
			}

			// Format and print
			cleanOutput := strings.ReplaceAll(result.Output, "\n", "\r\n")
			if cleanOutput != "" {
				term.Write([]byte(cleanOutput + "\r\n"))
			}
			
			// Re-print prompt is handled by next loop iteration
		}()
		
		// Wait for the AI response before looping back to prompt?
		// In a simple loop, yes. The 'go func' above is actually risky if we want synchronous output.
		// Let's do it synchronously for now to keep the prompt tidy.
		
		reqBody, _ := json.Marshal(BrainRequest{
			SessionID: "session-123",
			Command:   rawCmd,
			Cwd:       state.CurrentDir,
			History:   state.History,
		})

		resp, err := http.Post("http://localhost:5000/hallucinate", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			term.Write([]byte("Remote Error.\r\n"))
		} else {
			var result BrainResponse
			json.NewDecoder(resp.Body).Decode(&result)
			resp.Body.Close()
			cleanOutput := strings.ReplaceAll(result.Output, "\n", "\r\n")
			if cleanOutput != "" {
				term.Write([]byte(cleanOutput + "\r\n"))
			}
		}
	}
}

// Helper to handle "cd .." and relative paths
func resolvePath(current, target string) string {
	if target == "/" {
		return "/"
	}
	if target == ".." {
		// Move up one level
		if current == "/" {
			return "/"
		}
		lastSlash := strings.LastIndex(current, "/")
		if lastSlash == 0 {
			return "/" // Parent of /root is /
		}
		return current[:lastSlash]
	}
	if strings.HasPrefix(target, "/") {
		return target // Absolute path
	}
	// Relative path
	if current == "/" {
		return "/" + target
	}
	return current + "/" + target
}