package patcher

import (
	"encoding/json"
	"io/ioutil"
	"github.com/iain17/logger"
	"github.com/iain17/patcher/scanner"
	"os"
	"container/list"
	"fmt"
	"bufio"
)

type Patch struct {
	Signature string `json:"signature"`
	Write string `json:"write"`
	Address int64 `json:"address"`
	AddressHex string `json:"addressHex"`//Just for convience. Not read.
}

type Patcher struct {
	patchPath string
	inputPath string
	outputPath string
	patches map[string]*Patch
}

func New(patchPath string, inputPath, outputPath string) (*Patcher, error) {
	instance := &Patcher {
		patchPath: patchPath,
		inputPath: inputPath,
		outputPath: outputPath,
	}
	err := instance.load()
	if err != nil {
		return nil, err
	}
	return instance, nil
}

//Loads the json of patchPath
func (p *Patcher) load() error {
	file, err := ioutil.ReadFile(p.patchPath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, &p.patches); err != nil {
		return err
	}
	return nil
}

//Finds all the addresses of the
func (p *Patcher) Find() error {
	for name, patch := range p.patches {
		logger.Infof("Finding %s with sig %s", name, patch.Signature)
		address, err := scanner.Scan(patch.Signature, p.inputPath)
		if err != nil {
			if patch.Address == int64(0) {
				logger.Warningf("Failed to find sig %s. %v", patch.Signature, err)
			}
			patch.Address = 0
			continue
		}

		if patch.Address != address {
			logger.Infof("Found %s on %s", name, address)
		}
		patch.Address = address
		patch.AddressHex = fmt.Sprintf("%#08x", address)//Because it reads easier.
	}
	return nil
}

func (p *Patcher) getPatches() map[int64]*list.List {
	patches := map[int64]*list.List{}
	for name, patch := range p.patches {
		//Either we didn't find it or we have nothing to write.
		if patch.Address == int64(0) || patch.Write == "" {
			continue
		}
		bytes := scanner.ByteSeqToByteList(patch.Write)
		if bytes.Len() == 0 {
			logger.Warningf("'%s' contained no good instructions", patch.Write)
			continue
		}
		patches[patch.Address] = bytes
		logger.Infof("%s added to patchset.", name)
	}
	return patches
}

//Patches the memory
func (p *Patcher) Patch() error {
	//Open original and patched file.
	file, err := os.Open(p.inputPath)
	if err != nil {
		return err
	}
	defer file.Close()
	patched, err := os.Create(p.outputPath)
	if err != nil {
		return err
	}
	defer patched.Close()

	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)

	//Get a list of patches
	patches := p.getPatches()
	if len(patches) == 0{
		return nil
	}

	memory := int64(0)
	var patch *list.Element
	for scanner.Scan() {
		//Up the memory address
		memory += int64( len(scanner.Bytes()) )

		//Write original or patched?
		if patch != nil {
			logger.Debugf("Writing %X instead of %X on %#08x", patch.Value, scanner.Bytes(), memory)
			patched.Write([]byte{patch.Value.(byte)})
			patch = patch.Next()
		} else {
			//original
			patched.Write(scanner.Bytes())
		}

		//Do we have a patch for this piece of memory?
		if patches[memory] != nil {
			patch = patches[memory].Front()
		}

	}
	return nil
}

//Saves the p.patches to the input path
func (p *Patcher) Save() error {
	data, err := json.MarshalIndent(p.patches, "", "    ")
	if err != nil {
		return err
	}
	os.Remove(p.patchPath)
	ioutil.WriteFile(p.patchPath, data, 0644)
	return nil
}