package pipeline

import (
	"fmt"
	"plex_monitor/internal/database"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// PipelineCollection is the name of the collection that stores the pipelines
	PipelineCollection = "pipelines"

	// StatusSuccess is the status for a successful step.
	StatusSuccess = "success"

	// StatusFailed is the status for a failed step.
	StatusFailed = "failed"

	// StatusSkipped is the status for a skipped step.
	StatusSkipped = "skipped"

	// StatusPending is the status for a pending step.
	StatusPending = "pending"

	// StatusUnknown is the status for an unknown step.
	StatusUnknown = "unknown"
)

// Pipeline is the struct that represents all pipelines in the application.
type Pipeline struct {
	// ID is the ID of the pipeline.
	ID string `bson:"_id" json:"id"`

	// Name is the name of the pipeline.
	Name string `bson:"name" json:"name"`

	// Steps is the list of steps in the pipeline.
	Steps []Step `bson:"steps" json:"steps"`

	// CurrentStep is the current step in the pipeline.
	CurrentStep int `bson:"current_step" json:"current_step"`

	// Metadata is the metadata for the pipeline.
	Metadata map[string]interface{} `bson:"metadata" json:"metadata"`
}

// Step is the struct that represents a step in the pipeline.
type Step struct {
	// Name is the name of the step.
	Name string `bson:"name" json:"name"`

	// Key is the key for the step.
	Key string `bson:"key" json:"key"`

	// Status is the status of the step.
	Status string `bson:"status" json:"status"`

	// CompletedAt is the time that the step was completed.
	CompletedAt *time.Time `bson:"completed_at" json:"completed_at"`

	// StartedAt is the time that the step was started.
	StartedAt time.Time `bson:"started_at" json:"started_at"`

	// Function is the function that is executed for the step.
	Function func() error `bson:"-" json:"-"`
}

// CreatePipeline creates a new pipeline in the database.
func CreatePipeline(data Pipeline) (id string, err error) {
	result, err := database.DB.Collection(PipelineCollection).InsertOne(database.Ctx, data)
	return result.InsertedID.(string), err
}

// GetPipelineByID returns a pipeline by ID.
func GetPipelineByID(id string) (*Pipeline, error) {
	var pipeline Pipeline

	err := database.DB.Collection(PipelineCollection).FindOne(database.Ctx, bson.M{"_id": id}).Decode(&pipeline)
	if err != nil {
		return &Pipeline{}, err
	}

	return &pipeline, nil
}

// GeneratePipelineID generates a new pipeline ID based on the request user name, the media type, and the media title.
func GeneratePipelineID(mediaType, mediaTitle string) string {
	return fmt.Sprintf("%so_o%s", mediaType, mediaTitle)
}

// NewPipeline returns a new pipeline.
func NewPipeline(name string, logger *logrus.Entry) *Pipeline {
	return &Pipeline{
		Name: name,
	}
}

// AddStep adds a step to the pipeline.
func (p *Pipeline) AddStep(name, key string, function func() error) {
	p.Steps = append(p.Steps, Step{
		Name:      name,
		Key:       key,
		Status:    StatusPending,
		Function:  function,
		StartedAt: time.Now(),
	})
}

// GetStepByKey returns a step by key.
// Returns the step, step index, and error
func (p *Pipeline) GetStepByKey(key string) (*Step, int, error) {
	for idx, step := range p.Steps {
		if step.Key == key {
			return &step, idx, nil
		}
	}

	return nil, -1, fmt.Errorf("step with key %s does not exist", key)
}

// RunStep runs a specific step in the pipeline by index.
func (p *Pipeline) RunStep(key string) error {
	logrus.Infof("Running step %s", key)
	step, idx, err := p.GetStepByKey(key)
	if err != nil {
		return err
	}

	// Check if the step has already been completed
	if step.Status == StatusSuccess {
		// no-op
		return nil
	}

	// Set the current step
	p.CurrentStep = idx

	// Run the step
	err = step.Function()
	if err != nil {
		step.Status = StatusFailed
		return err
	}

	// Set the step as completed, and set the completed time, then save the step in the pipeline
	completedTime := time.Now()
	step.Status = StatusSuccess
	step.CompletedAt = &completedTime
	p.Steps[idx] = *step

	// Save the pipeline to the database
	err = p.Save()
	if err != nil {
		return err
	}

	logrus.Infof("Finished running step %s", key)

	return nil
}

// MarkStepAsSkipped marks a step as skipped.
func (p *Pipeline) MarkStepAsSkipped(key string) error {
	step, idx, err := p.GetStepByKey(key)
	if err != nil {
		return err
	}

	// Set the step as skipped, then save the step in the pipeline
	step.Status = StatusSkipped
	p.Steps[idx] = *step

	// Set the current step
	p.CurrentStep = idx

	// Save the pipeline to the database
	err = p.Save()
	if err != nil {
		return err
	}

	return nil
}

// AddMetadata adds metadata to the pipeline.
func (p *Pipeline) AddMetadata(key string, value interface{}) {
	if p.Metadata == nil {
		p.Metadata = make(map[string]interface{})
	}
	p.Metadata[key] = value
}

// Save saves the pipeline to the database.
func (p *Pipeline) Save() error {
	// Save the pipeline to the database.
	opts := options.Update().SetUpsert(true)
	_, err := database.DB.Collection(PipelineCollection).UpdateOne(database.Ctx, bson.M{"_id": p.ID}, bson.M{"$set": p}, opts)
	if err != nil {
		return err
	}

	return nil
}
