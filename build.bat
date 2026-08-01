@rem 关闭命令回显，且当前行也不显示（@ 符号抑制该行自身的回显），使输出更简洁
@echo off

@rem 创建一个局部环境，确保变量只在这个批处理文件中有效
@setlocal

@title MinIO Build

@rem 指定 Go 版本
set PATH=D:\program\go1.26.4.windows-amd64\bin;C:\msys64\mingw64\bin;C:\msys64\usr\bin;%PATH%

@rem 启用 CGO
set CGO_ENABLED=1

@rem 操作系统
for /f "delims=" %%i in ('go env GOOS') do set OS=%%i
@rem CPU 架构
for /f "delims=" %%i in ('go env GOARCH') do set ARCH=%%i
@rem 当前目录
set CURR_DIR=%cd%
@rem 程序目录
set PROGRAM_DIR=%CURR_DIR%
@rem 输出名称
set OUT_NAME=minio.exe
@rem 输出目录
set OUT_DIR=%CURR_DIR%\dist\minio-%OS%-%ARCH%
@rem 输出路径
set OUT_PATH=%OUT_DIR%\%OUT_NAME%

echo OS      : %OS%
echo ARCH    : %ARCH%
echo OutName : %OUT_NAME%
echo OutDir  : %OUT_DIR%

@rem 如果目录不存在则创建
if not exist "%OUT_DIR%" (
    mkdir "%OUT_DIR%"
)

@rem 拷贝文件
@rem 隐藏无用输出：> nul（标准输出），2> nul（错误输出）
copy /Y "%PROGRAM_DIR%\start.bat" "%OUT_DIR%\" > nul 2> nul

@rem 构建
cd "%PROGRAM_DIR%" && go build -ldflags="-s -w" -o "%OUT_PATH%"

@rem 压缩可执行文件
@rem https://upx.github.io
::upx --best --backup "%OUT_PATH%"
::upx --brute --backup "%OUT_PATH%"

@endlocal
@pause