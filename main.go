package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type Manifest struct {
	Permissions     []string `json:"permissions"`
	HostPermissions []string `json:"host_permissions"`
}

type Finding struct {
	File    string
	Line    int
	Pattern string
	Content string
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func cloneRepo(repoURL, destDir string) error {
	fmt.Printf("[*] Clonando repositório: %s\n", repoURL)
	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		fmt.Printf("[*] Removendo diretório existente: %s\n", destDir)
		os.RemoveAll(destDir)
	}
	cmd := exec.Command("git", "clone", repoURL, destDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("erro ao clonar: %v", err)
	}
	fmt.Println("[+] Repositório clonado com sucesso.")
	return nil
}

func parseManifest(extDir string) *Manifest {
	manifestPath := filepath.Join(extDir, "manifest.json")
	file, err := os.Open(manifestPath)
	if err != nil {
		fmt.Println("[-] manifest.json não encontrado.")
		return nil
	}
	defer file.Close()

	var manifest Manifest
	if err := json.NewDecoder(file).Decode(&manifest); err != nil {
		fmt.Println("[-] Erro ao parsear manifest.json.")
		return nil
	}
	return &manifest
}

func scanFile(path, extDir string, patterns []*regexp.Regexp, findingsChan chan<- []Finding, wg *sync.WaitGroup) {
	defer wg.Done()
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("[-] Erro ao ler o arquivo %s: %v\n", path, err)
		return
	}
	defer file.Close()

	var localFindings []Finding
	relPath, _ := filepath.Rel(extDir, path)
	scanner := bufio.NewScanner(file)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		for _, pattern := range patterns {
			if pattern.MatchString(line) {
				content := strings.TrimSpace(line)
				if len(content) > 100 {
					content = content[:100]
				}
				localFindings = append(localFindings, Finding{
					File:    relPath,
					Line:    lineNum,
					Pattern: pattern.String(),
					Content: content,
				})
			}
		}
		lineNum++
	}
	if len(localFindings) > 0 {
		findingsChan <- localFindings
	}
}

func analyzeJSFiles(extDir string) []Finding {
	suspiciousStrings := []string{
		`eval\(`,
		`fetch\(`,
		`XMLHttpRequest`,
		`chrome\.tabs\.executeScript`,
		`document\.write`,
		`innerHTML`,
		`chrome\.storage`,
	}

	var patterns []*regexp.Regexp
	for _, p := range suspiciousStrings {
		patterns = append(patterns, regexp.MustCompile(p))
	}

	findingsChan := make(chan []Finding, 100)
	var wg sync.WaitGroup

	err := filepath.Walk(extDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".js") {
			wg.Add(1)
			go scanFile(path, extDir, patterns, findingsChan, &wg)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("[-] Erro ao buscar arquivos: %v\n", err)
	}

	go func() {
		wg.Wait()
		close(findingsChan)
	}()

	var allFindings []Finding
	for f := range findingsChan {
		allFindings = append(allFindings, f...)
	}
	return allFindings
}

func generateReport(extName string, manifest *Manifest, findings []Finding, reportDir string) error {
	os.MkdirAll(reportDir, os.ModePerm)
	reportPath := filepath.Join(reportDir, fmt.Sprintf("relatorio-%s.md", extName))

	file, err := os.Create(reportPath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	defer writer.Flush()

	writer.WriteString(fmt.Sprintf("# Relatório final - %s\n\n", extName))
	writer.WriteString("## Resumo executivo\n\n")
	writer.WriteString(fmt.Sprintf("- **Extensão analisada:** %s\n", extName))
	writer.WriteString("- **Objetivo da análise:** Análise estática automatizada (Otimizada em Go)\n\n")

	if manifest != nil {
		writer.WriteString("## Análise do Manifest\n\n")
		writer.WriteString("### Permissões Solicitadas\n")
		if len(manifest.Permissions) > 0 {
			for _, p := range manifest.Permissions {
				writer.WriteString(fmt.Sprintf("- `%s`\n", p))
			}
		} else {
			writer.WriteString("- Nenhuma permissão explícita.\n")
		}

		if len(manifest.HostPermissions) > 0 {
			writer.WriteString("\n### Host Permissions\n")
			for _, hp := range manifest.HostPermissions {
				writer.WriteString(fmt.Sprintf("- `%s`\n", hp))
			}
		}
	} else {
		writer.WriteString("## Análise do Manifest\n\n")
		writer.WriteString("- Manifesto não encontrado ou inválido.\n")
	}

	writer.WriteString("\n## Achados no Código JavaScript\n\n")
	if len(findings) > 0 {
		writer.WriteString("| Arquivo | Linha | Padrão Encontrado | Trecho |\n")
		writer.WriteString("|---------|-------|-------------------|--------|\n")
		for _, f := range findings {
			cleanContent := strings.ReplaceAll(f.Content, "|", "\\|")
			writer.WriteString(fmt.Sprintf("| `%s` | %d | `%s` | `%s` |\n", f.File, f.Line, f.Pattern, cleanContent))
		}
	} else {
		writer.WriteString("- Nenhum padrão suspeito identificado nos arquivos `.js` analisados.\n")
	}

	writer.WriteString("\n## Conclusão Inicial\n\n")
	writer.WriteString("- *Requer revisão manual baseada nos achados acima.*\n")

	fmt.Printf("[+] Relatório gerado em: %s\n", reportPath)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: shield_analyzer <URL_ou_Diretorio>")
		return
	}
	target := os.Args[1]

	workDir := filepath.Join("src", "codigo-analisado", "temp_ext")
	extName := "extensao-local"
	targetDir := target

	if isURL(target) {
		parts := strings.Split(target, "/")
		extName = strings.TrimSuffix(parts[len(parts)-1], ".git")
		if err := cloneRepo(target, workDir); err != nil {
			fmt.Printf("Falha ao clonar repo: %v\n", err)
			return
		}
		targetDir = workDir
	} else {
		if _, err := os.Stat(targetDir); os.IsNotExist(err) {
			fmt.Printf("[-] Diretório local não encontrado: %s\n", targetDir)
			return
		}
		extName = filepath.Base(filepath.Clean(targetDir))
	}

	fmt.Printf("\n[*] Iniciando análise otimizada para: %s\n", extName)

	manifest := parseManifest(targetDir)
	findings := analyzeJSFiles(targetDir)

	generateReport(extName, manifest, findings, "reports")

	if isURL(target) {
		fmt.Printf("[*] Limpando diretório temporário: %s\n", workDir)
		os.RemoveAll(workDir)
	}

	fmt.Println("[*] Análise concluída.")
}
