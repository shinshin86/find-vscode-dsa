package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

type ProjectInfo struct {
	ProjectName                    string `json:"projectName"`
	ProjectPath                    string `json:"projectPath"`
	VscodeWorkspaceStoragePath     string `json:"vscodeWorkspaceStoragePath"`
	WorkspaceRecommendationsIgnore bool   `json:"workspaceRecommendationsIgnore"`
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// fetch vscode project info
func (a *App) ProjectInfoList() ([]ProjectInfo, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	var path string
	if runtime.GOOS == "windows" {
		appdataPath, exists := os.LookupEnv("APPDATA")

		if !exists {
			return nil, fmt.Errorf("APPDATA environment variable not found")
		}

		path = filepath.Join(appdataPath, "\\Code\\User\\workspaceStorage")
	} else {
		path = filepath.Join(home, "Library/Application Support/Code/User/workspaceStorage")
	}

	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var projects []ProjectInfo

	for _, dir := range dirs {
		if dir.IsDir() {
			workspaceJSONPath := filepath.Join(path, dir.Name(), "workspace.json")
			data, err := ioutil.ReadFile(workspaceJSONPath)
			if err != nil {
				log.Println("Could not read file:", workspaceJSONPath, "Error:", err)
				continue
			}

			// read workspace.json (VSCode workspaceStorage)
			var jsonContent map[string]interface{}
			err = json.Unmarshal(data, &jsonContent)
			if err != nil {
				log.Println("Could not parse JSON from file:", workspaceJSONPath, "Error:", err)
				continue
			}

			projectPath, ok := jsonContent["folder"].(string)
			if !ok {
				log.Println("Could not find 'projectPath' key in JSON from file:", workspaceJSONPath)
				continue
			}

			vscodeWorkspaceStoragePathAbs, err := filepath.Abs(filepath.Join(path, dir.Name()))
			if err != nil {
				log.Println("Could not get absolute path for:", dir.Name(), "Error:", err)
				continue
			}

			// DB access
			dbPath := filepath.Join(vscodeWorkspaceStoragePathAbs, "state.vscdb")
			db, err := sql.Open("sqlite3", dbPath)
			if err != nil {
				log.Println("Could not open database at:", dbPath, "Error:", err)
				continue
			}

			rows, err := db.Query("SELECT value FROM ItemTable WHERE key = 'extensionsAssistant/workspaceRecommendationsIgnore';")
			if err != nil {
				log.Println("Could not query database at:", dbPath, "Error:", err)
				continue
			}
			defer rows.Close()

			var value string
			var boolVal bool
			var parseBoolErr error
			for rows.Next() {
				err = rows.Scan(&value)
				if err != nil {
					log.Println("Could not read row from database at:", dbPath, "Error:", err)
					continue
				}
			}

			if value != "" {
				boolVal, parseBoolErr = strconv.ParseBool(value)
				if parseBoolErr != nil {
					log.Println("Could not convert string to bool, Error:", err)
					return nil, err
				}
			} else {
				boolVal = false
			}

			projectInfo := ProjectInfo{
				ProjectName:                    filepath.Base(projectPath),
				ProjectPath:                    projectPath,
				VscodeWorkspaceStoragePath:     vscodeWorkspaceStoragePathAbs,
				WorkspaceRecommendationsIgnore: boolVal,
			}

			projects = append(projects, projectInfo)
		}
	}

	return projects, nil
}

func updateIgnoreSetting(project ProjectInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	dbPath := filepath.Join(project.VscodeWorkspaceStoragePath, "state.vscdb")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Could not open database at:", dbPath, "Error:", err)
		return
	}

	stmt, err := db.Prepare(`UPDATE ItemTable SET 
		"value" = ? WHERE key = "extensionsAssistant/workspaceRecommendationsIgnore"`)
	if err != nil {
		log.Println("Could not prepare SQL statement. Error:", err)
		return
	}
	defer stmt.Close()

	intBoolValue := 0
	if project.WorkspaceRecommendationsIgnore {
		intBoolValue = 1
	}
	_, err = stmt.Exec(intBoolValue)
	if err != nil {
		log.Println("Could not execute SQL statement. Error:", err)
		return
	}
}

func (a *App) UpdateWorkspaceRecommendationsIgnore(projectList []ProjectInfo) {
	var wg sync.WaitGroup

	for _, project := range projectList {
		wg.Add(1)
		go updateIgnoreSetting(project, &wg)
	}

	wg.Wait()
}

func (a *App) OpenDir(path string) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("start", path)
	} else {
		cmd = exec.Command("open", path)
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Failed open directory:  %v", err)
	}

	return nil
}
