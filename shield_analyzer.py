import argparse
import os
import subprocess
import json
import re
from pathlib import Path
from urllib.parse import urlparse

def is_url(url):
    try:
        result = urlparse(url)
        return all([result.scheme, result.netloc])
    except ValueError:
        return False

import shutil

def clone_repo(url, dest_dir):
    print(f"[*] Clonando repositório: {url}")
    if os.path.exists(dest_dir):
        print(f"[*] Removendo diretório existente: {dest_dir}")
        shutil.rmtree(dest_dir)
    subprocess.run(['git', 'clone', url, dest_dir], check=True)
    print("[+] Repositório clonado com sucesso.")

def parse_manifest(ext_dir):
    manifest_path = Path(ext_dir) / 'manifest.json'
    if not manifest_path.exists():
        print("[-] manifest.json não encontrado.")
        return None

    with open(manifest_path, 'r', encoding='utf-8') as f:
        try:
            return json.load(f)
        except json.JSONDecodeError:
            print("[-] Erro ao parsear manifest.json.")
            return None

def analyze_js_files(ext_dir):
    suspicious_patterns = [
        'eval\\(',
        'fetch\\(',
        'XMLHttpRequest',
        'chrome.tabs.executeScript',
        'document.write',
        'innerHTML',
        'chrome.storage'
    ]
    findings = []

    for path in Path(ext_dir).rglob('*.js'):
        try:
            with open(path, 'r', encoding='utf-8') as f:
                lines = f.readlines()
                for i, line in enumerate(lines):
                    for pattern in suspicious_patterns:
                        if re.search(pattern, line):
                            findings.append({
                                'file': str(path.relative_to(ext_dir)),
                                'line': i + 1,
                                'pattern': pattern,
                                'content': line.strip()[:100]
                            })
        except Exception as e:
            print(f"[-] Erro ao ler o arquivo {path}: {e}")

    return findings

def generate_report(ext_name, manifest_data, js_findings, report_dir):
    os.makedirs(report_dir, exist_ok=True)
    report_path = Path(report_dir) / f"relatorio-{ext_name}.md"

    with open(report_path, 'w', encoding='utf-8') as f:
        f.write(f"# Relatório final - {ext_name}\n\n")
        f.write("## Resumo executivo\n\n")
        f.write(f"- **Extensão analisada:** {ext_name}\n")
        f.write("- **Objetivo da análise:** Análise estática automatizada\n\n")

        if manifest_data:
            f.write("## Análise do Manifest\n\n")
            f.write("### Permissões Solicitadas\n")
            permissions = manifest_data.get('permissions', [])
            host_permissions = manifest_data.get('host_permissions', [])

            if permissions:
                for p in permissions:
                    f.write(f"- `{p}`\n")
            else:
                f.write("- Nenhuma permissão explícita.\n")

            if host_permissions:
                f.write("\n### Host Permissions\n")
                for hp in host_permissions:
                    f.write(f"- `{hp}`\n")
        else:
            f.write("## Análise do Manifest\n\n")
            f.write("- Manifesto não encontrado ou inválido.\n")

        f.write("\n## Achados no Código JavaScript\n\n")
        if js_findings:
            f.write("| Arquivo | Linha | Padrão Encontrado | Trecho |\n")
            f.write("|---------|-------|-------------------|--------|\n")
            for finding in js_findings:
                clean_content = finding['content'].replace('|', '\\|')
                f.write(f"| `{finding['file']}` | {finding['line']} | `{finding['pattern']}` | `{clean_content}` |\n")
        else:
            f.write("- Nenhum padrão suspeito identificado nos arquivos `.js` analisados.\n")

        f.write("\n## Conclusão Inicial\n\n")
        f.write("- *Requer revisão manual baseada nos achados acima.*\n")

    print(f"[+] Relatório gerado em: {report_path}")

def main():
    parser = argparse.ArgumentParser(description="SHIELD_EXTENSION - Analisador Estático de Extensões")
    parser.add_argument("target", help="URL do repositório GitHub da extensão ou caminho local para o diretório da extensão.")
    args = parser.parse_args()

    target = args.target
    work_dir = "src/codigo-analisado/temp_ext"
    ext_name = "extensao-local"

    if is_url(target):
        ext_name = target.split('/')[-1].replace('.git', '')
        clone_repo(target, work_dir)
        target_dir = work_dir
    else:
        target_dir = target
        if os.path.exists(target_dir):
            ext_name = os.path.basename(os.path.normpath(target_dir))
            if not ext_name:
                ext_name = "extensao-local"
        else:
            print(f"[-] Diretório local não encontrado: {target_dir}")
            return

    print(f"\n[*] Iniciando análise para: {ext_name}")

    manifest_data = parse_manifest(target_dir)
    js_findings = analyze_js_files(target_dir)

    generate_report(ext_name, manifest_data, js_findings, "reports")

    if is_url(target) and os.path.exists(work_dir):
        print(f"[*] Limpando diretório temporário: {work_dir}")
        shutil.rmtree(work_dir)

    print("[*] Análise concluída.")

if __name__ == "__main__":
    main()
