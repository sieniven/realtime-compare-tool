package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/sieniven/realtime-compare-tool/compare"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func main() {
	app := cli.NewApp()
	app.Name = "realtime-comparator"
	app.Action = run
	app.Flags = compare.DefaultFlags
	if err := app.Run(os.Args); err != nil {
		_, printErr := fmt.Fprintln(os.Stderr, err)
		if printErr != nil {
			fmt.Printf("Fprintln error: %v\n", printErr)
		}
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	logger := log.Default()
	configFilePath := ctx.String(compare.ConfigFlag.Name)
	if configFilePath != "" {
		if err := setFlagsFromConfigFile(ctx, configFilePath, logger); err != nil {
			logger.Printf("failed setting config flags from yaml/toml file, err: %v\n", err)
			return err
		}
	}

	compareCfg := compare.NewCompareConfig(ctx)
	service, err := compare.NewCompareService(compareCfg, logger)
	if err != nil {
		logger.Printf("failed creating compare service, err: %v\n", err)
		return err
	}

	service.Start(ctx.Context)
	return nil
}

func setFlagsFromConfigFile(ctx *cli.Context, filePath string, logger *log.Logger) error {
	fileConfig := make(map[string]interface{})
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, fileConfig)
	if err != nil {
		return err
	}

	for key, value := range fileConfig {
		if ctx.IsSet(key) {
			continue
		}

		if err := setFlag(ctx, key, value, logger); err != nil {
			return err
		}
	}
	return nil
}

func setFlag(ctx *cli.Context, key string, value interface{}, logger *log.Logger) error {
	isSlice := reflect.ValueOf(value).Kind() == reflect.Slice
	if isSlice {
		return setMultiValueFlag(ctx, key, value, logger)
	}
	return setSingleValueFlag(ctx, key, value, logger)
}

func setMultiValueFlag(ctx *cli.Context, key string, value interface{}, logger *log.Logger) error {
	sliceInterface := value.([]interface{})
	slice := make([]string, len(sliceInterface))
	for i, v := range sliceInterface {
		slice[i] = fmt.Sprintf("%v", v)
	}

	return setFlagInContext(ctx, key, strings.Join(slice, ","), logger)
}

func setSingleValueFlag(ctx *cli.Context, key string, value interface{}, logger *log.Logger) error {
	return setFlagInContext(ctx, key, fmt.Sprintf("%v", value), logger)
}

func setFlagInContext(ctx *cli.Context, key, value string, logger *log.Logger) error {
	if err := ctx.Set(key, value); err != nil {
		return handleFlagError(key, value, err, logger)
	}
	return nil
}

func handleFlagError(key, value string, err error, logger *log.Logger) error {
	errUnknownFlag := fmt.Errorf("no such flag -%s", key)
	if err.Error() == errUnknownFlag.Error() {
		logger.Printf("failed setting %s flag with value=%s, error=%s", key, value, err)
		return nil
	}

	return fmt.Errorf("failed setting %s flag with value=%s, error=%s", key, value, err)
}
