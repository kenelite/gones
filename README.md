# Gones - NES Emulator

一个用 Go 语言编写的 NES 模拟器，基于 Ebiten 游戏引擎。

## 功能特性

- ✅ 基本的 NES 硬件模拟
- ✅ CPU (6502) 模拟
- ✅ PPU (图形处理单元) 模拟
- ✅ APU (音频处理单元) 模拟
- ✅ 内存总线模拟
- ✅ ROM 加载和解析
- ✅ 键盘和手柄输入支持
- ✅ 基本的图形渲染

## 系统要求

- Go 1.24.3 或更高版本
- macOS (已测试) / Linux / Windows

## 安装和运行

1. 克隆仓库：
```bash
git clone https://github.com/kenelite/gones.git
cd gones
```

2. 安装依赖：
```bash
go mod tidy
```

3. 运行模拟器：
```bash
go run .
```

## 使用方法

1. 启动程序后，会弹出文件选择对话框
2. 选择一个 .nes 格式的 ROM 文件
3. 模拟器会自动加载并运行游戏
4. 使用以下按键控制：
   - **Z**: A 按钮
   - **X**: B 按钮
   - **A**: Select 按钮
   - **S**: Start 按钮
   - **方向键**: 控制方向
   - **ESC**: 切换菜单

## 支持的游戏

目前支持 Mapper 0 的 NES 游戏，包括：
- Super Mario Bros.
- Duck Hunt
- 其他使用 Mapper 0 的游戏

## 项目结构

```
gones/
├── cmd/           # 命令行工具
├── core/          # 核心模拟器组件
│   ├── apu/      # 音频处理单元
│   ├── bus/      # 内存总线
│   ├── cpu/      # CPU 模拟
│   ├── input/    # 输入处理
│   ├── ppu/      # 图形处理单元
│   └── rom/      # ROM 加载和解析
├── internal/      # 内部包
│   ├── assets/   # 资源文件
│   └── ui/       # 用户界面
└── main.go        # 主程序入口
```

## 开发状态

这是一个正在开发中的项目。目前实现了：

- 基本的 NES 硬件架构
- ROM 加载和解析
- 简单的图形渲染
- 输入处理

## 已知问题

- 音频功能尚未完全实现
- 某些复杂的游戏可能无法正常运行
- 性能优化仍在进行中

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 致谢

- [Ebiten](https://ebiten.org/) - 2D 游戏引擎
- [6502.org](http://6502.org/) - 6502 处理器文档
- NES 开发社区
