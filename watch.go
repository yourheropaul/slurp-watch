package watch

import (
	"os"
	"time"

	"github.com/omeid/slurp"
	"github.com/omeid/slurp/tools/glob"
)

type fileNameChan chan string
type errorChan chan error

func watchFile(filePath string) error {

	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

func monitorFile(p string, c fileNameChan, e errorChan) {

	go func() {
		for {
			// Wait for the file to change. If there's an error,
			// put it into the error channel
			err := watchFile(p)
			if err != nil {
				e <- err
				break
			}

			// On change, return the file path and loop
			c <- p
		}
	}()
}

func Watch(c *slurp.C, task func(string), globs ...string) {

	files, err := glob.Glob(globs...)

	if err != nil {
		c.Error(err)
		return
	}

	// Create some channels
	f := make(fileNameChan)
	e := make(errorChan)

	for matchpair := range files {
		monitorFile(matchpair.Name, f, e)
	}

	go func() {
		for {
			select {
			case fn := <-f:
				task(fn)

			case err := <-e:
				c.Error(err)
			}
		}
	}()
}
