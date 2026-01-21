#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUT_DIR="${OUT_DIR:-${ROOT}/dist}"
HOST_GOOS="$(go env GOOS)"

mkdir -p "${OUT_DIR}"

echo "[1/3] 构建 Linux CLI..."
go build -o "${OUT_DIR}/sunshine-proxy" "${ROOT}/cmd/proxy"

echo "[2/3] 构建 Windows GUI..."
if [[ "${HOST_GOOS}" == "windows" ]]; then
	go build -o "${OUT_DIR}/sunshine-proxy.exe" "${ROOT}/cmd/proxy-gui"
elif [[ "${BUILD_WINDOWS_GUI:-0}" == "1" ]]; then
	GOOS=windows GOARCH=amd64 go build -o "${OUT_DIR}/sunshine-proxy.exe" "${ROOT}/cmd/proxy-gui"
else
	echo "跳过 Windows GUI（非 Windows 环境）。如需交叉编译请设置 BUILD_WINDOWS_GUI=1 并准备对应工具链。"
fi

echo "[3/3] 构建 Windows CLI..."
GOOS=windows GOARCH=amd64 go build -o "${OUT_DIR}/sunshine-proxy-cli.exe" "${ROOT}/cmd/proxy"

echo "构建完成: ${OUT_DIR}"
