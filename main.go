package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"syscall"
)

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/ext" {
		fmt.Printf("current pid: %d\n", syscall.Getgid())

		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200M --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
		os.Exit(-1)
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWNET,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		fmt.Printf("%v\n", cmd.Process.Pid)

		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)
		os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"), []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
		os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)
	}

	cmd.Process.Wait()
}
