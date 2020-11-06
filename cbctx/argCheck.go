package cbctx

// CheckMinArgs checks if the argument count is less than the given count
func (ctx *Ctx) CheckMinArgs(c int) (err error) {
	if len(ctx.Args) < c {
		return &ErrorNotEnoughArgs{
			NumRequiredArgs: c,
			SuppliedArgs:    len(ctx.Args),
		}
	}
	return nil
}

// CheckRequiredArgs checks if the arg count is exactly the given count
func (ctx *Ctx) CheckRequiredArgs(c int) (err error) {
	if len(ctx.Args) != c {
		if len(ctx.Args) > c {
			return &ErrorTooManyArguments{
				MaxArgs:      c,
				SuppliedArgs: len(ctx.Args),
			}
		} else {
			return &ErrorNotEnoughArgs{
				NumRequiredArgs: c,
				SuppliedArgs:    len(ctx.Args),
			}
		}
	}
	return nil
}

// CheckArgRange checks if the number of arguments is within the given range
func (ctx *Ctx) CheckArgRange(min, max int) (err error) {
	if len(ctx.Args) > max {
		return &ErrorTooManyArguments{
			MaxArgs:      max,
			SuppliedArgs: len(ctx.Args),
		}
	}
	if len(ctx.Args) < min {
		return &ErrorNotEnoughArgs{
			NumRequiredArgs: min,
			SuppliedArgs:    len(ctx.Args),
		}
	}
	return nil
}
