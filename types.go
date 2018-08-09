package run4labelled

type (
	// Configuration contains all data required to run the program
	Configuration struct {
		Label    string   // Label is the name of the file to look for
		Excludes []string // Excludes is a list of directory names to ignore
		Run      string   // Run is executed in each labelled directory

		sendChannel <-chan Execute // sendChannel is used to tell executor to xecute other procexx
	}

	// Execute is sent to executor with relevant information
	Execute struct {
		Directory string
		Command   []string
	}
)
