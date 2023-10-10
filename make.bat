if not exist mkdir bin
go build -o bin\skimmer.exe cmd\skimmer\skimmer.go
.\bin\skimmer -version