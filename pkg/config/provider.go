package config

import (
	"io/ioutil"
)

func NewFilesProvider(files []*File) Provider {
	return &FilesProvider{
		files: files,
	}
}

type File struct {
	Name string
	Path string
}

type FilesProvider struct {
	files   []*File
	buffers map[string][]byte
}

func (f *FilesProvider) Name() string {
	return "file"
}

func (f *FilesProvider) Load() error {
	return f.load()
}

func (f *FilesProvider) load() error {
	if f.buffers == nil {
		f.buffers = make(map[string][]byte)
	}
	for _, file := range f.files {
		buf, err := ioutil.ReadFile(file.Path)
		if err != nil {
			return err
		}
		// 支持环境变量
		f.buffers[file.Name] = []byte(ExpandEnv(string(buf)))
	}
	return nil
}

func (f *FilesProvider) Reload() error {
	return f.load()
}

func (f *FilesProvider) Read(file string) ([]byte, error) {
	return f.buffers[file], nil
}

func (f *FilesProvider) Watch(callback ProviderCallback) {
}
