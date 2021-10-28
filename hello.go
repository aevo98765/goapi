package main


func main() {
	poi := Protein{"vps35", "9606"}
	requestBytes := poi.BiogridAPIRequest()
	proteomeResults := poi.CollectingDataFromAPIRequest(requestBytes)
	results := poi.CollectPOIInteractors(proteomeResults)
	totalHits := 0
	for x := range results {
		println(x)
		totalHits++
	}
	println(totalHits)
}
