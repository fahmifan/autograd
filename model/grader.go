package model

// GradingArg ..
type GradingArg struct {
	SourceCodePath string
	Expecteds      []string
	Inputs         []string
}

// GradingResult ..
type GradingResult struct {
	Outputs  []string
	Corrects []bool
}

// Score ..
func (g *GradingResult) Score() (sum int64) {
	for _, c := range g.Corrects {
		if c == true {
			sum++
		}
	}

	return (sum * 100 / int64(len(g.Corrects)))
}

// GraderEngine describe a grader engine
type GraderEngine interface {
	Grade(arg *GradingArg) (*GradingResult, error)
	Compile(inputPath string) (outPath string, err error)
	Run(source, input string) (out string, err error)
	Remove(source string) error
}

// GraderUsecase ..
type GraderUsecase interface {
	GradeBySubmission(submissionID string) error
}
