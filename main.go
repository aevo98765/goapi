package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"sync"
)

func main() {
	lambda.Start(LambdaHandler)
}

//type LambdaEvent struct {
//	Gene_name string `json:"gene_name"`
//	Organism_id string `json:"organism_id"`
//}


type LambdaResponse struct {
	Message string `json:"message"`
	Body string `json:"body"`
}


func LambdaHandler(event interface{}) (LambdaResponse, error) {
	fmt.Println("Triggered")
	fmt.Println(event)
	data, ok := event.(map[string]interface{})
	fmt.Println(data["queryStringParameters"])
	inputs, ok := data["queryStringParameters"].(map[string]interface{})
	fmt.Println(inputs)
	fmt.Println(inputs["gene_name"])
	fmt.Println(inputs["organism_id"])
	fmt.Println(ok)
	parentProtein := Protein{}
	parentProtein.Gene_name = fmt.Sprintf("%v", inputs["gene_name"])
	parentProtein.OrganismID = fmt.Sprintf("%v", inputs["organism_id"])
	fmt.Println("assigned values")
	fmt.Println(parentProtein.Gene_name)
	fmt.Println(parentProtein.OrganismID)
	requestBytes := parentProtein.BiogridAPIRequest()
	proteomeJSONResults := parentProtein.CollectingDataFromAPIRequest(requestBytes)
	parentProtein.GenerateProteome(proteomeJSONResults)

	finalData := FinalData{}
	finalData.ChildData = make(map[string]Protein)

	var wg sync.WaitGroup

	for resultProteinName := range parentProtein.Proteome {
		wg.Add(1)
		go childProteomeAddedToFinalData(resultProteinName, parentProtein, finalData, &wg)
	}
	wg.Wait()
	childDataArray := make([]Protein, 0)
	for _, s := range finalData.ChildData {
		childDataArray = append(childDataArray, s)
	}

	finalData.ChildDataArray = childDataArray

	jsonResults, err := json.Marshal(finalData)
	if err != nil {
		fmt.Println(err)
		return LambdaResponse {
			Message: fmt.Sprintf("Failed"),
			Body: "nil",
		}, nil
	}

	jsonString := string(jsonResults)

	var results FinalData
	json.Unmarshal(jsonResults, &results)

	for _, child := range results.ChildData {
		fmt.Println(child.Gene_name)
		fmt.Println(child.NumberOfHitsInCommonWithParent)
		fmt.Println(child.ProportionOfHitsInCommonWithParent)
		fmt.Println(child.HitsInCommonWithParent)
		fmt.Println()

	}


	return LambdaResponse {
		Message: fmt.Sprintf("Success"),
		Body: jsonString,
	}, nil
}

func childProteomeAddedToFinalData(resultProteinName string, parentProtein Protein, finalData FinalData, wg *sync.WaitGroup) {
	childProtein := Protein{}
	childProtein.Gene_name = resultProteinName
	childProtein.OrganismID = parentProtein.getOrganismID()
	requestBytes := childProtein.BiogridAPIRequest()
	proteomeJSONResults := childProtein.CollectingDataFromAPIRequest(requestBytes)
	childProtein.GenerateProteome(proteomeJSONResults)
	proteinsInCommon := make([]string, 0)
	if len(childProtein.Proteome) > 0 {
		childProtein.NumberOfHitsInCommonWithParent = 0
		for childResultProtein := range childProtein.Proteome {
			if parentProtein.Proteome[childResultProtein] == true {
				childProtein.NumberOfHitsInCommonWithParent++
				proteinsInCommon = append(proteinsInCommon, childResultProtein)

			}
		}
		childProtein.HitsInCommonWithParent = proteinsInCommon
		childProtein.ProportionOfHitsInCommonWithParent = float32(childProtein.NumberOfHitsInCommonWithParent) / float32(len(childProtein.Proteome))
		if len(proteinsInCommon) > 0 {
			finalData.ChildData[childProtein.Gene_name] = childProtein
		}
	}
	wg.Done()
}