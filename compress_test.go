package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"sync"
	"testing"

	"github.com/klauspost/compress/s2"
	"github.com/stretchr/testify/assert"
)

var b1, b2, b3, b4, b5, b6 []byte

func init() {
	var err error

	b1, err = base64.StdEncoding.DecodeString("gqZVc2VySUTAsVByb2ZpbGVQaWN0dXJlVVJMxFdodHRwOi8vbG9jYWxob3N0OjgwODAvZXh0ZXJuYWwvYXp1cmUvcHJvZmlsZS82Nzg1YmM1Mi0wYmM4LTQxMDctYjMyNC02NzYwYzE0NTJkZmYvaW1hZ2U=")
	if err != nil {
		panic(err)
	}

	b2, err = base64.StdEncoding.DecodeString("gqZVc2VySUTAsVByb2ZpbGVQaWN0dXJlVVJMxFdodHRwOi8vbG9jYWxob3N0OjgwODAvZXh0ZXJuYWwvYXp1cmUvcHJvZmlsZS8xM2YzYzc3Yy0wZGUwLTRmODUtOGEzMy1mYjUzMzFjNmU3NzcvaW1hZ2U=")
	if err != nil {
		panic(err)
	}

	b3, err = base64.StdEncoding.DecodeString("k4eiSUS7Um5WdWJua2djR1Z5YldsemMybHZiaUJzZFd3plVzZXJJRMClRW1haWy1ZXhhbXBsZTFAdHBhLWdyb3VwLmF0pE5hbWWoQmF1IE1heGylUm9sZXOSpXdyaXRlpHJlYWSxUHJvZmlsZVBpY3R1cmVVUkzZV2h0dHA6Ly9sb2NhbGhvc3Q6ODA4MC9leHRlcm5hbC9henVyZS9wcm9maWxlLzY3ODViYzUyLTBiYzgtNDEwNy1iMzI0LTY3NjBjMTQ1MmRmZi9pbWFnZa1Jc1RQQUVtcGxveWVlw4eiSUS2WXpwa2J5QnViM1FnYzJodmR5QnRaUaZVc2VySUTApUVtYWlst3NvbWUtb3duZXJAdHBhLWdyb3VwLmF0pE5hbWWwT3duZXIgb2YgQ29tcGFueaVSb2xlc5Glb3duZXKxUHJvZmlsZVBpY3R1cmVVUkzZV2h0dHA6Ly9sb2NhbGhvc3Q6ODA4MC9leHRlcm5hbC9henVyZS9wcm9maWxlLzEzZjNjNzdjLTBkZTAtNGY4NS04YTMzLWZiNTMzMWM2ZTc3Ny9pbWFnZa1Jc1RQQUVtcGxveWVlwoeiSUS4WVNCemFHRnlaWEJ2YVc1MElHZHliM1Z3plVzZXJJRMClRW1haWzApE5hbWW9TWVtYmVycyBvZiBTb21lIEdyb3VwIENvbXBhbnmlUm9sZXORpHJlYWSxUHJvZmlsZVBpY3R1cmVVUkzArUlzVFBBRW1wbG95ZWXC")
	if err != nil {
		panic(err)
	}

	b4, err = base64.StdEncoding.DecodeString("gqZVc2VySUTZJDc2YTdkMTc0LTUyZDktNGZjOC05ZTM4LWFiOWU4MWEyNTgxZLFQcm9maWxlUGljdHVyZVVSTMRQaHR0cDovL2xvY2FsaG9zdDo4MDgwL2Fzc2V0cy91cGxvYWRzL2ltYWdlcy81M2Y0Y2NmMC1iNWUwLTQwNWMtODcyMy04OTIzMDI2MDlmOGU=")
	if err != nil {
		panic(err)
	}

	b5, err = base64.StdEncoding.DecodeString("gqZVc2VySUTAsVByb2ZpbGVQaWN0dXJlVVJMxFdodHRwOi8vbG9jYWxob3N0OjgwODAvZXh0ZXJuYWwvYXp1cmUvcHJvZmlsZS82MTA1YjRmZC0zZDg4LTRkMGUtYmI0Ny0wOWY2MDBiNTBiZWYvaW1hZ2U=")
	if err != nil {
		panic(err)
	}

	b6, err = base64.StdEncoding.DecodeString("koeiSUTZME9HRmlabUZpWVRNdE5EQmlNQzAwWmpkakxXSmxNR1F0TlRJNU5qVXhaVGczTlRNeKZVc2VySUTZJDc2YTdkMTc0LTUyZDktNGZjOC05ZTM4LWFiOWU4MWEyNTgxZKVFbWFpbLZleHRfcG9ydDFAdHBhLWdyb3VwLmF0pE5hbWWoU29tZSBHdXmlUm9sZXORpXdyaXRlsVByb2ZpbGVQaWN0dXJlVVJM2VBodHRwOi8vbG9jYWxob3N0OjgwODAvYXNzZXRzL3VwbG9hZHMvaW1hZ2VzLzUzZjRjY2YwLWI1ZTAtNDA1Yy04NzIzLTg5MjMwMjYwOWY4Za1Jc1RQQUVtcGxveWVlw4eiSUS4WVNCemFHRnlaWEJ2YVc1MElHZHliM1Z3plVzZXJJRMClRW1haWy0c29tZVVzZXJAZXhhbXBsZS5jb22kTmFtZblTb21lIFVzZXIgd2l0aCBQZXJtaXNzaW9upVJvbGVzkaV3cml0ZbFQcm9maWxlUGljdHVyZVVSTNlXaHR0cDovL2xvY2FsaG9zdDo4MDgwL2V4dGVybmFsL2F6dXJlL3Byb2ZpbGUvNjEwNWI0ZmQtM2Q4OC00ZDBlLWJiNDctMDlmNjAwYjUwYmVmL2ltYWdlrUlzVFBBRW1wbG95ZWXC")
	if err != nil {
		panic(err)
	}
}

