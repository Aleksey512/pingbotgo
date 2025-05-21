package ping

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

func GeneratePingReport(servers map[string]string) string {
	var wg sync.WaitGroup
	results := make(map[string]string, len(servers))
	mu := sync.Mutex{}

	for name, ip := range servers {
		wg.Add(1)
		go func(name, ip string) {
			defer wg.Done()
			_, status := pingServer(ip)
			mu.Lock()
			results[name] = status
			mu.Unlock()
		}(name, ip)
	}

	wg.Wait()

	var report strings.Builder
	for name, status := range results {
		report.WriteString(fmt.Sprintf("%s - %s\n", name, status))
	}

	return report.String()
}

func pingServer(ip string) (bool, string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-n", "4", ip)
	} else {
		cmd = exec.Command("ping", "-c", "4", ip)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Sprintf("❌ Сервер %s недоступен: %v", ip, err)
	}

	if strings.Contains(string(output), "timeout") {
		return false, fmt.Sprintf("❌ Сервер %s недоступен (timeout)", ip)
	}

	return true, fmt.Sprintf("✅ Сервер %s доступен", ip)
}
