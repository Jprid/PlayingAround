package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	fmt.Println("\nRetrieving Current user")
	userRequest, err := http.NewRequest("GET", "https://api.Github.com/repos/jprid/jprid.github.io/actions/runs?per_page=5&page=1", nil)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	// add username and token to headers
	resp, err := http.DefaultClient.Do(userRequest)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	defer resp.Body.Close()

	actionList := &ActionListResponse{}
	populateActionRunList(actionList, resp.Body)
	fmt.Printf("%9s : %9s : %28s:\n", "ID", "Status", "CreatedAt")
	for _, value := range actionList.WorkflowRuns {
		fmt.Printf("%d : %s : %s\n", value.ID, value.Status, value.Timestamp)
	}
}

// ActionListResponse represents the list of actions returned from the actions/runs GET API available on your public repositories from github
type ActionListResponse struct {
	TotalCount   int                 `json:"total_count"`
	WorkflowRuns []ActionWorkflowRun `json:"workflow_runs"`
}

type ActionWorkflowRun struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"created_at"`
	Status    string    `json:"status"`
}

// getPasswordFromStdin prompts user to write their password in stdin
func getPasswordFromStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input your password")
	text, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalln("failed to read pasword")
	}
	text = strings.Replace(text, "\r\n", "", -1)
	return text, nil
}

func populateActionRunList(runlist *ActionListResponse, responseBody io.ReadCloser) {
	enc := json.NewDecoder(responseBody)
	if err := enc.Decode(runlist); err != nil {
		log.Fatalf("error, can't decode response\n%v\n", err)
	}
}
