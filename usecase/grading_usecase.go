package usecase

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/miun173/autograd/grader"
	"github.com/sirupsen/logrus"
)

// GradingUsecase ...
type GradingUsecase struct {
	grader grader.Grader
}

// NewGradingUsecase :nodoc:
func NewGradingUsecase() *GradingUsecase {
	return &GradingUsecase{}
}

type assignment struct {
	id        int64
	inputURL  string
	outputURL string
	input     []string
	output    []string
}

// GradeSubmission ..
func (g *GradingUsecase) GradeSubmission(submissionID int64) (err error) {
	srcCodePath, err := g.getSubmissionSrcCodeByID(submissionID)
	if err != nil {
		err = fmt.Errorf("unable to get submission %d: %w", submissionID, err)
		logrus.Error(err)
		return
	}

	asg, err := g.getAssignmentBySubmissionID(submissionID)
	if err != nil {
		err = fmt.Errorf("unable to get assignment for submission %d: %w", submissionID, err)
		logrus.Error(err)
		return
	}

	_, _, err = g.grader.Grade(srcCodePath, asg.input, asg.output)
	return
}

func (g *GradingUsecase) getAssignmentBySubmissionID(submissionID int64) (assignment, error) {
	return g.getAssignment(0)
}

func (g *GradingUsecase) getSubmissionSrcCodeByID(id int64) (srcCodePath string, err error) {
	return
}

func (g *GradingUsecase) getAssignment(assignmentID int64) (asg assignment, err error) {
	// download source code to local filepath
	inputPath := ""
	inputPath, asg.input, err = g.downloadAndParse(asg.inputURL)
	if err != nil {
		return
	}
	defer removeFile(inputPath)

	outputPath := ""
	outputPath, asg.output, err = g.downloadAndParse(asg.outputURL)
	if err != nil {
		return
	}
	defer removeFile(outputPath)

	return
}

func (g *GradingUsecase) downloadAndParse(srcURL string) (filePath string, results []string, err error) {
	filePath, err = download(srcURL, nil)
	if err != nil {
		err = fmt.Errorf("unable to download from %s: %w", srcURL, err)
		return
	}

	results, err = g.parseFilePerLine(filePath)
	if err != nil {
		err = fmt.Errorf("unable to parse from %v: %w", filePath, err)
		return
	}

	return
}

// read file and append each line in slices of string
func (g *GradingUsecase) parseFilePerLine(filePath string) ([]string, error) {
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

// if dest is nil, the file will be download
// into temp with md5 hashed srcURL as file name
func download(srcURL string, dest *string) (outputPath string, err error) {
	resp, err := http.Get(srcURL)
	if err != nil {
		err = fmt.Errorf("unable to get from url: %w", err)
		return
	}
	defer resp.Body.Close()

	if dest == nil {
		hash := md5.New()
		hashed := hash.Sum([]byte(srcURL))
		outputPath = fmt.Sprintf(path.Join(os.TempDir(), string(hashed)))
	}

	out, err := os.Create(outputPath)
	if err != nil {
		err = fmt.Errorf("unable to create output file: %w", err)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		err = fmt.Errorf("unable to write response to file: %w", err)
		return
	}

	return
}

func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logrus.Error(err)
	}
}
