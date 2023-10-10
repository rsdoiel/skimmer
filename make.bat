if not exist mkdir bin
go build -o bin\skimmer.exe cmd\skimmer\skimmer.go
.\bin\skimmer -version
go build -o bin\skim2md.exe cmd\skim2md\skim2md.go
.\bin\skim2md -version