package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type ProjectConfig struct {
	ProjectName  string // 项目名称同时也是品牌名称
	PackageName  string
	TemplatePath string
	TargetPath   string
}

// 查找qor5boot项目根目录
func findQor5bootRoot() (string, error) {
	// 首先尝试当前目录
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("无法获取当前工作目录: %v", err)
	}

	// 检查当前目录是否是qor5boot项目
	if isQor5bootRoot(cwd) {
		return cwd, nil
	}

	// 如果不是，查找二进制所在目录
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("无法获取可执行文件路径: %v", err)
	}

	execDir := filepath.Dir(execPath)

	// 检查可执行文件所在目录是否是qor5boot项目
	if isQor5bootRoot(execDir) {
		return execDir, nil
	}

	// 如果都不是，尝试向上找父目录
	dir := cwd
	for i := 0; i < 5; i++ { // 最多向上查找5层
		dir = filepath.Dir(dir)
		if isQor5bootRoot(dir) {
			return dir, nil
		}
	}

	return "", fmt.Errorf("找不到qor5boot项目根目录，请确保在qor5boot项目内或其子目录中运行该工具")
}

// 检查目录是否是qor5boot项目根目录
func isQor5bootRoot(dir string) bool {
	// 检查go.mod文件是否存在
	modPath := filepath.Join(dir, "go.mod")
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		return false
	}

	// 检查go.mod文件内容是否包含qor5boot
	content, err := os.ReadFile(modPath)
	if err != nil {
		return false
	}

	return strings.Contains(string(content), "qor5boot")
}

func main() {
	fmt.Println("QOR5Boot 项目初始化工具")
	fmt.Println("========================")

	// 查找qor5boot项目根目录
	qor5bootRoot, err := findQor5bootRoot()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("使用模板目录: %s\n", qor5bootRoot)

	config := getProjectConfig(qor5bootRoot)

	// 验证项目路径不存在或为空
	if dirExists(config.TargetPath) {
		files, err := os.ReadDir(config.TargetPath)
		if err != nil {
			fmt.Printf("无法读取目标目录: %v\n", err)
			return
		}
		if len(files) > 0 {
			fmt.Printf("错误: 目标目录 %s 已存在且不为空\n", config.TargetPath)
			return
		}
	}

	fmt.Println("开始创建新项目...")

	// 创建项目目录
	if err := os.MkdirAll(config.TargetPath, 0755); err != nil {
		fmt.Printf("创建项目目录失败: %v\n", err)
		return
	}

	// 复制并处理项目文件
	if err := processProject(config); err != nil {
		fmt.Printf("处理项目文件失败: %v\n", err)
		return
	}

	// 初始化 Git 仓库
	initGitRepo(config.TargetPath)

	fmt.Printf("\n项目 %s 创建成功！\n", config.ProjectName)
	fmt.Printf("项目路径: %s\n", config.TargetPath)
	fmt.Println("\n您可以通过以下命令开始开发:")
	fmt.Printf("  cd %s\n", config.ProjectName)
	fmt.Println("  cp dev_env.example dev_env")
	fmt.Println("  # 编辑 dev_env 设置环境变量")
	fmt.Println("  ./dev.sh")
}

func getProjectConfig(qor5bootRoot string) ProjectConfig {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("请输入项目名称（将同时用作品牌名称）: ")
	projectName := readLine(reader)

	fmt.Print("请输入包名称（例如: github.com/yourname/" + projectName + "）: ")
	packageName := readLine(reader)
	if packageName == "" {
		// 默认使用项目名作为包名的最后一部分
		packageName = "github.com/yourname/" + projectName
	}

	// 获取包名最后一部分作为目录名
	packageParts := strings.Split(packageName, "/")
	lastPart := packageParts[len(packageParts)-1]

	// 询问目标目录
	fmt.Printf("请输入目标目录（按Enter使用默认值: 当前目录/%s）: ", lastPart)
	targetDir := readLine(reader)
	var targetPath string
	if targetDir == "" {
		// 默认为当前工作目录加上包名最后一部分
		cwd, _ := os.Getwd()
		targetPath = filepath.Join(cwd, lastPart)
	} else {
		// 如果用户输入了目标目录，使用绝对路径
		if filepath.IsAbs(targetDir) {
			targetPath = filepath.Join(targetDir, lastPart)
		} else {
			absPath, _ := filepath.Abs(targetDir)
			targetPath = filepath.Join(absPath, lastPart)
		}
	}

	return ProjectConfig{
		ProjectName:  projectName,
		PackageName:  packageName,
		TemplatePath: qor5bootRoot,
		TargetPath:   targetPath,
	}
}

