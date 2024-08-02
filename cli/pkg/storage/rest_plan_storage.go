package storage

import (
	"encoding/json"
	"fmt"
)

type PlanStorageRest struct {
	Endpoint string
	Method   string
	PrURl    string
}

type RestPlanBody struct {
	PRURL  string `json:"pr_url"`
	TfFile string `json:"tffile"`
}

func (psg *PlanStorageRest) PlanExists(artifactName string, storedPlanFilePath string) (bool, error) {

	return false, nil
}

func (psr *PlanStorageRest) StorePlanFile(fileContents []byte, artifactName string, fileName string) error {

	terraformPlanJson := string(fileContents)
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	jsonPayload := RestPlanBody{
		PRURL:  psr.PrURl,
		TfFile: terraformPlanJson,
	}
	jsonBody, err := json.Marshal(jsonPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	//doRequest(psr.Method, psr.Endpoint, headers, jsonBody)
	return nil
}

func (psg *PlanStorageRest) RetrievePlan(localPlanFilePath string, artifactName string, storedPlanFilePath string) (*string, error) {
	return nil, fmt.Errorf("unable to read data from Rest Endpoint")
}

func (psg *PlanStorageRest) DeleteStoredPlan(artifactName string, storedPlanFilePath string) error {
	return fmt.Errorf("unable to delete data on Rest Endpoint")
}
