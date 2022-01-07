package command

import (
	"fmt"
	"io/ioutil"

	"github.com/goccy/go-yaml"

	"github.com/rozifus/foref/cmd"
	"github.com/rozifus/foref/cmd/runner"
	"github.com/rozifus/foref/pkg/source"
)


type ImportCmd struct {
	Namespace  string   `kong:"flag,short='n',default='DEFAULT',help='Which folder namespace to use.'"`
	Path  string   		`kong:"flag,short='p',help='Which folder path to use.'"`
	IdentifierFiles []string `kong:"arg,optional,type='path'"`
}

func (cmd *ImportCmd) Run(ctx *cmd.Context) error {
	path, err := runner.GetNamespacePath(cmd.Namespace, cmd.Path)
	if err != nil {
		return err
	}

	gDatas := make([]*ForefData, 0)
	for _, gFile := range cmd.IdentifierFiles {
		gDatas = append(gDatas, readForefFile(gFile))
	}

	identifiers := make([]string, 0)
	for _, gData := range gDatas {
		identifiers = append(identifiers, gData.Identifiers...)
	}

	sourceCtx := &source.Context{
		NamespacePath: path,
		Source: "github.com", //TODO: source from flag
	}

	runner.CollectIdentifiers(sourceCtx, identifiers)

	return nil
}

func readForefFile(path string) *ForefData {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to open file: %v", path)
		fmt.Printf("%v", err)
		return nil
	}

	var gd ForefData
	if err = yaml.Unmarshal(data, &gd); err != nil {
		fmt.Printf("Failed to parse yaml in: %v", path)
		fmt.Printf("%v", err)
		return nil
	}

	return &gd
}

