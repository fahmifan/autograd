package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetReportCaller(true)
	log.SetLevel(log.DebugLevel)
}

func compile(filePath string) (outPath string, err error) {
	outPath = path.Join(fmt.Sprintf("%s.bin", filePath))
	args := strings.Split(fmt.Sprintf("%s -o %s", filePath, outPath), " ")
	cmd := exec.Command("g++", args...)
	bt, err := cmd.CombinedOutput()
	if err != nil {
		log.WithFields(log.Fields{
			"output":  string(bt),
			"args":    args,
			"outPath": outPath,
		}).Error(err)
	}

	return
}

func run(source, input string) (out string) {
	inputs := strings.Split(input, " ")
	input = strings.Join(inputs, "\n")
	cmd := exec.Command(source, inputs...)

	var buffOut bytes.Buffer
	var buffErr bytes.Buffer

	cmd.Stdin = bytes.NewBuffer([]byte(input))
	cmd.Stdout = &buffOut
	cmd.Stderr = &buffErr

	err := cmd.Run()
	if err != nil {
		return
	}

	out = strings.TrimSpace(string(buffOut.Bytes()))
	return
}

func remove(source string) {
	if err := os.Remove(source); err != nil {
		log.Println("error : ", err)
	}
}

func grade(source, input, expected string) (result string) {
	outPath, err := compile(source)
	if err != nil {
		log.WithField("source", source).Error(err)
		return ""
	}

	defer remove(outPath)
	out := run(outPath, input)
	if expected != out {
		result = "NAY"
		return
	}

	result = "AYE"

	return
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
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	subs, err := findSubmissionsInDir(path.Join(cwd, "submission"))
	if err != nil {
		log.Error(err)
		return
	}

	// read input
	inputDir := path.Join(cwd, "input")
	inputs, err := findFilesInDir(inputDir)
	if err != nil {
		log.Error(err)
		return
	}

	// read output
	outputDir := path.Join(cwd, "output")
	outputs, err := findFilesInDir(outputDir)
	if err != nil {
		log.Error(err)
		return
	}

	for testCode := range inputs {
		outs, err := readFile(outputs[testCode])
		if err != nil {
			log.Fatal(err)
		}

		ins, err := readFile(inputs[testCode])
		if err != nil {
			log.Fatal(err)
		}

		if len(outs) != len(ins) {
			log.Error("error : unmatch input %d & ouput file %d\n", len(ins), len(outs))
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
			for i := 0; i < len(outs); i++ {
				g := grade(src.Path, ins[i], outs[i])
				fmt.Printf("- case %d %s\n", i+1, g)
			}
			fmt.Println()
		}
		fmt.Println()
	}
}
