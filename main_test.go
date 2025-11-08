package main

import (
	"context"
	"sync"
	"testing"
	"time"
)

// TestContextCancel memastikan cancel menghentikan proses download
func TestContextCancel(t *testing.T) {
	files := []string{"a", "b", "c"} // Membuat file baru
	ctx, cancel := context.WithCancel(context.Background()) // Membuat cancel
	defer cancel() // Memastikan memanggil Cancel setelah semua selesai

	ch := make(chan string, len(files)) // Membuat chanel, dan data dari files
	var wg sync.WaitGroup //Membuat variable wg untuk waitgroup

	for _, f := range files {
		wg.Add(1) // menambahkan perhitungan 
		go downloadFile(ctx, &wg, f, ch)
	}

	// Batalkan lebih cepat dari waktu download
	time.Sleep(100 * time.Millisecond)
	cancel()

	go func() {
		wg.Wait()
		close(ch)
	}()

	// Ambil hasil dari channel
	results := []string{}
	for msg := range ch {
		results = append(results, msg)
	}

	// Pastikan semua hasil mengandung kata "dibatalkan"
	for _, msg := range results {
		if msg == "" {
			t.Errorf("Pesan kosong diterima")
		}
		if msg[:2] != "❌" {
			t.Errorf("Ditemukan file yang tidak dibatalkan: %s", msg)
		}
	}
}
