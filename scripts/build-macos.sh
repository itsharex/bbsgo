#!/bin/bash
# macOS GUI 应用构建脚本
# 在 M1 Mac 上运行，生成 bbsgo-macos-arm64.zip

set -e

VERSION=${1:-"v1.0.1"}
APP_NAME="bbsgo"
ARCH="arm64"
OUTPUT="bbsgo-macos-${ARCH}.zip"

echo "=== 构建 macOS GUI 应用 ==="
echo "版本: $VERSION"
echo "架构: $ARCH"

# 1. 构建 Go GUI 二进制
echo ""
echo "[1/4] 编译 Go GUI 二进制..."
cd server
CGO_ENABLED=1 GOOS=darwin GOARCH=${ARCH} go build -tags "gui sqlite_fts5" -ldflags="-s -w" -o bbsgo-macos
cd ..

# 2. 创建 .app 目录结构
echo ""
echo "[2/4] 创建 .app 目录结构..."
rm -rf "${APP_NAME}.app"
mkdir -p "${APP_NAME}.app/Contents/MacOS"
mkdir -p "${APP_NAME}.app/Contents/Resources"

# 3. 复制二进制并设置可执行权限
echo ""
echo "[3/4] 配置应用包..."
cp server/bbsgo-macos "${APP_NAME}.app/Contents/MacOS/${APP_NAME}"
chmod +x "${APP_NAME}.app/Contents/MacOS/${APP_NAME}"

# 4. 生成 ICNS 图标（使用 site/src/assets/bbsgo.png）
ICON_SOURCE="site/src/assets/bbs.png"
if [ -f "$ICON_SOURCE" ]; then
    echo "生成 ICNS 图标（使用 $ICON_SOURCE）..."
    ICONSET_DIR="${APP_NAME}.app/Contents/Resources/icon.iconset"
    mkdir -p "$ICONSET_DIR"
    
    sips -z 16 16     "$ICON_SOURCE" --out "$ICONSET_DIR/icon_16x16.png"      >/dev/null 2>&1
    sips -z 32 32     "$ICON_SOURCE" --out "$ICONSET_DIR/icon_16x16@2x.png"   >/dev/null 2>&1
    sips -z 32 32     "$ICON_SOURCE" --out "$ICONSET_DIR/icon_32x32.png"      >/dev/null 2>&1
    sips -z 64 64     "$ICON_SOURCE" --out "$ICONSET_DIR/icon_32x32@2x.png"   >/dev/null 2>&1
    sips -z 128 128   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_128x128.png"    >/dev/null 2>&1
    sips -z 256 256   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_128x128@2x.png" >/dev/null 2>&1
    sips -z 256 256   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_256x256.png"    >/dev/null 2>&1
    sips -z 512 512   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_256x256@2x.png" >/dev/null 2>&1
    sips -z 512 512   "$ICON_SOURCE" --out "$ICONSET_DIR/icon_512x512.png"    >/dev/null 2>&1
    sips -z 1024 1024 "$ICON_SOURCE" --out "$ICONSET_DIR/icon_512x512@2x.png" >/dev/null 2>&1
    
    iconutil -c icns "$ICONSET_DIR" -o "${APP_NAME}.app/Contents/Resources/${APP_NAME}.icns"
    rm -rf "$ICONSET_DIR"
else
    echo "警告: 未找到 $ICON_SOURCE，跳过图标生成"
fi

# 5. 创建 Info.plist（包含图标配置）
echo ""
echo "[4/4] 创建 Info.plist..."
cat > "${APP_NAME}.app/Contents/Info.plist" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleExecutable</key>
    <string>bbsgo</string>
    <key>CFBundleIdentifier</key>
    <string>com.bbsgo.macos</string>
    <key>CFBundleName</key>
    <string>bbsgo</string>
    <key>CFBundleDisplayName</key>
    <string>bbsgo</string>
    <key>CFBundleIconFile</key>
    <string>bbsgo</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>LSMinimumSystemVersion</key>
    <string>11.0</string>
    <key>LSUIElement</key>
    <true/>
    <key>NSHighResolutionCapable</key>
    <true/>
</dict>
</plist>
EOF

# 6. 打包成 zip
echo ""
echo "打包成 zip..."
rm -f "$OUTPUT"
codesign --force --deep --sign - "${APP_NAME}.app"
zip -r "$OUTPUT" "${APP_NAME}.app"

echo ""
echo "=== 构建完成 ==="
echo "输出文件: $OUTPUT"
echo "大小: $(du -h "$OUTPUT" | cut -f1)"
echo ""
echo "上传到 GitHub Release:"
echo "  gh release upload $VERSION $OUTPUT --clobber"