var bufPool = sync.Pool{
	New: newBuffer,
}

func newBuffer() interface{} {
	return &bytes.Buffer{}
}

// getBuffer retrieves a buffer from the buffer pool
func getBuffer() *bytes.Buffer {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	return buf
}

// putBuffer back into the buffer pool
func putBuffer(buf *bytes.Buffer) {
	bufPool.Put(buf)
}

func cycle(buf []byte) ([]byte, error) {
	// lenCompressed := s2.MaxEncodedLen(len(buf)) + 1 + 4
	// compressed := make([]byte, lenCompressed)
	compressed := s2.Encode(nil, buf)
	// compressed = append(compressed, 0x1)

	intermediate := make([]byte, len(compressed))
	copy(intermediate, compressed)

	// intermediate = intermediate[:len(intermediate)-1]

	lenIntermediate, err := s2.DecodedLen(intermediate)
	if err != nil {
		return nil, err
	}

	uncompressed := make([]byte, lenIntermediate)

	if _, err = s2.Decode(uncompressed, intermediate); err != nil {
		return nil, err
	}

	return uncompressed, nil
}

func TestS3CompressBufferOutOfBounds(t *testing.T) {

	runs := 500 // v1.11.4 corrupts at 500
	size := 20  // v1.11.4 corrupts at 20

	corruptedMap := make(map[string]map[string][]uint64)
	corruptedMapMu := sync.Mutex{}

	var wg sync.WaitGroup

	for run := 0; run < runs; run++ {
		wg.Add(1)

		go func(t assert.TestingT, run int) {
			defer wg.Done()

			buf := getBuffer()
			defer putBuffer(buf)

			corrupted := make(map[string][]uint64)
			for i := 1; i < size; i++ {

				key := fmt.Sprintf("%d", i)
				arr := make([]uint64, i)
				for it := range arr {
					arr[it] = uint64(65280 + it)
					buf.WriteString(fmt.Sprintf("%v", uint64(65280+it)))
				}

				corrupted[key] = arr
				corruptedMapMu.Lock()
				corruptedMap[fmt.Sprintf("%d", run)] = corrupted
				corruptedMapMu.Unlock()
			}

			{
				buf, err := cycle(b1)
				assert.NoError(t, err)
				assert.Equal(t, b1, buf)
			}

			{
				buf, err := cycle(b2)
				assert.NoError(t, err)
				assert.Equal(t, b2, buf)
			}

			{
				buf, err := cycle(b3)
				assert.NoError(t, err)
				assert.Equal(t, b3, buf)
			}

			{
				buf, err := cycle(b4)
				assert.NoError(t, err)
				assert.Equal(t, b4, buf)
			}

			{
				buf, err := cycle(b5)
				assert.NoError(t, err)
				assert.Equal(t, b5, buf)
			}

			{
				buf, err := cycle(b6)
				assert.NoError(t, err)
				assert.Equal(t, b6, buf)
			}
		}(t, run)

	}

	wg.Wait()

	for run := 0; run < runs; run++ {

		corruptedMapMu.Lock()
		corrupted := corruptedMap[fmt.Sprintf("%d", run)]

		for i := 1; i < size; i++ {
			key := fmt.Sprintf("%d", i)
			arr := corrupted[key]
			for it, val := range arr {
				assert.Equal(t, val, uint64(65280+it), fmt.Sprintf("corrupted[%v][%d]", key, it))
			}
		}

		corruptedMapMu.Unlock()
	}

}
