@echo off
title Stop Blog Services
echo Stopping blog services...
for /f "tokens=5" %%p in ('netstat -ano ^| findstr /c:":5173" ^| findstr /c:"LISTENING"') do taskkill /F /PID %%p 2>nul
for /f "tokens=5" %%p in ('netstat -ano ^| findstr /c:":8080" ^| findstr /c:"LISTENING"') do taskkill /F /PID %%p 2>nul
echo Done. Ports 5173 and 8080 released.
pause
