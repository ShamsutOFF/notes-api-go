package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"notes-api/internal/domain"

	"github.com/spf13/cobra"
)

const (
	baseURL = "http://localhost:8080/api/notes"
)

var (
	rootCmd = &cobra.Command{
		Use:   "client",
		Short: "Notes CLI Client",
		Long:  "A CLI client for interacting with the Notes API",
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd, listCmd, getCmd, updateCmd, deleteCmd)
}

var createCmd = &cobra.Command{
	Use:   "create [title] [content]",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		content := strings.Join(args[1:], " ")

		note := domain.CreateNoteRequest{
			Title:   title,
			Content: content,
		}

		data, err := json.Marshal(note)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}

		resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(data))
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, func(body []byte) {
			var createdNote domain.Note
			if err := json.Unmarshal(body, &createdNote); err != nil {
				fmt.Printf("Error parsing response: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Note created successfully with ID: %d\n", createdNote.ID)
			printNoteDetail(createdNote)
		}, http.StatusCreated)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")

		resp, err := http.Get(baseURL)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, func(body []byte) {
			var notes []domain.Note
			if err := json.Unmarshal(body, &notes); err != nil {
				fmt.Printf("Error parsing response: %v\n", err)
				os.Exit(1)
			}

			if len(notes) == 0 {
				fmt.Println("No notes found.")
				return
			}

			if limit > 0 && limit < len(notes) {
				notes = notes[:limit]
			}

			switch format {
			case "json":
				output, _ := json.MarshalIndent(notes, "", "  ")
				fmt.Println(string(output))
			case "table":
				printNotesTable(notes)
			default:
				printNotesTable(notes)
			}
		}, http.StatusOK)
	},
}

var getCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get note by ID",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil || id <= 0 {
			fmt.Println("Error: invalid note ID")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/%d", baseURL, id)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, func(body []byte) {
			var note domain.Note
			if err := json.Unmarshal(body, &note); err != nil {
				fmt.Printf("Error parsing response: %v\n", err)
				os.Exit(1)
			}

			printNoteDetail(note)
		}, http.StatusOK)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [id] [title] [content]",
	Short: "Update a note",
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil || id <= 0 {
			fmt.Println("Error: invalid note ID")
			os.Exit(1)
		}

		title := args[1]
		content := strings.Join(args[2:], " ")

		note := domain.UpdateNoteRequest{
			Title:   title,
			Content: content,
		}

		data, err := json.Marshal(note)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/%d", baseURL, id)
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(data))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		handleResponse(resp, func(body []byte) {
			var updatedNote domain.Note
			if err := json.Unmarshal(body, &updatedNote); err != nil {
				fmt.Printf("Error parsing response: %v\n", err)
				os.Exit(1)
			}

			fmt.Println("Note updated successfully:")
			printNoteDetail(updatedNote)
		}, http.StatusOK)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a note",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil || id <= 0 {
			fmt.Println("Error: invalid note ID")
			os.Exit(1)
		}

		url := fmt.Sprintf("%s/%d", baseURL, id)
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			os.Exit(1)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error sending request: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			fmt.Printf("Note with ID %d deleted successfully\n", id)
		} else {
			handleResponse(resp, nil, http.StatusOK)
		}
	},
}

func handleResponse(resp *http.Response, successHandler func([]byte), expectedStatus int) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode != expectedStatus {
		var errorResp map[string]string
		if err := json.Unmarshal(body, &errorResp); err == nil {
			fmt.Printf("Error: %s\n", errorResp["error"])
		} else {
			fmt.Printf("Error: HTTP %d - %s\n", resp.StatusCode, string(body))
		}
		os.Exit(1)
	}

	if successHandler != nil {
		successHandler(body)
	}
}

func printNotesTable(notes []domain.Note) {
	fmt.Printf("Total notes: %d\n\n", len(notes))
	fmt.Println("ID  | Title                          | Created At          | Updated At")
	fmt.Println("----|--------------------------------|---------------------|---------------------")

	for _, note := range notes {
		title := note.Title
		if len(title) > 30 {
			title = title[:27] + "..."
		}

		fmt.Printf("%-4d| %-30s | %-19s | %-19s\n",
			note.ID,
			title,
			formatTime(note.CreatedAt),
			formatTime(note.UpdatedAt))
	}
}

func printNoteDetail(note domain.Note) {
	fmt.Println("=== Note Details ===")
	fmt.Printf("ID:         %d\n", note.ID)
	fmt.Printf("Title:      %s\n", note.Title)
	fmt.Printf("Content:    %s\n", note.Content)
	fmt.Printf("Created:    %s\n", formatTime(note.CreatedAt))
	fmt.Printf("Updated:    %s\n", formatTime(note.UpdatedAt))
	fmt.Println("===================")
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func init() {
	listCmd.Flags().StringP("format", "f", "table", "Output format (table, json)")
	listCmd.Flags().IntP("limit", "l", 0, "Limit number of notes to display")
}
