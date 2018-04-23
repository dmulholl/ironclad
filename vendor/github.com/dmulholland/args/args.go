/*
    Package args is a minimalist argument-parsing library designed for
    building elegant command-line interfaces.
*/
package args


import (
    "fmt"
    "os"
    "strings"
    "strconv"
    "unicode"
    "sort"
)


// Package version.
const Version = "1.0.0"


// Print a message to stderr and exit with a non-zero error code.
func exit(msg string) {
    fmt.Fprintf(os.Stderr, "Error: %v.\n", msg)
    os.Exit(1)
}


// -----------------------------------------------------------------------------
// Options
// -----------------------------------------------------------------------------


// Internal type for storing option data.
type option struct {
    optiontype string
    found bool
    bools []bool
    strings []string
    ints []int
    floats []float64
    boolfb bool
    stringfb string
    intfb int
    floatfb float64
}


// Attempts to set the value of an option by parsing a string argument.
func (opt *option) trySet(arg string) {
    switch opt.optiontype {

    case "string":
        opt.strings = append(opt.strings, arg)

    case "int":
        intvalue, err := strconv.ParseInt(arg, 0, 0)
        if err != nil {
            exit(fmt.Sprintf("cannot parse '%v' as an integer", arg))
        }
        opt.ints = append(opt.ints, int(intvalue))

    case "float":
        floatvalue, err := strconv.ParseFloat(arg, 64)
        if err != nil {
            exit(fmt.Sprintf("cannot parse '%v' as a float", arg))
        }
        opt.floats = append(opt.floats, floatvalue)
    }
}


// -----------------------------------------------------------------------------
// ArgStream
// -----------------------------------------------------------------------------


// Makes a slice of string arguments available as a stream.
type argstream struct {
    args []string
    index int
    length int
}


// Initialize a new argstream instance.
func newArgStream(args []string) *argstream {
    return &argstream{
        args: args,
        index: 0,
        length: len(args),
    }
}


// Returns the next argument from the stream.
func (stream *argstream) next() string {
    stream.index += 1
    return stream.args[stream.index - 1]
}


// Returns true if the stream contains at least one more element.
func (stream *argstream) hasNext() bool {
    return stream.index < stream.length
}


// -----------------------------------------------------------------------------
// ArgParser
// -----------------------------------------------------------------------------


// An ArgParser instance is responsible for storing registered options and
// commands. Note that every registered command recursively receives an
// ArgParser instance of its own.
type ArgParser struct {

    // Help text for the application or command.
    Helptext string

    // Application version.
    Version string

    // Stores option instances indexed by option name.
    options map[string]*option

    // Stores command sub-parser instances indexed by command name.
    commands map[string]*ArgParser

    // Stores positional arguments parsed from the input array.
    arguments []string

    // Registered callback function for a command parser.
    callback func(*ArgParser)

    // Stores the command name, if a command is found while parsing.
    command string

    // Stores a reference to a command parser's parent.
    parent *ArgParser
}


// NewParser initializes a new ArgParser instance.
func NewParser() *ArgParser {
    return &ArgParser {
        options: make(map[string]*option),
        commands: make(map[string]*ArgParser),
        arguments: make([]string, 0),
    }
}


// -----------------------------------------------------------------------------
// ArgParser: register options.
// -----------------------------------------------------------------------------


/// NewFlag registers a flag (a boolean option) with a default value of false.
// Flag options take no arguments but are either present (true) or absent
// (false).
func (parser *ArgParser) NewFlag(name string) {
    opt := &option{}
    opt.optiontype = "flag"
    opt.bools = make([]bool, 0)
    for _, alias := range strings.Split(name, " ") {
        parser.options[alias] = opt
    }
}


// NewString registers a string option, optionally specifying a fallback value.
func (parser *ArgParser) NewString(name string, fallback ...string) {
    opt := &option{}
    opt.optiontype = "string"
    opt.strings = make([]string, 0)
    if len(fallback) > 0 {
        opt.stringfb = fallback[0]
    }
    for _, alias := range strings.Split(name, " ") {
        parser.options[alias] = opt
    }
}


// NewInt registers an integer option, optionally specifying a fallback value.
func (parser *ArgParser) NewInt(name string, fallback ...int) {
    opt := &option{}
    opt.optiontype = "int"
    opt.ints = make([]int, 0)
    if len(fallback) > 0 {
        opt.intfb = fallback[0]
    }
    for _, alias := range strings.Split(name, " ") {
        parser.options[alias] = opt
    }
}


