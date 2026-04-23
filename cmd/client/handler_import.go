package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	tea "charm.land/bubbletea/v2"
)

type ImportResponse struct {
	CueID string `json:"cue_id"`
}

type ImportStatus struct {
	Code     int        `json:"code"`
	Message  string     `json:"message"`
	Error    string     `json:"error"`
	Progress int        `json:"progress"`
	Missing  [][]string `json:"missing"`
}

type importStatusMsg struct {
	status   string
	err      string
	cue      string
	progress int
	missing  [][]string
	isDone   bool
	isError  bool
}

func (m *model) saveMissingCards() string {
	if len(m.filepicker.missing) == 0 {
		return ""
	}

	dir := filepath.Dir(m.filepicker.selectedFile)
	filename := fmt.Sprintf("missing_cards_%s.txt", time.Now().Format("20060102_150405"))
	path := filepath.Join(dir, filename)

	f, err := os.Create(path)
	if err != nil {
		return fmt.Sprintf("failed to create missing cards file: %v", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	for _, line := range m.filepicker.missing {
		_ = writer.Write(line)
	}
	writer.Flush()

	return fmt.Sprintf("Missing cards saved to: %s/%s", dir, filename)
}

func (m *model) doImportCollectionCmd() tea.Cmd {
	return func() tea.Msg {
		file, err := os.Open(m.filepicker.selectedFile)
		if err != nil {
			return importStatusMsg{err: err.Error(), isError: true, isDone: true}
		}
		defer file.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		part, err := writer.CreateFormFile("file", m.filepicker.selectedFile)
		if err != nil {
			return importStatusMsg{err: err.Error(), isError: true, isDone: true}
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return importStatusMsg{err: err.Error(), isError: true, isDone: true}
		}

		writer.Close()

		req, err := http.NewRequest("POST", website+"/collections/import?format=dragon_shield", body)
		if err != nil {
			return importStatusMsg{err: err.Error(), isError: true, isDone: true}
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Authorization", "Bearer "+m.jwtToken)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return importStatusMsg{err: err.Error(), isError: true, isDone: true}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusAccepted {
			type errResp struct {
				Error string `json:"error"`
			}
			var e errResp
			_ = json.NewDecoder(resp.Body).Decode(&e)
			if e.Error == "" {
				e.Error = resp.Status
			}
			return importStatusMsg{err: e.Error, isError: true, isDone: true}
		}

		var result ImportResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return importStatusMsg{err: "failed to parse response", isError: true, isDone: true}
		}

		return importStatusMsg{status: "Import started...", cue: result.CueID}
	}
}

func (m *model) checkImportStatusCmd(cue string) tea.Cmd {
	return func() tea.Msg {
		url := website + "/cue/" + cue

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return importStatusMsg{err: err.Error(), isError: true, isDone: true}
		}

		req.Header.Set("Authorization", "Bearer "+m.jwtToken)

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return importStatusMsg{err: err.Error(), isError: true, isDone: true}
		}
		defer resp.Body.Close()

		var status ImportStatus
		if resp.StatusCode == 404 {
			return importStatusMsg{err: "Status not found", isError: true, isDone: true}
		}

		if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
			return importStatusMsg{err: "Failed to read status", isError: true, isDone: true}
		}

		effectiveCode := status.Code
		if effectiveCode == 0 {
			effectiveCode = resp.StatusCode
		}

		msg := status.Message
		if msg == "" && status.Error != "" {
			msg = status.Error
		}

		if effectiveCode >= 400 {
			return importStatusMsg{
				err:      msg,
				isError:  true,
				isDone:   true,
				progress: status.Progress,
				missing:  status.Missing,
			}
		}

		if effectiveCode == 201 {
			return importStatusMsg{
				status:   msg,
				isDone:   true,
				progress: status.Progress,
				missing:  status.Missing,
			}
		}

		return importStatusMsg{
			status:   msg,
			cue:      cue,
			progress: status.Progress,
			missing:  status.Missing,
		}
	}
}

func waitTick() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg{}
	})
}

type tickMsg struct{}
