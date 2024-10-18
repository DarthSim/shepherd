package main

import (
	"os"
	"path/filepath"

	"github.com/urfave/cli"

	_ "github.com/DarthSim/godotenv/autoload"
)

const version = "1.1.0"

func main() {
	var (
		conf hivemindConfig
		err  error
	)

	app := cli.NewApp()

	app.Name = "Hivemind"
	app.HelpName = "hivemind"
	app.Usage = "The mind to rule processes of your development environment"
	app.Description = "Hivemind is a process manager for Procfile-based applications"
	app.Author = "Sergey \"DarthSim\" Alexandrovich"
	app.Email = "darthsim@gmail.com"
	app.Version = version
	app.ArgsUsage = "[procfile] (Use '-' to read from stdin, Procfile path can be also set with $HIVEMIND_PROCFILE)"
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "title, w", EnvVar: "HIVEMIND_TITLE", Usage: "specify a title of the application", Destination: &conf.Title},
		cli.StringFlag{Name: "processes, l", EnvVar: "HIVEMIND_PROCESSES", Usage: "specify process names to launch. Divide names with comma", Destination: &conf.ProcNames},
		cli.IntFlag{Name: "port, p", EnvVar: "HIVEMIND_PORT,PORT", Usage: "specify a port to use as the base", Value: 5000, Destination: &conf.PortBase},
		cli.IntFlag{Name: "port-step, P", EnvVar: "HIVEMIND_PORT_STEP", Usage: "specify a step to increase port number", Value: 100, Destination: &conf.PortStep},
		cli.StringFlag{Name: "root, d", EnvVar: "HIVEMIND_ROOT", Usage: "specify a working directory of application. Default: directory containing the Procfile", Destination: &conf.Root},
		cli.IntFlag{Name: "timeout, t", EnvVar: "HIVEMIND_TIMEOUT", Usage: "specify the amount of time (in seconds) processes have to shut down gracefully before being brutally killed", Value: 5, Destination: &conf.Timeout},
		cli.BoolFlag{Name: "no-prefix", EnvVar: "HIVEMIND_NO_PREFIX", Usage: "process names will not be printed if the flag is specified", Destination: &conf.NoPrefix},
		cli.BoolFlag{Name: "print-timestamps, T", EnvVar: "HIVEMIND_PRINT_TIMESTAMPS", Usage: "timestamps will be printed if the flag is specified", Destination: &conf.PrintTimestamps},
		cli.BoolFlag{Name: "exit-with-highest-exit-code, e", EnvVar: "HIVEMIND_EXIT_WITH_HIGHEST_EXIT_CODE", Usage: "exit hivemind with highest exit code from all processes", Destination: &conf.ExitWithHighest},
	}

	app.Action = func(c *cli.Context) error {
		switch c.NArg() {
		case 0:
			if path := os.Getenv("HIVEMIND_PROCFILE"); len(path) > 0 {
				conf.Procfile = path
			} else {
				conf.Procfile = "./Procfile"
			}
		case 1:
			conf.Procfile = c.Args().First()
		default:
			fatal("Specify a single procfile")
		}

		if conf.Timeout < 1 {
			fatal("Timeout should be greater than 0")
		}

		if len(conf.Root) == 0 {
			conf.Root = filepath.Dir(conf.Procfile)
		}

		conf.Root, err = filepath.Abs(conf.Root)
		fatalOnErr(err)

		exitCode := newHivemind(conf).Run()
		if exitCode > 0 && conf.ExitWithHighest {
			return cli.NewExitError("At least one process failed", exitCode)
		}

		return nil
	}

	app.Run(os.Args)
}
