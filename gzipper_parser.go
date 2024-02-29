package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type Work struct {
	FilePath string
}

func FileNameGen(dir string, pattern *regexp.Regexp) <-chan Work {
	jobs := make(chan Work)
	go func() {
		defer close(jobs)
		filepath.Walk(dir, func(path string, d fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				res := pattern.Find([]byte(d.Name()))
				if len(res) > 0 {
					fullPath, _ := filepath.Abs(path)
					jobs <- Work{FilePath: fullPath}
				}
			}
			return nil
		})
	}()
	return jobs
}

func compress(jobs <-chan Work) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				err := zipFile(job.FilePath, job.FilePath+".gz")
				if err != nil {
					fmt.Println("compression error:", err)
				}
			}
		}()
	}
}

func zipFile(sourcePath, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	gzw := gzip.NewWriter(target)
	defer gzw.Close()

	_, err = io.Copy(gzw, source)
	if err != nil {
		return err
	}
	return nil
}
