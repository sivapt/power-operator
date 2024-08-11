package main

import (
    "context"
    "fmt"
    "log"

    "google.golang.org/grpc"
    podresourcesapi "k8s.io/kubelet/pkg/apis/podresources/v1"
)

const (
    socketPath = "/var/lib/kubelet/pod-resources/kubelet.sock"
)

func main() {
    conn, err := grpc.Dial(
        fmt.Sprintf("unix://%s", socketPath),
        grpc.WithInsecure(),
        grpc.WithBlock(),
    )
    if err != nil {
        log.Fatalf("Failed to connect to Kubelet socket: %v", err)
    }
    defer conn.Close()

    client := podresourcesapi.NewPodResourcesListerClient(conn)

    listPodResources(client)
}

func listPodResources(client podresourcesapi.PodResourcesListerClient) {
    req := &podresourcesapi.ListPodResourcesRequest{}
    resp, err := client.List(context.Background(), req)
    if err != nil {
        log.Fatalf("Failed to list pod resources: %v", err)
    }

    for _, pod := range resp.PodResources {
        fmt.Printf("Pod Name: %s, Namespace: %s\n", pod.Name, pod.Namespace)
        for _, container := range pod.Containers {
            fmt.Printf("  Container Name: %s\n", container.Name)
            for _, device := range container.Devices {
                fmt.Printf("    Resource Name: %s, Device IDs: %v\n", device.ResourceName, device.DeviceIds)
            }
            if len(container.CpuIds) > 0 {
                fmt.Printf("    CPU IDs: %v\n", container.CpuIds)
            } else {
                fmt.Printf("    CPU IDs: NULL\n")
            }
        }
    }
}
