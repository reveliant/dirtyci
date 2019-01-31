package server

import "sort"

type Repositories []Repository

func (repositories Repositories) SetDefaults(defaults Repository) {
	for i, _ := range repositories {
		repositories[i].SetDefaults(defaults);
	}
}

func (repositories Repositories) Sort() {
	sort.Sort(repositories)
}

func (repositories Repositories) Len() int {
	return len(repositories)
}

func (repositories Repositories) Less(i, j int) bool {
	return repositories[i].RemoteUrl < repositories[j].RemoteUrl
}

func (repositories Repositories) Swap(i, j int) {
	var rep Repository
	rep = repositories[i]
	repositories[i] = repositories[j]
	repositories[j] = rep
}

func (repositories Repositories) Search(reponame string) *Repository {
	var idx int

	idx = sort.Search(repositories.Len(), func (i int) (bool) {
		return repositories[i].RemoteUrl >= reponame
	})

	if idx >= repositories.Len() {
		return nil
	}

	return &(repositories[idx])
}
