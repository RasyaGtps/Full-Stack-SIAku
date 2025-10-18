@echo off
echo.
echo ================================================
echo   RESET WhatsApp Session
echo ================================================
echo.
echo Menghapus session lama...

if exist auth_info (
    rmdir /s /q auth_info
    echo ✓ Folder auth_info dihapus
) else (
    echo - Folder auth_info tidak ada
)

if exist data.json (
    del /f data.json
    echo ✓ File data.json dihapus
) else (
    echo - File data.json tidak ada
)

if exist baileys_store.json (
    del /f baileys_store.json
    echo ✓ File baileys_store.json dihapus
) else (
    echo - File baileys_store.json tidak ada
)

echo.
echo ✅ Reset selesai!
echo.
echo Jalankan: npm start
echo.
pause

