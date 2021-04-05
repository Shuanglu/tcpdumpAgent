package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	//"syscall"

	"github.com/Shuanglu/containerTcpdump/pkg/containerdtcpdump"
	"github.com/Shuanglu/containerTcpdump/pkg/dockertcpdump"
	gopacket "github.com/Shuanglu/containerTcpdump/pkg/tcpdump"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"github.com/vishvananda/netns"
)

const (
	CLONE_NEWUTS  = 0x04000000   /* New utsname group? */
	CLONE_NEWIPC  = 0x08000000   /* New ipcs */
	CLONE_NEWUSER = 0x10000000   /* New user namespace */
	CLONE_NEWPID  = 0x20000000   /* New pid namespace */
	CLONE_NEWNET  = 0x40000000   /* New network namespace */
	CLONE_IO      = 0x80000000   /* Get io context */
	CLONE_NEWNS   = 0x20000      /* New mount namespace */
	bindMountPath = "/run/netns" /* Bind mount path for named netns */
)

type containerOperation interface {
	RetrieveNetns(UID string) (string, error)
}

type targetContainer struct {
	UID string `json:"UID"`
}

type targetContainers struct {
	Containers      []targetContainer `json:"Containers"`
	Runtime         string            `json:"Runtime,omitempty"`
	RuntimeEndpoint string            `json:"RuntimeEndpoint,omitempty"`
	Duration        int               `json:"Duration"`
}

func tcpdump(workerGroup *sync.WaitGroup, podNetns string, podId string, duration int) error {
	defer workerGroup.Done()
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	//tContainerJson, err := cli.ContainerInspect(ctx, container.ID)
	//tContainerJson.State.Pid
	//nsHandle, _ := netns.GetFromPid(containerPid)
	log.Info(fmt.Sprintf("%s", podNetns))
	podNetns = strings.TrimSuffix(podNetns, "\n")
	podNetns = strings.TrimSuffix(podNetns, "\"")
	podNetns = strings.TrimPrefix(podNetns, "\"")
	log.Info(fmt.Sprintf("%s", podNetns))
	if podNetns != "" {
		nsHandle, err := netns.GetFromPath(podNetns)
		if err != nil {
			log.Warn(fmt.Sprintf("Failed to get ns handle for pod %q due to %q", podId, err))
			return err
		}
		netns.Set(nsHandle)
		log.Info(fmt.Sprintf("Entering the network namespace: %q", podNetns))

	}
	path := "/tmp/" + podId + ".cap"
	err := gopacket.Capture(path, duration)
	if err != nil {
		log.Warn("There was error while capturing the requests of pod %q", podId)
	}
	return nil
}

//err := os.Remove("/tmp/" + tContainerName + ".cap")
/*cmd := exec.Command("/bin/bash", "-c", "timeout "+duration+" tcpdump -i any -w /tmp/"+podId+".cap")
	err = cmd.Start()
	log.Info("Tcpdump command starts")
	err = cmd.Wait()
	if err != nil {
		if fmt.Sprintf("%s", err) == "exit status 124" {
			log.Info(fmt.Sprintf("tcpdump command for the container %q completed successfully", podId))
		} else {
			log.Info(fmt.Sprintf("tcpdump command for the container %q didn't complete successfully due to %q", podId, err))
			return err
		}
	}
	defer runtime.UnlockOSThread()
	//originPid := os.Getpid()
	cmd := exec.Command("/bin/bash", "-c", "nsenter -t 1 -m -u")
	//err := cmd.Run()
	//if err != nil {
	//	log.Info(fmt.Sprintf("Failed to switch to the init namespace due to %q", err))
	//} else {
	//	log.Info("Succesfully switch to the init namespace")
	//}
	//originNS, _ := netns.GetFromPath("/proc/" + strconv.Itoa(originPid) + "/ns/pid")
	//defer originNS.Close()
	//nsPIDFD, err := unix.Open("/proc/1/ns/pid", unix.O_RDONLY|unix.O_CLOEXEC, 0)
	//nsUTSFD, err := unix.Open("/proc/1/ns/uts", unix.O_RDONLY|unix.O_CLOEXEC, 0)
	//nsMNTFD, err := unix.Open("/proc/1/ns/mnt", unix.O_RDONLY|unix.O_CLOEXEC, 0)
	//nsNETFD, err := unix.Open("/proc/1/ns/net", unix.O_RDONLY|unix.O_CLOEXEC, 0)
	//nsHandle, err := netns.GetFromPath("/proc/1/ns/pid")
	//log.Info(strconv.Itoa(int(nsUTSFD)))
	//err = unix.Setns(nsPIDFD, CLONE_NEWPID)
	//if err != nil {
	//	log.Fatal(fmt.Sprintf("Failed to set the 'PID' namespace due to %q", err))
	//}
	//err = unix.Setns(nsUTSFD, CLONE_NEWUTS)
	//if err != nil {
	//	log.Fatal(fmt.Sprintf("Failed to set the 'UTS' namespace due to %q", err))
	//}
	//err = unix.Setns(nsNETFD, CLONE_NEWNET)
	//hostname, _ := os.Hostname()
	//if err != nil {
	//	log.Fatal(fmt.Sprintf("Failed to set the 'net' namespace due to %q", err))
	//}
	//err = unix.Setns(nsMNTFD, CLONE_NEWNS)
	//hostname, _ := os.Hostname()
	//if err != nil {
	//	log.Info(hostname)
	//	log.Fatal(fmt.Sprintf("Failed to set the 'mount' namespace due to %q", err))
	//}

	//log.Info(fmt.Sprintf(strconv.Itoa(unix.Getpid())))

	//if err != nil {
	//	log.Fatal("Failed to enter the namespace of the pid 1 on node %q", hostname)
	//}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to initialize the docker client on node %q. EXIT now.", hostname))
	}
	containers, err := cli.ContainerList(context.TODO(), types.ContainerListOptions{})
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to list running docker containers on node %q due to error %q. EXIT now.", hostname, err))
	}
	containerPids := make(map[string]int)
	for _, container := range containers {
		containerInspect, err := cli.ContainerInspect(context.TODO(), container.ID)
		if err != nil {
			log.Info(fmt.Sprintf("Failed to get detail of the container %q", container.ID))
		}
		containerPids[containerInspect.Name] = containerInspect.State.Pid
		//log.Info()
	}
	//unix.Setns(int(originNS), CLONE_NEWPID)
	//unix.Setns(int(originNS), CLONE_NEWNS)
	//cmd = exec.Command("/bin/bash", "-c", "exit")
	//cmd.Run()
	return containerPids, nil
}*/

