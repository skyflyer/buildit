package build

import (
	"conf"
	"log"
	"os"

	"path"

	"io"

	"util"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func debug(args ...interface{}) {
	if false {
		args = append([]interface{}{"[DEBUG]"}, args...)
		log.Println(args...)
	}
}

// CheckoutFiles checks out files from the git repo given the reference
func CheckoutFiles(r *git.Repository, ref *plumbing.Reference) error {
	var res error
	c, err := r.Commit(ref.Hash())
	util.Check(err)

	log.Printf("Got commit: %s by %s on %s\r\n", c.Hash, c.Author.Name, c.Author.When)
	files, err := c.Files()
	util.Check(err)
	files.ForEach(func(f *object.File) error {
		debug("File:", f.Name, f.Mode)
		if f.Mode&os.ModeType != 0 {
			// ignore special files like symlinks for now
			return nil
		}

		dir := path.Join(conf.BuildDir, path.Dir(f.Name))
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			res = err
			return err
		}

		err = createFile(f)
		if err != nil {
			res = err
			return err
		}

		return nil
	})
	return res
}

func createFile(f *object.File) error {
	filepath := path.Join(conf.BuildDir, f.Name)
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, f.Mode)
	if err != nil {
		return err
	}
	debug("Writing file", filepath, "with size", f.Size, "and type", f.Type().String())
	rdr, err := f.Blob.Reader()
	if err != nil {
		return err
	}
	defer rdr.Close()

	n, err := io.Copy(file, rdr)
	if err != nil {
		return err
	}
	if n != f.Size {
		debug("Copied only", n, "bytes")
	}
	return nil
}

func listCommits(r *git.Repository) {
	ci, err := r.Commits()
	util.Check(err)

	ci.ForEach(func(c *object.Commit) error {
		log.Print("Commit: ", c.Hash)
		log.Print("-- ", c.Committer.When)
		log.Print("-- ", c.Message)
		return nil
	})
}

func listRefs(r *git.Repository) {
	refs, err := r.References()
	util.Check(err)

	refs.ForEach(func(ref *plumbing.Reference) error {
		log.Print("Ref:", ref, ref.IsBranch(), ref.IsRemote())
		return nil
	})
}
