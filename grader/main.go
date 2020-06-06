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
	args := strings.Split(fmt.Sprintf("-o %s %s", outPath, filePath), " ")
	cmd := exec.Command("g++", args...)
	bt, err := cmd.CombinedOutput()
	if err != nil {
		log.WithField("output", string(bt)).Error(err)
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
		log.Error(err)
		return ""
	}

	defer remove(outPath)
	out := run(outPath, input)
	if expected != out {
		log.Debugf("out=%s expected=%s\n", out, expected)
		result = "NAY"
		return
	}

	result = "AYE"

	return
}

func findFilesInDir(dir string) (map[string]string, error) {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	files := make(map[string]string)
	for _, d := range dirs {
		if d.IsDir() {
			continue
		}

		files[d.Name()] = path.Join(dir, d.Name())
	}

	return files, nil
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
		// lines
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

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// read submission
	submissionDir := path.Join(cwd, "submission")
	submissions, err := findFilesInDir(submissionDir)
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

	for _, sourcePath := range submissions {
		for k2 := range outputs {
			outs, err := readFile(outputs[k2])
			if err != nil {
				log.Fatal(err)
			}

			ins, err := readFile(inputs[k2])
			if err != nil {
				log.Fatal(err)
			}

			if len(outs) != len(ins) {
				log.Error("error : unmatch input %d & ouput file %d\n", len(ins), len(outs))
				return
			}

			for i := 0; i < len(outs); i++ {
				g := grade(sourcePath, ins[i], outs[i])
				fmt.Println(g)
			}
		}
	}
}
