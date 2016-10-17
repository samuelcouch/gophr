package github

import (
	"log"
	"sort"
	"sync"
	"time"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/model"
)

func fetchPackageVersions(metadata godocMetadata) ([]string, error) {
	refs, err := common.FetchRefs(metadata.author, metadata.repo)
	if err != nil {
		return nil, err
	}

	// Make sure the candidates are sorted first
	candidates := refs.Candidates
	sort.Sort(candidates)

	var versions []string
	for _, version := range candidates {
		versions = append(versions, version.String())
	}

	return versions, nil
}

func buildPackageModels(godocMetadataList []godocMetadata, awesomeGoIndex map[string]bool) ([]*models.PackageModel, error) {
	nbConcurrentGet := 20
	metadataChan := make(chan godocMetadata, nbConcurrentGet)
	packageModelChan := make(chan *models.PackageModel, nbConcurrentGet)

	var wg sync.WaitGroup
	for i := 0; i < nbConcurrentGet; i++ {
		wg.Add(1)
		go func() {
			for metadata := range metadataChan {
				log.Printf("Fetching package versions for \"%s\"", metadata.githubURL)
				packageVersions, err := fetchPackageVersions(metadata)
				if err == nil {
					log.Printf("\"%s\" versions retrieved successfully", metadata.githubURL)
					_, isAwesome := awesomeGoIndex[metadata.githubURL]
					packageModel, err := models.NewPackageModelForInsert(
						metadata.author,
						true,
						metadata.repo,
						packageVersions,
						"godoc.org/"+metadata.githubURL,
						time.Now(),
						isAwesome,
						metadata.description,
						0,
					)
					if err == nil {
						packageModelChan <- packageModel
					} else {
						packageModelChan <- nil
					}
				} else {
					log.Printf("\"%s\" versions failed to retrieve successfully", metadata.githubURL)
					packageModelChan <- nil
				}
			}
			log.Println("wait group done")
			wg.Done()
		}()
	}

	var packageModels []*models.PackageModel
	go func() {
		wg.Add(1)
		for i := 0; i < len(godocMetadataList); i++ {
			log.Printf("Waiting for %d out of %d of godocMetadataList length \n", i+1, len(godocMetadataList))
			packageModel := <-packageModelChan
			log.Println("Received one package from packageModelChan")
			if packageModel != nil {
				packageModels = append(packageModels, packageModel)
				log.Println("Appending package to packageModels slice")
			}
		}
		log.Println("wait group done")
		wg.Done()
	}()

	for _, metadata := range godocMetadataList {
		metadataChan <- metadata
		log.Println("Queuing package into metadataChan")
	}

	close(metadataChan)

	wg.Wait()
	log.Println("Done waiting")

	close(packageModelChan)

	return packageModels, nil
}