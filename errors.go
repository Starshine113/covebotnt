package main

import "fmt"

type errorNoPermissions struct {
	missingPerms string
}

func (e *errorNoPermissions) Error() string {
	return fmt.Sprintf("You are missing the following permission(s) to run this command: `%s`", e.missingPerms)
}

type errorNotEnoughArgs struct {
	numRequiredArgs int
	suppliedArgs    int
}

func (e *errorNotEnoughArgs) Error() string {
	return fmt.Sprintf("This command requires `%v` arguments, but you supplied `%v` arguments.", e.numRequiredArgs, e.suppliedArgs)
}

type errorMissingRequiredArgs struct {
	requiredArgs string
	missingArgs  string
}

func (e *errorMissingRequiredArgs) Error() string {
	return fmt.Sprintf("You are missing one or more required arguments.\nExpected `%v`, missing `%v`.", e.requiredArgs, e.missingArgs)
}

type errorNoDMs struct{}

func (e *errorNoDMs) Error() string {
	return fmt.Sprintf("This command cannot be run in DMs.")
}

type errorTooManyArguments struct {
	maxArgs      int
	suppliedArgs int
}

func (e *errorTooManyArguments) Error() string {
	return fmt.Sprintf("This command requires `%v` arguments, but you supplied `%v` arguments.", e.maxArgs, e.suppliedArgs)
}
