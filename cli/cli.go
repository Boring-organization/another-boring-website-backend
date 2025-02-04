package cli

import (
	"flag"
	"reflect"
)

type Flag string

const (
	GenerateSchema Flag = "GenerateSchema"
)

const (
	BoolType   int = iota
	StringType     = iota
	IntType        = iota
)

type FlagKeys = []string

type FlagParams struct {
	keys         FlagKeys
	hint         string
	defaultValue interface{}
	dataType     int
}

var FlagParamValues = map[Flag]FlagParams{
	GenerateSchema: {keys: []string{"generate-schema", "gs"}, hint: "Flag to generate schema or not", defaultValue: false, dataType: BoolType},
}

var FlagValues = map[Flag]any{}

var isParsed = false

func parseArgs() {
	for flagKeys, flagParams := range FlagParamValues {

		switch flagParams.dataType {
		case BoolType:
			for _, flagKey := range flagParams.keys {
				interfaceValue, ok := interface{}(flag.Bool(flagKey, flagParams.defaultValue.(bool), flagParams.hint)).(any)

				if ok {
					FlagValues[flagKeys] = interfaceValue
					break
				}
			}

		case StringType:
			for _, flagKey := range flagParams.keys {
				FlagValues[flagKeys] = interface{}(flag.String(flagKey, flagParams.defaultValue.(string), flagParams.hint)).(any)
				if FlagValues[flagKeys] != nil {
					break
				}
			}

		case IntType:
			for _, flagKey := range flagParams.keys {
				FlagValues[flagKeys] = interface{}(flag.Int(flagKey, flagParams.defaultValue.(int), flagParams.hint)).(*any)
				if FlagValues[flagKeys] != nil {
					break
				}
			}
		}
	}

	flag.Parse()

	isParsed = true
}

func GetArgs(flagName Flag) any {
	if !isParsed {
		parseArgs()
	}

	if len(FlagValues) == 0 {
		return nil
	}

	valuePointer := FlagValues[flagName]
	if valuePointer == nil {
		return nil
	}

	boolPointer, ok := valuePointer.(interface{}).(*bool)
	if ok && reflect.TypeOf(*boolPointer) == reflect.TypeOf(true) {
		return *boolPointer
	}

	intPointer, ok := valuePointer.(interface{}).(*int)
	if ok && reflect.TypeOf(*intPointer) == reflect.TypeOf(0) {
		return *intPointer
	}

	stringPointer, ok := valuePointer.(interface{}).(*string)
	if ok && reflect.TypeOf(*stringPointer) == reflect.TypeOf("") {
		return *stringPointer
	}

	return nil
}
