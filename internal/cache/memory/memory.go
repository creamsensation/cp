package memory

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Client interface {
	Get(key string) string
	Set(key string, value string, expiration time.Duration) error
	Exists(key string) bool
	Destroy(key string) error
}

type memory struct {
	sync.RWMutex
	data map[string]memoryData
	dir  string
}

type memoryData struct {
	Value      string    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

const (
	jsonSuffix    = ".json"
	watchInterval = time.Second
)

func New() Client {
	m := &memory{
		data: make(map[string]memoryData),
		dir:  getDir(),
	}
	go m.load()
	go m.watch()
	return m
}

func (m *memory) Get(key string) string {
	m.Lock()
	data, ok := m.data[key]
	if !ok {
		return ""
	}
	m.Unlock()
	return data.Value
}

func (m *memory) Set(key string, value string, expiration time.Duration) error {
	data := memoryData{
		Value:      value,
		Expiration: time.Now().Add(expiration),
	}
	m.Lock()
	m.data[key] = data
	m.Unlock()
	if err := m.setTempFile(key, data); err != nil {
		return err
	}
	return nil
}

func (m *memory) Exists(key string) bool {
	_, ok := m.data[key]
	return ok
}

func (m *memory) Destroy(key string) error {
	delete(m.data, key)
	if err := m.deleteTempFile(key); err != nil {
		return err
	}
	return nil
}

func (m *memory) load() {
	if _, err := os.Stat(m.dir); os.IsNotExist(err) {
		return
	}
	if err := filepath.Walk(
		m.dir, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() || !strings.HasSuffix(info.Name(), jsonSuffix) {
				return nil
			}
			key := strings.TrimSuffix(info.Name(), jsonSuffix)
			fbts, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			var data memoryData
			if err := json.Unmarshal(fbts, &data); err != nil {
				return err
			}
			m.data[key] = data
			return nil
		},
	); err != nil {
		log.Fatalln(err)
	}
}

func (m *memory) watch() {
	ticker := time.NewTicker(watchInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			t := time.Now()
			expired := make([]string, 0)
			m.RLock()
			for key, data := range m.data {
				if t.Before(data.Expiration) {
					continue
				}
				expired = append(expired, key)
			}
			m.RUnlock()
			if len(expired) > 0 {
				m.Lock()
				for _, key := range expired {
					if err := m.Destroy(key); err != nil {
						log.Fatalln(err)
					}
				}
				expired = nil
				m.Unlock()
			}
		}
	}
}

func (m *memory) setTempFile(key string, data memoryData) error {
	if _, err := os.Stat(m.dir); os.IsNotExist(err) {
		if err := os.MkdirAll(m.dir, os.ModePerm); err != nil {
			return err
		}
	}
	path := fmt.Sprintf("%s/%s%s", m.dir, key, jsonSuffix)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	dbts, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, dbts, 0444); err != nil {
		return err
	}
	return nil
}

func (m *memory) deleteTempFile(key string) error {
	if _, err := os.Stat(m.dir); os.IsNotExist(err) {
		return nil
	}
	path := fmt.Sprintf("%s/%s%s", m.dir, key, jsonSuffix)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if err := os.Remove(path); err != nil {
			return err
		}
	}
	return nil
}

func getDir() string {
	tmpDir := os.TempDir()
	if strings.HasSuffix(tmpDir, "/") {
		tmpDir = strings.TrimSuffix(tmpDir, "/")
	}
	return fmt.Sprintf("%s/.creamsensation/memory", tmpDir)
}
