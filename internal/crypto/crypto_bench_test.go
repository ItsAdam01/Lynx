package crypto

import (
	"crypto/rand"
	"os"
	"testing"
)

func BenchmarkHashFile_1MB(b *testing.B) {
	// Create a 1MB dummy file
	tmpfile, err := os.CreateTemp("", "bench-1mb")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	data := make([]byte, 1024*1024)
	rand.Read(data)
	tmpfile.Write(data)
	tmpfile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashFile(tmpfile.Name())
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkHashFile_10MB(b *testing.B) {
	// Create a 10MB dummy file
	tmpfile, err := os.CreateTemp("", "bench-10mb")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	data := make([]byte, 10*1024*1024)
	rand.Read(data)
	tmpfile.Write(data)
	tmpfile.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := HashFile(tmpfile.Name())
		if err != nil {
			b.Fatal(err)
		}
	}
}
