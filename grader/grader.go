package grader

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/miun173/autograd/model"
	"github.com/sirupsen/logrus"
)

// SubmissionUsecase ..
type SubmissionUsecase interface {
	FindByID(ctx context.Context, id int64) (*model.Submission, error)
	UpdateGradeByID(ctx context.Context, id, grade int64) error
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
func NewGrader(c Compiler, opts ...Option) *Grader {
	g := &Grader{
		compiler: c,
	}

	for _, opt := range opts {
		opt(g)
	}

	return g
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
	submission, err := g.submisisonUsecase.FindByID(context.Background(), submissionID)
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

	_, corrects, err := g.Grade(srcCodePath, asg.input, asg.output)
	if err != nil {
		logrus.Error(err)
		return err
	}

	err = g.submisisonUsecase.UpdateGradeByID(context.Background(), submissionID, calcCorrects(corrects))
	if err != nil {
		logrus.Error(err)
		return err
	}

	return
}

func calcCorrects(corrects []bool) (sum int64) {
	for _, c := range corrects {
		if c == true {
			sum++
		}
	}

	return (sum / int64(len(corrects))) * 100
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

	if res == nil {
		return nil, errors.New("unable to find assignment")
	}

	asg = &assignment{}

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

	theURL, err := url.Parse(srcURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse url %v", err)
	}

	paths := strings.Split(theURL.Path, "/")
	if dest == nil {
		outputPath = path.Join(os.TempDir(), paths[len(paths)-1])
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
