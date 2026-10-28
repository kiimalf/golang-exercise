package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"api-students/app/model"
	"api-students/config"
	"api-students/helper"
)

func fire(n int, job func(i int) string) map[string]int {
	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan string, n)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int, out chan<- string) {
			defer wg.Done()
			<-start
			out <- job(i)
		}(i, results)
	}

	close(start)
	wg.Wait()
	close(results)

	tally := map[string]int{}
	for r := range results {
		tally[r]++
	}
	return tally
}

func main() {
	supervisorID := flag.Int("supervisor", 3, "ID dosen pembimbing target (kuota 3)")
	flag.Parse()

	config.LoadEnv()
	jwtSecret := config.GetEnv("JWT_SECRET", "super-secret-key-yang-sangat-panjang-dan-aman")
	jwtManager := helper.NewJWTManager(jwtSecret, "api-students", 1*time.Hour)

	baseURL := "http://localhost:" + config.GetEnv("APP_PORT", "3000") + "/api/v1"
	client := &http.Client{Timeout: 10 * time.Second}

	adminToken, err := jwtManager.GenerateAccess(model.User{
		ID:       1,
		Username: "admin",
		Role:     "admin",
		IsActive: true,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("30 mahasiswa serentak ditetapkan ke Pembimbing ID %d (kuota 3)...\n", *supervisorID)
	start := time.Now()

	tally := fire(30, func(i int) string {
		studentID := i + 1 // Menguji mahasiswa ID 1 sampai 30
		payload := fmt.Sprintf(`{"student_id": %d}`, studentID)

		req, _ := http.NewRequest("POST", fmt.Sprintf("%s/supervisors/%d/assign", baseURL, *supervisorID), strings.NewReader(payload))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return "ERROR"
		}
		defer resp.Body.Close()

		if resp.StatusCode == 201 {
			return "201 OK"
		}
		var body map[string]any
		json.NewDecoder(resp.Body).Decode(&body)
		if code, ok := body["code"].(string); ok {
			return fmt.Sprintf("%d %s", resp.StatusCode, code)
		}
		return fmt.Sprintf("%d", resp.StatusCode)
	})

	fmt.Println("Hasil request:")
	for k, v := range tally {
		fmt.Printf("  %-25s %d\n", k, v)
	}
	fmt.Printf("Waktu total: %v\n", time.Since(start))
}
