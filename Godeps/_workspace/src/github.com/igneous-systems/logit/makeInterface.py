
#A simple script for automating/avoiding bugs in the creation of the full cross
#product of debug callers.
targets = []
targets.append(["", "Global"])
targets.append(["(log *Logger) ", "log"])


types=[]
types.append("Finest")
types.append("Fine")
types.append("Debug")
types.append("Trace")
types.append("Info")
types.append("Warning")
types.append("Error")
types.append("Critical")

funcs=[]


func="\
// Utility for #UPPERREP log messages\n\
// This behaves like Logf but with the #UPPERREP log level.\n\
func #TARGET1#NORMALREPf(format string, args ...interface{}) {\n\
\t#TARGET2.intLogf(#UPPERREP, format, args...)\n\
}"
funcs.append(func)

func="\
// Utility for #UPPERREP log messages\n\
// This behaves like Logln but with the #UPPERREP log level.\n\
func #TARGET1#NORMALREPln(args ...interface{}) {\n\
\t#TARGET2.intLogln(#UPPERREP, args...)\n\
}"
funcs.append(func)

func="\
// Utility for #UPPERREP log messages\n\
// This behaves like Logc but with the #UPPERREP log level.\n\
func #TARGET1#NORMALREPc(closure func() string) {\n\
\t#TARGET2.intLogc(#UPPERREP, closure)\n\
}"
funcs.append(func)
first = True

print("// Interface functions for the Global and Local loggers for each type.")
print("// This was auto generated with makeInterface.py.\n")
print "package logit\n"
for target in targets:
    for typev in types:
        for func in funcs:
            if not first:
                print ""
            first = False
            rep = func
            rep = rep.replace("#TARGET1", target[0])
            rep = rep.replace("#TARGET2", target[1])
            rep = rep.replace("#NORMALREP", typev)
            rep = rep.replace("#UPPERREP", typev.upper())
            print rep
