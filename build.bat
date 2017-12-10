
REM TODO: Get the release version using Git tag
set VERSION=0.0.0

REM Build back end binary
REM Set linker flags
set LDFLAGS=-X main.version=%VERSION%
IF "%1" == "pd" (
    set LDFLAGS=%LDFLAGS% -X main.env=production -X main.opSys=win -s
)
go build -ldflags "%LDFLAGS%"

REM Build the client code
cd ui
IF "%1" == "pd" (
    set NODE_ENV=production
)
call npm run build
IF "%1" == "pd" (
    call npm run dist
)
cd ..
set NODE_ENV=

REM Package together
IF "%1" == "pd" (
    mkdir gringotts-%VERSION%-win
    copy gringotts.exe gringotts-%VERSION%-win
    copy "ui\dist\gringotts-%VERSION%.exe" gringotts-%VERSION%-win
)
