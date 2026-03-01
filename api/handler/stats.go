package handler

/*
#ifdef __APPLE__
#include <mach/mach.h>

static mach_port_t my_task_self() {
    return mach_task_self();
}

static int get_mem(uint64_t *vsize, uint64_t *rss) {
    task_vm_info_data_t info;
    mach_msg_type_number_t count = TASK_VM_INFO_COUNT;

    kern_return_t kerr = task_info(
        my_task_self(),
        TASK_VM_INFO,
        (task_info_t)&info,
        &count
    );

    if (kerr != KERN_SUCCESS) return -1;

    *vsize = info.virtual_size;
    *rss   = info.phys_footprint;

    return 0;
}
#else
// Placeholder for Linux/other OS to prevent "undefined" errors
#include <stdint.h>
static int get_mem(uint64_t *vsize, uint64_t *rss) {
    *vsize = 0;
    *rss = 0;
    return -1;
}
#endif
*/
import "C"
import (
    "encoding/json"
    "net/http"
    "os"
    "runtime"
    "strconv"
    "strings"
    "sync"
    "syscall"
    "time"
)

type Stats struct {
    RSSMB        float64 `json:"rss_mb"`
    VMSMB        float64 `json:"vms_mb"`
    HeapAllocMB  float64 `json:"heap_alloc_mb"`
    HeapSysMB    float64 `json:"heap_sys_mb"`
    GoSysMB      float64 `json:"go_sys_mb"`
    CPUPercent   float64 `json:"cpu_percent"`
    NumGoroutine int     `json:"num_goroutine"`
    NumCPU       int     `json:"num_cpu"`
    UptimeSec    float64 `json:"uptime_sec"`
    Platform     string  `json:"platform"`
}

var (
    lastCPUTime float64
    lastTime    time.Time
    statsMu     sync.Mutex
)

func GetGoMem() (heapAlloc, heapSys, goSys float64) {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    const mb = 1024 * 1024

    return float64(m.Alloc) / mb,
        float64(m.HeapSys) / mb,
        float64(m.Sys) / mb
}

func getPlatformMemDarwin() (rssMB, vmsMB float64) {
    var vsize C.uint64_t
    var rss C.uint64_t

    if C.get_mem(&vsize, &rss) != 0 {
        return 0, 0
    }

    const mb = 1024 * 1024
    return float64(rss) / mb, float64(vsize) / mb
}

func getPlatformMemLinux() (rssMB, vmsMB float64) {
    var data, readErr = os.ReadFile("/proc/self/statm")
    if readErr != nil {
        return 0, 0
    }

    var fields []string = strings.Fields(string(data))
    if len(fields) < 2 {
        return 0, 0
    }

    var pageSize float64 = float64(os.Getpagesize())
    var rssPages, rssParseErr = strconv.ParseFloat(fields[1], 64)
    if rssParseErr != nil {
        return 0, 0
    }
    var vmsPages, vmsParseErr = strconv.ParseFloat(fields[0], 64)
    if vmsParseErr != nil {
        return 0, 0
    }

    const mb = 1024 * 1024

    return (rssPages * pageSize) / mb,
        (vmsPages * pageSize) / mb
}

func GetPlatformMem() (rssMB, vmsMB float64) {
    switch runtime.GOOS {
    case "darwin":
        return getPlatformMemDarwin()
    case "linux":
        return getPlatformMemLinux()
    default:
        return 0, 0
    }
}

func getProcessCPUTime() float64 {
    var r syscall.Rusage
    syscall.Getrusage(syscall.RUSAGE_SELF, &r)

    sec := float64(r.Utime.Sec) + float64(r.Stime.Sec)
    usec := float64(r.Utime.Usec+r.Stime.Usec) / 1e6

    return sec + usec
}

func GetCPUPercent() float64 {
    statsMu.Lock()
    defer statsMu.Unlock()

    var now time.Time = time.Now()
    var cpuTime float64 = getProcessCPUTime()

    if lastTime.IsZero() {
        lastTime = now
        lastCPUTime = cpuTime
        return 0
    }

    var deltaCPU float64 = cpuTime - lastCPUTime
    var deltaTime float64 = now.Sub(lastTime).Seconds()

    lastCPUTime = cpuTime
    lastTime = now

    if deltaTime < 0.5 {
        return 0
    }

    return (deltaCPU / deltaTime) / float64(runtime.NumCPU()) * 100
}

func HandleStats(w http.ResponseWriter, r *http.Request) {
    var rss, vms = GetPlatformMem()
    var heapAlloc, heapSys, goSys = GetGoMem()

    w.Header().Set("Content-Type", "application/json")

    var t, err = time.Parse(os.Getenv("UPTIME"), os.Getenv("UPTIME"))
    if err != nil {
        HttpError(w, err, http.StatusInternalServerError)
    }

    if err := json.NewEncoder(w).Encode(Stats{
        RSSMB:        rss,
        VMSMB:        vms,
        HeapAllocMB:  heapAlloc,
        HeapSysMB:    heapSys,
        GoSysMB:      goSys,
        CPUPercent:   GetCPUPercent(),
        NumGoroutine: runtime.NumGoroutine(),
        NumCPU:       runtime.NumCPU(),
        UptimeSec:    time.Since(t).Seconds(),
        Platform:     runtime.GOOS,
    }); err != nil {
        HttpError(w, err, http.StatusInternalServerError)
    }
}

