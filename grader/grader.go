package grader

import (
	"bufio"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/miun173/autograd/model"
	"github.com/sirupsen/logrus"
)

// SubmissionUsecase ..
type SubmissionUsecase interface {
	FindByID(id int64) (*model.Submission, error)
}

type assignment struct {
	input  []string
	output []string
}

// AssignmentUsecase ..
type AssignmentUsecase interface {
	FindByID(ctx context.Context, id int64) (*model.Assignment, error)
}

// Grader implements worker.GraderUsecase
type Grader struct {
	compiler          Compiler
	submisisonUsecase SubmissionUsecase
	assignmentUsecase AssignmentUsecase
}

// NewGrader ..
func NewGrader(c Compiler) *Grader {
	return &Grader{
		compiler: c,
	}
}

// Grade ..
func (g *Grader) Grade(srcCodePath string, inputs, expecteds []string) (outputs []string, corrects []bool, err error) {
	outPath, err := g.compiler.Compile(srcCodePath)
	if err != nil {
		logrus.WithField("source", srcCodePath).Error(err)
		return
	}

	defer func() {
		if err := g.compiler.Remove(outPath); err != nil {
			logrus.WithField("path", outPath).Error(err)
		}
	}()

	if len(inputs) != len(expecteds) {
		return
	}

	for i := range inputs {
		input := inputs[i]
		expected := expecteds[i]

		out, err := g.compiler.Run(outPath, input)
		if err != nil {
			logrus.Error(err)
			return outputs, corrects, err
		}

		outputs = append(outputs, out)
		corrects = append(corrects, out == expected)
	}

	return
}

// GradeSubmission find the submission source code and call Grade
func (g *Grader) GradeSubmission(submissionID int64) (err error) {
	submission, err := g.submisisonUsecase.FindByID(submissionID)
	if err != nil {
		err = fmt.Errorf("unable to get submission %d: %w", submissionID, err)
		logrus.Error(err)
		return
	}

	srcCodePath, err := g.download(submission.FileURL, nil)
	if err != nil {
		err = fmt.Errorf("unable to download submission %d: %w", submissionID, err)
		logrus.Error(err)
		return
	}

	asg, err := g.getAssignment(submission.AssignmentID)
	if err != nil {
		err = fmt.Errorf("unable to get assignment for submission %d: %w", submissionID, err)
		logrus.Error(err)
		return
	}

	_, _, err = g.Grade(srcCodePath, asg.input, asg.output)
	return
}

func (g *Grader) getSubmissionSrcCodeByID(id int64) (srcCodePath string, err error) {
	return
}

// find the model.Assignment from usecase
// then download the input & output code to local path
func (g *Grader) getAssignment(assignmentID int64) (asg *assignment, err error) {
	res, err := g.assignmentUsecase.FindByID(context.Background(), assignmentID)
	if err != nil {
		return nil, err
	}

	// download source code to local filepath
	inputPath := ""
	inputPath, asg.input, err = g.parseAssignmentIO(res.CaseInputFileURL)
	if err != nil {
		return
	}
	defer g.removeFile(inputPath)

	outputPath := ""
	outputPath, asg.output, err = g.parseAssignmentIO(res.CaseOutputFileURL)
	if err != nil {
		return
	}
	defer g.removeFile(outputPath)

	return
}

// use for downloading & parsing the assingment's input & output code
func (g *Grader) parseAssignmentIO(srcURL string) (filePath string, results []string, err error) {
	filePath, err = g.download(srcURL, nil)
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

// download srcURL into local file
// if dest is nil, the file will be download into temp with md5 hashed srcURL as file name
func (g *Grader) download(srcURL string, dest *string) (outputPath string, err error) {
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

// read file and append each line in slices of string
func (g *Grader) parseFilePerLine(filePath string) ([]string, error) {
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

func (g *Grader) removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logrus.Error(err)
	}
}