// NewFloat registers a floating-point option, optionally specifying a fallback
// value.
func (parser *ArgParser) NewFloat(name string, fallback ...float64) {
    opt := &option{}
    opt.optiontype = "float"
    opt.floats = make([]float64, 0)
    if len(fallback) > 0 {
        opt.floatfb = fallback[0]
    }
    for _, alias := range strings.Split(name, " ") {
        parser.options[alias] = opt
    }
}


// -----------------------------------------------------------------------------
// ArgParser: retrieve option values.
// -----------------------------------------------------------------------------


// Print an informative error message if a dev typos an option name.
func (parser *ArgParser) getOpt(name string) *option {
    if opt, found := parser.options[name]; found {
        return opt
    }
    panic(fmt.Sprintf("args: '%s' is not a registered option name", name))
}


// Found returns true if the specified option was found while parsing.
func (parser *ArgParser) Found(name string) bool {
    return parser.getOpt(name).found
}


// GetFlag returns the value of the specified boolean option.
func (parser *ArgParser) GetFlag(name string) bool {
    opt := parser.getOpt(name)
    if len(opt.bools) > 0 {
        return opt.bools[len(opt.bools) - 1]
    } else {
        return opt.boolfb
    }
}


// GetString returns the value of the specified string option.
func (parser *ArgParser) GetString(name string) string {
    opt := parser.getOpt(name)
    if len(opt.strings) > 0 {
        return opt.strings[len(opt.strings) - 1]
    } else {
        return opt.stringfb
    }
}


// GetInt returns the value of the specified integer option.
func (parser *ArgParser) GetInt(name string) int {
    opt := parser.getOpt(name)
    if len(opt.ints) > 0 {
        return opt.ints[len(opt.ints) - 1]
    } else {
        return opt.intfb
    }
}


// GetFloat returns the value of the specified floating-point option.
func (parser *ArgParser) GetFloat(name string) float64 {
    opt := parser.getOpt(name)
    if len(opt.floats) > 0 {
        return opt.floats[len(opt.floats) - 1]
    } else {
        return opt.floatfb
    }
}


// LenList returns the length of the specified option's list of values.
func (parser *ArgParser) LenList(name string) int {
    opt := parser.getOpt(name)
    switch opt.optiontype {
        case "flag":
            return len(opt.bools)
        case "string":
            return len(opt.strings)
        case "int":
            return len(opt.ints)
        case "float":
            return len(opt.floats)
    }
    return 0
}


// GetStringList returns the specified option's list of values.
func (parser *ArgParser) GetStringList(name string) []string {
    return parser.getOpt(name).strings
}


// GetIntList returns the specified option's list of values.
func (parser *ArgParser) GetIntList(name string) []int {
    return parser.getOpt(name).ints
}


// GetFloatList returns the specified option's list of values.
func (parser *ArgParser) GetFloatList(name string) []float64 {
    return parser.getOpt(name).floats
}


// -----------------------------------------------------------------------------
// ArgParser: positional arguments.
// -----------------------------------------------------------------------------


// HasArgs returns true if the parser has found one or more positional
// arguments.
func (parser *ArgParser) HasArgs() bool {
    return len(parser.arguments) > 0
}


// NumArgs returns the number of positional arguments.
func (parser *ArgParser) NumArgs() int {
    return len(parser.arguments)
}


// GetArg returns the positional argument at the specified index.
func (parser *ArgParser) GetArg(index int) string {
    return parser.arguments[index]
}


// GetArgs returns the positional arguments as a slice of strings.
func (parser *ArgParser) GetArgs() []string {
    return parser.arguments
}


// GetArgsAsInts attempts to parse and return the positional arguments as a
// slice of integers. The application will exit with an error message if any
// of the arguments cannot be parsed as an integer.
func (parser *ArgParser) GetArgsAsInts() []int {
    ints := make([]int, 0)
    for _, strArg := range parser.arguments {
        intArg, err := strconv.ParseInt(strArg, 0, 0)
        if err != nil {
            exit(fmt.Sprintf("cannot parse '%v' as an integer", strArg))
        }
        ints = append(ints, int(intArg))
    }
    return ints
}


