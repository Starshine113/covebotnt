package main

import "fmt"

type errorNoPermissions struct {
	missingPerms string
}

type errorNotEnoughArgs struct {
	numRequiredArgs int
	suppliedArgs    int
}

func (e *errorNoPermissions) Error() string {
	return fmt.Sprintf("You are missing the following permission(s) to run this command: `%s`", e.missingPerms)
}

func (e *errorNotEnoughArgs) Error() string {
	return fmt.Sprintf("This command requires `%v` arguments, but you supplied `%v` arguments.", e.numRequiredArgs, e.suppliedArgs)
}
