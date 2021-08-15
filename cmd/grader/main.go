package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"github.com/fahmifan/autograd/grader"
	"github.com/fahmifan/autograd/model"
)

func init() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}

// findFilesInDir return map[fileName]filePath
func findFilesInDir(dir string) (map[string]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	mapFiles := make(map[string]string)
	for _, f := range files {
		if f.IsDir() {
			continue
		}

		mapFiles[f.Name()] = path.Join(dir, f.Name())
	}

	return mapFiles, nil
}

// read file and append each line in slices of string
func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lines := make([]string, 0)
	for {
		line, err := reader.ReadString(byte('\n'))
		if err != nil && err != io.EOF {
			return nil, err
		}

		if err == io.EOF {
			break
		}

		lines = append(lines, strings.TrimSpace(line))
	}

	return lines, err
}

type subm struct {
	Name   string
	Path   string
	Code   string
	UserID string
}

// parse code from fileName
func (s *subm) parseCode() {
	if s == nil {
		return
	}

	ss := strings.Split(s.Name, "-")
	if len(ss) != 3 {
		return
	}

	s.Code = ss[1]
	s.UserID = ss[0]
}

func findSubmissionsInDir(dir string) (subs map[string][]subm, err error) {
	var submissions map[string]string
	submissions, err = findFilesInDir(dir)
	if err != nil {
		log.Error(err)
		return
	}

	subs = make(map[string][]subm)
	for fname, fpath := range submissions {
		sub := subm{Name: fname, Path: fpath}
		sub.parseCode()
		subs[sub.Code] = append(subs[sub.Code], sub)
	}

	return subs, nil
}

func main() {
	grad := grader.New(grader.TypeCPP)

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	subs, err := findSubmissionsInDir(path.Join(cwd, "example/submission"))
	if err != nil {
		log.Error(err)
		return
	}

	// read input
	inputDir := path.Join(cwd, "example/input")
	mapInputs, err := findFilesInDir(inputDir)
	if err != nil {
		log.Error(err)
		return
	}

	// read output
	outputDir := path.Join(cwd, "example/output")
	outputs, err := findFilesInDir(outputDir)
	if err != nil {
		log.Error(err)
		return
	}

	for testCode := range mapInputs {
		expecteds, err := readFile(outputs[testCode])
		if err != nil {
			log.Fatal(err)
		}

		inputs, err := readFile(mapInputs[testCode])
		if err != nil {
			log.Fatal(err)
		}

		if len(expecteds) != len(inputs) {
			log.Errorf("error : unmatch input %d & ouput file %d\n", len(inputs), len(expecteds))
			return
		}

		if _, ok := subs[testCode]; !ok {
			log.Debug("not found >>> ", testCode, " >>> ", subs)
			continue
		}

		source := subs[testCode]
		fmt.Printf("%s:\n---\n", testCode)
		for _, src := range source {
			fmt.Printf("%s:\n", src.UserID)
			result, err := grad.Grade(&model.GradingArg{
				SourceCodePath: src.Path,
				Expecteds:      expecteds,
				Inputs:         inputs,
			})
			if err != nil {
				logrus.Error(err)
				return
			}

			fmt.Printf("outputs: %v | corrects: %v | score: %d\n\n", result.Outputs, result.Corrects, result.Score())
		}
	}
}
