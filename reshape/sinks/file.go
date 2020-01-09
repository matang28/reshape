package sinks

import (
	"fmt"
	"github.com/matang28/reshape/reshape"
	"github.com/matang28/reshape/reshape/etc"
	"github.com/matang28/reshape/reshape/serde"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type FileSinkConfig struct {
	PathToFolder   string
	FileNamePrefix string

	RotateTimeout *time.Duration
	RotateSizeMb  int

	Serializer serde.Serializer
	LineBreak  string
}

type FileSink struct {
	config FileSinkConfig

	file     *os.File
	fileSize int64

	queue  chan string
	errors chan error
	mutex  sync.Mutex
}

func NewFileSink(config FileSinkConfig, errors chan error) *FileSink {
	ensureDefaults(&config)
	fs := &FileSink{
		config:   config,
		file:     nil,
		queue:    make(chan string),
		fileSize: 0,
		mutex:    sync.Mutex{},
		errors:   errors,
	}

	if err := fs.openOrCreate(); err != nil {
		panic(err)
	}

	go fs.pullFromQueue()
	return fs
}

func (this *FileSink) Dump(objects ...interface{}) error {
	sb := strings.Builder{}

	for _, o := range objects {
		str, err := this.config.Serializer(o)
		if err != nil {
			return err
		}
		sb.WriteString(str)
		sb.WriteString(this.config.LineBreak)
	}

	go func() {
		this.queue <- sb.String()
	}()
	return nil
}

func (this *FileSink) Close() error {
	close(this.queue)
	return this.file.Close()
}

func (this *FileSink) pullFromQueue() {
	var timeout <-chan time.Time
	if this.config.RotateTimeout != nil {
		timeout = time.Tick(*this.config.RotateTimeout)
	} else {
		timeout = make(chan time.Time, 1)
	}

	for {
		select {
		case line, ok := <-this.queue:
			if !ok {
				return
			}
			if this.shouldRotateSize(line) {
				this.rotate()
			}

			if err := this.writeLine(line); err != nil {
				reshape.Report(reshape.NewSinkError(err), this.errors)
			}

		case <-timeout:
			this.rotate()
		}
	}
}

func (this *FileSink) shouldRotateSize(line string) bool {
	if this.config.RotateSizeMb <= 0 {
		return false
	}

	futureSize := etc.Bytes2Mb(int64(len(line)) + this.fileSize)
	if futureSize > this.config.RotateSizeMb {
		return true
	}

	return false
}

func (this *FileSink) writeLine(line string) error {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	nBytes, err := this.file.WriteString(line)
	if err != nil {
		return err
	}

	this.fileSize += int64(nBytes)
	return nil
}

func (this *FileSink) rotate() {
	this.mutex.Lock()
	defer this.mutex.Unlock()

	if err := this.file.Close(); err != nil {
		reshape.Report(reshape.NewUnrecoverableError(err), this.errors)
	}

	if err := this.openOrCreate(); err != nil {
		reshape.Report(reshape.NewUnrecoverableError(err), this.errors)
	}
}

func (this *FileSink) openOrCreate() error {
	if err := os.MkdirAll(this.config.PathToFolder, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filepath.Join(this.config.PathToFolder, this.generateName()))
	if err != nil {
		return err
	}
	this.file = file
	return nil
}

func (this *FileSink) generateName() string {
	return fmt.Sprintf("%s-%d.log", this.config.FileNamePrefix, time.Now().Unix())
}

func ensureDefaults(config *FileSinkConfig) {
	if config.PathToFolder == "" {
		config.PathToFolder = filepath.Join(os.TempDir(), "reshape")
	}
	if config.FileNamePrefix == "" {
		config.FileNamePrefix = "file_sink"
	}
	if config.Serializer == nil {
		config.Serializer = serde.JsonSerializer
	}
	if config.LineBreak == "" {
		config.LineBreak = "\n"
	}
}
