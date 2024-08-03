package backend

import (
	"github.com/diggerhq/digger/cli/pkg/core/execution"
	"github.com/diggerhq/digger/libs/orchestrator/scheduler"
	"time"
)

type RunDetails struct {
	Id        int    `json:"id"`
	StartedAt string `json:"startedAt"`
	EndedAt   string `json:"endedAt"`
	Status    string `json:"status"`
	Command   string `json:"command"`
	Output    string `json:"output"`
}
type Api interface {
	ReportProject(repo string, projectName string, configuration string) error
	ReportProjectRun(repo string, projectName string, startedAt time.Time, endedAt time.Time, status string, command string, output string) (RunDetails, error)
	ReportProjectJobStatus(repo string, projectName string, jobId string, status string, timestamp time.Time, summary *execution.DiggerExecutorPlanResult, PrCommentUrl string, terraformOutput string) (*scheduler.SerializedBatch, error)
}
