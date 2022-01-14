package basecli

import (
	"errors"
	"io/ioutil"

	"github.com/go-playground/validator/v10"
	"github.com/iter8-tools/iter8/base"
	"github.com/iter8-tools/iter8/base/log"
	"sigs.k8s.io/yaml"
)

// Experiment type that includes a list of runnable tasks derived from the experiment spec
type Experiment struct {
	tasks []base.Task
	*base.Experiment
}

// ExpIO enables interacting with experiment spec and result stored externally
type ExpIO interface {
	// ReadSpec reads the experiment spec
	ReadSpec() ([]base.TaskSpec, error)
	// ReadResult reads the experiment results
	ReadResult() (*base.ExperimentResult, error)
	// WriteResult writes the experimeent results
	WriteResult(r *Experiment) error
}

const (
	experimentSpecPath   = "experiment.yaml"
	experimentResultPath = "result.yaml"
)

// Build an experiment
func Build(withResult bool, expio ExpIO) (*Experiment, error) {
	e := &Experiment{
		Experiment: &base.Experiment{},
	}
	var err error
	// read it in
	log.Logger.Trace("build started")
	e.Tasks, err = expio.ReadSpec()
	if err != nil {
		return nil, err
	}
	e.InitResults()
	if withResult {
		e.Result, err = expio.ReadResult()
		if err != nil {
			return nil, err
		}
	}

	err = e.buildTasks()
	if err != nil {
		return nil, err
	}

	return e, err
}

// build experiment tasks
func (e *Experiment) buildTasks() error {
	for _, t := range e.Tasks {
		if (t.Task == nil || len(*t.Task) == 0) && (t.Run == nil) {
			log.Logger.Error("invalid task found without a task name or a run command")
			return errors.New("invalid task found without a task name or a run command")
		}

		var err error
		var task base.Task

		validate := validator.New()
		err = validate.Struct(t)
		if err != nil {
			log.Logger.WithStackTrace(err.Error()).Error("invalid task specification")
			return err
		}

		// this is a run task
		if t.Run != nil {
			task, err = base.MakeRun(&t)
			e.tasks = append(e.tasks, task)
			if err != nil {
				return err
			}
		} else {
			// this is some other task
			switch *t.Task {
			case base.CollectTaskName:
				task, err = base.MakeCollect(&t)
				e.tasks = append(e.tasks, task)
			case base.AssessTaskName:
				task, err = base.MakeAssess(&t)
				e.tasks = append(e.tasks, task)
			default:
				log.Logger.Error("unknown task: " + *t.Task)
				return errors.New("unknown task: " + *t.Task)
			}

			if err != nil {
				return err
			}
		}
	}
	return nil
}

//FileExpIO enables reading and writing experiment spec and result files
type FileExpIO struct{}

// SpecFromBytes reads experiment spec from bytes
func SpecFromBytes(b []byte) ([]base.TaskSpec, error) {
	e := []base.TaskSpec{}
	err := yaml.Unmarshal(b, &e)
	if err != nil {
		log.Logger.WithStackTrace(err.Error()).Error("unable to unmarshal experiment spec")
		return nil, err
	}
	return e, err
}

// ResultFromBytes reads experiment result from bytes
func ResultFromBytes(b []byte) (*base.ExperimentResult, error) {
	r := &base.ExperimentResult{}
	err := yaml.Unmarshal(b, r)
	if err != nil {
		log.Logger.WithStackTrace(err.Error()).Error("unable to unmarshal experiment result")
		return nil, err
	}
	return r, err
}

// ReadSpec reads experiment spec from file
func (f *FileExpIO) ReadSpec() ([]base.TaskSpec, error) {
	b, err := ioutil.ReadFile(experimentSpecPath)
	if err != nil {
		log.Logger.WithStackTrace(err.Error()).Error("unable to read experiment spec")
		return nil, errors.New("unable to read experiment spec")
	}
	return SpecFromBytes(b)
}

// ReadResult reads experiment result from file
func (f *FileExpIO) ReadResult() (*base.ExperimentResult, error) {
	b, err := ioutil.ReadFile(experimentResultPath)
	if err != nil {
		log.Logger.WithStackTrace(err.Error()).Error("unable to read experiment result")
		return nil, errors.New("unable to read experiment result")
	}
	return ResultFromBytes(b)
}

// WriteResult writes experiment result to file
func (f *FileExpIO) WriteResult(r *Experiment) error {
	rBytes, _ := yaml.Marshal(r.Result)
	err := ioutil.WriteFile(experimentResultPath, rBytes, 0664)
	if err != nil {
		log.Logger.WithStackTrace(err.Error()).Error("unable to write experiment result")
		return err
	}
	return err
}

// Completed returns true if the experiment is complete
// if the result stanza is missing, this function returns false
func (exp *Experiment) Completed() bool {
	if exp != nil {
		if exp.Result != nil {
			if exp.Result.NumCompletedTasks == len(exp.Tasks) {
				return true
			}
		}
	}
	return false
}

// NoFailure returns true if no task int he experiment has failed
// if the result stanza is missing, this function returns false
func (exp *Experiment) NoFailure() bool {
	if exp != nil {
		if exp.Result != nil {
			if !exp.Result.Failure {
				return true
			}
		}
	}
	return false
}