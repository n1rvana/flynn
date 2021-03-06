package client

import (
	"fmt"
	"io"
	"net/url"

	"github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/opts"
	flag "github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/pkg/mflag"
	"github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/pkg/parsers"
	"github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/registry"
	"github.com/flynn/flynn/Godeps/_workspace/src/github.com/docker/docker/utils"
)

func (cli *DockerCli) CmdImport(args ...string) error {
	cmd := cli.Subcmd("import", "URL|- [REPOSITORY[:TAG]]", "Create an empty filesystem image and import the contents of the\ntarball (.tar, .tar.gz, .tgz, .bzip, .tar.xz, .txz) into it, then\noptionally tag it.", true)
	flChanges := opts.NewListOpts(nil)
	cmd.Var(&flChanges, []string{"c", "-change"}, "Apply Dockerfile instruction to the created image")
	cmd.Require(flag.Min, 1)

	utils.ParseFlags(cmd, args, true)

	var (
		v          = url.Values{}
		src        = cmd.Arg(0)
		repository = cmd.Arg(1)
	)

	v.Set("fromSrc", src)
	v.Set("repo", repository)
	for _, change := range flChanges.GetAll() {
		v.Add("changes", change)
	}
	if cmd.NArg() == 3 {
		fmt.Fprintf(cli.err, "[DEPRECATED] The format 'URL|- [REPOSITORY [TAG]]' has been deprecated. Please use URL|- [REPOSITORY[:TAG]]\n")
		v.Set("tag", cmd.Arg(2))
	}

	if repository != "" {
		//Check if the given image name can be resolved
		repo, _ := parsers.ParseRepositoryTag(repository)
		if err := registry.ValidateRepositoryName(repo); err != nil {
			return err
		}
	}

	var in io.Reader

	if src == "-" {
		in = cli.in
	}

	return cli.stream("POST", "/images/create?"+v.Encode(), in, cli.out, nil)
}