func main() {
	var runtime string
	var runtimeEndpoint string
	//ctx := context.Background()
	log.SetFormatter(&log.JSONFormatter{})

	filePath := flag.String("parameter-file", "/mnt/containerTcpdump/containers.json", "path of the parameter file")
	file, err := ioutil.ReadFile(*filePath)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to load the parameter file(%q) on node due to %q", *filePath, err))
	}
	containerJson := targetContainers{}
	err = json.Unmarshal([]byte(file), &containerJson)
	tContainers := containerJson.Containers
	if containerJson.Runtime != "" {
		runtime = containerJson.Runtime
	}
	if containerJson.RuntimeEndpoint != "" {
		runtimeEndpoint = containerJson.RuntimeEndpoint
	}
	duration := containerJson.Duration
	err = json.Unmarshal([]byte(file), &duration)
	//tempDuration, _ := strconv.Atoi(duration)
	//tempDuration = tempDuration + 2
	//duration = strconv.Itoa(tempDuration)
	//duration = duration + 2

	//containerPids := make(map[string]int)
	//containerPids, err = getContainerPIDs()
	containerMap := make(map[string]string)
	switch runtime {
	case "docker":
		dockerClient := dockertcpdump.DockerRuntime{}
		for _, tContainer := range tContainers {
			tContainerNetns, err := dockerClient.RetrieveNetns(tContainer.UID)
			if err != nil {
				log.Warn(fmt.Sprintf("Failed to retrieve the PID of the container %q", tContainer.UID))
				continue
			} else {
				containerMap[tContainer.UID] = tContainerNetns
			}
		}
	case "containerd":
		containerdClient := containerdtcpdump.RuntimeClient{}
		containerdClient.RuntimeClientInit("unix:///run/containerd/containerd.sock")
		for _, tContainer := range tContainers {
			tContainerNetns, err := containerdClient.RetrieveNetns(tContainer.UID)
			if err != nil {
				log.Warn(fmt.Sprintf("Failed to retrieve the network namespace of the container %q", tContainer.UID))
				continue
			} else {
				containerMap[tContainer.UID] = tContainerNetns
			}
		}
	default:
		log.Fatal(fmt.Sprintf("Not supported runtime: %q", runtime))
	}
	log.Info(fmt.Sprintf("%s", containerMap))

	if len(containerMap) > 0 {
		var workerGroup sync.WaitGroup
		for podId, podNetns := range containerMap {
			workerGroup.Add(1)
			log.Info(fmt.Sprintf("network namespace: %q", podNetns))
			go tcpdump(&workerGroup, podNetns, podId, duration)
		}
		workerGroup.Wait()
	} else {
		log.Fatal("Failed to retrieve PIDs of all the containers")
	}
	cmd := exec.Command("touch", "/tmp/containerTcpdumpComplete")
	cmd.Run()

}