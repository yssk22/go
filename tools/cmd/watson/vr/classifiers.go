package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	"github.com/olekukonko/tablewriter"
	vr "github.com/speedland/go/services/watson/visualrecognition"
	"github.com/urfave/cli"
)

var prepareClassifier = cli.Command{
	Name:     "prepare-classifier",
	Usage:    "prepare face files from a directory to pass create-classifier",
	Category: "classifiers",
	Flags: []cli.Flag{
		cli.IntFlag{
			Name:  "numfiles, n",
			Usage: "the number of files sent for face detection (must be less than 15)",
			Value: maxFaceDetectionFiles,
		},
	},
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		args := c.Args()
		if len(args) != 2 {
			return fmt.Errorf("must specify a source directory and an output directory")
		}
		rootDir := args[0]
		outputDir := args[1]
		if err := os.MkdirAll(outputDir, os.FileMode(0755)); err != nil {
			return err
		}
		rootInfo, err := os.Stat(rootDir)
		if err != nil {
			return err
		}
		if !rootInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", rootDir)
		}
		sources, err := collectSources(rootDir)
		if err != nil {
			return err
		}
		size := c.Int("numfiles")
		if size < 0 || size > 15 {
			return fmt.Errorf("numfiles must be >0 or <15")
		}
		_, err = extractFacesFromSources(sources, outputDir, client, size)
		if err != nil {
			if resperr, ok := err.(*vr.ErrorResponse); ok {
				if resperr.Code == 413 {
					return fmt.Errorf("too many files found to extract. Use --numfiles to reduce the request size")
				}
			}
			return err
		}
		return nil
	},
}

var createClassifier = cli.Command{
	Name:     "create-classifier",
	Usage:    "create a new custom classifier from a directory",
	Category: "classifiers",
	Flags:    []cli.Flag{},
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		args := c.Args()
		if len(args) != 2 {
			return fmt.Errorf("must specify a classifier name and a source directory")
		}
		name := args[0]
		rootDir := args[1]
		rootInfo, err := os.Stat(rootDir)
		if err != nil {
			return err
		}
		if !rootInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", rootDir)
		}
		sources, err := collectSources(rootDir)
		if err != nil {
			return err
		}
		positives, negatives := buildExamples(sources)
		defer func() {
			for _, a := range positives {
				a.Close()
			}
			if negatives != nil {
				negatives.Close()
			}
		}()
		log.Printf("Creating classifier %s....", name)
		resp, err := client.CreateClassifier(context.Background(), name, positives, negatives)
		if err != nil {
			return err
		}
		fmt.Printf("%s created\n", resp.ClassifierID)
		return nil
	},
}

var listClassifiers = cli.Command{
	Name:     "list-classifiers",
	Usage:    "list all classifiers",
	Category: "classifiers",
	Flags:    []cli.Flag{},
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		resp, err := client.ListClassifiers(context.Background(), true)
		if err != nil {
			return err
		}

		table := tablewriter.NewWriter(os.Stdout)
		header := []string{"ID", "Name", "Status", "Created"}
		table.SetHeader(header)
		for _, classifier := range resp.Classifiers {
			row := []string{
				classifier.ClassifierID,
				classifier.Name,
				string(classifier.Status),
				formatTime(classifier.Created),
			}
			table.Append(row)
		}
		table.Render()
		return nil
	},
}

var showClassifier = cli.Command{
	Name:     "show-classifier",
	Usage:    "show details of a clasifier",
	Category: "classifiers",
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		for _, id := range c.Args() {
			fmt.Printf("%s:\n", id)
			resp, err := client.GetClassifier(context.Background(), id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Printf("\tname: %s\n", resp.Name)
			if resp.Status == vr.StatusFailed || resp.Status == vr.StatusUnavailable {
				fmt.Printf("\tstatus: %s (%s)\n", resp.Status, resp.Explanation)
			} else {
				fmt.Printf("\tstatus: %s\n", resp.Status)
			}
			fmt.Printf("\tclasses:\n")
			for _, class := range resp.Classes {
				fmt.Printf("\t\t%s\n", class.Class)
			}
		}
		return nil
	},
}

var updateClassifier = cli.Command{
	Name:     "update-classifier",
	Usage:    "update a classifier",
	Category: "classifiers",
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		args := c.Args()
		if len(args) != 2 {
			return fmt.Errorf("must specify a classifier id and a source directory")
		}
		id := args[0]
		rootDir := args[1]
		rootInfo, err := os.Stat(rootDir)
		if err != nil {
			return err
		}
		if !rootInfo.IsDir() {
			return fmt.Errorf("%s is not a directory", rootDir)
		}
		sources, err := collectSources(rootDir)
		if err != nil {
			return err
		}
		positives, negatives := buildExamples(sources)
		defer func() {
			for _, a := range positives {
				a.Close()
			}
			if negatives != nil {
				negatives.Close()
			}
		}()
		log.Printf("Updating classifier %s....", id)
		resp, err := client.UpdateClassifier(context.Background(), id, positives, negatives)
		if err != nil {
			return err
		}
		fmt.Printf("%s updated\n", resp.ClassifierID)
		return nil
	},
}

var deleteClassifier = cli.Command{
	Name:     "delete-classifier",
	Usage:    "delete a classifier",
	Category: "classifiers",
	Action: func(c *cli.Context) error {
		client, err := NewClient(c)
		if err != nil {
			return err
		}
		for _, id := range c.Args() {
			log.Printf("Deleting %s ...", id)
			_, err := client.DeleteClassifier(context.Background(), id)
			if err != nil {
				fmt.Printf("%s", err)
			} else {
				fmt.Printf("%s deleted\n", id)
			}
		}
		return nil
	},
}
