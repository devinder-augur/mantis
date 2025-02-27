package utils

import (
	"time"

	"github.com/diggerhq/digger/cli/pkg/core/backend"
	"github.com/diggerhq/digger/cli/pkg/core/execution"
	"github.com/diggerhq/digger/libs/orchestrator/scheduler"

	"github.com/diggerhq/digger/libs/orchestrator"
)

type MockTerraform struct {
	commands []string
}

func (tf *MockTerraform) Apply() (string, string, error) {
	tf.commands = append(tf.commands, "apply")
	return "", "", nil
}

func (tf *MockTerraform) Plan() (bool, string, string, error) {
	tf.commands = append(tf.commands, "plan")
	return true, "", "", nil
}

type MockPolicyChecker struct {
}

func (t MockPolicyChecker) CheckAccessPolicy(ciService orchestrator.OrgService, prService *orchestrator.PullRequestService, SCMOrganisation string, SCMrepository string, projectName string, projectDir string, command string, prNumber *int, requestedBy string, planPolicyViolations []string) (bool, error) {
	return false, nil
}

func (t MockPolicyChecker) CheckPlanPolicy(SCMrepository string, SCMOrganisation string, projectname string, projectDir string, planOutput string) (bool, []string, error) {
	return false, nil, nil
}

func (t MockPolicyChecker) CheckDriftPolicy(SCMOrganisation string, SCMrepository string, projectname string) (bool, error) {
	return true, nil
}

type MockPullRequestManager struct {
	ChangedFiles []string
	Teams        []string
	Approvals    []string
}

func (t MockPullRequestManager) GetUserTeams(organisation string, user string) ([]string, error) {
	return t.Teams, nil
}

func (t MockPullRequestManager) GetChangedFiles(prNumber int) ([]string, error) {
	return t.ChangedFiles, nil
}
func (t MockPullRequestManager) PublishComment(prNumber int, comment string) (*orchestrator.Comment, error) {
	return nil, nil
}

func (t MockPullRequestManager) ListIssues() ([]*orchestrator.Issue, error) {
	return nil, nil
}

func (t MockPullRequestManager) PublishIssue(title string, body string) (int64, error) {
	return 0, nil
}

func (t MockPullRequestManager) SetStatus(prNumber int, status string, statusContext string) error {
	return nil
}

func (t MockPullRequestManager) GetCombinedPullRequestStatus(prNumber int) (string, error) {
	return "", nil
}

func (t MockPullRequestManager) GetApprovals(prNumber int) ([]string, error) {
	return t.Approvals, nil
}

func (t MockPullRequestManager) MergePullRequest(prNumber int) error {
	return nil
}

func (t MockPullRequestManager) IsMergeable(prNumber int) (bool, error) {
	return true, nil
}

func (t MockPullRequestManager) IsMerged(prNumber int) (bool, error) {
	return false, nil
}

func (t MockPullRequestManager) DownloadLatestPlans(prNumber int) (string, error) {
	return "", nil
}

func (t MockPullRequestManager) IsClosed(prNumber int) (bool, error) {
	return false, nil
}

func (t MockPullRequestManager) GetComments(prNumber int) ([]orchestrator.Comment, error) {
	return []orchestrator.Comment{}, nil
}

func (t MockPullRequestManager) EditComment(prNumber int, commentId interface{}, comment string) error {
	return nil
}

func (t *MockPullRequestManager) CreateCommentReaction(id interface{}, reaction string) error {
	return nil
}

func (t MockPullRequestManager) GetBranchName(prNumber int) (string, string, error) {
	return "", "", nil
}

func (t MockPullRequestManager) SetOutput(prNumber int, key string, value string) error {
	return nil
}

type MockPlanStorage struct {
}

func (t *MockPlanStorage) StorePlanFile(fileContents []byte, artifactName string, fileName string) error {
	return nil
}

func (t MockPlanStorage) RetrievePlan(localPlanFilePath string, artifactName string, storedPlanFilePath string) (*string, error) {
	return nil, nil
}

func (t MockPlanStorage) DeleteStoredPlan(artifactName string, storedPlanFilePath string) error {
	return nil
}

func (t MockPlanStorage) PlanExists(artifactName string, storedPlanFilePath string) (bool, error) {
	return false, nil
}

type MockBackendApi struct {
}

func (t MockBackendApi) ReportProject(namespace string, projectName string, configuration string) error {
	return nil
}

func (t MockBackendApi) ReportProjectRun(repo string, projectName string, startedAt time.Time, endedAt time.Time, status string, command string, output string) (backend.RunDetails, error) {
	return backend.RunDetails{}, nil
}

func (t MockBackendApi) ReportProjectJobStatus(repo string, projectName string, jobId string, status string, timestamp time.Time, summary *execution.DiggerExecutorPlanResult, PrCommentUrl string, terraformOutput string) (*scheduler.SerializedBatch, error) {
	return nil, nil
}
