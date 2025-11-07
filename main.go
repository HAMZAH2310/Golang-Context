package main

import (
	"context"
	"math/rand"
	"fmt"
	"sync"
	"time"
)

func downloadFile(ctx context.Context, wg *sync.WaitGroup, fileName string, ch chan string) {
	defer wg.Done()

	duration := time.Duration(rand.Intn(2000)+500) * time.Millisecond

	select{
	case <-time.After(duration):
		ch <-fmt.Sprintf("%s selesai di unduh %s", fileName, duration)
	case <-ctx.Done():
		ch <- fmt.Sprintf("%s dibatalkan", fileName)
		return
	}
}


func main(){
	rand.New(rand.NewSource(time.Now().UnixNano()))

	files := []string{
		"file-A",
		"file-B",
		"File-C",
		"File-D",
		"File-F",
	}

	parent:= context.Background()
	ctx, cancel := context.WithCancel(parent)
	defer cancel() // Di cancel ketika semua function selesai

	ch := make(chan string, len(files))
	var wg sync.WaitGroup

	for _, f := range files{
		wg.Add(1)
		go downloadFile(ctx,&wg,f,ch)
	}

	go func() {
		time.Sleep(2 *time.Second)
		fmt.Println("User Membatalkan Proses Download") // Kelamaan menunggu download
		cancel()
	}()

	go func(){
		wg.Wait()
		close(ch)
	}()

	for msg:= range ch{
		fmt.Println(msg)
	}

	fmt.Println("Proses Selesai")
}	