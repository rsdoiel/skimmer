call make.bat
if not exist mkdir %userprofile%\bin
move bin\skimmer.exe %userprofile%\bin\
set PATH=%PATH%;%userprofile%\bin
echo OFF
echo This has installed skimmer in %userprofile%\bin
echo This directory needs to be in your PATH
echo You can do this with the "set" command
echo
echo        set PATH=%PATH%;%userprofile%\bin
echo ON