// GetArgsAsFloats attempts to parse and return the positional arguments as a
// slice of floats. The application will exit with an error message if any
// of the arguments cannot be parsed as a float.
func (parser *ArgParser) GetArgsAsFloats() []float64 {
    floats := make([]float64, 0)
    for _, strArg := range parser.arguments {
        floatArg, err := strconv.ParseFloat(strArg, 64)
        if err != nil {
            exit(fmt.Sprintf("cannot parse '%v' as a float", strArg))
        }
        floats = append(floats, floatArg)
    }
    return floats
}


// -----------------------------------------------------------------------------
// ArgParser: commands.
// -----------------------------------------------------------------------------


// NewCmd registers a command, its help text, and its associated callback
// function. The callback function should accept the command's ArgParser
// instance as its sole agument and should have no return value.
func (parser *ArgParser) NewCmd(name, helptext string, callback func(*ArgParser)) *ArgParser {
    cmdParser := NewParser()
    cmdParser.Helptext = helptext
    cmdParser.callback = callback
    cmdParser.parent = parser
    for _, alias := range strings.Split(name, " ") {
        parser.commands[alias] = cmdParser
    }
    return cmdParser
}


// HasCmd returns true if the parser has found a command.
func (parser *ArgParser) HasCmd() bool {
    return parser.command != ""
}


// GetCmd returns the command name, if the parser has found a command.
func (parser *ArgParser) GetCmdName() string {
    return parser.command
}


// GetCmdParser returns the command's parser instance, if the parser has found
// a command.
func (parser *ArgParser) GetCmdParser() *ArgParser {
    return parser.commands[parser.command]
}


// GetParent returns a command parser's parent parser.
func (parser *ArgParser) GetParent() *ArgParser {
    return parser.parent
}


// -----------------------------------------------------------------------------
// ArgParser: parse arguments.
// -----------------------------------------------------------------------------


// Parse a stream of string arguments.
func (parser *ArgParser) parseStream(stream *argstream) {
    parsing := true
    isFirstArg := true

    // Loop while we have arguments to process.
    for stream.hasNext() {
        arg := stream.next()

        // If parsing has been turned off, simply add the argument to the
        // list of positionals.
        if !parsing {
            parser.arguments = append(parser.arguments, arg)
            continue
        }

        // If we encounter a -- argument, turn off option-parsing.
        if arg == "--" {
            parsing = false
            continue
        }

        // Is the argument a long-form option or flag?
        if strings.HasPrefix(arg, "--") {
            parser.parseLongOption(arg[2:], stream)
            continue
        }

        // Is the argument a short-form option or flag?
        if strings.HasPrefix(arg, "-") {
            if arg == "-" || unicode.IsDigit([]rune(arg)[1]) {
                parser.arguments = append(parser.arguments, arg)
            } else {
                parser.parseShortOption(arg[1:], stream)
            }
            continue
        }

        // Is the argument a registered command?
        if cmdParser, found := parser.commands[arg]; found {
            parser.command = arg
            cmdParser.parseStream(stream)
            cmdParser.callback(cmdParser)
            continue
        }

        // Is the argument the automatic 'help' command?
        if isFirstArg && arg == "help" {
            if stream.hasNext() {
                name := stream.next()
                if cmdParser, ok := parser.commands[name]; ok {
                    cmdParser.ExitHelp()
                } else {
                    exit(fmt.Sprintf("'%v' is not a recognised command", name))
                }
            } else {
                exit("the help command requires an argument")
            }
        }

        // If we get here, we have a positional argument.
        parser.arguments = append(parser.arguments, arg)
        isFirstArg = false
    }
}


// ParseArgs parses a slice of string arguments.
func (parser *ArgParser) ParseArgs(args []string) {
    parser.parseStream(newArgStream(args))
}


// Parse parses the application's command line arguments.
func (parser *ArgParser) Parse() {
    parser.ParseArgs(os.Args[1:])
}


// Parse a long-form option, i.e. an option beginning with a double dash.
func (parser *ArgParser) parseLongOption(arg string, stream *argstream) {

    // Do we have an option of the form --name=value?
    if strings.Contains(arg, "=") {
        parser.parseEqualsOption("--", arg)
        return
    }

    // Is the argument a registered option name?
    if opt, found := parser.options[arg]; found {
        opt.found = true
        if opt.optiontype == "flag" {
            opt.bools = append(opt.bools, true)
        } else if stream.hasNext() {
            opt.trySet(stream.next())
        } else {
            exit(fmt.Sprintf("missing argument for --%v", arg))
        }
        return
    }

    // Is the argument an automatic --help flag?
    if arg == "help" && parser.Helptext != "" {
        parser.ExitHelp()
    }

    // Is the argument an automatic --version flag?
    if arg == "version" && parser.Version != "" {
        parser.ExitVersion()
    }

    // The argument is not a recognised option name.
    exit(fmt.Sprintf("--%v is not a recognised option", arg))
}


