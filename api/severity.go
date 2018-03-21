package api

//go:generate go-enum -f=$GOFILE

// Severity of linter issue
// ENUM(
// Info,
// Warning,
// Error
// )
type Severity int8
