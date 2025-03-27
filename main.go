package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type CommandRequest struct {
	Name string `json:"name"`
}

type CommandResponse struct {
	Output string `json:"output"`
	Error  string `json:"error"`
}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	var req CommandRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "请求体解析失败", http.StatusBadRequest)
		return
	}

	name := req.Name
	command := fmt.Sprintf(`rm -f /srv/%s/tmp/login.png && sed -i 's/"gewechat_token": "[^"]*"/"gewechat_token": ""/' /srv/%s/config.json && docker-compose -f /srv/%s/docker-compose.yml restart %s`, name, name, name, name)
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()

	log.Print(string(output))
	resp := CommandResponse{
		Output: "sucess",
	}
	if err != nil {
		resp.Error = err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/restart", restartHandler)
	fmt.Println("Run :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
