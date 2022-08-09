package names

import (
	"bufio"
	"math/rand"
	"os"
)

/**
 * NOTE: Defined as an interface so we can use more complex methods in the
 * future(markov chain, etc).
 */
type NameGenerator interface {
	GenerateName() string
}

type SelectorNameGenerator struct {
	names []string
}

func NewSelectorNameGenerator(inputFile string) *SelectorNameGenerator {
	nameGenerator := &SelectorNameGenerator{
		names: make([]string, 0),
	}

	readFile, err := os.Open(inputFile)
	defer readFile.Close()

	if err != nil {
		panic(err.Error())
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	for fileScanner.Scan() {
		inputName := fileScanner.Text()
		nameGenerator.names = append(nameGenerator.names, inputName)
	}

	return nameGenerator
}

func (generator *SelectorNameGenerator) GenerateName() string {
	return generator.names[rand.Intn(len(generator.names))]
}
