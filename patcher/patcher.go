package patcher

import (
	"encoding/json"
	"io/ioutil"
)

type Patcher struct {
	inputPath string
	outputPath string
	patches map[string]string
}

func New(patchPath string, inputPath, outputPath string) (*Patcher, error) {
	instance := &Patcher {
		inputPath: inputPath,
		outputPath: outputPath,
	}
	file, err := ioutil.ReadFile(patchPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(file, &instance.patches); err != nil {
		panic(err)
	}
	return instance, nil
}

func (p *Patcher) Patch() {
	for find, replace := range p.patches {
		println(find)
		println(replace)
	}
}