call make.bat
if not exist mkdir %userprofile%\bin
copy /Y bin\*.exe %userprofile%\bin\
set PATH=%PATH%;%userprofile%\bin
@echo     
@echo This has installed skimmer in %userprofile%\bin
@echo This directory needs to be in your PATH
@echo You can do this with the "set" command
@echo
@echo        set PATH=%PATH%;%userprofile%\bin
@echo      