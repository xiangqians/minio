@rem 关闭命令回显，且当前行也不显示（@ 符号抑制该行自身的回显），使输出更简洁
@echo off

@rem 创建一个局部环境，确保变量只在这个批处理文件中有效
@setlocal

@title MinIO

minio.exe server ./data --address ":9000" --console-address ":9001"

endlocal
pause