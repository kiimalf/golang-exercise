package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"

	"latihan-fiber/app/model"
	"latihan-fiber/config"
	"latihan-fiber/helper"
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
	mode := flag.String("mode", "crowd", "mode pengujian: crowd | credits | swap")
	flag.Parse()

	config.LoadEnv()
	jwtSecret := config.GetEnv("JWT_SECRET", "super-secret-key-yang-sangat-panjang-dan-aman")
	jwtManager := helper.NewJWTManager(jwtSecret, "praktikum-backend", 1*time.Hour)

	baseURL := "http://localhost:" + config.GetEnv("APP_PORT", "3000") + "/api/v1"
	client := &http.Client{Timeout: 10 * time.Second}

	// Buat token mandiri untuk mhs001 - mhs060
	tokens := make([]string, 61)
	for i := 1; i <= 60; i++ {
		u := model.User{
			ID:       i, // ID pada database uji beban
			Username: fmt.Sprintf("mhs%03d", i),
			Email:    fmt.Sprintf("mhs%03d@student.unair.ac.id", i),
			Role:     "user",
			IsActive: true,
		}
		t, _ := jwtManager.GenerateAccess(u)
		tokens[i] = t
	}

	start := time.Now()

	switch *mode {
	case "crowd":
		fmt.Println("60 mahasiswa berebut kelas IF399 (kuota 5)...")
		// Mengambil ID untuk kelas IF399 (asumsi ID 13)
		tally := fire(60, func(i int) string {
			req, _ := http.NewRequest("POST", baseURL+"/courses/13/enroll", nil)
			req.Header.Set("Authorization", "Bearer "+tokens[i+1])
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
		fmt.Printf("hasil request:\n")
		for k, v := range tally {
			fmt.Printf("  %-25s %d\n", k, v)
		}
		fmt.Printf("waktu total: %v\n", time.Since(start))

	case "credits":
		fmt.Println("10 mahasiswa mendaftar 12 kelas serentak (batas 24 SKS)...")
		tally := fire(120, func(i int) string {
			mhsIdx := (i % 10) + 1
			courseID := (i / 10) + 1 // ID 1 sampai 12
			req, _ := http.NewRequest("POST", fmt.Sprintf("%s/courses/%d/enroll", baseURL, courseID), nil)
			req.Header.Set("Authorization", "Bearer "+tokens[mhsIdx])
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
		fmt.Printf("hasil request:\n")
		for k, v := range tally {
			fmt.Printf("  %-25s %d\n", k, v)
		}
		fmt.Printf("waktu total: %v\n", time.Since(start))

	case "swap":
		fmt.Println("20 pasang mahasiswa pindah kelas berlawanan arah...")
		tally := fire(40, func(i int) string {
			mhsIdx := i + 1
			fromID, toID := 1, 2
			if i%2 == 1 {
				fromID, toID = 2, 1
			}
			payload, _ := json.Marshal(map[string]int{
				"from_course_id": fromID,
				"to_course_id":   toID,
			})
			req, _ := http.NewRequest("POST", baseURL+"/enrollments/swap", bytes.NewBuffer(payload))
			req.Header.Set("Authorization", "Bearer "+tokens[mhsIdx])
			req.Header.Set("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				return "ERROR"
			}
			defer resp.Body.Close()
			return fmt.Sprintf("%d OK", resp.StatusCode)
		})
		fmt.Printf("hasil request:\n")
		for k, v := range tally {
			fmt.Printf("  %-25s %d\n", k, v)
		}
		fmt.Printf("waktu total: %v\n", time.Since(start))
	}
}
