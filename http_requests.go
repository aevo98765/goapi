package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ProteomeFullResults map[string]ProteomeJSONStructure

type Proteome map[string]bool

type ProteomeJSONStructure struct {
	BiogridInteractionID   int    `json:"BIOGRID_INTERACTION_ID"`
	EntrezGeneA            string `json:"ENTREZ_GENE_A"`
	EntrezGeneB            string `json:"ENTREZ_GENE_B"`
	BiogridIDA             int    `json:"BIOGRID_ID_A"`
	BiogridIDB             int    `json:"BIOGRID_ID_B"`
	SystematicNameA        string `json:"SYSTEMATIC_NAME_A"`
	SystematicNameB        string `json:"SYSTEMATIC_NAME_B"`
	OfficialSymbolA        string `json:"OFFICIAL_SYMBOL_A"`
	OfficialSymbolB        string `json:"OFFICIAL_SYMBOL_B"`
	SynonymsA              string `json:"SYNONYMS_A"`
	SynonymsB              string `json:"SYNONYMS_B"`
	ExperimentalSystem     string `json:"EXPERIMENTAL_SYSTEM"`
	ExperimentalSystemType string `json:"EXPERIMENTAL_SYSTEM_TYPE"`
	PubmedAuthor           string `json:"PUBMED_AUTHOR"`
	PubmedID               int    `json:"PUBMED_ID"`
	OrganismA              int    `json:"ORGANISM_A"`
	OrganismB              int    `json:"ORGANISM_B"`
	Throughput             string `json:"THROUGHPUT"`
	Quantitation           string `json:"QUANTITATION"`
	Modification           string `json:"MODIFICATION"`
	OntologyTerms          struct {
	Num332623 struct {
	OntologyTermID string `json:"ONTOLOGY_TERM_ID"`
	Name           string `json:"NAME"`
	TypeID         int    `json:"TYPE_ID"`
	TypeName       string `json:"TYPE_NAME"`
	Qualifiers     struct {
} `json:"QUALIFIERS"`
	Flag             string `json:"FLAG"`
	Desc             string `json:"DESC"`
	ID               int    `json:"ID"`
	OntologyCategory string `json:"ONTOLOGY_CATEGORY"`
} `json:"332623"`
} `json:"ONTOLOGY_TERMS"`
	Qualifications string `json:"QUALIFICATIONS"`
	Tags           string `json:"TAGS"`
	Sourcedb       string `json:"SOURCEDB"`
}


type FinalData struct {
	ParentName string `json:"parent_name"`
	ChildData map[string]Protein `json:"child_data"`
	ChildDataArray []Protein `json:"child_data_array,omitempty"`
}

type Protein struct {
	Gene_name string `json:"gene_name"`
	OrganismID string `json:"organism_id"`
	Proteome map[string]bool `json:"parent_proteome,omitempty"`
	NumberOfHitsInCommonWithParent int `json:"number_of_hits_in_common_with_parent,omitempty"`
	ProportionOfHitsInCommonWithParent float32 `json:"proportion_of_hits_in_common_with_parent"`
	HitsInCommonWithParent []string `json:"hits_in_common_with_parent,omitempty"`
}

func (protein *Protein) getOrganismID() string {
	return protein.OrganismID

}

func (protein *Protein) getAllInfo() {
	fmt.Println(protein.Gene_name)
	fmt.Println(protein.OrganismID)
	fmt.Println(protein.NumberOfHitsInCommonWithParent)
	fmt.Println(protein.Proteome)
	fmt.Println(protein.HitsInCommonWithParent)
}

func (protein *Protein) BiogridAPIRequest() []byte {
	biogrid_access_key := "160c9d5820c29e02dca794cacdbdee5f"
	fmt.Println("Just before requests")
	protein.getAllInfo()
	url := "https://webservice.thebiogrid.org/interactions?searchNames=true&geneList=" + protein.Gene_name + "&includeInteractors=true&format=json&max=1000&includeInteractorInteractions=false&taxId=" + protein.OrganismID + "&accesskey=" + biogrid_access_key


	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error with the response from Biogrid")
	}
	fmt.Println(response.StatusCode)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	return body
}

func (protein *Protein) CollectingDataFromAPIRequest(body []byte) ProteomeFullResults {
	var proteomeFullResults ProteomeFullResults
	json.Unmarshal(body, &proteomeFullResults)
	return proteomeFullResults
}

func (protein *Protein) GenerateProteome(proteomeFullResults ProteomeFullResults) {
	proteome := make(Proteome)
	for interactor := range proteomeFullResults {
		proteome[strings.ToUpper(proteomeFullResults[interactor].OfficialSymbolA)] = true
		proteome[strings.ToUpper(proteomeFullResults[interactor].OfficialSymbolB)] = true
	}

	protein.Proteome = proteome
}