func readLine(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func processProject(config ProjectConfig) error {
	// 要忽略的目录和文件
	ignores := []string{
		".git",
		"vendor",
		"qor5boot", // 编译后的二进制文件
		"cmd/newboot",
		"cmd/data-resetor",
		"cmd/publisher",
	}

	// 获取原包名
	oldPackageName := getOldPackageName(config.TemplatePath)
	if oldPackageName == "" {
		fmt.Println("警告: 无法确定原始包名，使用默认值 'github.com/naokij/qor5boot'")
		oldPackageName = "github.com/naokij/qor5boot"
	}

	fmt.Printf("正在从包 '%s' 复制到 '%s'\n", oldPackageName, config.PackageName)

	// 遍历并处理文件
	return filepath.Walk(config.TemplatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 获取相对路径
		relPath, err := filepath.Rel(config.TemplatePath, path)
		if err != nil {
			return err
		}

		// 跳过根目录
		if relPath == "." {
			return nil
		}

		// 检查是否在忽略列表中
		for _, ignore := range ignores {
			if strings.HasPrefix(relPath, ignore) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		// 处理目标路径
		targetPath := filepath.Join(config.TargetPath, relPath)

		// 如果是目录，创建对应的目录
		if info.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		// 读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// 处理文件内容
		newContent := processFileContent(string(content), relPath, oldPackageName, config)

		// 确保目标目录存在
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}

		// 写入新文件，保持原有权限
		return os.WriteFile(targetPath, []byte(newContent), info.Mode())
	})
}

func processFileContent(content, relPath, oldPackageName string, config ProjectConfig) string {
	// 替换包名
	newContent := strings.ReplaceAll(content, oldPackageName, config.PackageName)

	// 根据文件类型进行特殊处理
	switch {
	case filepath.Ext(relPath) == ".go":
		// 处理 Go 文件
		newContent = processGoFile(newContent, config)
	case strings.HasSuffix(relPath, ".yml") || strings.HasSuffix(relPath, ".yaml"):
		// 处理 YAML 文件（部署配置等）
		newContent = processYamlFile(newContent, config)
	case strings.HasSuffix(relPath, ".j2"):
		// 处理 Jinja2 模板文件
		newContent = processJ2Template(newContent, config)
	case strings.HasSuffix(relPath, ".sh"):
		// 处理 Shell 脚本
		newContent = processShellScript(newContent, config)
	}

	return newContent
}

func processGoFile(content string, config ProjectConfig) string {
	// 处理品牌名称等
	return content
}

func processYamlFile(content string, config ProjectConfig) string {
	// 替换 YAML 文件中的应用名称
	re := regexp.MustCompile(`app_name:\s*qor5boot`)
	content = re.ReplaceAllString(content, fmt.Sprintf("app_name: %s", config.ProjectName))

	// 替换服务名称
	namePattern := regexp.MustCompile(`name:\s*qor5boot`)
	content = namePattern.ReplaceAllString(content, fmt.Sprintf("name: %s", config.ProjectName))

	// 替换二进制文件名称
	binaryPattern := regexp.MustCompile(`app_binary:\s*qor5boot`)
	content = binaryPattern.ReplaceAllString(content, fmt.Sprintf("app_binary: %s", config.ProjectName))

	return content
}

func processJ2Template(content string, config ProjectConfig) string {
	// 这里不需要替换，因为模板中使用的是变量如 {{ app_name }}
	return content
}

func processShellScript(content string, config ProjectConfig) string {
	// 替换部署脚本中的应用名称
	re := regexp.MustCompile(`qor5boot`)
	return re.ReplaceAllString(content, config.ProjectName)
}

func getOldPackageName(templatePath string) string {
	modPath := filepath.Join(templatePath, "go.mod")
	content, err := os.ReadFile(modPath)
	if err != nil {
		return ""
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 2 {
			return parts[1]
		}
	}
	return ""
}

func initGitRepo(targetPath string) {
	cmd := exec.Command("git", "init")
	cmd.Dir = targetPath
	if err := cmd.Run(); err != nil {
		fmt.Printf("警告: 无法初始化Git仓库: %v\n", err)
		return
	}
	fmt.Println("Git 仓库初始化成功")
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}
