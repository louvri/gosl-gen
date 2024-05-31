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

var templates = []string{"helper", "key", "model", "query", "repository", "service_reader", "service_writer", "modify_body", "modify_request", "search_body", "search_request"}

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
	var servicePath string
	if tmp, ok := r.config["$SERVICE_PATH"].(string); ok {
		servicePath = tmp
	}

	var requestPath string

	if tmp, ok := r.config["$REQUEST_PATH"].(string); ok {
		if servicePath != "" && requestPath == "" {
			fmt.Printf("Error: invalid request path\n")
			return errors.New("request is mandatory if service is generated")
		}
		requestPath = tmp
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
			data = strings.ReplaceAll(data, "$PROJECT_PATH", r.config["$PROJECT_PATH"].(string))
			return writer.Write(data)
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		err = tomlTemplate.Read(func(data string) error {
			if template == "repository" {
				data = strings.ReplaceAll(data, "$DB_CONNECTION_STRING", r.config["$DB_CONNECTION_STRING"].(string))
				data = strings.ReplaceAll(data, "$DB_SCHEMA", fmt.Sprintf("%v", r.config["$DB_SCHEMA"]))
				data = strings.ReplaceAll(data, "$DB_INCLUDE_TABLES", fmt.Sprintf("%v", r.config["$DB_INCLUDE_TABLES"]))
				data = strings.ReplaceAll(data, "$DB_EXCLUDE_TABLES", fmt.Sprintf("%v", r.config["$DB_EXCLUDE_TABLES"]))
				data = strings.ReplaceAll(data, "$DB_TYPE", r.config["$DB_TYPE"].(string))
				data = strings.ReplaceAll(data, "$WORKDIR_PATH", workdirPath)

				data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/%s.go\" =\"%s.gotmpl\"", repoPath, template, template))
			} else {
				data = strings.ReplaceAll(data, "$DB_CONNECTION_STRING", r.config["$DB_CONNECTION_STRING"].(string))
				data = strings.ReplaceAll(data, "$DB_SCHEMA", fmt.Sprintf("%v", r.config["$DB_SCHEMA"]))
				data = strings.ReplaceAll(data, "$DB_INCLUDE_TABLES", fmt.Sprintf("%v", r.config["$DB_INCLUDE_TABLES"]))
				data = strings.ReplaceAll(data, "$DB_EXCLUDE_TABLES", fmt.Sprintf("%v", r.config["$DB_EXCLUDE_TABLES"]))
				data = strings.ReplaceAll(data, "$DB_TYPE", r.config["$DB_TYPE"].(string))
				data = strings.ReplaceAll(data, "$WORKDIR_PATH", workdirPath)
				data = strings.ReplaceAll(data, "$REPOSITORY_PATH", repoPath)
				data = strings.ReplaceAll(data, "$MODEL_PATH", modelPath)
				data = strings.ReplaceAll(data, "$TEMPLATE", template)
				switch template {
				case "modify_body":
					{
						if requestPath != "" {
							data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/{{.Table}}/modify/body.go\" =\"%s.gotmpl\"", requestPath, template))
						}
					}
				case "modify_request":
					{
						if requestPath != "" {
							data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/{{.Table}}/modify/request.go\" =\"%s.gotmpl\"", requestPath, template))
						}
					}
				case "search_body":
					{
						if requestPath != "" {
							data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/{{.Table}}/search/body.go\" =\"%s.gotmpl\"", requestPath, template))
						}
					}
				case "search_request":
					{
						if requestPath != "" {
							data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/{{.Table}}/search/request.go\" =\"%s.gotmpl\"", requestPath, template))
						}
					}
				case "service_writer":
					{
						if servicePath != "" {
							data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/{{.Table}}/writer/service.go\" =\"%s.gotmpl\"", servicePath, template))
						}
					}
				case "service_reader":
					{
						if servicePath != "" {
							data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/{{.Table}}/reader/service.go\" =\"%s.gotmpl\"", servicePath, template))
						}
					}
				case "key":
					{
						data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/key/%s.go\" =\"%s.gotmpl\"", modelPath, template, template))
					}
				case "model":
					{
						data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/{{.Table}}/%s.go\" =\"%s.gotmpl\"", modelPath, template, template))
					}
				default:
					{
						data = strings.ReplaceAll(data, "$GENERATE_PATH", fmt.Sprintf("\"%s/{{.Table}}/%s.go\" =\"%s.gotmpl\"", repoPath, template, template))
					}
				}
			}
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
	if r.config["$WORKDIR_PATH"] != nil {
		command := fmt.Sprintf(`cd %s && go mod tidy`, r.config["$WORKDIR_PATH"].(string))
		if err := run("bash", "-c", command); err != nil {
			return err
		}
		fmt.Printf("go mod tidy\n")
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
	if r.config["$PROJECT_PATH"] == nil {
		return errors.New("$PROJECT_PATH is mandatory")
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
		return err
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	bytes, err := io.ReadAll(stdOut)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Error: %v\n", err)
		if exitError, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exit code is %d\n", exitError.ExitCode())
		}
		return err
	}
	fmt.Println(string(bytes))
	return nil
}
