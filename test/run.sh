echo "generating model and repository\n"

# install gosl-gen 
go install github.com/louvri/gosl-gen@latest


# initiate config 
~/go/bin/gosl-gen init -c gosl_config/config.json

# generate model and repository
~/go/bin/gosl-gen gen -c gosl_config/config.json

# tidy module
go mod tidy
