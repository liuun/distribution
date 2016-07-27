package allower

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"

	distribution "github.com/docker/distribution"
	"github.com/docker/distribution/context"
	// registrymiddleware "github.com/docker/distribution/registry/middleware/registry"
	repositorymiddleware "github.com/docker/distribution/registry/middleware/repository"
)

type AllowerConfiguration struct {
	Whitelist Whitelist `yaml:"whitelist"`
}

type Repository string
type Image string

type Whitelist struct {
	Repositories []Repository `yaml:"repositories"`
	Images       []Image      `yaml:"images"`
}

func (a *AllowerConfiguration) checkRepository(name string) bool {
	repoName := Repository(strings.Split(name, "/")[0])
	for _, r := range a.Whitelist.Repositories {
		if r == repoName {
			return true
		}
	}

	imageName := Image(name)
	for _, i := range a.Whitelist.Images {
		if i == imageName {
			return true
		}
	}

	return false
}

func readConfigurationFromOptions(options map[string]interface{}) (*AllowerConfiguration, error) {
	p, ok := options["path"]
	if !ok {
		return nil, fmt.Errorf("no allower config provided")
	}
	cPath, ok := p.(string)
	if !ok {
		return nil, fmt.Errorf("path must be a string")
	}

	cBytes, err := ioutil.ReadFile(cPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %s", err)
	}
	c := AllowerConfiguration{}
	yaml.Unmarshal(cBytes, &c)
	return &c, nil
}

func checkRepository(ctx context.Context, repository distribution.Repository, options map[string]interface{}) (distribution.Repository, error) {
	name := repository.Named().Name()
	fmt.Println("Hello " + name)

	conf, err := readConfigurationFromOptions(options)
	if err != nil {
		return nil, err
	}

	if !conf.checkRepository(name) {
		return nil, errors.New("Image not allowed to be pulled")
	}

	return repository, nil
}

// func checkRegistry(ctx context.Context, registry distribution.Namespace, options map[string]interface{}) (distribution.Namespace, error) {
// 	fmt.Println("Hello repo")
// 	return registry, nil
// }

func init() {
	repositorymiddleware.Register("allower", repositorymiddleware.InitFunc(checkRepository))
	// registrymiddleware.Register("allower", registrymiddleware.InitFunc(checkRegistry))

}
