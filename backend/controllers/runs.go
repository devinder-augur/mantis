package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/diggerhq/digger/backend/middleware"
	"github.com/diggerhq/digger/backend/models"
	"github.com/diggerhq/digger/backend/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

func RunsForProject(c *gin.Context) {
	currentOrg, exists := c.Get(middleware.ORGANISATION_ID_KEY)
	projectIdStr := c.Param("project_id")

	if projectIdStr == "" {
		c.String(http.StatusBadRequest, "ProjectId not specified")
		return
	}

	projectId, err := strconv.Atoi(projectIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid ProjectId")
		return
	}

	if !exists {
		c.String(http.StatusForbidden, "Not allowed to access this resource")
		return
	}

	var org models.Organisation
	err = models.DB.GormDB.Where("id = ?", currentOrg).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, fmt.Sprintf("Could not find organisation: %v", currentOrg))
		} else {
			c.String(http.StatusInternalServerError, "Unknown error occurred while fetching database")
		}
		return
	}

	project, err := models.DB.GetProject(uint(projectId))
	if err != nil {
		log.Printf("could not fetch project: %v", err)
		c.String(http.StatusInternalServerError, "Could not fetch project")
		return
	}

	if project.OrganisationID != org.ID {
		log.Printf("Forbidden access: not allowed to access projectID: %v logged in org: %v", project.OrganisationID, org.ID)
		c.String(http.StatusForbidden, "No access to this project")
		return

	}

	runs, err := models.DB.ListDiggerRunsForProject(project.Name, project.Repo.ID)
	if err != nil {
		log.Printf("could not fetch runs: %v", err)
		c.String(http.StatusInternalServerError, "Could not fetch runs")
		return
	}

	serializedRuns := make([]interface{}, 0)
	for _, run := range runs {
		serializedRun, err := run.MapToJsonStruct()
		if err != nil {
			log.Printf("could not unmarshal run: %v", err)
			c.String(http.StatusInternalServerError, "Could not unmarshal runs")
			return
		}
		serializedRuns = append(serializedRuns, serializedRun)
	}
	response := make(map[string]interface{})
	response["runs"] = serializedRuns
	c.JSON(http.StatusOK, response)
}

func GetRuns(c *gin.Context) {
	runs, done := models.DB.GetProjectRunsFromContext(c, middleware.ORGANISATION_ID_KEY)
	if !done {
		return
	}

	pageContext := services.GetMessages(c)
	maps.Copy(pageContext, gin.H{
		"Runs": runs,
	})
	c.JSON(http.StatusOK, gin.H{"data": runs})
}

func RunDetailsPage(c *gin.Context) {
	runId64, err := strconv.ParseUint(c.Param("runid"), 10, 32)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse project run id")
		return
	}
	runId := uint(runId64)
	run, ok := models.DB.GetProjectByRunId(c, runId, middleware.ORGANISATION_ID_KEY)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": run})
	// stateSyncOutput := ""
	// terraformPlanOutput := ""
	// runOutput := string(ansihtml.ConvertToHTMLWithClasses([]byte(run.Output), "terraform-output-", true))
	// runOutput = strings.Replace(runOutput, "  ", "&nbsp;&nbsp;", -1)
	// runOutput = strings.Replace(runOutput, "\n", "<br>\n", -1)

	// planIndex := strings.Index(runOutput, "Terraform used the selected providers to generate the following execution")
	// if planIndex != -1 {
	// 	stateSyncOutput = runOutput[:planIndex]
	// 	terraformPlanOutput = runOutput[planIndex:]

	// 	pageContext := services.GetMessages(c)
	// 	maps.Copy(pageContext, gin.H{
	// 		"Run":                      run,
	// 		"TerraformStateSyncOutput": template.HTML(stateSyncOutput),
	// 		"TerraformPlanOutput":      template.HTML(terraformPlanOutput),
	// 	})
	// 	c.HTML(http.StatusOK, "run_details.tmpl", pageContext)
	// } else {
	// 	pageContext := services.GetMessages(c)
	// 	maps.Copy(pageContext, gin.H{
	// 		"Run":       run,
	// 		"RunOutput": template.HTML(runOutput),
	// 	})
	// 	c.HTML(http.StatusOK, "run_details.tmpl", pageContext)
	// }
}
func RunDetails(c *gin.Context) {
	currentOrg, exists := c.Get(middleware.ORGANISATION_ID_KEY)
	runIdStr := c.Param("run_id")

	if runIdStr == "" {
		c.String(http.StatusBadRequest, "RunID not specified")
		return
	}

	runId, err := strconv.Atoi(runIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid RunId")
		return
	}

	if !exists {
		c.String(http.StatusForbidden, "Not allowed to access this resource")
		return
	}

	var org models.Organisation
	err = models.DB.GormDB.Where("id = ?", currentOrg).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, fmt.Sprintf("Could not find organisation: %v", currentOrg))
		} else {
			c.String(http.StatusInternalServerError, "Unknown error occurred while fetching database")
		}
		return
	}

	run, err := models.DB.GetDiggerRun(uint(runId))
	if err != nil {
		log.Printf("Could not fetch run: %v", err)
		c.String(http.StatusBadRequest, "Could not fetch run, please check that it exists")
	}
	if run.Repo.OrganisationID != org.ID {
		c.String(http.StatusForbidden, "Not allowed to access this resource")
		return
	}

	response, err := run.MapToJsonStruct()
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not unmarshall data")
		return
	}
	c.JSON(http.StatusOK, response)
}

func ApproveRun(c *gin.Context) {

	currentOrg, exists := c.Get(middleware.ORGANISATION_ID_KEY)
	runIdStr := c.Param("run_id")

	if runIdStr == "" {
		c.String(http.StatusBadRequest, "RunID not specified")
		return
	}

	runId, err := strconv.Atoi(runIdStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid RunId")
		return
	}

	if !exists {
		c.String(http.StatusForbidden, "Not allowed to access this resource")
		return
	}

	var org models.Organisation
	err = models.DB.GormDB.Where("id = ?", currentOrg).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.String(http.StatusNotFound, fmt.Sprintf("Could not find organisation: %v", currentOrg))
		} else {
			c.String(http.StatusInternalServerError, "Unknown error occurred while fetching database")
		}
		return
	}

	run, err := models.DB.GetDiggerRun(uint(runId))
	if err != nil {
		log.Printf("Could not fetch run: %v", err)
		c.String(http.StatusBadRequest, "Could not fetch run, please check that it exists")
	}
	if run.Repo.OrganisationID != org.ID {
		c.String(http.StatusForbidden, "Not allowed to access this resource")
		return
	}

	if run.Status != models.RunPendingApproval {
		log.Printf("Run status not ready for approval: %v", run.ID)
		c.String(http.StatusBadRequest, "Approval not possible for run (%v) because status is %v", run.ID, run.Status)
		return
	}

	if run.IsApproved == false {
		run.ApprovalAuthor = "a_user"
		run.IsApproved = true
		run.ApprovalDate = time.Now()
		err := models.DB.UpdateDiggerRun(run)
		if err != nil {
			log.Printf("Could update run: %v", err)
			c.String(http.StatusInternalServerError, "Could not update approval")
		}
	} else {
		log.Printf("Run has already been approved")
	}

	response, err := run.MapToJsonStruct()
	if err != nil {
		c.String(http.StatusInternalServerError, "Could not unmarshall data")
		return
	}
	c.JSON(http.StatusOK, response)
}
