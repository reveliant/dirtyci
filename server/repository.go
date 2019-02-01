package server

import (
	"log"
	"gopkg.in/libgit2/git2go.v27"
)

func certificateCheckCallback(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
	if (cert.Kind == git.CertificateX509) && !valid {
		return git.ErrCertificate
	}
	return git.ErrOk
}

func transferProgressCallback(stats git.TransferProgress) git.ErrorCode {
	if (stats.ReceivedObjects == stats.TotalObjects) {
		log.Printf("Resolving %d deltas\n", stats.TotalDeltas)
	} else if (stats.TotalObjects > 0) {
		log.Printf(
			"Received %d/%d objects (%d) in %zu bytes\n",
			stats.ReceivedObjects, stats.TotalObjects,
			stats.IndexedObjects, stats.ReceivedBytes,
		)
	}

	log.Printf("%s / %s\n", stats.ReceivedObjects, stats.TotalObjects)
	return git.ErrOk
}

type Repository struct {
	RemoteUrl       string  `yaml:"remoteUrl"`
	LocalUrl        string  `yaml:"localUrl"`
	RemoteName      string  `yaml:"remoteName"`
	RemoteBranch    string  `yaml:"remoteBranch"`
	LocalBranch     string  `yaml:"localBranch"`
	PublicKeyPath   string  `yaml:"publicKeyPath"`
	PrivateKeyPath  string  `yaml:"privateKeyPath"`
}

func (repo *Repository) GetCredentialsCallback() git.CredentialsCallback {
	return func (url string, username string, allowedTypes git.CredType) (git.ErrorCode, *git.Cred) {
		var ret int
		var cred git.Cred

		if (allowedTypes & git.CredTypeSshKey) != 0 {
			ret, cred = git.NewCredSshKey("git", repo.PublicKeyPath, repo.PrivateKeyPath, "")
		}

		return git.ErrorCode(ret), &cred
	}
}

func (repo *Repository) SetDefaults(defaults Repository) {
	repo.SetDefault(&repo.RemoteName, defaults.RemoteName)
	repo.SetDefault(&repo.RemoteBranch, defaults.RemoteBranch)
	repo.SetDefault(&repo.LocalBranch, defaults.LocalBranch)
	repo.SetDefault(&repo.PublicKeyPath, defaults.PublicKeyPath)
	repo.SetDefault(&repo.PrivateKeyPath, defaults.PrivateKeyPath)
}

func (repo *Repository) SetDefault(key *string, defaultValue string) {
	if *key == "" {
		*key = defaultValue
	}
}

func (repo *Repository) Pull() error {
	var remote *git.Remote
	var remoteRef *git.Reference
	var remoteRefSpecs = make([]string, 1);
	remoteRefSpecs[0] = "refs/remotes/" + repo.RemoteName + "/" + repo.RemoteBranch
	var fetchOptions = git.FetchOptions{
		RemoteCallbacks: git.RemoteCallbacks{
			CredentialsCallback: repo.GetCredentialsCallback(),
			CertificateCheckCallback: certificateCheckCallback,
			TransferProgressCallback: transferProgressCallback,
		},
	}

	// Opening repository
	var gitrepo, err = git.OpenRepository(repo.LocalUrl)
	if err != nil {
		return err
	}
	log.Printf("[%s] Opened repository\n", repo.LocalUrl)

	// Getting remote
	remote, err = gitrepo.Remotes.Lookup(repo.RemoteName)
	if err != nil {
		log.Printf("[%s] Looked up remote '%s' failed!\n", repo.LocalUrl, repo.RemoteName)
		return err
	}
	log.Printf("[%s] Looked up remote '%s'\n", repo.LocalUrl, repo.RemoteName)

	// Fetching from remote
	err = remote.Fetch(remoteRefSpecs, &fetchOptions, "")
	if err != nil {
		log.Printf("[%s] Fetching %s/%s failed!\n", repo.LocalUrl, repo.RemoteName, repo.RemoteBranch)
		return err
	}
	log.Printf("[%s] Fetched %s/%s\n", repo.LocalUrl, repo.RemoteName, repo.RemoteBranch)
	remote.Free()

	// Get remote head
	remoteRef, err = gitrepo.References.Lookup(remoteRefSpecs[0])
	if err != nil {
		return err
	}
	var mergeRemoteHeads = make([]*git.AnnotatedCommit, 1)
	mergeRemoteHeads[0], err = gitrepo.AnnotatedCommitFromRef(remoteRef)
	if err != nil {
		return err
	}
	remoteRef.Free()

	// Merging remote into local head
	err = gitrepo.Merge(mergeRemoteHeads, nil, nil)
	if err != nil {
		log.Printf("[%s] Merging %s/%s into %s failed!\n", repo.LocalUrl, repo.RemoteName, repo.RemoteBranch, repo.LocalBranch)
		return err
	}
	mergeRemoteHeads[0].Free()
	log.Printf("[%s] Merged %s/%s into %s\n", repo.LocalUrl, repo.RemoteName, repo.RemoteBranch, repo.LocalBranch)

	gitrepo.Free()
	return nil
}
