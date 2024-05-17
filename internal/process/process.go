package process

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/louvri/gosl-gen/internal/file/config"
	"github.com/louvri/gosl-gen/internal/file/json"
)

const buildPath = "/tmp/.gen_gosl_build"

var templates = []string{"helper", "key", "model", "query", "repository"}

var instance Runner

type Runner interface {
	Initialize(cfg string) error
	Generate(cfg string) error
	IsInitiated() error
}

func New() Runner {
	if instance == nil {
		instance = &runner{}
	}
	return instance
}

type runner struct {
	config map[string]interface{}
}

func (r *runner) Initialize(path string) error {
	fmt.Println("read config")
	err := r.getConfig(path)
	if err != nil {
		return err
	}
	err = run("go", "get", "gnorm.org/gnorm")
	if err != nil {
		return err
	}
	fmt.Println("gnorm.org is installed")

	//build template
	workdirPath, ok := r.config["$WORKDIR_PATH"].(string)
	if !ok {
		fmt.Printf("Error: invalid workdir path\n")
		return errors.New("invalid workdir path")
	}
	repoPath, ok := r.config["$REPOSITORY_PATH"].(string)
	if !ok {
		fmt.Printf("Error: invalid repository path\n")
		return errors.New("invalid repository path")
	}
	modelPath, ok := r.config["$MODEL_PATH"].(string)
	if !ok {
		fmt.Printf("Error: invalid model path\n")
		return errors.New("invalid model path")
	}
	var tmp strings.Builder
	max := len(modelPath)
	for i, r := range repoPath {
		if i < max && rune(modelPath[i]) == r {
			tmp.WriteRune(r)
		} else {
			break
		}
	}
	rootPath := tmp.String()
	if len(rootPath) > 0 {
		rootPath = rootPath[:len(rootPath)-1]
	}

	for _, template := range templates {
		var reader config.File
		var writer config.File
		tomlTemplate := config.New("template/config.toml", config.Read)
		tomlBuild := config.New(fmt.Sprintf("%s/%s/config.toml", buildPath, template), config.Write)
		reader = config.New("template/"+template+".gotmpl", config.Read)
		writer = config.New(fmt.Sprintf("%s/%s/%s.gotmpl", buildPath, template, template), config.Write)
		defer func() {
			writer.Close()
			reader.Close()
			tomlTemplate.Close()
			tomlBuild.Close()
		}()
		err := reader.Read(func(data string) error {
			data = strings.ReplaceAll(data, "$REPOSITORY_PATH", repoPath)
			data = strings.ReplaceAll(data, "$MODEL_PATH", modelPath)
			return writer.Write(data)
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		err = tomlTemplate.Read(func(data string) error {
			data = strings.ReplaceAll(data, "$DB_CONNECTION_STRING", r.config["$DB_CONNECTION_STRING"].(string))
			data = strings.ReplaceAll(data, "$DB_SCHEMA", fmt.Sprintf("%v", r.config["$DB_SCHEMA"]))
			data = strings.ReplaceAll(data, "$DB_INCLUDE_TABLES", fmt.Sprintf("%v", r.config["$DB_INCLUDE_TABLES"]))
			data = strings.ReplaceAll(data, "$DB_EXCLUDE_TABLES", fmt.Sprintf("%v", r.config["$DB_EXCLUDE_TABLES"]))
			data = strings.ReplaceAll(data, "$DB_TYPE", r.config["$DB_TYPE"].(string))
			data = strings.ReplaceAll(data, "$WORKDIR_PATH", workdirPath)
			data = strings.ReplaceAll(data, "$REPOSITORY_PATH", repoPath)
			data = strings.ReplaceAll(data, "$MODEL_PATH", modelPath)
			data = strings.ReplaceAll(data, "$TEMPLATE", template)
			data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("%v/%v", rootPath, template))
			return tomlBuild.Write(data)
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		fmt.Printf("- %s/%s copied\n", buildPath, template)
	}
	fmt.Println("initiated")
	return nil
}

func (r *runner) Generate(path string) error {
	fmt.Println("read config")
	err := r.getConfig(path)
	if err != nil {
		return err
	}
	for _, template := range templates {
		command := fmt.Sprintf(`cd %s/%s && gnorm gen -c ./config.toml`, buildPath, template)
		if err := run("bash", "-c", command); err != nil {
			return err
		}
		fmt.Printf("- generating %s/%s \n", buildPath, template)
	}
	fmt.Println("generated")
	return nil
}

func (r *runner) IsInitiated() error {
	command := fmt.Sprintf(`ls -al %s | wc -l`, buildPath)
	cmd := exec.Command("bash", "-c", command)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("%v \n", err)
		return err
	}
	re := regexp.MustCompile("[0-9]+")
	numbers := re.FindAllString(string(stdout), -1)
	for _, number := range numbers {
		n, _ := strconv.ParseInt(number, 10, 64)
		if n == 0 {
			return errors.New("gen-gosl not initiated")
		}
	}
	return nil
}

func (r *runner) getConfig(path string) error {
	reader := json.New(path)
	if tmp, err := reader.Object(); err != nil {
		return err
	} else {
		r.config = tmp
	}
	if r.config["$DB_CONNECTION_STRING"] == nil {
		return errors.New("$DB_CONNECTION_STRING is mandatory")
	}
	if r.config["$DB_SCHEMA"] == nil {
		return errors.New("$DB_SCHEMA is mandatory")
	} else if arr, ok := r.config["$DB_SCHEMA"].([]interface{}); !ok {
		return errors.New("$DB_SCHEMA value should be array")
	} else {
		var output strings.Builder
		output.WriteRune('[')
		for i, item := range arr {
			if i > 0 {
				output.WriteRune(',')
			}
			output.WriteRune('"')
			output.WriteString(item.(string))
			output.WriteRune('"')
		}
		output.WriteRune(']')
		r.config["$DB_SCHEMA"] = output.String()
	}
	if arr, ok := r.config["$DB_INCLUDE_TABLES"].([]interface{}); !ok {
		r.config["$DB_INCLUDE_TABLES"] = "[]"
	} else {
		var output strings.Builder
		output.WriteRune('[')
		for i, item := range arr {
			if i > 0 {
				output.WriteRune(',')
			}
			output.WriteRune('"')
			output.WriteString(item.(string))
			output.WriteRune('"')
		}
		output.WriteRune(']')
		r.config["$DB_INCLUDE_TABLES"] = output.String()
	}
	if arr, ok := r.config["$DB_EXCLUDE_TABLES"].([]interface{}); !ok {
		r.config["$DB_EXCLUDE_TABLES"] = "[]"
	} else {
		var output strings.Builder
		output.WriteRune('[')
		for i, item := range arr {
			if i > 0 {
				output.WriteRune(',')
			}
			output.WriteRune('"')
			output.WriteString(item.(string))
			output.WriteRune('"')
		}
		output.WriteRune(']')
		r.config["$DB_EXCLUDE_TABLES"] = output.String()
	}
	if r.config["$WORKDIR_PATH"] == nil {
		return errors.New("$WORKDIR_PATH is mandatory")
	}
	if r.config["$REPOSITORY_PATH"] == nil {
		return errors.New("$REPOSITORY_PATH is mandatory")
	}
	if r.config["$MODEL_PATH"] == nil {
		return errors.New("$MODEL_PATH is mandatory")
	}
	if r.config["$DB_TYPE"] == nil {
		return errors.New("$DB_TYPE is mandatory")
	}
	return nil
}

func run(cli string, args ...string) error {

	cmd := exec.Command(cli, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	bytes, err := io.ReadAll(stdOut)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code is %d\n", exitError.ExitCode())
		}
	}
	fmt.Println(string(bytes))
	return nil
}
