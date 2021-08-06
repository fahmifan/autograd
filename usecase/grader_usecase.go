package usecase

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/fahmifan/autograd/grader"
	"github.com/fahmifan/autograd/model"
	"github.com/sirupsen/logrus"
)

// make sure implement interface
var _ model.GraderUsecase = (*GraderUsecase)(nil)

type assignment struct {
	inputs    []string
	expecteds []string
}

// GraderUsecase implements worker.GraderUsecase
type GraderUsecase struct {
	submisisonUsecase model.SubmissionUsecase
	assignmentUsecase model.AssignmentUsecase
}

// NewGraderUsecase ..
func NewGraderUsecase(s model.SubmissionUsecase, a model.AssignmentUsecase) *GraderUsecase {
	g := &GraderUsecase{
		submisisonUsecase: s,
		assignmentUsecase: a,
	}
	return g
}

// GradeBySubmission find the submission source code and call Grade
func (s *GraderUsecase) GradeBySubmission(submissionID string) (err error) {
	submission, err := s.submisisonUsecase.FindByID(context.Background(), submissionID)
	if err != nil {
		err = fmt.Errorf("unable to get submission %s: %w", submissionID, err)
		logrus.Error(err)
		return
	}

	srcCodePath, err := s.download(submission.FileURL, nil)
	if err != nil {
		err = fmt.Errorf("unable to download submission %s: %w", submissionID, err)
		logrus.Error(err)
		return
	}

	asg, err := s.findAssignment(submission.AssignmentID)
	if err != nil {
		err = fmt.Errorf("unable to get assignment for submission %s: %w", submissionID, err)
		logrus.Error(err)
		return
	}

	grader := grader.New(grader.TypeCPP)
	res, err := grader.Grade(&model.GradingArg{
		SourceCodePath: srcCodePath,
		Expecteds:      asg.expecteds,
		Inputs:         asg.inputs,
	})
	if err != nil {
		logrus.WithField("srcCodePath", srcCodePath).Error(err)
		return err
	}

	err = s.submisisonUsecase.UpdateGradeByID(context.Background(), submissionID, res.Score())
	if err != nil {
		logrus.WithField("submissionID", submissionID).Error(err)
		return err
	}

	return
}

// find the model.Assignment from usecase
// then download the input & output code to local path
func (s *GraderUsecase) findAssignment(assignmentID string) (asg *assignment, err error) {
	res, err := s.assignmentUsecase.FindByID(context.Background(), assignmentID)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, errors.New("unable to find assignment")
	}

	asg = &assignment{}

	// download source code to local filepath
	inputPath := ""
	inputPath, asg.inputs, err = s.parseAssignmentIO(res.CaseInputFileURL)
	if err != nil {
		return
	}
	defer s.removeFile(inputPath)

	outputPath := ""
	outputPath, asg.expecteds, err = s.parseAssignmentIO(res.CaseOutputFileURL)
	if err != nil {
		return
	}
	defer s.removeFile(outputPath)

	return
}

// use for downloading & parsing the assingment's input & output code
func (s *GraderUsecase) parseAssignmentIO(srcURL string) (filePath string, results []string, err error) {
	filePath, err = s.download(srcURL, nil)
	if err != nil {
		err = fmt.Errorf("unable to download from %s: %w", srcURL, err)
		return
	}

	results, err = s.parseFilePerLine(filePath)
	if err != nil {
		err = fmt.Errorf("unable to parse from %v: %w", filePath, err)
		return
	}

	return
}

// download srcURL into local file
// if dest is nil, the file will be download into temp with md5 hashed srcURL as file name
func (s *GraderUsecase) download(srcURL string, dest *string) (outputPath string, err error) {
	logger := logrus.WithField("srcURL", srcURL)

	resp, err := http.Get(srcURL)
	if err != nil {
		logger.Error(err)
		err = fmt.Errorf("unable to get from url: %w", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bt, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
			return "", err
		}

		logger.WithField("body", string(bt)).Error(err)
		return "", err
	}

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
func (s *GraderUsecase) parseFilePerLine(filePath string) ([]string, error) {
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

func (s *GraderUsecase) removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		logrus.Error(err)
	}
}