// Parse a short-form option, i.e. an option beginning with a single dash.
func (parser *ArgParser) parseShortOption(arg string, stream *argstream) {

    // Do we have an option of the form -n=value?
    if strings.Contains(arg, "=") {
        parser.parseEqualsOption("-", arg)
        return
    }

    // We examine each character individually to support condensed options
    // with trailing arguments: -abc foo bar. If we don't recognise the
    // character as a registered option name, we check for an automatic
    // -h or -v flag before exiting.
    for _, char := range arg {
        name := string(char)
        if opt, found := parser.options[name]; found {
            opt.found = true
            if opt.optiontype == "flag" {
                opt.bools = append(opt.bools, true)
            } else if stream.hasNext() {
                opt.trySet(stream.next())
            } else {
                exit(fmt.Sprintf("missing argument for -%v", arg))
            }
        } else {
            if name == "h" && parser.Helptext != "" {
                parser.ExitHelp()
            } else if name == "v" && parser.Version != "" {
                parser.ExitVersion()
            } else {
                exit(fmt.Sprintf("-%v is not a recognised option", name))
            }
        }
    }
}


// Parse an option of the form --name=value or -n=value.
func (parser *ArgParser) parseEqualsOption(prefix string, arg string) {
    split := strings.SplitN(arg, "=", 2)
    name := split[0]
    value := split[1]

    // Do we have the name of a registered option?
    opt, found := parser.options[name]
    if !found {
        exit(fmt.Sprintf("%s%s is not a recognised option", prefix, name))
    }
    opt.found = true

    // Boolean flags should never contain an equals sign.
    if opt.optiontype == "flag" {
        exit(fmt.Sprintf("invalid format for boolean flag %s%s", prefix, name))
    }

    // Check that a value has been supplied.
    if value == "" {
        exit(fmt.Sprintf("missing argument for %s%s", prefix, name))
    }

    // Try to parse the argument as a value of the appropriate type.
    opt.trySet(value)
}


// -------------------------------------------------------------------------
// ArgParser: utilities.
// -------------------------------------------------------------------------


// ExitHelp prints the parser's help text, then exits.
func (parser *ArgParser) ExitHelp() {
    fmt.Println(strings.TrimSpace(parser.Helptext))
    os.Exit(0)
}


// ExitVersion prints the parser's version string, then exits.
func (parser *ArgParser) ExitVersion() {
    fmt.Println(strings.TrimSpace(parser.Version))
    os.Exit(0)
}


// String returns a string representation of the parser instance.
func (parser *ArgParser) String() string {
    lines := make([]string, 0)

    lines = append(lines, "Options:")
    if len(parser.options) > 0 {
        names := make([]string, 0, len(parser.options))
        for name := range parser.options {
            names = append(names, name)
        }
        sort.Strings(names)
        for _, name := range names {
            var values string
            opt := parser.options[name]
            switch opt.optiontype {
                case "flag":
                    values = fmt.Sprintf("(%v) %v", opt.boolfb, opt.bools)
                case "string":
                    values = fmt.Sprintf("(%v) %v", opt.stringfb, opt.strings)
                case "int":
                    values = fmt.Sprintf("(%v) %v", opt.intfb, opt.ints)
                case "float":
                    values = fmt.Sprintf("(%v) %v", opt.floatfb, opt.floats)
            }
            lines = append(lines, fmt.Sprintf("  %v: %v", name, values))
        }
    } else {
        lines = append(lines, "  [none]")
    }

    lines = append(lines, "\nArguments:")
    if len(parser.arguments) > 0 {
        for _, arg := range parser.arguments {
            lines = append(lines, fmt.Sprintf("  %v", arg))
        }
    } else {
        lines = append(lines, "  [none]")
    }

    lines = append(lines, "\nCommand:")
    if parser.HasCmd() {
        lines = append(lines, fmt.Sprintf("  %v", parser.GetCmdName()))
    } else {
        lines = append(lines, "  [none]")
    }

    return strings.Join(lines, "\n")
}
