package crouter

import "fmt"

// ErrorNoPermissions is called when the user doesn't have permission to execute the command
type ErrorNoPermissions struct {
	MissingPerms string
}

func (e *ErrorNoPermissions) Error() string {
	return fmt.Sprintf("You are missing the following permission(s) to run this command: `%s`", e.MissingPerms)
}

// ErrorNotEnoughArgs is called when not enough arguments were supplied
type ErrorNotEnoughArgs struct {
	NumRequiredArgs int
	SuppliedArgs    int
}

func (e *ErrorNotEnoughArgs) Error() string {
	return fmt.Sprintf("This command requires `%v` arguments, but you supplied `%v` arguments.", e.NumRequiredArgs, e.SuppliedArgs)
}

// ErrorMissingRequiredArgs is called when one or more required arguments is missing
type ErrorMissingRequiredArgs struct {
	RequiredArgs string
	MissingArgs  string
}

func (e *ErrorMissingRequiredArgs) Error() string {
	return fmt.Sprintf("You are missing one or more required arguments.\nExpected `%v`, missing `%v`.", e.RequiredArgs, e.MissingArgs)
}

// ErrorNoDMs is called when the command cannot be used in a DM
type ErrorNoDMs struct{}

func (e *ErrorNoDMs) Error() string {
	return fmt.Sprintf("This command cannot be run in DMs.")
}

// ErrorTooManyArguments is called when there are too many arguments
type ErrorTooManyArguments struct {
	MaxArgs      int
	SuppliedArgs int
}

func (e *ErrorTooManyArguments) Error() string {
	return fmt.Sprintf("This command requires `%v` arguments, but you supplied `%v` arguments.", e.MaxArgs, e.SuppliedArgs)
}

// ErrorNotACommand is called when the given string isn't a command.
type ErrorNotACommand struct{}

func (e *ErrorNotACommand) Error() string {
	return fmt.Sprintf("the provided input wasn't a command")
}
