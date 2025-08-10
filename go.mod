module github.com/yantianyv/AkashaTerminal

go 1.24.5

replace (
	github.com/yantianyv/AkashaTerminal => ./
	github.com/yantianyv/AkashaTerminal/internal/commands => ./internal/commands
	github.com/yantianyv/AkashaTerminal/internal/config => ./internal/config
	github.com/yantianyv/AkashaTerminal/internal/operations => ./internal/operations
	github.com/yantianyv/AkashaTerminal/internal/providers => ./internal/providers
	github.com/yantianyv/AkashaTerminal/internal/state => ./internal/state
	github.com/yantianyv/AkashaTerminal/internal/utils => ./internal/utils
	github.com/yantianyv/AkashaTerminal/pkg/types => ./pkg/types
)

require (
	github.com/fatih/color v1.18.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/spf13/cobra v1.9.1 // indirect
	github.com/spf13/pflag v1.0.7 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
