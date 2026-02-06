import os

# --- 项目配置 ---
LIB_SOURCES = ["src/SimpleCache.cpp"]  # 共享实现文件，链接到可执行文件
MAIN_TARGET = "src/main.cpp"           # 主程序
CLI_TARGET = "apps/cli.cpp"            # CLI 程序
PROJECT_NAME = "main"
CXX_STANDARD = 17
TEST_DIR = "tests"

# --- 扫描 tests 目录 ---
tests = [f for f in os.listdir(TEST_DIR) if f.endswith(".cpp")]
tests.sort()  # 保证顺序一致

# --- 拼接 CMakeLists 内容 ---
lines = [
    "cmake_minimum_required(VERSION 3.10)",
    f"project({PROJECT_NAME})",
    f"set(CMAKE_CXX_STANDARD {CXX_STANDARD})",
    "set(CMAKE_EXPORT_COMPILE_COMMANDS ON)",
    "",
    "include_directories(include)",
    "",
    "enable_testing()"
]

# --- 添加测试 ---
for test in tests:
    name = os.path.splitext(test)[0]
    test_path = os.path.join(TEST_DIR, test)
    srcs = " ".join([test_path] + LIB_SOURCES)
    lines.append(f"add_executable({name} {srcs})")
    lines.append(f"add_test(NAME {name} COMMAND {name})")
    lines.append("")

# --- 添加主程序 ---
main_srcs = " ".join([MAIN_TARGET] + LIB_SOURCES)
lines.append(f"add_executable({PROJECT_NAME} {main_srcs})")
lines.append("")

# --- 添加 CLI 程序 ---
cli_srcs = " ".join([CLI_TARGET] + LIB_SOURCES)
lines.append(f"add_executable(cli {cli_srcs})")
lines.append("")

# --- 写入文件 ---
cmake_content = "\n".join(lines)
with open("CMakeLists.txt", "w", encoding="utf-8") as f:
    f.write(cmake_content)

print("✅ CMakeLists.txt 已生成，包含以下测试：")
for t in tests:
    print("  -", t)
