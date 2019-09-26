package libcmd

func (cmd *Cmd) getOptVal(name string) interface{} {
	if opt := cmd.findOpt("-" + name); opt != nil {
		return opt.val.raw
	}

	return cmd.findOpt("--" + name).val.raw
}

// GetString returns the string pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetString(name string) *string {
	return cmd.getOptVal(name).(*string)
}

// GetBool returns the bool pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetBool(name string) *bool {
	return cmd.getOptVal(name).(*bool)
}

// GetInt returns the int pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetInt(name string) *int {
	return cmd.getOptVal(name).(*int)
}

// GetInt8 returns the int8 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetInt8(name string) *int8 {
	return cmd.getOptVal(name).(*int8)
}

// GetInt16 returns the int16 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetInt16(name string) *int16 {
	return cmd.getOptVal(name).(*int16)
}

// GetInt32 returns the int32 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetInt32(name string) *int32 {
	return cmd.getOptVal(name).(*int32)
}

// GetInt64 returns the int64 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetInt64(name string) *int64 {
	return cmd.getOptVal(name).(*int64)
}

// GetUint returns the uint pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetUint(name string) *uint {
	return cmd.getOptVal(name).(*uint)
}

// GetUint8 returns the uint8 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetUint8(name string) *uint8 {
	return cmd.getOptVal(name).(*uint8)
}

// GetUint16 returns the uint16 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetUint16(name string) *uint16 {
	return cmd.getOptVal(name).(*uint16)
}

// GetUint32 returns the uint32 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetUint32(name string) *uint32 {
	return cmd.getOptVal(name).(*uint32)
}

// GetUint64 returns the uint64 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetUint64(name string) *uint64 {
	return cmd.getOptVal(name).(*uint64)
}

// GetFloat32 returns the float32 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetFloat32(name string) *float32 {
	return cmd.getOptVal(name).(*float32)
}

// GetFloat64 returns the float64 pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetFloat64(name string) *float64 {
	return cmd.getOptVal(name).(*float64)
}

// GetChoice returns the string pointer used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetChoice(name string) *string {
	return cmd.GetCustom(name).(*choiceString).value
}

// GetCustom returns the CustomArg value used as value
// for the argument 'name' (you can use either the short or long name).
// If the argument does not exist, this routine panics.
func (cmd *Cmd) GetCustom(name string) CustomArg {
	return cmd.getOptVal(name).(CustomArg)
}
